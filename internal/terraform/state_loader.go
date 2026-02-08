package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// State represents a Terraform state
type State struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// Resource represents a Terraform resource
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Instances  []Instance             `json:"instances"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// Instance represents a resource instance
type Instance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// StateLoader loads Terraform state
type StateLoader struct {
	statePath string
}

// NewStateLoader creates a new state loader
func NewStateLoader(statePath string) *StateLoader {
	return &StateLoader{
		statePath: statePath,
	}
}

// LoadState loads the Terraform state from file
func (l *StateLoader) LoadState(ctx context.Context) (*State, error) {
	log.Debugf("Loading Terraform state from: %s", l.statePath)

	// Read state file
	data, err := os.ReadFile(l.statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	// Parse JSON
	var rawState map[string]interface{}
	if err := json.Unmarshal(data, &rawState); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	// Extract resources
	state := &State{
		Resources: []Resource{},
	}

	if version, ok := rawState["version"].(float64); ok {
		state.Version = int(version)
	}

	// Parse resources from state
	resources, ok := rawState["resources"].([]interface{})
	if !ok {
		log.Warn("No resources found in state file")
		return state, nil
	}

	for _, r := range resources {
		resourceMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		resource := Resource{
			Type:     getString(resourceMap, "type"),
			Name:     getString(resourceMap, "name"),
			Provider: getString(resourceMap, "provider"),
		}

		// Extract instances
		instances, ok := resourceMap["instances"].([]interface{})
		if ok && len(instances) > 0 {
			for _, inst := range instances {
				instMap, ok := inst.(map[string]interface{})
				if !ok {
					continue
				}

				attrs, ok := instMap["attributes"].(map[string]interface{})
				if ok {
					resource.Attributes = attrs
					resource.Instances = append(resource.Instances, Instance{
						Attributes: attrs,
					})
				}
			}
		}

		state.Resources = append(state.Resources, resource)
	}

	log.Infof("Loaded %d resources from Terraform state", len(state.Resources))
	return state, nil
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
