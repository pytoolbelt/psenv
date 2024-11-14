package parameterstore

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type ParameterStore struct {
	Client *ssm.Client
}

func New() (*ParameterStore, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config, %v", err)
	}

	return &ParameterStore{
		Client: ssm.NewFromConfig(cfg),
	}, nil
}

func BuildPutParameterInput(paramName, paramValue, keyId string, overwrite bool) *ssm.PutParameterInput {
	return &ssm.PutParameterInput{
		Name:      aws.String(paramName),
		Value:     aws.String(paramValue),
		Type:      types.ParameterTypeSecureString,
		Overwrite: aws.Bool(overwrite),
		KeyId:     aws.String(keyId),
	}
}

func BuildGetParamsByPathInput(path, next string, decrypt bool) *ssm.GetParametersByPathInput {
	return &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		WithDecryption: aws.Bool(decrypt),
		NextToken:      aws.String(next),
		MaxResults:     aws.Int32(10),
	}
}

func BuildDeleteParamsInput(names []string) *ssm.DeleteParametersInput {
	return &ssm.DeleteParametersInput{
		Names: names,
	}
}

//func (p *ParamStore) BuildDescribeParametersInput() *ssm.DescribeParametersInput {
//	return &ssm.DescribeParametersInput{
//		ParameterFilters: []types.ParameterStringFilter{
//			{
//				Key:    aws.String("Name"),
//				Option: aws.String("BeginsWith"),
//				Values: []string{p.SSMPath},
//			},
//		},
//	}
//}
//
//func (p *ParamStore) DescribeParameters() ([]string, error) {
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	var names []string
//
//	input := p.BuildDescribeParametersInput()
//
//	result, err := p.SSMClient.DescribeParameters(ctx, input)
//	if err != nil {
//		return names, fmt.Errorf("Error describing parameters: %s", err)
//	}
//	for _, param := range result.Parameters {
//		names = append(names, strings.TrimPrefix(*param.Name, p.SSMPath+"/"))
//	}
//
//	return names, nil
//}

func (p *ParameterStore) PutParameters(params map[string]string, keyId string, overwrite bool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for paramKey, paramValue := range params {
		params := BuildPutParameterInput(paramKey, paramValue, keyId, overwrite)
		result, err := p.Client.PutParameter(ctx, params)

		if err != nil {

			return fmt.Errorf("Error putting parameter %s: %s", paramKey, err)
		}
		fmt.Printf("Parameter added: %s Version: %d\n", *params.Name, result.Version)
	}
	return nil
}

func (p *ParameterStore) GetParameters(path string, decrypt bool) (map[string]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	next := ""
	params := make(map[string]string)

	for {
		input := BuildGetParamsByPathInput(path, next, decrypt)
		result, err := p.Client.GetParametersByPath(ctx, input)

		if err != nil {
			return nil, fmt.Errorf("error getting parameters: %s", err)
		}

		for _, param := range result.Parameters {
			params[*param.Name] = *param.Value
		}

		if result.NextToken == nil {
			break
		}
		next = *result.NextToken
	}
	return params, nil
}

func (p *ParameterStore) DeleteParameters(names []string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := BuildDeleteParamsInput(names)
	result, err := p.Client.DeleteParameters(ctx, input)

	if err != nil {
		return fmt.Errorf("Error deleting parameters: %s", err)
	}

	for _, deleted := range result.DeletedParameters {
		fmt.Printf("Parameter deleted: %s\n", deleted)
	}

	return nil
}

func CheckCredentials() error {
	_, err := New()
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config, %v", err)
	}
	return nil
}

//}
//
//func (p *ParamStore) ParseParameterName(name string) string {
//	parts := strings.Split(name, "/")
//	return parts[len(parts)-1]
//}
//
//func GetEnvNameFromSSMPath(path string) (string, error) {
//	// Split the path by '/'
//	parts := strings.Split(path, "/")
//
//	// Check if there are enough parts to get the second-to-last part
//	if len(parts) < 1 {
//		return "", errors.New("path does not have enough parts")
//	}
//
//	// Return the last part
//	return parts[len(parts)-1], nil
//}
//
//func SplitAndDeduplicatePaths(paths []string) []string {
//	uniquePaths := make(map[string]struct{})
//	for _, path := range paths {
//		// Split the path by '/' and remove the last part
//		parts := strings.Split(path, "/")
//		if len(parts) > 1 {
//			// Join the parts back together, excluding the last part
//			newPath := strings.Join(parts[:len(parts)-1], "/")
//			uniquePaths[newPath] = struct{}{}
//		}
//	}
//
//	// Convert the map keys to a slice
//	result := make([]string, 0, len(uniquePaths))
//	for path := range uniquePaths {
//		result = append(result, path)
//	}
//
//	return result
//}
//
//func FormatParamsAsEnv(params map[string]string) []string {
//	var env []string
//
//	for k, v := range params {
//		env = append(env, fmt.Sprintf("%s=%s", k, v))
//	}
//	return env
//}
