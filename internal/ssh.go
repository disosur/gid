package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GenerateSSHKey creates an ed25519 key at ~/.ssh/id_<alias>.
func GenerateSSHKey(alias, email string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	keyPath := filepath.Join(home, ".ssh", fmt.Sprintf("id_%s", alias))

	if _, err := os.Stat(keyPath); err == nil {
		return keyPath, fmt.Errorf("key already exists at %s", keyPath)
	}

	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", keyPath, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ssh-keygen failed: %w", err)
	}

	return keyPath, nil
}

// ReadPublicKey reads the .pub file for a given key path.
func ReadPublicKey(keyPath string) (string, error) {
	data, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// AddSSHHostBlock appends a Host block for the alias to ~/.ssh/config.
func AddSSHHostBlock(alias, keyPath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sshConfig := filepath.Join(home, ".ssh", "config")

	block := fmt.Sprintf(`
Host github.com-%s
    HostName github.com
    User git
    IdentityFile %s
    IdentitiesOnly yes
`, alias, keyPath)

	f, err := os.OpenFile(sshConfig, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(block)
	return err
}
