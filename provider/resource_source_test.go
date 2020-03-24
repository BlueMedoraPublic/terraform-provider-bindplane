package provider

import (
	"encoding/json"
	"testing"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

const fakeValidUUID = "abcdefAB-0123-4ABC-ab12-CDEF01234567"
const fakeValidJSON = "{\"key\":\"value\"}" // important that there are no spaces

func TestInitSource(t *testing.T) {
	d := makeSchema()
	x, err := initSource(d)
	if err != nil {
		t.Errorf("Expected initSource() to return nil, got: " + err.Error())
		return
	}

	if x.CollectionInterval != 20 {
		t.Errorf("Expected initSource() to return collection interval of 20")
	}

	if x.CollectorID != fakeValidUUID {
		t.Errorf("Expected initSource() to return collectorID " + fakeValidUUID)
	}

	if x.Credentials.Credentials != fakeValidUUID {
		t.Errorf("Expected initSource() to return credentialID " + fakeValidUUID)
	}

	if x.Name != "abc" {
		t.Errorf("Expected initSource() to return name 'abc'")
	}

	if x.SourceType != "abc" {
		t.Errorf("Expected initSource() to return source type 'abc'")
	}

	jsonBytes, err := json.Marshal(x.Configuration)
	if err != nil {
		t.Errorf("Expected configuration interface to marshal into json byte array")
	} else if string(jsonBytes) != fakeValidJSON {
		t.Errorf("Expected initSource() to return a configuration interface{} that can be converted back to the origonal json string. Got: " + string(jsonBytes))
	}
}

func TestConfDiff(t *testing.T) {
	d := makeSchema()
	s, err := makeinitSource()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	x, err := confDiff(d, s)
	if err != nil {
		t.Errorf("Expected confDiff to return a nil error, got: " + err.Error())
		return
	}

	if x != true {
		t.Errorf("Expected confDiff to return true, got false, despite the state and config struct being identical")
	}
}

func makeSchema() *schema.ResourceData {
	r := resourceSource()
	d := r.Data(nil)
	d.Set("collection_interval", 20)
	d.Set("collector_id", fakeValidUUID)
	d.Set("name", "abc")
	d.Set("source_type", "abc")
	d.Set("credential_id", fakeValidUUID)
	d.Set("configuration", fakeValidJSON)
	return d
}

func makeinitSource() (sdk.SourceConfigGet, error) {
	x := sdk.SourceConfigGet{}
	x.CollectionInterval = 20
	x.Collector.ID = fakeValidUUID
	x.Name = "abc"
	x.SourceType.Name = "abc"
	x.Credentials = append(x.Credentials, sdk.Credential{ID: fakeValidUUID})

	// error is not checked because it will be checked by
	// the function caller
	err := json.Unmarshal([]byte(fakeValidJSON), &x.Configuration)
	return x, err
}
