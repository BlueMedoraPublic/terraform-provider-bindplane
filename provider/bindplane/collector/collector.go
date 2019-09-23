package collector

import (
	"github.com/pkg/errors"
	"time"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"
)

// Max time to wait for a collector to deploy in seconds
const collectorTimeout = int64(300)

/*
WaitForAPI checks the BindPlane API every ten seconds
until a collector with the name 'collectorName' is found. The
collector id is returned.
*/
func WaitForAPI(collectorName string) (string, error) {
	bp, err := common.New()
	if err != nil {
		return "", errors.Wrap(err, "Not attempting to import collector, failed to run sdk.bindplane.New()")
	}

	startTime := time.Now().Unix()

	for {
		collectors, err := bp.GetCollectors()
		if err != nil {
			return "", errors.Wrap(err, "sdk.GetCollectors() returned an error")
		}

		for _, collector := range collectors {
			if collector.Name == collectorName {
				return collector.ID, nil
			}
		}

		time.Sleep(10 * time.Second)
		timeCurrent := time.Now().Unix()
		if (timeCurrent - startTime) > collectorTimeout {
			msg := "Timeout exceeded for collector creation: " + "CollectorName: " + collectorName
			return "", errors.New(msg)
		}
	}
}

// Read returns an error if a collector does not exist
func Read(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}

	_, err = bp.GetCollector(id)
	return err
}

// Delete deletes a collector from the API
func Delete(id string) error {
	bp, err := common.New()
	if err != nil {
		return err
	}
	return bp.DeleteCollector(id)
}
