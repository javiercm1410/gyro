package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

// Left here figure out how to pass any object

func DisplayData(outputType, path string, value any) {
	switch outputType {
	case "json":
		jsonOutput(value)
	case "file":
		fileOutput(value, path)
	default:
		log.Error("Invalid output type")
	}

}

func jsonOutput(value any) error {
	marshaled, err := json.MarshalIndent(value, "", "   ")
	if err != nil {
		log.Fatalf("marshaling error: %s", err)
		return err
	}
	fmt.Println(string(marshaled))
	return nil
}

func fileOutput(value any, path string) error {
	marshaled, err := json.MarshalIndent(value, "", "   ")
	if err != nil {
		log.Fatalf("marshaling error: %s", err)
		return err
	}

	err = os.WriteFile(path, marshaled, 0644)
	if err != nil {
		log.Fatalf("WriteFile error: %s", err)
		return err
	}

	log.Info("Output created on %s", path)
	return nil
}
