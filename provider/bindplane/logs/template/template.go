package template

import (
	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"
)

// Create creates a bindplane log template and
// returns the id
func Create(t sdk.LogTemplate) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", err
	}

	x, err := bp.CreateLogTemplate(t)
	return x.ID, err
}

// Read returns nil if the log template exists
func Read(id string) (sdk.LogTemplate, error) {
	bp, err := common.New()
	if err != nil {
		return sdk.LogTemplate{}, err
	}
	return bp.GetLogTemplate(id)
}

// Delete returns nil if the log template is deleted
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DeleteLogTemplate(id)
}
