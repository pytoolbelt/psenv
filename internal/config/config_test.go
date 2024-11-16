package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestFile(t *testing.T, filename, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	require.NoError(t, err)
}

func TestLoadProjectConfig(t *testing.T) {
	// Create the test file
	createTestFile(t, "psenv-project.yml", `
default: dev
environments:
  - base
  - dev
  - prod
  - test
prefix: /path/to/params
project: foobar
`)
	// Run the test
	_, err := LoadProjectConfig()
	require.NoError(t, err)

	// Clean up the generated test file
	RemoveTestFiles(t, "psenv-project.yml")
}

func RemoveTestFiles(t *testing.T, files ...string) {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			t.Errorf("Failed to remove file %s: %v", file, err)
		}
	}
}

func TestProjectConfig_PrintTable(t *testing.T) {
	projectConfig := &ProjectConfig{
		Default:      "dev",
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
	}
	projectConfig.PrintTable()
}

func TestProjectConfig_GetEnvironmentPath(t *testing.T) {
	projectConfig := &ProjectConfig{
		Default:      "dev",
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
	}
	path := projectConfig.GetEnvironmentPath("dev")
	require.Equal(t, "/path/to/params/foobar/dev", path)
}

func TestProjectConfig_GetBasePath(t *testing.T) {
	projectConfig := &ProjectConfig{
		Default:      "dev",
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
	}
	basePath := projectConfig.GetBasePath()
	require.Equal(t, "/path/to/params/foobar/", basePath)
}

func TestProjectConfig_HasEnvironment(t *testing.T) {
	projectConfig := &ProjectConfig{
		Default:      "dev",
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
	}
	require.True(t, projectConfig.HasEnvironment("dev"))
	require.False(t, projectConfig.HasEnvironment("staging"))
}

func TestProjectConfig_Save(t *testing.T) {
	projectConfig := &ProjectConfig{
		Default:      "dev",
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
	}
	err := projectConfig.Save()
	require.NoError(t, err)

	// Clean up the generated test file
	RemoveTestFiles(t, ProjectConfigFile)
}

func TestProjectConfig_RemoveEnvironment(t *testing.T) {
	projectConfig := &ProjectConfig{
		Default:      "dev",
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
	}
	projectConfig.RemoveEnvironment("dev")
	require.False(t, projectConfig.HasEnvironment("dev"))
}

func TestSecretsConfig_GetEnvironmentPath(t *testing.T) {
	secretsConfig := &SecretsConfig{
		Project: "foobar",
		Prefix:  "/path/to/params",
		Environments: map[string]map[string]string{
			"dev": {
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
	}
	path := secretsConfig.GetEnvironmentPath("dev")
	require.Equal(t, "/path/to/params/foobar/dev", path)
}

func TestSecretsConfig_GetEnvironmentParams(t *testing.T) {
	secretsConfig := &SecretsConfig{
		Project: "foobar",
		Prefix:  "/path/to/params",
		Environments: map[string]map[string]string{
			"dev": {
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
	}
	params := secretsConfig.GetEnvironmentParams("dev")
	expected := map[string]string{
		"/path/to/params/foobar/dev/KEY1": "value1",
		"/path/to/params/foobar/dev/KEY2": "value2",
	}
	require.Equal(t, expected, params)
}

func TestSecretsConfig_ClearEnvironment(t *testing.T) {
	secretsConfig := &SecretsConfig{
		Project: "foobar",
		Prefix:  "/path/to/params",
		Environments: map[string]map[string]string{
			"dev": {
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
	}
	secretsConfig.ClearEnvironment("dev")
	require.Empty(t, secretsConfig.Environments["dev"])
}

func TestSecretsConfig_ClearEnvironments(t *testing.T) {
	secretsConfig := &SecretsConfig{
		Project: "foobar",
		Prefix:  "/path/to/params",
		Environments: map[string]map[string]string{
			"dev": {
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
	}
	secretsConfig.ClearEnvironments()
	require.Empty(t, secretsConfig.Environments)
}

func TestSecretsConfig_UpdateSecretsConfigFromParameters(t *testing.T) {
	secretsConfig := &SecretsConfig{
		Environments: make(map[string]map[string]string),
	}
	params := map[string]string{
		"/path/to/params/foobar/dev/KEY1": "value1",
		"/path/to/params/foobar/dev/KEY2": "value2",
	}
	err := secretsConfig.UpdateSecretsConfigFromParameters(params)
	require.NoError(t, err)
	require.Equal(t, "value1", secretsConfig.Environments["dev"]["KEY1"])
	require.Equal(t, "value2", secretsConfig.Environments["dev"]["KEY2"])
}

func TestSecretsConfig_Save(t *testing.T) {
	secretsConfig := &SecretsConfig{
		Project: "foobar",
		Prefix:  "/path/to/params",
		Environments: map[string]map[string]string{
			"dev": {
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
	}
	err := secretsConfig.Save()
	require.NoError(t, err)

	// Clean up the generated test file
	RemoveTestFiles(t, SecretsConfigFile)
}

func TestLoadSecretsConfig(t *testing.T) {
	// Create the test file
	createTestFile(t, "psenv-secrets.yml", `
project: foobar
prefix: /path/to/params
environments:
  dev:
    KEY1: value1
    KEY2: value2
`)

	// Run the test
	_, err := LoadSecretsConfig()
	require.NoError(t, err)

	// Clean up the generated test file
	RemoveTestFiles(t, "psenv-secrets.yml")
}

func TestCreateNewProjectConfigFile(t *testing.T) {
	_, err := CreateNewProjectConfigFile()
	require.NoError(t, err)
	RemoveTestFiles(t, ProjectConfigFile)
}

func TestCreateNewSecretsConfigFile(t *testing.T) {
	_, err := CreateNewSecretsConfigFile()
	require.NoError(t, err)
	RemoveTestFiles(t, SecretsConfigFile)
}
