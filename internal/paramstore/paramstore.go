package paramstore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type ParamStore struct {
	SSMClient *ssm.Client
	SSMPath   string
}

func NewParamStore(ssmPath string) (*ParamStore, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config, %v", err)
	}

	return &ParamStore{
		SSMClient: ssm.NewFromConfig(cfg),
		SSMPath:   ssmPath,
	}, nil
}

func (p *ParamStore) FormatParamName(name string) string {
	return fmt.Sprintf("%s/%s", p.SSMPath, name)
}

func (p *ParamStore) BuildPutParamInput(name, value string, overwrite bool) *ssm.PutParameterInput {
	return &ssm.PutParameterInput{
		Name:      aws.String(p.FormatParamName(name)),
		Value:     aws.String(value),
		Type:      types.ParameterTypeSecureString,
		Overwrite: aws.Bool(overwrite),
	}
}

func (p *ParamStore) BuildGetParamsByPathInput(next string, decrypt bool) *ssm.GetParametersByPathInput {
	return &ssm.GetParametersByPathInput{
		Path:           aws.String(p.SSMPath),
		WithDecryption: aws.Bool(decrypt),
		NextToken:      aws.String(next),
		MaxResults:     aws.Int32(10),
	}
}

func (p *ParamStore) BuildDeleteParamsInput(names []string) *ssm.DeleteParametersInput {

	for i, name := range names {
		names[i] = p.FormatParamName(name)
	}

	return &ssm.DeleteParametersInput{
		Names: names,
	}
}

func (p *ParamStore) BuildDescribeParametersInput() *ssm.DescribeParametersInput {
	return &ssm.DescribeParametersInput{
		ParameterFilters: []types.ParameterStringFilter{
			{
				Key:    aws.String("Name"),
				Option: aws.String("BeginsWith"),
				Values: []string{p.SSMPath},
			},
		},
	}
}

func (p *ParamStore) DescribeParameters() ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var names []string

	input := p.BuildDescribeParametersInput()

	result, err := p.SSMClient.DescribeParameters(ctx, input)
	if err != nil {
		return names, fmt.Errorf("Error describing parameters: %s", err)
	}
	for _, param := range result.Parameters {
		names = append(names, strings.TrimPrefix(*param.Name, p.SSMPath+"/"))
	}

	return names, nil
}

func (p *ParamStore) PutParameters(params map[string]string, overwrite bool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for k, v := range params {
		params := p.BuildPutParamInput(k, v, overwrite)
		result, err := p.SSMClient.PutParameter(ctx, params)
		if err != nil {
			return fmt.Errorf("Error putting parameter %s: %s", k, err)
		}
		fmt.Printf("Parameter added: %s Version: %d\n", *params.Name, result.Version)
	}
	return nil
}

func (p *ParamStore) GetParameters(decrypt bool) (map[string]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	next := ""
	params := make(map[string]string)

	for {
		input := p.BuildGetParamsByPathInput(next, decrypt)
		result, err := p.SSMClient.GetParametersByPath(ctx, input)

		if err != nil {
			return nil, fmt.Errorf("Error getting parameters: %s", err)
		}

		for _, param := range result.Parameters {
			n := p.ParseParameterName(*param.Name)
			params[n] = *param.Value
		}

		if result.NextToken == nil {
			break
		}
		next = *result.NextToken
	}
	return params, nil
}

func (p *ParamStore) DeleteParameters(names []string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := p.BuildDeleteParamsInput(names)
	result, err := p.SSMClient.DeleteParameters(ctx, input)

	if err != nil {
		return fmt.Errorf("Error deleting parameters: %s", err)
	}

	for _, deleted := range result.DeletedParameters {
		fmt.Printf("Parameter deleted: %s\n", deleted)
	}

	return nil
}

func (p *ParamStore) ParseParameterName(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}

func GetEnvNameFromSSMPath(path string) (string, error) {
	// Split the path by '/'
	parts := strings.Split(path, "/")

	// Check if there are enough parts to get the second-to-last part
	if len(parts) < 1 {
		return "", errors.New("path does not have enough parts")
	}

	// Return the last part
	return parts[len(parts)-1], nil
}

func SplitAndDeduplicatePaths(paths []string) []string {
	uniquePaths := make(map[string]struct{})
	for _, path := range paths {
		// Split the path by '/' and remove the last part
		parts := strings.Split(path, "/")
		if len(parts) > 1 {
			// Join the parts back together, excluding the last part
			newPath := strings.Join(parts[:len(parts)-1], "/")
			uniquePaths[newPath] = struct{}{}
		}
	}

	// Convert the map keys to a slice
	result := make([]string, 0, len(uniquePaths))
	for path := range uniquePaths {
		result = append(result, path)
	}

	return result
}

func FormatParamsAsEnv(params map[string]string) []string {
	var env []string

	for k, v := range params {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}
