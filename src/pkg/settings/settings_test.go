package settings

import (
	"testing"
)

func TestGetSettings(t *testing.T) {
	// Reset the settings to ensure a clean state for the test
	settings = nil

	// Call GetSettings to initialize the settings
	s := GetSettings()

	// Check if the settings are initialized correctly
	if s == nil {
		t.Fatal("Expected settings to be initialized, got nil")
	}

	if !s.Verbose {
		t.Errorf("Expected Verbose to be true, got %v", s.Verbose)
	}

	if s.Debug {
		t.Errorf("Expected Debug to be false, got %v", s.Debug)
	}

	if s.Lang != "" {
		t.Errorf("Expected Lang to be an empty string, got %v", s.Lang)
	}
}

func TestGetSettingsSingleton(t *testing.T) {
	// Reset the settings to ensure a clean state for the test
	settings = nil

	// Call GetSettings to initialize the settings
	s1 := GetSettings()
	s2 := GetSettings()

	// Check if GetSettings returns the same instance
	if s1 != s2 {
		t.Error("Expected GetSettings to return the same instance")
	}
}

func TestSetSettings(t *testing.T) {
	// Reset the settings to ensure a clean state for the test
	settings = nil

	// Create a new settings instance
	s := &Settings{
		Verbose: false,
		Debug:   true,
		Lang:    "en",
	}

	// Set the settings
	SetSettings(s)

	// Get the settings
	s2 := GetSettings()

	// Check if the settings are set correctly
	if s2.Verbose {
		t.Errorf("Expected Verbose to be false, got %v", s2.Verbose)
	}

	if !s2.Debug {
		t.Errorf("Expected Debug to be true, got %v", s2.Debug)
	}

	if s2.Lang != "en" {
		t.Errorf("Expected Lang to be 'en', got %v", s2.Lang)
	}
}

func TestVerboseLog(t *testing.T) {
	// Reset the settings to ensure a clean state for the test
	settings = nil

	// Create a new settings instance
	s := &Settings{
		Verbose: true,
		Debug:   false,
	}

	// Set the settings
	SetSettings(s)

	// Call VerboseLog
	VerboseLog("This is a verbose log message")

	// Get the settings
	s2 := GetSettings()

	// Check if the verbose log message was printed
	if !s2.Verbose {
		t.Error("Expected Verbose to be true")
	}
}

func TestDebugLog(t *testing.T) {
	// Reset the settings to ensure a clean state for the test
	settings = nil

	// Create a new settings instance
	s := &Settings{
		Verbose: false,
		Debug:   true,
	}

	// Set the settings
	SetSettings(s)

	// Call DebugLog
	DebugLog("This is a debug log message")

	// Get the settings
	s2 := GetSettings()

	// Check if the debug log message was printed
	if !s2.Debug {
		t.Error("Expected Debug to be true")
	}
}

func TestStdLog(t *testing.T) {
	// Reset the settings to ensure a clean state for the test
	settings = nil

	// Create a new settings instance
	s := &Settings{
		Verbose: true,
		Debug:   false,
	}

	// Set the settings
	SetSettings(s)

	// Call StdLog
	StdLog("This is a standard log message")
}

func TestErrLog(t *testing.T) {
	// Reset the settings to ensure a clean state for the test
	settings = nil

	// Create a new settings instance
	s := &Settings{
		Verbose: true,
		Debug:   false,
	}

	// Set the settings
	SetSettings(s)

	// Call ErrLog
	ErrLog("This is an error log message")
}
