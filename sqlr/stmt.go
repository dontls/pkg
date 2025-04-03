package sqlr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func _reflectValue(v any) reflect.Value {
	value := reflect.ValueOf(v)
	// json.Unmarshal returns errors for these
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

// 反射生成slice的基本类型对象
func BaseModel(dest any) any {
	value := _reflectValue(dest)
	if value.Type().Kind() == reflect.Slice {
		return reflect.New(value.Type().Elem()).Interface()
	}
	return nil
}

// 返回值 s: sprintf类型格式; n: 字段名
func fieldsFmt(v reflect.Value) (s string, ns []string) {
	ref := v.Type()
	n := ref.NumField()
	for i := 0; i < n; i++ {
		fld := ref.Field(i)
		t := fld.Type.Kind()
		switch {
		case t <= reflect.Uint64:
			s += "%d,"
			ns = append(ns, fld.Name)
		case t == reflect.Float32 || t == reflect.Float64:
			s += "%f,"
			ns = append(ns, fld.Name)
		case t == reflect.String:
			s += "'%s',"
			ns = append(ns, fld.Name)
		case t == reflect.Struct:
			if fld.Anonymous { // 结构体嵌套
				s1, n1 := fieldsFmt(v.Field(i))
				s += s1
				ns = append(ns, n1...)
			} else {
				s += "'%s'," // 默认json格式
				ns = append(ns, fld.Name)
			}
		default:
			panic(fld.Type.Name() + " is not supported")
		}
	}
	return
}

func fieldVal(v reflect.Value) any {
	t := v.Kind()
	switch {
	case t > reflect.Invalid && t < reflect.Uint:
		return v.Int()
	case t > reflect.Int64 && t < reflect.Uintptr:
		return v.Uint()
	case t == reflect.Float32 || t == reflect.Float64:
		return v.Float()
	case t == reflect.String:
		return v.String()
	default:
		return nil
	}
}

func fieldAddr(v reflect.Value) (r []any) {
	//如果传入指针则拿到指针指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	n := v.NumField()
	for i := 0; i < n; i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct && v.Type().Field(i).Anonymous {
			r = append(r, fieldAddr(field)...)
		} else {
			r = append(r, field.Addr().Interface())
		}
	}
	return
}

// 获取结构体实例字段的地址
func SqlValuesAddr(s any) []any {
	value := _reflectValue(s)
	//参数必须是结构体
	if value.Kind() != reflect.Struct {
		return nil
	}
	return fieldAddr(value)
}

func SqlValueFmt(v any) string {
	s, _ := fieldsFmt(_reflectValue(v))
	if s != "" {
		s = strings.TrimRight(s, ",")
		s = "(" + s + ")"
	}
	return s
}

func SqlValueNames(v any) string {
	_, ns := fieldsFmt(_reflectValue(v))
	if len(ns) > 1 {
		s := strings.Join(ns, ", ")
		return "(" + s + ")"
	}
	return ""
}

func _sqlValue(v reflect.Value) (r []any) {
	n := v.NumField()
	for i := 0; i < n; i++ {
		fd := v.Field(i)
		if fd.Kind() == reflect.Struct {
			if v.Type().Field(i).Anonymous {
				r = append(r, _sqlValue(fd)...)
			} else {
				b, _ := json.Marshal(fd.Interface())
				r = append(r, string(b)) // 默认json格式
			}
		} else {
			if val := fieldVal(fd); val != nil {
				r = append(r, val)
			}
		}
	}
	return
}

func SqlValues(v any, handle func(int) string) string {
	value := _reflectValue(v)
	sqls := ""
	switch value.Type().Kind() {
	case reflect.Slice:
		svf := SqlValueFmt(value.Index(0).Interface()) + ","
		for i := 0; i < value.Len(); i++ {
			if handle != nil {
				sqls += handle(i)
			}
			sqls += fmt.Sprintf(svf, _sqlValue(value.Index(i))...)
		}
	case reflect.Struct:
		svf := SqlValueFmt(v) + ","
		if handle != nil {
			sqls += handle(-1)
		}
		sqls += fmt.Sprintf(svf, _sqlValue(value)...)
	}
	return strings.TrimRight(sqls, ",")
}
