package provider

import (
    "testing"
)

func TestDataSourceAgentInstallCMDRead(t *testing.T) {
    if _, err := installCmd("all"); err == nil {
        t.Errorf("expected installCmD('all') to return an error, got nil")
    }
}
