package utils

import (
	"os"
	"testing"
)

func TestMergeDataFromManifests(t *testing.T) {
	tests := []struct {
		name      string
		manifests []YamlDoc
		expected  EnvVarObject
	}{
		{
			name:      "should return an empty object if no manifests are provided",
			manifests: []YamlDoc{},
			expected:  map[string]string{},
		},
		{
			name: "should merge data from multiple manifests",
			manifests: []YamlDoc{
				{Data: map[string]string{"key1": "value1"}},
				{Data: map[string]string{"key2": "value2"}},
				{Data: map[string]string{"key3": "value3"}},
			},
			expected: EnvVarObject{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "should merge stringData from multiple manifests",
			manifests: []YamlDoc{
				{StringData: map[string]string{"key1": "value1"}},
				{StringData: map[string]string{"key2": "value2"}},
				{StringData: map[string]string{"key3": "value3"}},
			},
			expected: EnvVarObject{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "should merge both data and stringData from multiple manifests",
			manifests: []YamlDoc{
				{Data: map[string]string{"key1": "value1"}},
				{StringData: map[string]string{"key2": "value2"}},
				{Data: map[string]string{"key3": "value3"}, StringData: map[string]string{"key4": "value4"}},
			},
			expected: EnvVarObject{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			name: "should override data with stringData if keys are the same",
			manifests: []YamlDoc{
				{Data: map[string]string{"key1": "value1"}},
				{StringData: map[string]string{"key1": "overrideValue"}},
			},
			expected: EnvVarObject{
				"key1": "overrideValue",
			},
		},
		{
			name: "should handle manifests with undefined data or stringData",
			manifests: []YamlDoc{
				{Data: map[string]string{"key1": "value1"}},
				{StringData: nil},
				{Data: nil},
				{StringData: map[string]string{"key2": "value2"}},
			},
			expected: EnvVarObject{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeDataFromManifests(tt.manifests)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("Expected %s for key %s, got %s", expectedValue, key, result[key])
				}
			}
		})
	}
}

func TestGenerateEnvFile(t *testing.T) {
	testFilePath := ".env.test"

	// defer func() {
	// 	os.Remove(testFilePath)
	// }()

	tests := []struct {
		name            string
		envObject       EnvVarObject
		expectedContent string
	}{
		{
			name:            "should handle an empty envObject",
			envObject:       EnvVarObject{},
			expectedContent: "",
		},
		{
			name: "should generate a .env file with the correct content",
			envObject: EnvVarObject{
				"KEY1": "value1",
				"KEY2": "value2",
				"KEY3": "value3",
			},
			expectedContent: `KEY1='value1'
KEY2='value2'
KEY3='value3'
`,
		},
		{
			name: "should handle special characters in values",
			envObject: EnvVarObject{
				"KEY1": "value with spaces",
				"KEY2": "value_with_underscores",
				"KEY3": "value-with-dashes",
				"KEY4": "value=with=equal",
			},
			expectedContent: `KEY1='value with spaces'
KEY2='value_with_underscores'
KEY3='value-with-dashes'
KEY4='value=with=equal'
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateEnvFile(tt.envObject, testFilePath); err != nil {
				t.Fatalf("Failed to generate env file: %v", err)
			}

			content, err := os.ReadFile(testFilePath)
			if err != nil {
				t.Fatalf("Failed to read env file: %v", err)
			}

			if string(content) != tt.expectedContent {
				t.Errorf("Expected content %s, got %s", tt.expectedContent, string(content))
			}
			os.Remove(testFilePath)

		})
	}
}
