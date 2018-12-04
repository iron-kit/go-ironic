package hptypes

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	pb "github.com/golang/protobuf/ptypes/struct"
	"gopkg.in/mgo.v2/bson"
)

func EncodeToStruct(m map[string]interface{}) *pb.Struct {
	mSize := len(m)

	if mSize == 0 {
		return nil
	}

	fields := make(map[string]*pb.Value, mSize)

	for k, val := range m {
		fields[k] = encodeStructValue(val)
	}

	return &pb.Struct{
		Fields: fields,
	}
}

func DecodeToMap(s *pb.Struct) map[string]interface{} {
	m := make(map[string]interface{})
	if s == nil {
		return m
	}

	for k, v := range s.Fields {
		m[k] = decodeStructValue(v)
	}

	return m
}

func decodeStructValue(v *pb.Value) interface{} {
	if v == nil {
		return nil
	}
	switch k := v.Kind.(type) {
	case *pb.Value_NullValue:
		return nil
	case *pb.Value_BoolValue:
		return k.BoolValue
	case *pb.Value_StringValue:

		s := k.StringValue
		if bson.IsObjectIdHex(s) {
			return bson.ObjectIdHex(s)
		}

		spattern := regexp.MustCompile(`^ObjectIdHex\(\"(\w+)\"\)$`)

		if spattern.MatchString(s) {
			sub := spattern.FindStringSubmatch(s)
			return bson.ObjectIdHex(sub[len(sub)-1])
		}

		cpattern := regexp.MustCompile(`^objectId:(\w+)$`)
		if cpattern.MatchString(s) {
			// fmt.Println("math objectId:x")
			sub := cpattern.FindStringSubmatch(s)
			// fmt.Println("find submatch", sub)
			// fmt.Println
			return bson.ObjectIdHex(sub[len(sub)-1])
		}

		// 如果是时间
		if parsedTime, err := time.Parse(time.RFC3339, s); err == nil {
			return parsedTime
		}
		// s := time.RFC3339
		// timePattern := regexp.MustCompile(`^time:(.+)$`)
		// if timePattern.MatchString(s) {
		// 	sub := timePattern.FindStringSubmatch(s)
		// 	if t, err := time.Parse(time.RFC1123Z, sub[len(sub)-1]); err != nil {
		// 		return t
		// 	}
		// 	return time.Now()
		// }
		return s
	case *pb.Value_NumberValue:
		return k.NumberValue
	case *pb.Value_StructValue:
		return DecodeToMap(k.StructValue)
	case *pb.Value_ListValue:
		res := make([]interface{}, len(k.ListValue.Values))
		for index, val := range k.ListValue.Values {
			res[index] = decodeStructValue(val)
		}
		return res
	default:
		panic("unknown kind")
	}
}

func encodeStructValue(v interface{}) *pb.Value {
	val := v
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return &pb.Value{
			Kind: &pb.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: v,
			},
		}
	case time.Time:
		t := val.(time.Time)
		return &pb.Value{
			Kind: &pb.Value_StringValue{
				// StringValue: fmt.Sprintf("time:%s", t.Format(time.RFC1123Z)),
				StringValue: fmt.Sprintf("%s", t.Format(time.RFC3339)),
			},
		}
	case bson.ObjectId:
		bv := val.(bson.ObjectId)
		return &pb.Value{
			Kind: &pb.Value_StringValue{
				StringValue: fmt.Sprintf("objectId:%x", string(bv)),
			},
		}
	// case string:
	// 	fmt.Println("is string")
	// 	return &pb.Value{
	// 		Kind: &pb.Value_StringValue{
	// 			StringValue: v,
	// 		},
	// 	}
	case error:
		return &pb.Value{
			Kind: &pb.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	default:
		// Fallback to reflection for other types
		return toValue(reflect.ValueOf(v))
	}
}

func toValue(v reflect.Value) *pb.Value {
	switch v.Kind() {
	case reflect.Bool:
		return &pb.Value{
			Kind: &pb.Value_BoolValue{
				BoolValue: v.Bool(),
			},
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v.Int()),
			},
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: float64(v.Uint()),
			},
		}
	case reflect.Float32, reflect.Float64:
		return &pb.Value{
			Kind: &pb.Value_NumberValue{
				NumberValue: v.Float(),
			},
		}
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return toValue(reflect.Indirect(v))
	case reflect.Array, reflect.Slice:
		size := v.Len()
		if size == 0 {
			return nil
		}
		values := make([]*pb.Value, size)

		for i := 0; i < size; i++ {
			vk := v.Index(i)
			if vk.Kind() == reflect.Interface && vk.Elem().Kind() == reflect.Map {
				vk = vk.Elem()
				// values[i] = toValue()
			}
			values[i] = toValue(vk)
		}
		return &pb.Value{
			Kind: &pb.Value_ListValue{
				ListValue: &pb.ListValue{
					Values: values,
				},
			},
		}
	case reflect.Struct:
		t := v.Type()
		size := v.NumField()
		if size == 0 {
			return nil
		}
		fields := make(map[string]*pb.Value, size)
		for i := 0; i < size; i++ {
			name := t.Field(i).Name
			// Better way?
			if len(name) > 0 && 'A' <= name[0] && name[0] <= 'Z' {
				fields[name] = toValue(v.Field(i))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &pb.Value{
			Kind: &pb.Value_StructValue{
				StructValue: &pb.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		fields := make(map[string]*pb.Value, len(keys))
		for _, k := range keys {
			if k.Kind() == reflect.String {
				fields[k.String()] = toValue(v.MapIndex(k))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &pb.Value{
			Kind: &pb.Value_StructValue{
				StructValue: &pb.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.String:
		return &pb.Value{
			Kind: &pb.Value_StringValue{
				StringValue: v.String(),
			},
		}
	default:
		// Last resort
		return &pb.Value{
			Kind: &pb.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}
