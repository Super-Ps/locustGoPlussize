package runfunc

import (
	"reflect"
	"runtime"
)

func CallOnly(f interface{}, args []interface{}) {
	fun := reflect.ValueOf(f)
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	fun.Call(in)
}

func CallReturn(f interface{}, args []interface{}) []reflect.Value {
	fun := reflect.ValueOf(f)
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	return fun.Call(in)
}

func GetName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
