package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	iam "github.com/javiercm1410/rotator/pkg/providers/aws"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

func DisplayData(outputFormat, path string, stale int, value []iam.UserData) {
	switch outputFormat {
	case "json":
		jsonOutput(value)
	case "file":
		fileOutput(value, path)
	case "table":
		headers, data, err := processTableData(value)
		if err != nil {
			log.Error("Could not process table data")
		}
		// fmt.Println("%v", len(data))
		tableOutput(headers, data, stale)
		// tableOutput(headers, data[:4])
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

	log.Infof("Output created on %s", path)
	return nil
}

func processTableData(value []iam.UserData) ([]string, [][]string, error) {
	var headers []string
	if reflect.TypeOf(value) == reflect.TypeOf([]iam.UserData{}) {
		headers = append(headers, []string{"UserName", "KeyId", "CreateDate", "KeyStatus", "LastUsedTime", "LastUsedService"}...)
	}

	var data [][]string
	for _, item := range value {
		// Type assert each item to UserAccessKeyData
		if user, ok := item.(iam.UserAccessKeyData); ok {
			for _, key := range user.Keys {
				var lastUsedTime string
				createDate := key.CreateDate.Format("2006-01-02 15:04:05")
				if key.LastUsedTime.IsZero() {
					lastUsedTime = "n/a"
				} else {
					lastUsedTime = key.LastUsedTime.Format("2006-01-02 15:04:05")

				}

				row := []string{
					user.UserName,
					*key.Id,
					createDate,
					string(key.KeyStatus),
					lastUsedTime,
					key.LastUsedService,
				}
				data = append(data, row)
			}
		}
	}

	return headers, data, nil
}

func tableOutput(headers []string, data [][]string, stale int) {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(130).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}

			even := row%2 == 0

			switch col {
			case 2: // Date column
				if row < len(data) && col < len(data[row]) { // Ensure bounds
					dateStr := data[row][col]
					parsedDate, err := time.Parse("2006-01-02 15:04:05", dateStr) // Match your date format
					if err == nil {                                               // If the date parsing is successful
						if time.Since(parsedDate).Hours() > float64(stale)*24 { // More than 90 days
							return baseStyle.Foreground(lipgloss.Color("#BA5F75")) // Red
						} else if time.Since(parsedDate).Hours() > float64(stale-10)*24 {
							return baseStyle.Foreground(lipgloss.Color("#FCFF5F")) // Yellow
						} else if even {
							return baseStyle.Foreground(lipgloss.Color("245"))
						}
					}
				}
				return baseStyle.Foreground(lipgloss.Color("252"))
			}

			if even {
				return baseStyle.Foreground(lipgloss.Color("245"))
			}
			return baseStyle.Foreground(lipgloss.Color("252"))
		})
	// fmt.Println("%v", len(data))

	fmt.Println(t)
}
