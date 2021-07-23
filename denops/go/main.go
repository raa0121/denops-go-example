package main

import (
	"reflect"
	"syscall/js"
)

/**
 * for go 1.17 (Generics)
 *
 * func main() {
 * 	c := make(chan struct{}, 0)
 * 	js.Global().Set("GoAdd", js.FuncOf(addJs))
 * 	<-c
 * }
 * type Addable interface {
 * 	type int, int8, int16, int32, int64,
 * 		uint, uint8, uint16, uint32, uint64, uintptr,
 * 		float32, float64, complex64, complex128,
 * 		string
 * }
 *
 * func add[T Addable](a, b T) T {
 *     return a + b
 * }
 * func addJs(this js.Value, args []js.Value) interface{} {
 * 	if args[0].Type() == js.TypeNumber &&
 * 	   args[1].Type() == js.TypeNumber {
 * 		return js.ValueOf(add(args[0].Float(), args[1].Float()))
 * 	}
 * 	if args[0].Type() == js.TypeString &&
 * 	   args[1].Type() == js.TypeString {
 * 		return js.ValueOf(add(args[0].String(), args[1].String()))
 * 	}
 * 	return nil
 * }
 **/

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("GoAddIntNoWrap", js.FuncOf(addIntNoWrap))
	js.Global().Set("GoAdd", js.FuncOf(wrap(add)))
	js.Global().Set("GoAddInt", js.FuncOf(wrap(addInt)))
	js.Global().Set("GoAddFloat", js.FuncOf(wrap(addFloat)))
	js.Global().Set("GoAddString", js.FuncOf(wrap(addFloat)))
	<-c
}

func addIntNoWrap(this js.Value, args []js.Value) interface{} {
	return js.ValueOf(args[0].Int() + args[1].Int())
}

func add(lhs, rhs interface{}) interface{} {
	lf := reflect.ValueOf(lhs)
	rf := reflect.ValueOf(rhs)
	if lf.Type().Kind() == reflect.Int &&
	   rf.Type().Kind() == reflect.Int {
		if l, ok := lf.Interface().(int); ok {
			if r, ok := rf.Interface().(int); ok {
				return l + r
			}
		}
	}
	if lf.Type().Kind() == reflect.Float64 &&
	   rf.Type().Kind() == reflect.Float64 {
		if l, ok := lf.Interface().(float64); ok {
			if r, ok := rf.Interface().(float64); ok {
				return l + r
			}
		}
	}
	if lf.Type().Kind() == reflect.String &&
	   rf.Type().Kind() == reflect.String {
		if l, ok := lf.Interface().(string); ok {
			if r, ok := rf.Interface().(string); ok {
				return l + r
			}
		}
	}
	return nil
}

func addInt(lhs, rhs int) int {
	return lhs + rhs
}

func addFloat(lhs, rhs float64) float64 {
	return lhs + rhs
}

func addString(lhs, rhs string) string {
	return lhs + rhs
}

func wrap(f interface{}) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		rf := reflect.ValueOf(f)
		rt := rf.Type()
		rargs := []reflect.Value{}
		for i := 0; i < rt.NumIn(); i++ {
			var arg reflect.Value
			switch args[i].Type() {
			case js.TypeUndefined:
				arg = reflect.Zero(reflect.TypeOf(nil)).Convert(rt.In(i))
			case js.TypeNull:
				arg = reflect.Zero(reflect.TypeOf(nil)).Convert(rt.In(i))
			case js.TypeBoolean:
				arg = reflect.ValueOf(args[i].Bool()).Convert(rt.In(i))
			case js.TypeNumber:
				arg = reflect.ValueOf(args[i].Float()).Convert(rt.In(i))
			case js.TypeString:
				arg = reflect.ValueOf(args[i].String()).Convert(rt.In(i))
			case js.TypeSymbol:
				arg = reflect.ValueOf(args[i].String()).Convert(rt.In(i))
			case js.TypeObject:
				arg = reflect.ValueOf(args[i].JSValue()).Convert(rt.In(i))
			case js.TypeFunction:
				arg = reflect.ValueOf(args[i].JSValue()).Convert(rt.In(i))
			}
			rargs = append(rargs, arg)
		}
		ret := rf.Call(rargs)
		if len(ret) > 0 {
			return ret[0].Interface()
		}
		return nil
	}
}
