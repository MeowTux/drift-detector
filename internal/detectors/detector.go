package detectors

import (
	"context"

	"github.com/MeowTux/drift-detector/internal/drift"
	"github.com/MeowTux/drift-detector/internal/terraform"
)

// Detector is the interface for cloud provider drift detectors
type Detector interface {
	// Name returns the detector name
	Name() string
	
	// Detect performs drift detection on resources
	Detect(ctx context.Context, state *terraform.State) ([]drift.DriftItem, error)
}
