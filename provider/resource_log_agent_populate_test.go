package provider

import (
	"testing"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
)

func TestUniqueAgantFail(t *testing.T) {
	name := "bob"
	a := []sdk.LogAgent{}
	a = append(a, sdk.LogAgent{Name: "bob"})
	a = append(a, sdk.LogAgent{Name: "bob"})
	if err := uniqueAgent(name, a); err == nil {
		t.Errorf("expected uniqueAgent() to return an error when given multiple agents with the same name")
	}
}

func TestUniqueAgantPass(t *testing.T) {
	name := "bob"
	a := []sdk.LogAgent{}
	a = append(a, sdk.LogAgent{Name: name})
	a = append(a, sdk.LogAgent{Name: "bob2"})

	if err := uniqueAgent(name, a); err != nil {
		t.Errorf("expected uniqueAgent() to return a nil error when given unique agents: " + err.Error())
	}
}
