package config

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"os"
	"slices"
	"strings"
)

const ProjectConfigFile = "psenv-project.yml"
const SecretsConfigFile = "psenv-secrets.yml"

type ProjectConfig struct {
	Default      string   `yaml:"default"`
	Environments []string `yaml:"environments"`
	Prefix       string   `yaml:"prefix"`
	Project      string   `yaml:"project"`
}

// PrintTable prints the project config as a table to the terminal
func (c *ProjectConfig) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"Prefix", "Project", "Default", "Environments", "Path"})

	table.Append([]string{c.Prefix, c.Project, c.Default, c.Environments[0], c.GetEnvironmentPath(c.Environments[0])})

	for _, env := range c.Environments[1:] {
		table.Append([]string{"", "", "", env, c.GetEnvironmentPath(env)})
	}

	table.Render()
}

// GetEnvironmentPath returns the path to the environment
func (c *ProjectConfig) GetEnvironmentPath(env string) string {
	i := slices.Index(c.Environments, env)
	if i == -1 {
		return ""
	}
	return c.Prefix + "/" + c.Project + "/" + env
}

func (c *ProjectConfig) GetBasePath() string {
	return fmt.Sprintf("%s/%s/", c.Prefix, c.Project)
}

func (c *ProjectConfig) HasEnvironment(env string) bool {
	return slices.Contains(c.Environments, env)
}

// *************** Secrets Config ***************

type SecretsConfig struct {
	Project      string                       `yaml:"project"`
	Prefix       string                       `yaml:"prefix"`
	Environments map[string]map[string]string `yaml:"environments"`
}

func (c *SecretsConfig) GetEnvironmentPath(env string) string {
	return fmt.Sprintf("%s/%s/%s", c.Prefix, c.Project, env)
}

func (c *SecretsConfig) GetEnvironmentParams(env string) map[string]string {
	keys := make(map[string]string)
	path := c.GetEnvironmentPath(env)
	for k, v := range c.Environments[env] {
		keys[path+"/"+strings.ToUpper(k)] = v
	}
	return keys
}

func (c *SecretsConfig) ClearEnvironment(env string) {
	delete(c.Environments, env)
}

func (c *SecretsConfig) ClearEnvironments() {
	c.Environments = make(map[string]map[string]string)
}

// UpdateSecretsConfigFromParameters updates the secrets configuration from the
// parameters passed in as received from the parameter store
func (c *SecretsConfig) UpdateSecretsConfigFromParameters(params map[string]string) error {
	for k, v := range params {
		parts := strings.Split(k, "/")

		if len(parts) < 3 {
			return fmt.Errorf("invalid parameter name: %s", k)
		}

		env := parts[len(parts)-2]
		key := strings.ToUpper(parts[len(parts)-1])
		prefix := strings.Join(parts[0:len(parts)-3], "/")
		project := parts[len(parts)-3]

		c.Prefix = prefix
		c.Project = project
		if _, ok := c.Environments[env]; !ok {
			c.Environments[env] = make(map[string]string)
		}
		c.Environments[env][key] = v
	}
	return nil
}

func (c *SecretsConfig) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(SecretsConfigFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// *************** utility functions ***************

// LoadProjectConfig loads the project configuration from the project configuration yml file
func LoadProjectConfig() (*ProjectConfig, error) {
	var projectConfig ProjectConfig

	data, err := os.ReadFile(ProjectConfigFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &projectConfig)
	if err != nil {
		return nil, err
	}

	return &projectConfig, nil
}

// LoadSecretsConfig loads the secrets configuration from the secrets configuration yml file
func LoadSecretsConfig() (*SecretsConfig, error) {
	var secretsConfig SecretsConfig

	data, err := os.ReadFile(SecretsConfigFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &secretsConfig)
	if err != nil {
		return nil, err
	}

	return &secretsConfig, nil
}

// CreateNewProjectConfigFile creates a new project configuration file with template data
func CreateNewProjectConfigFile() (*ProjectConfig, error) {
	templateData := ProjectConfig{
		Default:      "dev",
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
	}

	data, err := yaml.Marshal(&templateData)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(ProjectConfigFile, data, 0644)
	if err != nil {
		return nil, err
	}

	return &templateData, nil
}

// Create a new psenv-secrets.yml file with template values
func CreateNewSecretsConfigFile() (*SecretsConfig, error) {
	templateData := SecretsConfig{
		Project: "foobar",
		Prefix:  "/path/to/params",
		Environments: map[string]map[string]string{
			"dev": {
				"KEY1": "value1",
				"KEY2": "value2",
			},
			"prod": {
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
	}

	data, err := yaml.Marshal(&templateData)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(SecretsConfigFile, data, 0644)
	if err != nil {
		return nil, err
	}

	return &templateData, nil
}
