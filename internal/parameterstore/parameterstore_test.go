package parameterstore

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockSSMClient struct {
	mock.Mock
}

func (m *MockSSMClient) DescribeParameters(ctx context.Context, input *ssm.DescribeParametersInput, opts ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*ssm.DescribeParametersOutput), args.Error(1)
}

func (m *MockSSMClient) PutParameter(ctx context.Context, input *ssm.PutParameterInput, opts ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*ssm.PutParameterOutput), args.Error(1)
}

func (m *MockSSMClient) GetParametersByPath(ctx context.Context, input *ssm.GetParametersByPathInput, opts ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*ssm.GetParametersByPathOutput), args.Error(1)
}

func (m *MockSSMClient) DeleteParameters(ctx context.Context, input *ssm.DeleteParametersInput, opts ...func(*ssm.Options)) (*ssm.DeleteParametersOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*ssm.DeleteParametersOutput), args.Error(1)
}

func NewMockParameterStore() *ParameterStore {
	return &ParameterStore{
		Client: &MockSSMClient{},
	}
}

func TestDescribeParametersReturnsNames(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("DescribeParameters", mock.Anything, mock.Anything).Return(&ssm.DescribeParametersOutput{
		Parameters: []types.ParameterMetadata{
			{Name: aws.String("param1")},
			{Name: aws.String("param2")},
		},
	}, nil)

	names, err := ps.DescribeParameters("/path")
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"param1", "param2"}, names)
}

func TestDescribeParametersHandlesError(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("DescribeParameters", mock.Anything, mock.Anything).Return((*ssm.DescribeParametersOutput)(nil), errors.New("error"))

	names, err := ps.DescribeParameters("/path")
	require.Error(t, err)
	require.Empty(t, names)
}

func TestPutParametersSuccessfully(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("PutParameter", mock.Anything, mock.Anything).Return(&ssm.PutParameterOutput{
		Version: *aws.Int64(1),
	}, nil)

	err := ps.PutParameters(map[string]string{"param1": "value1"}, "keyId", true)
	require.NoError(t, err)
}

func TestPutParametersHandlesError(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("PutParameter", mock.Anything, mock.Anything).Return((*ssm.PutParameterOutput)(nil), errors.New("error"))

	err := ps.PutParameters(map[string]string{"param1": "value1"}, "keyId", true)
	require.Error(t, err)
}

func TestGetParametersReturnsValues(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("GetParametersByPath", mock.Anything, mock.Anything).Return(&ssm.GetParametersByPathOutput{
		Parameters: []types.Parameter{
			{Name: aws.String("param1"), Value: aws.String("value1")},
			{Name: aws.String("param2"), Value: aws.String("value2")},
		},
	}, nil)

	params, err := ps.GetParameters("/path", true)
	require.NoError(t, err)
	require.Equal(t, map[string]string{"param1": "value1", "param2": "value2"}, params)
}

func TestGetParametersHandlesError(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("GetParametersByPath", mock.Anything, mock.Anything).Return((*ssm.GetParametersByPathOutput)(nil), errors.New("error"))

	params, err := ps.GetParameters("/path", true)
	require.Error(t, err)
	require.Empty(t, params)
}

func TestDeleteParametersSuccessfully(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("DeleteParameters", mock.Anything, mock.Anything).Return(&ssm.DeleteParametersOutput{
		DeletedParameters: []string{"param1", "param2"},
	}, nil)

	err := ps.DeleteParameters([]string{"param1", "param2"})
	require.NoError(t, err)
}

func TestDeleteParametersHandlesError(t *testing.T) {
	mockClient := new(MockSSMClient)
	ps := &ParameterStore{Client: mockClient}

	mockClient.On("DeleteParameters", mock.Anything, mock.Anything).Return((*ssm.DeleteParametersOutput)(nil), errors.New("error"))

	err := ps.DeleteParameters([]string{"param1", "param2"})
	require.Error(t, err)
}
