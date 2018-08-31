package assistant

import (
	"reflect"
)

type PBSupport struct {
	value interface{}
}

func (p *PBSupport) D(d interface{}) *PBSupport {
	if p.value == nil {
		p.value = d
	}

	return p
}

func (p *PBSupport) ToPB(out interface{}) {
	outv := reflect.ValueOf(out)

	if p.value == nil {
		panic("please init value use d func ")
	}

	if outv.Type().Kind() != reflect.Ptr {
		panic("need a ptr")
	} else {
		outv = outv.Elem()
	}

	structInfo := getStructInfo(p.value)
	valv := reflect.ValueOf(p.value)
	if valv.Type().Kind() == reflect.Ptr {
		valv = valv.Elem()
	}
	for _, v := range structInfo.StructFields {
		var originVal reflect.Value
		if v.IsInline {
			originVal = valv.FieldByIndex(v.InlineIndex)
		} else {
			originVal = valv.Field(v.Index)
		}
		if v.IsNormal {
			// outv.Field(v.Index).
			outField := outv.FieldByName(v.Name)
			if outField.CanSet() {
				outField.Set(originVal)
			}
		}
	}
}

// func (p *PBSupport) SetFromPB(in interface{}) {
// 	inv := reflect.ValueOf(in)

// 	if p.value == nil {
// 		panic("please init value use D func")
// 	}
// 	outv := reflect.ValueOf(p.value)

// 	if inv.Type().Kind() == reflect.Ptr {
// 		inv = inv.Elem()
// 	}

// 	if outv.Type().Kind() == reflect.Ptr {
// 		outv = outv.Elem()
// 	}

// 	structInfo := getStructInfo(p.value)

// 	for _, v := range structInfo.StructFields {
// 		// the value of in field
// 		var field reflect.Value
// 		fieldVal := inv.FieldByName(v.Name)

// 		if v.IsInline {
// 			field = outv.FieldByIndex(v.InlineIndex)

// 		}

// 	}
// }
