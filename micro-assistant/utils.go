package assistant

import (
	"go/ast"
	"reflect"
	"strings"
	"sync"
	"time"
)

func getTypeName(t reflect.Type) string {
	// name := ""
	if t.Kind() == reflect.Ptr {
		return t.Elem().String()
	}

	return t.String()
}

func IsZero(v interface{}) bool {
	vv := reflect.ValueOf(v)

	return isZero(vv)
}

var typeTime = reflect.TypeOf(time.Time{})

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()
		if vt == typeTime {
			return v.Interface().(time.Time).IsZero()
		}
		for i := 0; i < v.NumField(); i++ {
			if vt.Field(i).PkgPath != "" && !vt.Field(i).Anonymous {
				continue // Private field
			}
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	}
	return false
}

type StructInfo struct {
	Type         reflect.Type
	StructFields []*StructField
}

type StructField struct {
	Name        string
	Index       int
	IsNormal    bool
	IsIgnored   bool
	IsInline    bool
	Tag         reflect.StructTag
	TagMap      map[string]string
	Struct      reflect.StructField
	Zero        reflect.Value
	InlineIndex []int
}

var structMap sync.Map

func parseTagConfig(tags reflect.StructTag) map[string]string {
	conf := map[string]string{}
	for index, str := range []string{tags.Get("bson"), tags.Get("monger")} {
		// bson
		if index == 0 {
			tags := strings.Split(str, ",")

			for _, tag := range tags {
				// if t == "" {
				// 	continue
				// }
				switch tag {
				case "inline":
					conf["INLINE"] = "true"
				case "omitempty":
					conf["OMITEMPTY"] = "true"
				default:
					conf["COLUMN"] = tag
				}
			}

			continue
		}
		tags := strings.Split(str, ",")
		for _, value := range tags {
			v := strings.Split(value, "=")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if k == "COLUMN" {
				continue
			}
			if len(v) >= 2 {
				conf[k] = strings.Join(v[1:], "=")
			} else {
				conf[k] = k
			}
		}
	}

	return conf
}

func getStructInfo(d interface{}) *StructInfo {
	var structInfo StructInfo
	if d == nil {
		return &structInfo
	}

	reflectV := reflect.ValueOf(d)
	reflectType := reflectV.Type()

	if reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	if reflectType.Kind() != reflect.Struct {
		return &structInfo
	}

	if v, found := structMap.Load(reflectType); found && v != nil {
		return v.(*StructInfo)
	}

	structInfo.Type = reflectType

	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {

			field := &StructField{
				Struct:   fieldStruct,
				Name:     fieldStruct.Name,
				Tag:      fieldStruct.Tag,
				TagMap:   parseTagConfig(fieldStruct.Tag),
				Zero:     reflect.New(fieldStruct.Type).Elem(),
				Index:    i,
				IsInline: false,
			}

			// hidden
			if _, found := field.TagMap["-"]; found {
				field.IsIgnored = true
			} else if v, foundInline := field.TagMap["INLINE"]; foundInline && v == "true" {
				// the field is inline
				inlineFieldStruct := getStructInfo(reflect.New(fieldStruct.Type).Interface())

				for _, inlineField := range inlineFieldStruct.StructFields {
					inlineField.IsInline = true
					// inlineField.Index = []int{i, field.Index[0]}
					inlineField.InlineIndex = []int{i, inlineField.Index}
					structInfo.StructFields = append(structInfo.StructFields, inlineField)
				}
				continue
			} else {

				indirectType := fieldStruct.Type
				for indirectType.Kind() == reflect.Ptr {
					indirectType = indirectType.Elem()
				}

				fieldValue := reflect.New(indirectType).Interface()
				if _, isTime := fieldValue.(*time.Time); isTime {
					field.IsNormal = false
				} else {
					switch fieldStruct.Type.Kind() {
					case reflect.Slice:
						field.IsNormal = false
					case reflect.Struct:
						fallthrough
					case reflect.Ptr:
						field.IsNormal = false
					default:
						field.IsNormal = true
					}
				}
			}
			structInfo.StructFields = append(structInfo.StructFields, field)
		}
	}

	structMap.Store(reflectType, &structInfo)
	return &structInfo
}
