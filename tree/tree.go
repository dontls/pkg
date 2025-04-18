package tree

import (
	"reflect"
)

type Builder struct {
	len       int
	elem      reflect.Value
	ChildName string
}

func (o *Builder) addChildren(i, j int) {
	elem := o.elem.Index(i).FieldByName(o.ChildName)
	if !elem.CanSet() {
		return
	}
	elem.Set(reflect.Append(elem, o.elem.Index(j)))
}

func (o *Builder) findChildren(i int, h2 func(int, int) bool) {
	for j := 0; j < o.len; j++ {
		if !h2(i, j) {
			continue
		}
		o.findChildren(j, h2)
		o.addChildren(i, j)
	}
}

// h1 判断主节点， h2判断子节点
func (o *Builder) Do(obj any, h1 func(int) bool, h2 func(int, int) bool) any {
	o.elem = reflect.ValueOf(obj)
	if reflect.TypeOf(obj).Kind() != reflect.Slice || o.elem.Len() == 0 {
		return nil
	}
	if o.elem.Index(0).FieldByName(o.ChildName).Kind() != reflect.Slice {
		return nil
	}
	o.len = o.elem.Len()
	var r []any
	for i := 0; i < o.len; i++ {
		if !h1(i) {
			continue
		}
		o.findChildren(i, h2)
		r = append(r, o.elem.Index(i).Interface())
	}
	return r
}

// h1 判断主节点， h2判断子节点
func Slice(obj any, h1 func(int) bool, h2 func(int, int) bool) any {
	o := Builder{ChildName: "Children"}
	return o.Do(obj, h1, h2)
}
