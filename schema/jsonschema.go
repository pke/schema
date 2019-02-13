package schema

import (
	"encoding/json"
	"reflect"
)

func generateNumberJSONSchema(tags schemaTags) (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type": "number",
	}
	if value, ok := tags["description"]; ok {
		schema["description"] = value
	}
	return schema, nil
}

func generateStringJSONSchema(tags schemaTags) (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type": "string",
	}
	if value, ok := tags["description"]; ok {
		schema["description"] = value
	}
	if value, ok := tags["format"]; ok {
		schema["format"] = value
	}
	return schema, nil
}

func generateBooleanJSONSchema(tags schemaTags) (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type": "boolean",
	}
	if value, ok := tags["description"]; ok {
		schema["description"] = value
	}
	return schema, nil
}

func generateArrayJSONSchema(typ reflect.Type, tags schemaTags) (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type": "array",
	}
	if value, ok := tags["description"]; ok {
		schema["description"] = value
	}

	items, err := generateJSONSchema(typ.Elem(), tags.itemsTags())
	if err != nil {
		return nil, err
	}
	schema["items"] = items
	return schema, nil
}

func generateObjectJSONSchema(typ reflect.Type, tags schemaTags) (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type": "object",
	}
	if value, ok := tags["description"]; ok {
		schema["description"] = value
	}

	properties := map[string]interface{}{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := getJSONName(field)
		if fieldName == "" {
			continue
		}
		fieldTags := newSchemaTags(field.Tag)
		property, err := generateJSONSchema(field.Type, fieldTags)
		if err != nil {
			return nil, err
		}
		properties[fieldName] = property
	}
	schema["properties"] = properties
	return schema, nil
}

func generateJSONSchema(typ reflect.Type, tags schemaTags) (map[string]interface{}, error) {
	kind := typ.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return generateNumberJSONSchema(tags)
	case reflect.String:
		return generateStringJSONSchema(tags)
	case reflect.Bool:
		return generateBooleanJSONSchema(tags)
	case reflect.Slice:
		return generateArrayJSONSchema(typ, tags)
	case reflect.Struct:
		return generateObjectJSONSchema(typ, tags)
	default:
		return nil, &UnsupportedKindError{kind: kind}
	}
}

func GenerateJSONSchema() ([]byte, error) {
	typ := reflect.TypeOf(Schema{})
	schema, err := generateJSONSchema(typ, map[string]string{})
	if err != nil {
		return nil, err
	}
	schema["$schema"] = "http://json-schema.org/draft-07/schema#"
	schema["title"] = "Resumic resume Schema"
	return json.MarshalIndent(schema, "", "  ")
}
