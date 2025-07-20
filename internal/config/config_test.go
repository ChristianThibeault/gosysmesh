package config

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid username", "admin", false},
		{"valid with underscore", "admin_user", false},
		{"valid with dash", "admin-user", false},
		{"valid with numbers", "admin123", false},
		{"empty username", "", true},
		{"too long", "this_username_is_way_too_long_for_any_system", true},
		{"invalid characters", "admin@host", true},
		{"invalid space", "admin user", true},
		{"invalid special chars", "admin$user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		wantErr  bool
	}{
		{"valid hostname", "example.com", false},
		{"valid subdomain", "sub.example.com", false},
		{"valid IP", "192.168.1.1", false},
		{"valid single word", "localhost", false},
		{"empty hostname", "", true},
		{"invalid characters", "host_name", true},
		{"invalid hostname chars", "host@name", true},
		{"invalid with space", "host name", true},
		{"too long", "this-is-a-very-long-hostname-that-exceeds-the-maximum-allowed-length-for-a-hostname-which-should-be-rejected-by-our-validation-function-because-it-is-way-too-long-and-would-cause-problems-in-real-world-usage", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHostname(tt.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHostname() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantErr  bool
	}{
		{"valid path", "/home/user/.ssh/id_rsa", false},
		{"valid relative", "~/.ssh/id_rsa", false},
		{"empty path", "", true},
		{"path traversal", "/home/user/../../../etc/passwd", true},
		{"null byte", "/home/user/\x00malicious", true},
		{"newline", "/home/user/file\nname", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}