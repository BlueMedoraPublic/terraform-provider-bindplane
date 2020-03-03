package source

import (
	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"
)

// Create creates a bindplane log source config and
// returns the id
func Create(config sdk.LogSourceConfig) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", err
	}

	x, err := bp.CreateLogSourceConfig(config)
	if err != nil {
		return "", err
	}

	return x.ID, nil
}

// Read returns nil if the log source exists
func Read(id string) (sdk.LogSourceConfig, error) {
	bp, err := common.New()
	if err != nil {
		return sdk.LogSourceConfig{}, err
	}

	return bp.GetLogSourceConfig(id)
}

// Delete returns nil if the log source config is deleted
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DeleteLogSourceConfig(id)
}
