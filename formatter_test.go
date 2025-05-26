package plogger_test

import (
	"testing"

	"github.com/lthphuw/plogger"
)

func TestFieldKeySettersAndReset(t *testing.T) {
	// Save original values to restore after test
	originalMsg := plogger.FieldKeyMsg
	originalLevel := plogger.FieldKeyLevel
	originalTime := plogger.FieldKeyTime
	originalFile := plogger.FieldKeyFileCaller
	originalFunc := plogger.FieldKeyFuncCaller
	originalLine := plogger.FieldKeyLineCaller

	// Restore default values at the end
	defer func() {
		plogger.SetFieldKeyMsg(originalMsg)
		plogger.SetFieldKeyLevel(originalLevel)
		plogger.SetFieldKeyTimestamp(originalTime)
		plogger.SetFieldKeyFileCaller(originalFile)
		plogger.SetFieldKeyFuncCaller(originalFunc)
		plogger.SetFieldKeyLineCaller(originalLine)
	}()

	// Set custom values
	plogger.SetFieldKeyMsg("custom_msg")
	plogger.SetFieldKeyLevel("custom_level")
	plogger.SetFieldKeyTimestamp("custom_time")
	plogger.SetFieldKeyFileCaller("custom_file")
	plogger.SetFieldKeyFuncCaller("custom_func")
	plogger.SetFieldKeyLineCaller("custom_line")

	if plogger.FieldKeyMsg != "custom_msg" {
		t.Errorf("FieldKeyMsg = %s, want %s", plogger.FieldKeyMsg, "custom_msg")
	}
	if plogger.FieldKeyLevel != "custom_level" {
		t.Errorf("FieldKeyLevel = %s, want %s", plogger.FieldKeyLevel, "custom_level")
	}
	if plogger.FieldKeyTime != "custom_time" {
		t.Errorf("FieldKeyTime = %s, want %s", plogger.FieldKeyTime, "custom_time")
	}
	if plogger.FieldKeyFileCaller != "custom_file" {
		t.Errorf("FieldKeyFileCaller = %s, want %s", plogger.FieldKeyFileCaller, "custom_file")
	}
	if plogger.FieldKeyFuncCaller != "custom_func" {
		t.Errorf("FieldKeyFuncCaller = %s, want %s", plogger.FieldKeyFuncCaller, "custom_func")
	}
	if plogger.FieldKeyLineCaller != "custom_line" {
		t.Errorf("FieldKeyLineCaller = %s, want %s", plogger.FieldKeyLineCaller, "custom_line")
	}

	// Reset and test default values
	plogger.SetFieldKeyAsDefault()

	if plogger.FieldKeyMsg != plogger.DefaultFieldKeyMsg {
		t.Errorf("FieldKeyMsg = %s, want %s", plogger.FieldKeyMsg, plogger.DefaultFieldKeyMsg)
	}
	if plogger.FieldKeyLevel != plogger.DefaultFieldKeyLevel {
		t.Errorf("FieldKeyLevel = %s, want %s", plogger.FieldKeyLevel, plogger.DefaultFieldKeyLevel)
	}
	if plogger.FieldKeyTime != plogger.DefaultFieldKeyTime {
		t.Errorf("FieldKeyTime = %s, want %s", plogger.FieldKeyTime, plogger.DefaultFieldKeyTime)
	}
	if plogger.FieldKeyFileCaller != plogger.DefaultFieldKeyFileCaller {
		t.Errorf(
			"FieldKeyFileCaller = %s, want %s",
			plogger.FieldKeyFileCaller,
			plogger.DefaultFieldKeyFileCaller,
		)
	}
	if plogger.FieldKeyFuncCaller != plogger.DefaultFieldKeyFuncCaller {
		t.Errorf(
			"FieldKeyFuncCaller = %s, want %s",
			plogger.FieldKeyFuncCaller,
			plogger.DefaultFieldKeyFuncCaller,
		)
	}
	if plogger.FieldKeyLineCaller != plogger.DefaultFieldKeyLineCaller {
		t.Errorf(
			"FieldKeyLineCaller = %s, want %s",
			plogger.FieldKeyLineCaller,
			plogger.DefaultFieldKeyLineCaller,
		)
	}
}
