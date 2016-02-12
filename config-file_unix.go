// +build darwin freebsd linux netbsd openbsd

package nodester

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var homedir_cache string = ""


func configDir() (string, error) {
	dir, err := homeDir()
	if err != nil {
		return "", err
	}
    
	return filepath.Join(dir, ".nodester"), nil
}

func homeDir() (string, error) {
	if homedir_cache != "" {
		return homedir_cache, nil
	}
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		log.Printf("Detected home directory from env var: %s", home)
		homedir_cache = home
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output")
	}
	homedir_cache = result
	return result, nil
}
