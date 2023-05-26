package boot

import (
	"net/http"
	"strings"
)

// 路由模块
type Router struct {
	//请求方法+":"+节点名
	pattern string
	//当前模块所有的子节点
	childs []*Router
	//当前模块的处理方法
	fs []func(c *Context)
}

var NotFound = func(c *Context) {
	c.Resp.WriteHeader(http.StatusNotFound)
}

func (r *Router) New(routerName string, funcations ...func(c *Context)) *Router {
	if funcations == nil {
		funcations = make([]func(c *Context), 0)
	}
	if r.childs == nil {
		r.childs = []*Router{{
			pattern:    routerName,
			childs:     nil,
			fs: funcations,
		}}
	} else {
		r.childs = append([]*Router{{
			pattern:    routerName,
			childs:     nil,
			fs: funcations,
		}}, r.childs...)
	}
	return r.childs[0]
}

// 处理http请求
func (r *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//初始化上下文模块
	c := &Context{
		Msg:   nil,
		Req:   req,
		Resp:  resp,
		Param: make(map[string]string, 0),
		fs:  r.fs,
	}
	if strings.Contains(req.URL.Path, ":") {
		NotFound(c)
		return
	}
	//获取路径节点集
	patterns := strings.Split(req.URL.Path, "/")[1:]
	patterns[len(patterns)-1] = req.Method + ":" + patterns[len(patterns)-1]
	//遍历root的中间变量
	temp := r
	//遍历路径节点集
	for _, pattern := range patterns {
		//如果当前路由的子节点为空，直接执行NotFound
		if temp.childs == nil {
			NotFound(c)
			return
		}
		//遍历当前路由的子节点
		for _, child := range temp.childs {
			//符合匹配条件则添加funcations，并进行路径节点集的下次遍历
			if pattern == child.pattern {
				temp = child
				c.fs = append(c.fs, (child.fs)...)
				goto NEXT
			}
			//判断有无动态参数
			childParma, patternParam, ok := dynamicMatch(child.pattern, pattern)
			if !ok {
				continue
			}
			temp = child
			c.Param[childParma] = patternParam

			c.fs = append(c.fs, (child.fs)...)
			goto NEXT
		}
		//如果当前路由的子节点全部不匹配，则执行NotFound
		NotFound(c)
		return
	NEXT:
	}

	//执行请求方法
	c.execute()
}

// 注册路由
func (r *Router) addRouter(method string, pattern string, fs ...func(c *Context)) {
	var isDynamic bool
	//判断动态参数是否合法
	switch strings.Count(pattern, ":") {
	case 0:
	case 1:
		if pattern[len(pattern)-1] != ':' {
			panic("invaild pattern")
		}
		isDynamic = true
	default:
		panic("invaild pattern")
	}
	patterns := strings.Split(pattern, "/")[1:]
	patterns[len(patterns)-1] = method+":" + patterns[len(patterns)-1]
	temp := r
	//遍历路径节点集
	for i := 0; i < len(patterns); i++ {

		//如果当前路由的子节点为空，创建路由子节点
		if temp.childs == nil {
			if i == len(patterns)-1 {
				//当前是最后一个节点,创建新节点，然后将fs参数赋给funcations，最后return
				//初始化节点
				temp.childs = []*Router{
					{
						pattern:    patterns[i],
						childs:     nil,
						fs: fs,
					},
				}
				return
			} else {
				//当前不是最后一个节点，则初始化节点，temp向后移动，进行下次循环
				f := make([]func(c *Context), 0)
				//初始化节点
				temp.childs = []*Router{
					{
						pattern:    patterns[i],
						childs:     nil,
						fs: f,
					},
				}
				//temp后移
				temp = temp.childs[0]
				continue
			}
		}

		//如果当前路由的子节点全部不匹配，则执行创建路由子节点
		if i == len(patterns)-1 {
			//当前是最后一个节点
			
			//遍历当前路由的子节点
			for _, child := range temp.childs {
				//符合匹配条件，panic
				if patterns[i] == child.pattern || (child.pattern[len(child.pattern)-1] == ':' && isDynamic)  {
					panic("invalid pattern")
				}
			}

			if isDynamic {
				//当前是最后一个节点,创建新节点，然后将fs参数赋给funcations，最后return
				temp.childs = append(temp.childs, &Router{
					pattern:    patterns[i],
					childs:     nil,
					fs: fs,
				})
			} else {
				temp.childs = append([]*Router{{
					pattern: patterns[i],
					childs: nil,
					fs: fs,
				}}, temp.childs...)
			}
			return
		} else {
			//当前不是最后一个节点
			
			//遍历当前路由的子节点
			for _, child := range temp.childs {
				//符合匹配条件，temp后移，然后进行路径节点集的下次遍历
				if patterns[i] == child.pattern || child.pattern[len(child.pattern)-1] == ':'  {
					//temp后移
					temp = child
					//进入下次外层循环
					goto NEXT
				}
			}
			
			//未匹配, 初始化节点，temp向后移动，进行下次循环
			f := make([]func(c *Context), 0)
			//初始化节点
			temp.childs = []*Router{
				{
					pattern:    patterns[i],
					childs:     nil,
					fs: f,
				},
			}
			//temp后移
			temp = temp.childs[0]
			continue
		}
	NEXT:
	}
}


// 注册GET路由
func (r *Router) GET(pattern string, fs ...func(c *Context)) {
	r.addRouter("GET", pattern, fs...)
}

// 注册POST路由
func (r *Router) POST(pattern string, fs ...func(c *Context)) {
	r.addRouter("POST", pattern, fs...)
}

// 注册DELETE路由
func (r *Router) DELETE(pattern string, fs ...func(c *Context)) {
	r.addRouter("DELETE", pattern, fs...)
}

// 注册PATCH路由
func (r *Router) PATCH(pattern string, fs ...func(c *Context)) {
	r.addRouter("PATCH", pattern, fs...)
}

// 注册PUT路由
func (r *Router) PUT(pattern string, fs ...func(c *Context)) {
	r.addRouter("PUT", pattern, fs...)
}

// 注册OPTIONS路由
func (r *Router) OPTIONS(pattern string, fs ...func(c *Context)) {
	r.addRouter("OPTIONS", pattern, fs...)
}

// 注册HEAD路由
func (r *Router) HEAD(pattern string, fs ...func(c *Context)) {
	r.addRouter("HEAD", pattern, fs...)
}