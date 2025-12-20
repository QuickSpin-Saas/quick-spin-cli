package output

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// TableFormatter formats output as a table
type TableFormatter struct{}

// Format formats data as a table
func (f *TableFormatter) Format(data interface{}) error {
	// For single items, display as key-value pairs
	return f.formatKeyValue(data)
}

// FormatList formats a list as a table
func (f *TableFormatter) FormatList(data interface{}, headers []string) error {
	// Print headers
	fmt.Fprintln(os.Stdout, strings.Join(headers, "\t"))
	fmt.Fprintln(os.Stdout, strings.Repeat("-", len(strings.Join(headers, "\t"))))

	// Extract rows from data
	rows := f.extractRows(data)
	for _, row := range rows {
		fmt.Fprintln(os.Stdout, strings.Join(row, "\t"))
	}

	return nil
}

// FormatError formats an error in table format
func (f *TableFormatter) FormatError(err error) error {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	return nil
}

// formatKeyValue formats a single item as key-value pairs
func (f *TableFormatter) formatKeyValue(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			// Skip unexported fields
			if !value.CanInterface() {
				continue
			}

			// Get field name from json tag if available
			fieldName := field.Name
			if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
				// Extract just the field name from the tag (before comma)
				parts := strings.Split(tag, ",")
				if parts[0] != "" {
					fieldName = parts[0]
				}
			}

			fmt.Fprintf(os.Stdout, "%-20s %v\n", fieldName+":", value.Interface())
		}
	} else if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			fmt.Fprintf(os.Stdout, "%-20s %v\n", fmt.Sprintf("%v:", key.Interface()), value.Interface())
		}
	} else {
		fmt.Fprintf(os.Stdout, "%v\n", data)
	}

	return nil
}

// extractRows extracts rows from a slice of structs or maps
func (f *TableFormatter) extractRows(data interface{}) [][]string {
	var rows [][]string

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return rows
	}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		row := f.extractRow(item)
		if len(row) > 0 {
			rows = append(rows, row)
		}
	}

	return rows
}

// extractRow extracts a single row from a struct or map
func (f *TableFormatter) extractRow(v reflect.Value) []string {
	var row []string

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !field.CanInterface() {
				continue
			}
			row = append(row, fmt.Sprintf("%v", field.Interface()))
		}
	} else if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			row = append(row, fmt.Sprintf("%v", value.Interface()))
		}
	}

	return row
}
