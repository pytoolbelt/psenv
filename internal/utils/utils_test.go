package utils

import (
	"reflect"
	"testing"
)

func TestMergeLocalAndRemoteParams_AddsNewParams(t *testing.T) {
	localParams := map[string]string{"key1": "value1"}
	remoteParams := map[string]string{}

	expected := &Parameters{
		ToAdd:    map[string]string{"key1": "value1"},
		ToUpdate: map[string]string{},
		ToDelete: []string{},
	}

	result := MergeLocalAndRemoteParams(localParams, remoteParams)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMergeLocalAndRemoteParams_UpdatesExistingParams(t *testing.T) {
	localParams := map[string]string{"key1": "newValue"}
	remoteParams := map[string]string{"key1": "oldValue"}

	expected := &Parameters{
		ToAdd:    map[string]string{},
		ToUpdate: map[string]string{"key1": "newValue"},
		ToDelete: []string{},
	}

	result := MergeLocalAndRemoteParams(localParams, remoteParams)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMergeLocalAndRemoteParams_DeletesMissingParams(t *testing.T) {
	localParams := map[string]string{}
	remoteParams := map[string]string{"key1": "value1"}

	expected := &Parameters{
		ToAdd:    map[string]string{},
		ToUpdate: map[string]string{},
		ToDelete: []string{"key1"},
	}

	result := MergeLocalAndRemoteParams(localParams, remoteParams)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMergeLocalAndRemoteParams_NoChanges(t *testing.T) {
	localParams := map[string]string{"key1": "value1"}
	remoteParams := map[string]string{"key1": "value1"}

	expected := &Parameters{
		ToAdd:    map[string]string{},
		ToUpdate: map[string]string{},
		ToDelete: []string{},
	}

	result := MergeLocalAndRemoteParams(localParams, remoteParams)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestConvertParamsToEnvVars(t *testing.T) {
	params := map[string]string{"key1": "value1", "key2": "value2"}
	expected := []string{"key1=value1", "key2=value2"}

	result := ConvertParamsToEnvVars(params)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestConvertParamsToEnvVars_WithSlashes(t *testing.T) {
	params := map[string]string{"path/to/key1": "value1", "another/path/to/key2": "value2"}
	expected := []string{"key1=value1", "key2=value2"}

	result := ConvertParamsToEnvVars(params)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
