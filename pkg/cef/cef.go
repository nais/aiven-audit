package cef

import "fmt"

type CefRecord struct {
	// An integer which identifies the version of the CEF format. The current CEF version is `0`.
	version int8
	// Identify the device vendor
	device_vendor string
	// Identify device name
	device_product string
	// Identify the device version
	device_version string
	// Signature ID also known as _Device Event Class ID_ identifies the type of event reported.
	type_id string
	// Representing a human-readable and understandable description of the event.
	message string
	// Reflects the importance of the event.
	severity_id string
	// Event timestamp
	end int64
}

func (record CefRecord) cefString() string {
	return fmt.Sprintf(
		"CEF:%v|%v|%v|%v|%v|%v|%v|end=%v",
		record.version,
		record.device_vendor,
		record.device_product,
		record.device_version,
		record.type_id,
		record.message,
		record.severity_id,
		record.end)
}
