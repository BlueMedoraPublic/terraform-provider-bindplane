package source

import (
    "github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"

    "github.com/pkg/errors"
)

// Create creates a bindplane log source config and
// returns the id
func Create(payload []byte) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", err
	}

	x, err := bp.CreateLogSourceConfig(payload)
	if err != nil {
		return "", err
	}

	return x.ID, nil
}

// Read returns nil if the log source exists
func Read(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}

	_, err = bp.GetLogSourceConfig(id)
    if err != nil {
        return errors.Wrap(err, "hello!!")
    }
	return nil
}

// Delete returns nil if the log source config is deleted
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DeleteLogSourceConfig(id)
}
