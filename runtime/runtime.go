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

//Caution: RunFunc is very slow!!!
func RunFunc(f interface{}, args ...interface{}) []interface{} {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()
	fName := runtime.FuncForPC(fValue.Pointer()).Name()
	inValues := []reflect.Value{}
	for idx, arg := range args {
		argValue := reflect.ValueOf(arg)
		if argValue.IsValid() {
			argType := argValue.Type()
			if !argType.ConvertibleTo(fType.In(idx)) {
				panic(fmt.Sprintf("function %s require %s, but get %s", fName, fType.In(idx).Name(), argType.Name()))
			}
		} else {
			argValue = reflect.Zero(fType.In(idx))
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
	defer func() {
		fmt.Printf("%s:%v  cost %s\n", GetFuncName(f), args, date.DateDurationFrom(now))
	}()
	return RunFunc(f, args...)
}
