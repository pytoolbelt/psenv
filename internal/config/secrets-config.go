package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Project      string                       `mapstructure:"project"`
	Prefix       string                       `mapstructure:"prefix"`
	Environments map[string]map[string]string `mapstructure:"environments"`
}

type Environment struct {
	Name    string            `mapstructure:"name"`
	Project string            `mapstructure:"project"`
	Prefix  string            `mapstructure:"prefix"`
	Params  map[string]string `mapstructure:"params"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) GetEnvironment(envName string) (*Environment, error) {
	env, exists := cfg.Environments[envName]

	if !exists {
		return nil, fmt.Errorf("environment %s not found", envName)
	}

	return &Environment{
		Project: cfg.Project,
		Prefix:  cfg.Prefix,
		Params:  env,
		Name:    envName,
	}, nil
}

func (cfg *Config) SetEnvironment(env *Environment) {
	k := fmt.Sprintf("environments.%s", env.Name)
	viper.Set(k, env.Params)
	viper.Set("project", env.Project)
	viper.Set("prefix", env.Prefix)
}

func (cfg *Config) ClearEnvironments() {
	cfg.Environments = make(map[string]map[string]string)
	viper.Set("environments", make(map[string]string))
}

func (cfg *Config) Save() error {
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	return viper.WriteConfig()
}

func (cfg *Config) GetEnvironments() []Environment {
	var envs []Environment
	for k, v := range cfg.Environments {
		envs = append(envs, Environment{
			Name:    k,
			Project: cfg.Project,
			Prefix:  cfg.Prefix,
			Params:  v,
		})
	}
	return envs
}

func (cfg *Config) GetBasePath() string {
	return fmt.Sprintf("%s/%s/", cfg.Prefix, cfg.Project)
}

func RemoveVarNameFromPath(path string) string {
	// removes the last /var from the path and returns the new path
	parts := strings.Split(path, "/")
	return fmt.Sprintf("%s/%s", strings.Join(parts[:len(parts)-1], "/"), "")
}

func RemoveVarNameFromPaths(paths []string) []string {
	var newPaths []string
	for _, p := range paths {
		newPaths = append(newPaths, RemoveVarNameFromPath(p))
	}
	return newPaths
}

func NewEnvironmentFromPath(path string) *Environment {
	// path is in the format prefix/project/env
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return nil
	}
	return &Environment{
		Name:    parts[len(parts)-1],
		Project: parts[len(parts)-2],
		Prefix:  strings.Join(parts[:len(parts)-2], "/"),
	}
}

func (e *Environment) GetParams() map[string]string {
	var params = make(map[string]string)
	for k, v := range e.Params {
		key := strings.ToUpper(k) //fmt.Sprintf("%s/%s/%s", e.Prefix, e.Project, k)
		params[key] = v
	}
	return params
}

func (e *Environment) GetParamStorePath() string {

	return fmt.Sprintf("%s/%s/%s", e.Prefix, e.Project, e.Name)
}

func (e *Environment) GetParamsToAdd(eParams map[string]string) map[string]string {
	lParams := e.GetParams()
	var paramsToAdd = make(map[string]string)

	for k, v := range lParams {
		value, exists := eParams[k]

		if !exists {
			paramsToAdd[k] = v
			continue
		}

		if value != v {
			paramsToAdd[k] = v
		}

	}
	return paramsToAdd
}

func (e *Environment) GetParamsToDelete(rParams map[string]string) []string {
	var paramsToDelete []string
	lParams := e.GetParams()

	// check if the remote params are in the local params, if they are not, add them to the params to delete
	for k := range rParams {
		_, exists := lParams[k]
		if !exists {
			paramsToDelete = append(paramsToDelete, k)
		}
	}
	return paramsToDelete
}

func InitAndLoad() (*Config, error) {
	if err := Init(); err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error initializing config: %w", err)
	}
	return cfg, nil
}

// CreateNewSecretsConfigFile creates a new secrets config file with template data
func CreateNewSecretsConfigFile() (*Config, error) {
	cfg := Config{
		Project:      "my-project",
		Prefix:       "/my-prefix/",
		Environments: map[string]map[string]string{"dev": {"var1": "value1"}, "prod": {"var2": "value2"}, "base": {"var3": "value3"}, "test": {"var4": "value4"}},
	}

	viper.Set("project", cfg.Project)
	viper.Set("prefix", cfg.Prefix)
	viper.Set("environments", cfg.Environments)

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.SafeWriteConfigAs("psenv-secrets.yml"); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Init() error {
	viper.SetConfigName("psenv-secrets")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
