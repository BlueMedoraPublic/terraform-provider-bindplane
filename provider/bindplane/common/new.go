package common

import (
	"github.com/pkg/errors"
	"os"

	_ "github.com/BlueMedoraPublic/bpcli/bindplane/api" // force go mod to download api package
	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/BlueMedoraPublic/bpcli/util/uuid"
)

// New returns an initilized sdk.BindPlae object
func New() (*sdk.BindPlane, error) {
	if err := checkEnv(); err != nil {
		return nil, errors.Wrap(err, "Not attempting to initilize the bindplane sdk")
	}

	bp := new(sdk.BindPlane)
	bp.Init()

	if err := testConnection(bp); err != nil {
		return nil, err
	}

	return bp, nil
}

/*
checkEnv validates the environment variable BINDPLANE_API_KEY
removing this code would allow ~/.bpcli credentials file to be used,
however we currently want to enforce an environment variable in
order to avoid mistakes
*/
func checkEnv() error {
	apiKey := os.Getenv("BINDPLANE_API_KEY")
	if len(apiKey) == 0 {
		return errors.New("required environment variable BINDPLANE_API_KEY is not set")
	}

	if uuid.IsUUID(apiKey) == false {
		return errors.New("required environment variable BINDPLANE_API_KEY is set but does not appear to be a uuid")
	}

	return nil
}

func testConnection(bp *sdk.BindPlane) error {
	if _, err := bp.ListJobs(); err != nil {
		return errors.Wrap(err, "Test connection failed, could not list jobs. Is the API key correct?")
	}
	return nil
}
