package detectors

import (
	"context"
	"fmt"

	"github.com/MeowTux/drift-detector/internal/drift"
	"github.com/MeowTux/drift-detector/internal/terraform"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

// AWSDetector detects drift in AWS resources
type AWSDetector struct {
	regions   []string
	ec2Client *ec2.Client
	s3Client  *s3.Client
}

// NewAWSDetector creates a new AWS detector
func NewAWSDetector(regions []string) (*AWSDetector, error) {
	if len(regions) == 0 {
		regions = []string{"us-east-1"}
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &AWSDetector{
		regions:   regions,
		ec2Client: ec2.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
	}, nil
}

// Name returns the detector name
func (d *AWSDetector) Name() string {
	return "AWS"
}

// Detect performs drift detection
func (d *AWSDetector) Detect(ctx context.Context, state *terraform.State) ([]drift.DriftItem, error) {
	var drifts []drift.DriftItem

	log.Debugf("Detecting drift in AWS resources across %d regions", len(d.regions))

	for _, resource := range state.Resources {
		// Only check AWS resources
		if !isAWSResource(resource.Type) {
			continue
		}

		switch resource.Type {
		case "aws_instance":
			drift, err := d.checkEC2Instance(ctx, resource)
			if err != nil {
				log.Warnf("Error checking EC2 instance %s: %v", resource.Name, err)
				continue
			}
			if drift != nil {
				drifts = append(drifts, *drift)
			}

		case "aws_s3_bucket":
			drift, err := d.checkS3Bucket(ctx, resource)
			if err != nil {
				log.Warnf("Error checking S3 bucket %s: %v", resource.Name, err)
				continue
			}
			if drift != nil {
				drifts = append(drifts, *drift)
			}

		case "aws_security_group":
			drift, err := d.checkSecurityGroup(ctx, resource)
			if err != nil {
				log.Warnf("Error checking security group %s: %v", resource.Name, err)
				continue
			}
			if drift != nil {
				drifts = append(drifts, *drift)
			}
		}
	}

	return drifts, nil
}

func (d *AWSDetector) checkEC2Instance(ctx context.Context, resource terraform.Resource) (*drift.DriftItem, error) {
	instanceID, ok := resource.Attributes["id"].(string)
	if !ok {
		return nil, fmt.Errorf("instance ID not found")
	}

	// Describe the instance
	result, err := d.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return &drift.DriftItem{
			ResourceType: resource.Type,
			ResourceName: resource.Name,
			Provider:     "AWS",
			Severity:     "critical",
			Changes: []drift.Change{{
				Field:    "existence",
				Expected: "exists",
				Actual:   "deleted",
			}},
		}, nil
	}

	instance := result.Reservations[0].Instances[0]
	var changes []drift.Change

	// Check instance type
	expectedType, _ := resource.Attributes["instance_type"].(string)
	if string(instance.InstanceType) != expectedType {
		changes = append(changes, drift.Change{
			Field:    "instance_type",
			Expected: expectedType,
			Actual:   string(instance.InstanceType),
		})
	}

	// Check tags
	expectedTags, _ := resource.Attributes["tags"].(map[string]interface{})
	actualTags := make(map[string]string)
	for _, tag := range instance.Tags {
		if tag.Key != nil && tag.Value != nil {
			actualTags[*tag.Key] = *tag.Value
		}
	}

	for key, expectedVal := range expectedTags {
		if actualVal, ok := actualTags[key]; !ok || actualVal != expectedVal {
			changes = append(changes, drift.Change{
				Field:    fmt.Sprintf("tags.%s", key),
				Expected: expectedVal,
				Actual:   actualTags[key],
			})
		}
	}

	if len(changes) > 0 {
		return &drift.DriftItem{
			ResourceType: resource.Type,
			ResourceName: resource.Name,
			ResourceID:   instanceID,
			Provider:     "AWS",
			Severity:     determineSeverity(changes),
			Changes:      changes,
		}, nil
	}

	return nil, nil
}

func (d *AWSDetector) checkS3Bucket(ctx context.Context, resource terraform.Resource) (*drift.DriftItem, error) {
	bucketName, ok := resource.Attributes["bucket"].(string)
	if !ok {
		return nil, fmt.Errorf("bucket name not found")
	}

	var changes []drift.Change

	// Check bucket versioning
	versioning, err := d.s3Client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: &bucketName,
	})
	if err != nil {
		// Bucket might have been deleted
		return &drift.DriftItem{
			ResourceType: resource.Type,
			ResourceName: resource.Name,
			Provider:     "AWS",
			Severity:     "critical",
			Changes: []drift.Change{{
				Field:    "existence",
				Expected: "exists",
				Actual:   "deleted or inaccessible",
			}},
		}, nil
	}

	expectedVersioning, _ := resource.Attributes["versioning"].(map[string]interface{})
	if expectedVersioning != nil {
		enabled, _ := expectedVersioning["enabled"].(bool)
		if enabled && versioning.Status != "Enabled" {
			changes = append(changes, drift.Change{
				Field:    "versioning",
				Expected: "Enabled",
				Actual:   string(versioning.Status),
			})
		}
	}

	// Check encryption
	encryption, err := d.s3Client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
		Bucket: &bucketName,
	})

	expectedEncryption, _ := resource.Attributes["server_side_encryption_configuration"]
	if expectedEncryption != nil && err != nil {
		changes = append(changes, drift.Change{
			Field:    "encryption",
			Expected: "enabled",
			Actual:   "disabled",
		})
	} else if expectedEncryption == nil && encryption != nil {
		changes = append(changes, drift.Change{
			Field:    "encryption",
			Expected: "disabled",
			Actual:   "enabled",
		})
	}

	if len(changes) > 0 {
		return &drift.DriftItem{
			ResourceType: resource.Type,
			ResourceName: resource.Name,
			ResourceID:   bucketName,
			Provider:     "AWS",
			Severity:     determineSeverity(changes),
			Changes:      changes,
		}, nil
	}

	return nil, nil
}

