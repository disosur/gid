package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

// SetGitUser sets git global user.name and user.email.
func SetGitUser(name, email string) error {
	if err := gitConfig("user.name", name); err != nil {
		return err
	}
	return gitConfig("user.email", email)
}

// SetGitSSH sets git global core.sshCommand to use the given key.
func SetGitSSH(keyPath string) error {
	expanded := expandTilde(keyPath)
	cmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", expanded)
	return gitConfig("core.sshCommand", cmd)
}

// GetGitEmail reads the current git global user.email.
func GetGitEmail() (string, error) {
	out, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GetGitName reads the current git global user.name.
func GetGitName() (string, error) {
	out, err := exec.Command("git", "config", "--global", "user.name").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config --global %s: %s", key, string(out))
	}
	return nil
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := exec.Command("sh", "-c", "echo $HOME").Output()
		if err == nil {
			return strings.TrimSpace(string(home)) + path[1:]
		}
	}
	return path
}
