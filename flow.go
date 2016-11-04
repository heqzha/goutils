package goutils

import(
	"fmt"
	"github.com/heqzha/goutils/flow"
)


var (
	flowFactory *flow.Factory
	errNilFactory = fmt.Errorf("Flow factory is nil")
)

func FlowConfig(){
	flowFactory = new(flow.Factory)
}

func FlowNewLine(handlers ...flow.HandlerFunc) (int, error){
	if flowFactory == nil{
		return 0, errNilFactory
	}
	return flowFactory.NewLine(handlers...), nil
}

func FlowStart(i int, p flow.Params) error{
	if flowFactory == nil{
		return errNilFactory
	}
	flowFactory.Start(i, p)
	return nil
}

func FlowStop(i int) error {
	if flowFactory == nil{
		return errNilFactory
	}
	flowFactory.Stop(i)
	return nil
}

func FlowIsStop(i int) (bool, error){
	if flowFactory == nil{
		return false, errNilFactory
	}
	return flowFactory.IsStop(i), nil
}

func FlowAreAllStop() (bool, error){
	if flowFactory == nil{
		return false, errNilFactory
	}
	return flowFactory.AreAllStop(), nil
}

func FlowDestory() error{
	if flowFactory == nil{
		return errNilFactory
	}
	flowFactory.Destory()
	return nil
}
