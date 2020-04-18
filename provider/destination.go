package provider

import (
	"encoding/json"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/pkg/errors"
)

func defaultLogDestVersion(id string) (string, error) {
	p, err := bp.GetLogDestTypeParameters(id)
	if err != nil {
		return "", err
	}

	c := sdk.LogDestConfig{}
	if err := json.Unmarshal(p, &c); err != nil {
		return "", errors.Wrap(err, "unable to get default version for destination type: "+id+". "+fileIssueErr)
	}
	return c.DestinationVersion, nil
}
