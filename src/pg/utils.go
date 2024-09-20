package pg

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

func getSafeColName(name string) string {
	var sb strings.Builder
	re := regexp.MustCompile("[A-Z][a-z]*")

	components := re.FindAllString(name, -1)
	for i, component := range components {
		if i > 0 {
			sb.WriteRune('_')
		}

		sb.WriteString(strings.ToLower(component))
	}

	return sb.String()
}

func GetColumnNamesOf(obj interface{}) ([]string, error) {
	var content []string
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		return content, fmt.Errorf("type error: expected struct type but got %T", obj)
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		if name, ok := field.Tag.Lookup("dbNamne"); ok {
			content = append(content, name)
		} else {
			content = append(content, getSafeColName(field.Name))
		}
	}

	return content, nil
}

func BuildCreateString(schema string, name string, obj interface{}) (string, error) {
	var content string
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		return content, fmt.Errorf("type error: expected struct type but got %T", obj)
	}

	var sb strings.Builder
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		typeName, ok := field.Tag.Lookup("dbType")
		if !ok {
			continue
		}

		columnName := field.Name
		if name, ok := field.Tag.Lookup("dbName"); ok {
			columnName = name
		} else {
			columnName = getSafeColName(columnName)
		}

		if typeMod, ok := field.Tag.Lookup("dbMod"); ok {
			typeName += fmt.Sprintf("(%s)", typeMod)
		}

		if isPrimary, ok := field.Tag.Lookup("dbIsPrimary"); ok {
			if strings.ToLower(isPrimary) == "true" {
				typeName += " PRIMARY KEY"
			}
		}

		sb.WriteString(fmt.Sprintf("\n\t%s  %s,", columnName, typeName))
	}

	content = sb.String()
	contentLen := len(content)
	if contentLen < 1 {
		return content, fmt.Errorf("invalid arguments: failed to generate table, are you missing struct tags?")
	}

	content = content[:contentLen-1]
	content = fmt.Sprintf(
		"CREATE TABLE %s (%s\n);",
		pgx.Identifier{schema, name}.Sanitize(), content,
	)

	return content, nil
}

func FlattenRow(obj interface{}) ([]any, error) {
	var content []any
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		return content, fmt.Errorf("type error: expected struct type but got %T", obj)
	}

	rv := reflect.ValueOf(obj)
	rt = rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		item := rv.FieldByName(field.Name)
		if item.CanInterface() {
			value := item.Interface()
			content = append(content, value)
		}
	}

	return content, nil
}
