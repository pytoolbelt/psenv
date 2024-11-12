package config

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"os"
	"slices"
)

type ProjectConfig struct {
	Environments []string `mapstructure:"environments"`
	Prefix       string   `mapstructure:"prefix"`
	Project      string   `mapstructure:"project"`
	Default      string   `mapstructure:"default"`
}

func InitAndLoadProjectConfig() (*ProjectConfig, error) {
	var cfg ProjectConfig
	viper.SetConfigName("psenv-project")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func CreateNewProjectConfigFile() (*ProjectConfig, error) {

	templateData := ProjectConfig{
		Environments: []string{"base", "dev", "prod", "test"},
		Prefix:       "/path/to/params",
		Project:      "foobar",
		Default:      "dev",
	}
	viper.Set("environments", templateData.Environments)
	viper.Set("prefix", templateData.Prefix)
	viper.Set("project", templateData.Project)
	viper.Set("default", templateData.Default)

	if err := viper.SafeWriteConfigAs("psenv-project.yml"); err != nil {
		return nil, err
	}

	return &templateData, nil
}

func (c *ProjectConfig) HasEnvironment(env string) bool {
	return slices.Contains(c.Environments, env)
}

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

// PrintTable prints the project config as a table to the terminal
func (c *ProjectConfig) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"Project", "Prefix", "Default", "Environments", "Path"})

	table.Append([]string{c.Project, c.Prefix, c.Default, c.Environments[0], c.GetEnvironmentPath(c.Environments[0])})

	for _, env := range c.Environments[1:] {
		table.Append([]string{"", "", "", env, c.GetEnvironmentPath(env)})
	}

	table.Render()
}
