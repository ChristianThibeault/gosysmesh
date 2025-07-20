package remote

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// RunSSHCommandOpenSSH executes a command over SSH via the OpenSSH client with ProxyJump support
func RunSSHCommandOpenSSH(user, host string, port int, keyPath, proxyJump, command string) (string, error) {
	// Input validation
	if err := validateSSHParams(user, host, keyPath); err != nil {
		return "", fmt.Errorf("invalid SSH parameters: %w", err)
	}

	keyPath = expandTilde(os.ExpandEnv(keyPath))
	args := []string{
		"-i", keyPath,
		"-p", fmt.Sprintf("%d", port),
		fmt.Sprintf("%s@%s", user, host),
		"-o", "ConnectTimeout=10",
		"-o", "ServerAliveInterval=30",
		"-o", "ServerAliveCountMax=3",
		// Enable host key checking for security
		"-o", "StrictHostKeyChecking=yes",
	}

	if proxyJump != "" {
		if err := validateHostname(proxyJump); err != nil {
			return "", fmt.Errorf("invalid proxy jump host: %w", err)
		}
		// Use configurable port instead of hardcoded 5822
		proxy := fmt.Sprintf("%s@%s", user, proxyJump)
		args = append([]string{"-J", proxy}, args...)
	}

	// Validate and sanitize command before execution
	if err := validateCommand(command); err != nil {
		return "", fmt.Errorf("invalid command: %w", err)
	}
	args = append(args, command)

	cmd := exec.Command("ssh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ssh error: %v — stderr: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

// expandTilde replaces ~ with home directory
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return home + path[1:]
		}
	}
	return path
}

// validateSSHParams validates SSH connection parameters
func validateSSHParams(user, host, keyPath string) error {
	if user == "" {
		return errors.New("user cannot be empty")
	}
	if host == "" {
		return errors.New("host cannot be empty")
	}
	if keyPath == "" {
		return errors.New("SSH key path cannot be empty")
	}

	// Validate username format (alphanumeric, underscore, dash)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(user) {
		return errors.New("invalid username format")
	}

	// Validate hostname/IP format
	if err := validateHostname(host); err != nil {
		return fmt.Errorf("invalid host: %w", err)
	}

	return nil
}

// validateHostname validates hostname or IP address format
func validateHostname(host string) error {
	if host == "" {
		return errors.New("hostname cannot be empty")
	}

	// Allow hostname format (RFC 1123) or IP address
	hostnameRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	ipRegex := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)

	if !hostnameRegex.MatchString(host) && !ipRegex.MatchString(host) {
		return errors.New("invalid hostname or IP address format")
	}

	return nil
}

// validateCommand validates command for basic security
func validateCommand(command string) error {
	if command == "" {
		return errors.New("command cannot be empty")
	}

	// Block dangerous injection patterns while allowing legitimate shell operators
	dangerousPatterns := []string{
		"$(", "`", "\n", "\r", ";rm", ";wget", ";curl", ";sh", ";bash",
		"&&rm", "&&wget", "&&curl", "&&sh", "&&bash",
		"||rm", "||wget", "||curl", "||sh", "||bash",
	}

	commandLower := strings.ToLower(command)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(commandLower, pattern) {
			return fmt.Errorf("command contains dangerous pattern: %s", pattern)
		}
	}

	return nil
}

// BuildPsCommand safely constructs a ps command with validated user parameter
func BuildPsCommand(user string) (string, error) {
	if err := validateUsername(user); err != nil {
		return "", fmt.Errorf("invalid user for ps command: %w", err)
	}
	return fmt.Sprintf("ps -u %s -o pid,user,%%cpu,%%mem,stat,lstart,args --no-headers", user), nil
}

// BuildSystemStatsCommand returns the pre-approved system stats command
func BuildSystemStatsCommand() string {
	return `top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/%us,//'; free -m | awk 'NR==2{printf "%.0f %.0f", $3,$2}'; df -h / | awk 'NR==2{gsub(/[^0-9.]/, "", $3); gsub(/[^0-9.]/, "", $2); printf " %.1f %.1f", $3, $2}'`
}

// validateUsername validates username for command construction
func validateUsername(user string) error {
	if user == "" {
		return errors.New("username cannot be empty")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(user) {
		return errors.New("invalid username format")
	}
	return nil
}

