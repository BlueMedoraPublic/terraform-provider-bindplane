package credential

import (
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"
)

// Create creates a bindplane credential and returns the id
func Create(payload []byte) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", err
	}

	x, err := bp.CreateCredential(payload)
	if err != nil {
		return "", err
	}

	return x.ID, nil
}

// Read returns nil if the credential exists
func Read(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}

	_, err = bp.GetCredential(id)
	return err
}

// Delete returns nil if the credential is deleted
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DeleteCredential(id)
}
