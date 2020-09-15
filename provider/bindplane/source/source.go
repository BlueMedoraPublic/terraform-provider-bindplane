package source

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"

	"github.com/google/uuid"
)

// Result describes the response from the bindplane source
// create api
type Result struct {
	SourceID string
	JobID    string
}

// Create attempts to create a source repeatedly until it is
// created succefully or timeout is exceeded
func Create(bp *sdk.BindPlane, source sdk.SourceConfigCreate, timeout int) (Result, error) {
	r := Result{}

	if err := source.Validate(); err != nil {
		return r, errors.Wrap(err, "Not attempting to create source, validation failed")
	}

	config, err := buildConfig(source)
	if err != nil {
		return r, errors.Wrap(err, "Not attempting to create source, buildConfig() failed")
	}

	exists, _, err := sourceExists(bp, source.Name, source.SourceType)
	if err != nil {
		return r, errors.Wrap(err, "Failed to check if source already exists")
	}
	if exists {
		return r, errors.New("Cannot create source, already exists")
	}

	resp, err := bp.CreateSource(config)
	if err != nil {
		return r, errors.Wrap(err, string(config))
	}
	r.JobID = resp.JobID

	id, err := checkJob(bp, int64(timeout), r.JobID, source)
	if err != nil {
		return r, err
	}
	r.SourceID = id
	return r, nil
}

func checkJob(bp *sdk.BindPlane, timeout int64, jobID string, config sdk.SourceConfigCreate) (string, error) {
	end := time.Now().Unix() + timeout
	for {
		job, err := bp.GetJob(jobID)
		if err != nil {
			return "", errors.Wrap(err, "Failed to call sdk.GetJob() with job id: "+jobID)
		}

		complete, status, err := parseStatus(job)
		if err != nil {
			return "", err
		}

		if complete {
			return getSourceID(bp, jobID)
		} else {
			// TEMP fallback check if the source exists, even if the job is not in a "complete" state
			exists, sourceID, err := sourceExists(bp, config.Name, config.SourceType)
			if err != nil {
				return "", err
			}
			if exists {
				return sourceID, nil
			}
		}

		if time.Now().Unix() > end {
			j, err := json.Marshal(job)
			if err != nil {
				return "", errors.New("Timed out waiting for job to complete. Last known status: "+status+". Job id: "+jobID)
			}
			return "", errors.New("Timed out waiting for job to complete. Last known status: "+status+". Job: "+string(j))
		}

		time.Sleep(time.Second * 5)
	}
}

// ParseStatus returns true if job complete and an error
// if job status failed or unexpected
func parseStatus(job sdk.Job) (bool, string, error) {
	status := strings.ToLower(job.Status)
	if status == "complete" {
		return true, status, nil
	} else if status == "in progress" {
		return false, status, nil
	} else if status == "testing connection to source" {
		return false, status, nil
	} else if status == "queued for completion" {
		return false, status, nil
	} else if status == "failed" {
		return false, status, errors.Wrap(jobErr(job), "job: "+job.ID+" failed. "+job.Message)

	}
	return false, status, errors.Wrap(jobErr(job), "ParseStatus() failed to parse job id "+job.ID)
}

/*
getSourceID returns a source uuid from a job. This is not safe
to call unless you know the job was a source create api call,
and the source has been created succefully.
*/
func getSourceID(bp *sdk.BindPlane, jobID string) (string, error) {
	j, err := bp.GetJob(jobID)
	if err != nil {
		return "", errors.Wrap(err, "Attempted to get source ID from Job with ID: "+jobID)
	}

	// not safe, could fail if the source was not created
	id := j.Result.(map[string]interface{})["id"].(string)

	if _, err := uuid.Parse(id); err != nil {
		msg := "job id " + id + " is not a valid uuid. This is likey an issue with the provider or BindPlane. Please file an issue on Github."
		return "", errors.Wrap(err, msg)

	}
	return id, nil
}

func buildConfig(source sdk.SourceConfigCreate) ([]byte, error) {
	x, err := json.Marshal(source)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert source struct to []byte with json.Marshal")
	}
	return x, err
}

func jobErr(job sdk.Job) error {
	s := job.Status
	m := job.Message
	r, ok := job.Result.(string)
	if ok != true {
		return errors.New("status: " + s + " message: " + m)
	}
	return errors.New("status: " + s + " message: " + m + " result: " + r)
}

func sourceExists(bp *sdk.BindPlane, name, sourceType string) (bool, string, error) {
	sources, err := bp.GetSources()
	if err != nil {
		return false, "", err
	}

	for _, s := range sources {
		if sourceType == s.SourceType.ID {
			if name == s.Name {
				return true, s.ID, nil
			}
		}
	}
	return false, "", nil
}
