package flow

import(
	"fmt"
)

var(
	errNilFactory = fmt.Errorf("Flow factory is nil")
)

type FlowHandler struct{
	fct *Factory
}

func FlowNewHandler() *FlowHandler{
	return &FlowHandler{
		fct: new(Factory),
	}
}

func (f *FlowHandler)NewLine(handlers ...HandlerFunc) (int, error){
	if f.fct == nil{
		return 0, errNilFactory
	}
	return f.fct.NewLine(handlers...), nil
}

func (f *FlowHandler)Start(i int, p Params) error{
	if f.fct == nil{
		return errNilFactory
	}
	f.fct.Start(i, p)
	return nil
}

func (f *FlowHandler)Stop(i int) error {
	if f.fct == nil{
		return errNilFactory
	}
	f.fct.Stop(i)
	return nil
}

func (f *FlowHandler)IsStopped(i int) (bool, error){
	if f.fct == nil{
		return false, errNilFactory
	}
	return f.fct.IsStopped(i), nil
}

func (f *FlowHandler)AreAllStopped() (bool, error){
	if f.fct == nil{
		return false, errNilFactory
	}
	return f.fct.AreAllStopped(), nil
}

func (f *FlowHandler)Destory() error{
	if f.fct == nil{
		return errNilFactory
	}
	f.fct.Destory()
	return nil
}
