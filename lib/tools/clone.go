package tools

import (
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"strings"
)

type Registry struct {
	Repository  string `yaml:"repo"`
	Description string `yaml:"desc"`
}

const jettiInstallGitRepo = "https://github.com/snowmerak/jetti-install.git"

var tempDir string
var regDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	tempDir = homeDir + "/.jetti-cache/jetti-install"

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		panic(err)
	}

	regDir = tempDir + "/registry"
}

func ResetTempDir() error {
	if err := os.RemoveAll(tempDir); err != nil {
		return err
	}

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	return nil
}

func CloneIfNotExists() error {
	exists := true
	if _, err := git.PlainOpen(tempDir); err != nil {
		exists = false
	}

	if exists {
		return nil
	}

	if _, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: jettiInstallGitRepo,
	}); err != nil {
		return err
	}

	return nil
}

func GetRegistries() ([]string, error) {
	if err := CloneIfNotExists(); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(regDir)
	if err != nil {
		return nil, err
	}

	registries := make([]string, 0, len(files))
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".yaml") {
			registries = append(registries, strings.TrimSuffix(name, ".yaml"))
		}
	}

	return registries, nil
}

func GetRegistryInfo(registry string) (*Registry, error) {
	file, err := os.Open(regDir + "/" + registry + ".yaml")
	if err != nil {
		return nil, err
	}

	reg := new(Registry)

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(reg); err != nil {
		return nil, err
	}

	return reg, nil
}

func InstallRegistry(repository string, version string) error {
	cmd := exec.Command("go", "install", repository+"@"+version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
