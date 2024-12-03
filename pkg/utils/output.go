package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	iam "github.com/javiercm1410/gyro/pkg/providers/aws"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

const dateFormat = "2006-01-02 15:04:05"

// DisplayData processes and displays data in the specified format.
func DisplayData(outputFormat, path string, stale int, value []iam.UserData) {
	if len(value) == 0 {
		log.Warn("No data available to display")
		return
	}

	switch outputFormat {
	case "json":
		if err := jsonOutput(value); err != nil {
			log.Error("Failed to generate JSON output", "error", err)
		}
	case "file":
		if err := fileOutput(value, path); err != nil {
			log.Error("Failed to write data to file", "error", err)
		}
	case "table":
		headers, data, err := processTableData(value)
		if err != nil {
			log.Error("Failed to process table data", "error", err)
			return
		}
		tableOutput(headers, data, stale)
	default:
		log.Error("Generate output error", "Error", outputFormat)
	}
}

func jsonOutput(value any) error {
	marshaled, err := json.MarshalIndent(value, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	fmt.Println(string(marshaled))
	return nil
}

func fileOutput(value any, path string) error {
	marshaled, err := json.MarshalIndent(value, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(path, marshaled, 0644); err != nil {
		return fmt.Errorf("error writing to file %s: %w", path, err)
	}

	log.Infof("Output saved to %s", path)
	return nil
}

// func processTableData(value []iam.UserData) ([]string, [][]string, error) {
// 	if reflect.TypeOf(value) != reflect.TypeOf([]iam.UserData{}) {
// 		return nil, nil, errors.New("unexpected data type for table processing")
// 	}

// 	headers := []string{"UserName", "KeyId", "CreateDate", "KeyStatus", "LastUsedTime", "LastUsedService"}
// 	data := make([][]string, 0, len(value))

// 	for _, user := range value {
// 		for _, key := range user.Keys {
// 			createDate := key.CreateDate.Format(dateFormat)
// 			lastUsedTime := "n/a"
// 			if !key.LastUsedTime.IsZero() {
// 				lastUsedTime = key.LastUsedTime.Format(dateFormat)
// 			}

// 			row := []string{
// 				user.UserName,
// 				*key.Id,
// 				createDate,
// 				string(key.KeyStatus),
// 				lastUsedTime,
// 				key.LastUsedService,
// 			}
// 			data = append(data, row)
// 		}
// 	}
// 	return headers, data, nil
// }

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

func tableOutput(headers []string, data [][]string, age int) {
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

			if col == 2 { // CreateDate column
				dateStr := data[row][col]
				parsedDate, err := time.Parse(dateFormat, dateStr)
				if err == nil { // If the date parsing is successful
					ageHours := float64(age) * 24
					switch {
					case time.Since(parsedDate).Hours() > ageHours:
						return baseStyle.Foreground(lipgloss.Color("#BA5F75")) // Red
					case time.Since(parsedDate).Hours() > ageHours-10*24:
						return baseStyle.Foreground(lipgloss.Color("#FCFF5F")) // Yellow
					}
				}
			}

			if even {
				return baseStyle.Foreground(lipgloss.Color("245"))
			}
			return baseStyle.Foreground(lipgloss.Color("252"))
		})

	fmt.Println(t)
}
