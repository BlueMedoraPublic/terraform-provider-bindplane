package source

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
)

const fakeValidUUID = "abcdefAB-0123-4ABC-ab12-CDEF01234567"

func TestJobErrEmpty(t *testing.T) {
	job := sdk.Job{}
	if err := jobErr(job); err == nil {
		t.Errorf("Expected jobErr to always return an error, got nil")
	}
}

func TestJobErrResult(t *testing.T) {
	job := sdk.Job{
		Status:  "some status",
		Message: "some message",
		Result:  nil,
	}
	err := jobErr(job)
	if strings.Contains(err.Error(), "result") == true {
		t.Errorf("Expected jobErr to not include a result if the Job object did not have a result")
	}

	job.Result = "some result"
	err = jobErr(job)
	if strings.Contains(err.Error(), "result") == false {
		t.Errorf("Expected jobErr to include a result")
	}
}

func TestParseStatus(t *testing.T) {
	var (
		j   sdk.Job
		x   bool
		err error
	)

	j.Status = "complete"
	x, err = parseStatus(j)
	if x != true {
		t.Errorf("Expected parseStatus() to return true when job status is 'complete'")
	}
	if err != nil {
		t.Errorf("Expected parseStatus() to return a nil error when job status is 'complete'")
	}

	j.Status = "in progress"
	x, err = parseStatus(j)
	if x != false {
		t.Errorf("Expected parseStatus() to return false when job status is 'in progress'")
	}
	if err != nil {
		t.Errorf("Expected parseStatus() to return a nil error when job status is 'in progress'")
	}

	j.Status = "testing connection to source"
	x, err = parseStatus(j)
	if x != false {
		t.Errorf("Expected parseStatus() to return false when job status is 'testing connection to source'")
	}
	if err != nil {
		t.Errorf("Expected parseStatus() to return a nil error when job status is 'testing connection to source'")
	}

	j.Status = "queued for completion"
	x, err = parseStatus(j)
	if x != false {
		t.Errorf("Expected parseStatus() to return false when job status is 'queued for completion'")
	}
	if err != nil {
		t.Errorf("Expected parseStatus() to return a nil error when job status is 'queued for completion'")
	}

	j.Status = "failed"
	x, err = parseStatus(j)
	if x != false {
		t.Errorf("Expected parseStatus() to return false when job status is 'failed'")
	}
	if err == nil {
		t.Errorf("Expected parseStatus() to return an error when job status is 'failed'")
	}

}

func TestBuildConfig(t *testing.T) {
	source := newTestSourceConfig()
	if err := source.Validate(); err != nil {
		t.Errorf("Expected validation to pass: " + err.Error())
		return
	}
	jsonBytes, err := buildConfig(source)
	if err != nil {
		t.Errorf("Expected buildConfig to return nil error, got:" + err.Error())
		return
	}

	s := newBlankSourceConfig()
	if err := json.Unmarshal(jsonBytes, &s); err != nil {
		t.Errorf("Expected to unmarshal json to sdk.SourceConfigCreate struct. The json was built from a valid struct. " + err.Error())
		return
	}

	if err := s.Validate(); err != nil {
		t.Errorf(err.Error())
		return
	}
}

func newTestSourceConfig() sdk.SourceConfigCreate {
	source := sdk.SourceConfigCreate{}
	source.CollectionInterval = 60
	source.CollectorID = fakeValidUUID
	source.Name = "test"
	source.SourceType = "postgresql"
	source.Credentials.Credentials = fakeValidUUID
	source.Configuration = []byte("")
	return source
}

func newBlankSourceConfig() sdk.SourceConfigCreate {
	source := sdk.SourceConfigCreate{}
	return source
}
