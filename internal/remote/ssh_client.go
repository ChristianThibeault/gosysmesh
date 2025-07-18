package remote

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunSSHCommandOpenSSH executes a command over SSH via the OpenSSH client with ProxyJump support
func RunSSHCommandOpenSSH(user, host string, port int, keyPath, proxyJump, command string) (string, error) {
	keyPath = expandTilde(os.ExpandEnv(keyPath))
	args := []string{
		"-i", keyPath,
		"-p", fmt.Sprintf("%d", port),
		fmt.Sprintf("%s@%s", user, host),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}

	if proxyJump != "" {
		args = append([]string{"-J", proxyJump}, args...)
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

