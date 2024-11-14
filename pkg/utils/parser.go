package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// YamlDoc represents a YAML document structure
type YamlDoc struct {
	Data       map[string]string
	StringData map[string]string
}

// EnvVarObject represents a map of environment variables
type EnvVarObject map[string]string

// MergeDataFromManifests merges data from multiple YAML manifests
func MergeDataFromManifests(manifests []YamlDoc) EnvVarObject {
	envData := make(EnvVarObject)

	for _, manifest := range manifests {
		for k, v := range manifest.Data {
			envData[k] = v
		}
		for k, v := range manifest.StringData {
			envData[k] = v
		}
	}

	return envData
}

// GenerateEnvFile generates an environment file from the given EnvVarObject
func GenerateEnvFile(envObject EnvVarObject, filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		existingContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading existing file: %v", err)
		}

		scanner := bufio.NewScanner(strings.NewReader(string(existingContent)))
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) != 0 && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					if _, exists := envObject[key]; !exists {
						// Remove quotes from value
						envObject[key] = strings.Trim(value, "'\"")
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error scanning existing file: %w", err)
		}
	}

	var envContent strings.Builder
	for key, value := range envObject {
		_, err := fmt.Fprintf(&envContent, "%s='%s'\n", key, value)
		if err != nil {
			return fmt.Errorf("error writing to string builder: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(envContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Printf("dotenv file generated in %s!\n", filePath)
	return nil
}
