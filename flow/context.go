package flow

import (
	"math"
	"sync"
)

const abortIndex int = math.MaxInt32 / 2

type Params map[string]interface{}

func (p Params) copy() Params {
	var cp = make(Params)
	for k, v := range p {
		cp[k] = v
	}
	return cp
}

type Context struct {
	params Params
	line   *Line
	index  int

	Errors errors
	mutex  *sync.Mutex
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) Init(params Params, l *Line) {
	c.params = params
	c.line = l
	c.index = 0
	c.Errors = c.Errors[0:0]
	c.mutex = &sync.Mutex{}
	//Set first input params
	c.line.do(c.index, c)
}

func (c *Context) Start(){
	go c.line.start()
}

func (c *Context) reset() {
	c.params = make(Params)
	c.line = nil
	c.index = 0
	c.Errors = c.Errors[0:0]
	c.mutex = &sync.Mutex{}
}

// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This have to be used then the context has to be passed to a goroutine.
func (c *Context) Copy() *Context {
	var cp = *c
	cp.params = c.params.copy()
	cp.index = abortIndex
	cp.line = nil
	return &cp
}

// HandlerName returns the current handler of line's name. For example if the handler is "handleGetUsers()", this
// function will return "main.handleGetUsers"
func (c *Context) HandlerName() string {
	return nameOfFunction(c.line.getHandler(int(c.index)))
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in github.
func (c *Context) Next() {
	if !c.IsAborted() {
		n, end := c.line.next(c.index)
		if end {
			// Reach the end of line
			c.Abort()
			return
		}
		n.in(c)
	}
}

// Repeat should be used only inside middleware.
// It repeats workflow from process i, but keeps all parameters
func (c *Context) Repeat() {
	if !c.IsAborted() {
		c.index = 0
		c.line.do(c.index, c)
	}
}

//Jump should be used only inside middleware.
// It jump to process i, but keeps all parameters
func (c *Context) Jump(i int){
	if !c.IsAborted() {
		c.index = i
		c.line.do(c.index, c)
	}
}

// IsAborted returns true if the currect context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized. If the
// authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.line.stop()
	c.index = abortIndex
}

// AbortWithError calls `AbortWithStatus()` and `Error()` internally. This method stops the chain, writes the status code and
// pushes the specified error to `c.Errors`.
// See Context.Error() for more details.
func (c *Context) AbortWithError(err error) *Error {
	c.Abort()
	return c.Error(err)
}

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors
// and push them to a database together, print a log, or append it in the HTTP response.
func (c *Context) Error(err error) *Error {
	var parsedError *Error
	switch err.(type) {
	case *Error:
		parsedError = err.(*Error)
	default:
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}
	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusivelly for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	c.mutex.Lock()
	if c.params == nil {
		c.params = make(Params)
	}
	c.params[key] = value
	c.mutex.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	if c.params != nil {
		c.mutex.Lock()
		value, exists = c.params[key]
		c.mutex.Unlock()
	}
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value

	}
	panic("Key \"" + key + "\" does not exist")
}
