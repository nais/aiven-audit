package cef

import "testing"

func TestCefString(t *testing.T) {
	record := CefRecord{
		version:        0,
		device_vendor:  "Aiven",
		device_product: "dev-cluster",
		device_version: "1.0",
		type_id:        "service_update",
		message:        "Created ACL entry {'permission': 'readwrite', 'topic': 'test.devtopic', 'username': 'test.tester*'}",
		severity_id:    "INFO",
		end:            1625567291,
	}

	got := record.cefString()
	want := "CEF:0|Aiven|dev-cluster|1.0|service_update|Created ACL entry {'permission': 'readwrite', 'topic': 'test.devtopic', 'username': 'test.tester*'}|INFO|end=1625567291"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
