package runtime

import (
	"fmt"
	"github.com/heqzha/goutils/date"
	"reflect"
	"runtime"
	"time"
)

func GetFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func RunFunc(f interface{}, args ...interface{}) []interface{} {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()
	fName := runtime.FuncForPC(fValue.Pointer()).Name()
	inValues := []reflect.Value{}
	for idx, arg := range args {
		argValue := reflect.ValueOf(arg)
		argType := argValue.Type()
		if !argType.ConvertibleTo(fType.In(idx)) {
			panic(fmt.Sprintf("function %s require %s, but get %s", fName, fType.In(idx).Name(), argType.Name()))
		}
		inValues = append(inValues, argValue)
	}

	outValues := fValue.Call(inValues)
	out := []interface{}{}
	for _, v := range outValues {
		out = append(out, v.Interface())
	}
	return out
}

func PrintTimeCost(f interface{}, args ...interface{}) []interface{} {
	now := time.Now()
	defer fmt.Printf("%s cost %s\n", GetFuncName(f), date.DateDurationFrom(now))
	return RunFunc(f, args...)
}
