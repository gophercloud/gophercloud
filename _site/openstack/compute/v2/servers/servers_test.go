package servers

import (
	"reflect"
	"testing"
)

// This provides more fine-grained failures when Servers differ, because Server structs are too damn big to compare by eye.
// FIXME I should absolutely refactor this into a general-purpose thing in testhelper.
func equalServers(t *testing.T, expected Server, actual Server) {
	if expected.ID != actual.ID {
		t.Errorf("ID differs. expected=[%s], actual=[%s]", expected.ID, actual.ID)
	}
	if expected.TenantID != actual.TenantID {
		t.Errorf("TenantID differs. expected=[%s], actual=[%s]", expected.TenantID, actual.TenantID)
	}
	if expected.UserID != actual.UserID {
		t.Errorf("UserID differs. expected=[%s], actual=[%s]", expected.UserID, actual.UserID)
	}
	if expected.Name != actual.Name {
		t.Errorf("Name differs. expected=[%s], actual=[%s]", expected.Name, actual.Name)
	}
	if expected.Updated != actual.Updated {
		t.Errorf("Updated differs. expected=[%s], actual=[%s]", expected.Updated, actual.Updated)
	}
	if expected.Created != actual.Created {
		t.Errorf("Created differs. expected=[%s], actual=[%s]", expected.Created, actual.Created)
	}
	if expected.HostID != actual.HostID {
		t.Errorf("HostID differs. expected=[%s], actual=[%s]", expected.HostID, actual.HostID)
	}
	if expected.Status != actual.Status {
		t.Errorf("Status differs. expected=[%s], actual=[%s]", expected.Status, actual.Status)
	}
	if expected.Progress != actual.Progress {
		t.Errorf("Progress differs. expected=[%s], actual=[%s]", expected.Progress, actual.Progress)
	}
	if expected.AccessIPv4 != actual.AccessIPv4 {
		t.Errorf("AccessIPv4 differs. expected=[%s], actual=[%s]", expected.AccessIPv4, actual.AccessIPv4)
	}
	if expected.AccessIPv6 != actual.AccessIPv6 {
		t.Errorf("AccessIPv6 differs. expected=[%s], actual=[%s]", expected.AccessIPv6, actual.AccessIPv6)
	}
	if !reflect.DeepEqual(expected.Image, actual.Image) {
		t.Errorf("Image differs. expected=[%s], actual=[%s]", expected.Image, actual.Image)
	}
	if !reflect.DeepEqual(expected.Flavor, actual.Flavor) {
		t.Errorf("Flavor differs. expected=[%s], actual=[%s]", expected.Flavor, actual.Flavor)
	}
	if !reflect.DeepEqual(expected.Addresses, actual.Addresses) {
		t.Errorf("Addresses differ. expected=[%s], actual=[%s]", expected.Addresses, actual.Addresses)
	}
	if !reflect.DeepEqual(expected.Metadata, actual.Metadata) {
		t.Errorf("Metadata differs. expected=[%s], actual=[%s]", expected.Metadata, actual.Metadata)
	}
	if !reflect.DeepEqual(expected.Links, actual.Links) {
		t.Errorf("Links differs. expected=[%s], actual=[%s]", expected.Links, actual.Links)
	}
	if expected.KeyName != actual.KeyName {
		t.Errorf("KeyName differs. expected=[%s], actual=[%s]", expected.KeyName, actual.KeyName)
	}
	if expected.AdminPass != actual.AdminPass {
		t.Errorf("AdminPass differs. expected=[%s], actual=[%s]", expected.AdminPass, actual.AdminPass)
	}
}
