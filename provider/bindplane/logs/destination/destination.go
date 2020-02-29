package destination

import (
    "github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"

    "github.com/pkg/errors"
)

// Create creates a bindplane log destination config and
// returns the id
func Create(payload []byte) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", err
	}

	x, err := bp.CreateLogDestConfig(payload)
	if err != nil {
		return "", err
	}

	return x.ID, nil
}

// Read returns nil if the log destination exists
func Read(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}

	_, err = bp.GetLogDestConfig(id)
    if err != nil {
        return errors.Wrap(err, "hello!!")
    }
	return nil
}

// Delete returns nil if the log destination config is deleted
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DelLogDestConfig(id)
}
