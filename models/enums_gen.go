// Code generated by enum-generate.go DO NOT EDIT.

package models

import "gogs.mikescher.com/BlackForestBytes/goext/langext"
import "gogs.mikescher.com/BlackForestBytes/goext/enums"

const ChecksumEnumGenerator = "767fbd1b356150f73246cf075ea9a22a8e485784495e87d15d2e5d6fc0fde0e9" // GoExtVersion: 0.0.291

// ================================ JobStatus ================================
//
// File:       jobExecution.go
// StringEnum: true
// DescrEnum:  false
//

var __JobStatusValues = []JobStatus{
	JobStatusRunning,
	JobStatusSuccess,
	JobStatusFailed,
}

var __JobStatusVarnames = map[JobStatus]string{
	JobStatusRunning: "JobStatusRunning",
	JobStatusSuccess: "JobStatusSuccess",
	JobStatusFailed:  "JobStatusFailed",
}

func (e JobStatus) Valid() bool {
	return langext.InArray(e, __JobStatusValues)
}

func (e JobStatus) Values() []JobStatus {
	return __JobStatusValues
}

func (e JobStatus) ValuesAny() []any {
	return langext.ArrCastToAny(__JobStatusValues)
}

func (e JobStatus) ValuesMeta() []enums.EnumMetaValue {
	return JobStatusValuesMeta()
}

func (e JobStatus) String() string {
	return string(e)
}

func (e JobStatus) VarName() string {
	if d, ok := __JobStatusVarnames[e]; ok {
		return d
	}
	return ""
}

func (e JobStatus) Meta() enums.EnumMetaValue {
	return enums.EnumMetaValue{VarName: e.VarName(), Value: e, Description: nil}
}

func ParseJobStatus(vv string) (JobStatus, bool) {
	for _, ev := range __JobStatusValues {
		if string(ev) == vv {
			return ev, true
		}
	}
	return "", false
}

func JobStatusValues() []JobStatus {
	return __JobStatusValues
}

func JobStatusValuesMeta() []enums.EnumMetaValue {
	return []enums.EnumMetaValue{
		JobStatusRunning.Meta(),
		JobStatusSuccess.Meta(),
		JobStatusFailed.Meta(),
	}
}

// ================================ JobLogLevel ================================
//
// File:       jobLog.go
// StringEnum: true
// DescrEnum:  false
//

var __JobLogLevelValues = []JobLogLevel{
	JobLogLevelDebug,
	JobLogLevelInfo,
	JobLogLevelWarn,
	JobLogLevelError,
	JobLogLevelFatal,
}

var __JobLogLevelVarnames = map[JobLogLevel]string{
	JobLogLevelDebug: "JobLogLevelDebug",
	JobLogLevelInfo:  "JobLogLevelInfo",
	JobLogLevelWarn:  "JobLogLevelWarn",
	JobLogLevelError: "JobLogLevelError",
	JobLogLevelFatal: "JobLogLevelFatal",
}

func (e JobLogLevel) Valid() bool {
	return langext.InArray(e, __JobLogLevelValues)
}

func (e JobLogLevel) Values() []JobLogLevel {
	return __JobLogLevelValues
}

func (e JobLogLevel) ValuesAny() []any {
	return langext.ArrCastToAny(__JobLogLevelValues)
}

func (e JobLogLevel) ValuesMeta() []enums.EnumMetaValue {
	return JobLogLevelValuesMeta()
}

func (e JobLogLevel) String() string {
	return string(e)
}

func (e JobLogLevel) VarName() string {
	if d, ok := __JobLogLevelVarnames[e]; ok {
		return d
	}
	return ""
}

func (e JobLogLevel) Meta() enums.EnumMetaValue {
	return enums.EnumMetaValue{VarName: e.VarName(), Value: e, Description: nil}
}

func ParseJobLogLevel(vv string) (JobLogLevel, bool) {
	for _, ev := range __JobLogLevelValues {
		if string(ev) == vv {
			return ev, true
		}
	}
	return "", false
}

func JobLogLevelValues() []JobLogLevel {
	return __JobLogLevelValues
}

func JobLogLevelValuesMeta() []enums.EnumMetaValue {
	return []enums.EnumMetaValue{
		JobLogLevelDebug.Meta(),
		JobLogLevelInfo.Meta(),
		JobLogLevelWarn.Meta(),
		JobLogLevelError.Meta(),
		JobLogLevelFatal.Meta(),
	}
}