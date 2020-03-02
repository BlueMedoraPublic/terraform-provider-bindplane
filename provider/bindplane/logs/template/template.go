package template

import (
    "github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"
    "github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
)

// Create creates a bindplane log template and
// returns the id
func Create(t sdk.LogTemplateCreate) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", err
	}

	x, err := bp.CreateLogTemplate(t)
	return x.ID, err
}

// Read returns nil if the log template exists
func Read(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}

	_, err = bp.GetLogTemplate(id)
    return err
}

// Delete returns nil if the log template is deleted
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DeleteLogTemplate(id)
}
