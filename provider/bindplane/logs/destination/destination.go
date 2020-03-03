package destination

import (
    "github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"
    "github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
)

// Create creates a bindplane log destination config and
// returns the id
func Create(config sdk.LogDestConfig) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", err
	}

	x, err := bp.CreateLogDestConfig(config)
	if err != nil {
		return "", err
	}

	return x.ID, nil
}

// Read returns nil if the log destination exists
func Read(id string) (sdk.LogDestConfig, error) {
	bp, err := common.New()
	if err != nil {
		return sdk.LogDestConfig{}, err
	}

	return bp.GetLogDestConfig(id)
}

// Delete returns nil if the log destination config is deleted
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DelLogDestConfig(id)
}
