package detectors

import (
	"context"
	"fmt"

	"github.com/MeowTux/drift-detector/internal/drift"
	"github.com/MeowTux/drift-detector/internal/terraform"
	log "github.com/sirupsen/logrus"
)

// AzureDetector detects drift in Azure resources
type AzureDetector struct {
	subscriptionID string
}

// NewAzureDetector creates a new Azure detector
func NewAzureDetector(subscriptionID string) (*AzureDetector, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("Azure subscription ID is required")
	}

	return &AzureDetector{
		subscriptionID: subscriptionID,
	}, nil
}

// Name returns the detector name
func (d *AzureDetector) Name() string {
	return "Azure"
}

// Detect performs drift detection
func (d *AzureDetector) Detect(ctx context.Context, state *terraform.State) ([]drift.DriftItem, error) {
	var drifts []drift.DriftItem

	log.Debugf("Detecting drift in Azure resources for subscription: %s", d.subscriptionID)

	for _, resource := range state.Resources {
		// Only check Azure resources
		if !isAzureResource(resource.Type) {
			continue
		}

		switch resource.Type {
		case "azurerm_virtual_machine":
			drift, err := d.checkVirtualMachine(ctx, resource)
			if err != nil {
				log.Warnf("Error checking virtual machine %s: %v", resource.Name, err)
				continue
			}
			if drift != nil {
				drifts = append(drifts, *drift)
			}

		case "azurerm_storage_account":
			drift, err := d.checkStorageAccount(ctx, resource)
			if err != nil {
				log.Warnf("Error checking storage account %s: %v", resource.Name, err)
				continue
			}
			if drift != nil {
				drifts = append(drifts, *drift)
			}
		}
	}

	return drifts, nil
}

func (d *AzureDetector) checkVirtualMachine(ctx context.Context, resource terraform.Resource) (*drift.DriftItem, error) {
	// Implementation for Azure VM drift detection
	log.Debug("Checking Azure virtual machine:", resource.Name)
	
	// In production, this would:
	// 1. Use Azure SDK to get VM details
	// 2. Compare with Terraform state
	// 3. Return drift if found
	
	return nil, nil
}

func (d *AzureDetector) checkStorageAccount(ctx context.Context, resource terraform.Resource) (*drift.DriftItem, error) {
	// Implementation for Azure storage account drift detection
	log.Debug("Checking Azure storage account:", resource.Name)
	
	return nil, nil
}

func isAzureResource(resourceType string) bool {
	return len(resourceType) > 8 && resourceType[:8] == "azurerm_"
}