func (d *AWSDetector) checkSecurityGroup(ctx context.Context, resource terraform.Resource) (*drift.DriftItem, error) {
	sgID, ok := resource.Attributes["id"].(string)
	if !ok {
		return nil, fmt.Errorf("security group ID not found")
	}

	result, err := d.ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{sgID},
	})
	if err != nil {
		return &drift.DriftItem{
			ResourceType: resource.Type,
			ResourceName: resource.Name,
			Provider:     "AWS",
			Severity:     "critical",
			Changes: []drift.Change{{
				Field:    "existence",
				Expected: "exists",
				Actual:   "deleted",
			}},
		}, nil
	}

	if len(result.SecurityGroups) == 0 {
		return &drift.DriftItem{
			ResourceType: resource.Type,
			ResourceName: resource.Name,
			Provider:     "AWS",
			Severity:     "critical",
			Changes: []drift.Change{{
				Field:    "existence",
				Expected: "exists",
				Actual:   "deleted",
			}},
		}, nil
	}

	sg := result.SecurityGroups[0]
	var changes []drift.Change

	// Compare ingress rules
	expectedIngress, _ := resource.Attributes["ingress"].([]interface{})
	if len(sg.IpPermissions) != len(expectedIngress) {
		changes = append(changes, drift.Change{
			Field:    "ingress_rules_count",
			Expected: len(expectedIngress),
			Actual:   len(sg.IpPermissions),
		})
	}

	// Compare egress rules
	expectedEgress, _ := resource.Attributes["egress"].([]interface{})
	if len(sg.IpPermissionsEgress) != len(expectedEgress) {
		changes = append(changes, drift.Change{
			Field:    "egress_rules_count",
			Expected: len(expectedEgress),
			Actual:   len(sg.IpPermissionsEgress),
		})
	}

	if len(changes) > 0 {
		return &drift.DriftItem{
			ResourceType: resource.Type,
			ResourceName: resource.Name,
			ResourceID:   sgID,
			Provider:     "AWS",
			Severity:     "high",
			Changes:      changes,
		}, nil
	}

	return nil, nil
}

func isAWSResource(resourceType string) bool {
	return len(resourceType) > 4 && resourceType[:4] == "aws_"
}

func determineSeverity(changes []drift.Change) string {
	for _, change := range changes {
		if change.Field == "existence" {
			return "critical"
		}
		if change.Field == "encryption" || change.Field == "public_access" {
			return "high"
		}
	}
	return "medium"
}
