package detectors

import (
	"context"
	"fmt"

	"github.com/MeowTux/drift-detector/internal/drift"
	"github.com/MeowTux/drift-detector/internal/terraform"
	log "github.com/sirupsen/logrus"
)

// GCPDetector detects drift in GCP resources
type GCPDetector struct {
	projectID string
}

// NewGCPDetector creates a new GCP detector
func NewGCPDetector(projectID string) (*GCPDetector, error) {
	if projectID == "" {
		return nil, fmt.Errorf("GCP project ID is required")
	}

	return &GCPDetector{
		projectID: projectID,
	}, nil
}

// Name returns the detector name
func (d *GCPDetector) Name() string {
	return "GCP"
}

// Detect performs drift detection
func (d *GCPDetector) Detect(ctx context.Context, state *terraform.State) ([]drift.DriftItem, error) {
	var drifts []drift.DriftItem

	log.Debugf("Detecting drift in GCP resources for project: %s", d.projectID)

	for _, resource := range state.Resources {
		// Only check GCP resources
		if !isGCPResource(resource.Type) {
			continue
		}

		switch resource.Type {
		case "google_compute_instance":
			drift, err := d.checkComputeInstance(ctx, resource)
			if err != nil {
				log.Warnf("Error checking compute instance %s: %v", resource.Name, err)
				continue
			}
			if drift != nil {
				drifts = append(drifts, *drift)
			}

		case "google_storage_bucket":
			drift, err := d.checkStorageBucket(ctx, resource)
			if err != nil {
				log.Warnf("Error checking storage bucket %s: %v", resource.Name, err)
				continue
			}
			if drift != nil {
				drifts = append(drifts, *drift)
			}
		}
	}

	return drifts, nil
}

func (d *GCPDetector) checkComputeInstance(ctx context.Context, resource terraform.Resource) (*drift.DriftItem, error) {
	// Implementation for GCP compute instance drift detection
	// This is a placeholder - full implementation would use GCP SDK
	log.Debug("Checking GCP compute instance:", resource.Name)
	
	// In production, this would:
	// 1. Use compute.NewInstancesRESTClient()
	// 2. Get the instance details
	// 3. Compare with Terraform state
	// 4. Return drift if found
	
	return nil, nil
}

func (d *GCPDetector) checkStorageBucket(ctx context.Context, resource terraform.Resource) (*drift.DriftItem, error) {
	// Implementation for GCP storage bucket drift detection
	log.Debug("Checking GCP storage bucket:", resource.Name)
	
	return nil, nil
}

func isGCPResource(resourceType string) bool {
	return len(resourceType) > 7 && resourceType[:7] == "google_"
}
