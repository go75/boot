package boot

import "net/http"

/*
*
core框架初始的router
*/
func New(funcations ...func(c *Context)) *Router {
	if funcations == nil {
		funcations = make([]func(c *Context), 0)
	}
	return &Router{
		pattern:    "",
		childs:     nil,
		fs: funcations,
	}
}

func Default(funcations ...func(c *Context)) *Router {
	return New(append([]func(c *Context){Recovery}, funcations...)...)
}

/*
*
运行core框架
*/
func (r *Router)Run(addr string) error {
	return http.ListenAndServe(addr, r)
}