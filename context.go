package boot

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// 当前core框架的上下文
type Context struct {
	Req   *http.Request
	Resp  http.ResponseWriter
	Msg   interface{}
	fs  []func(c *Context)
	Param map[string]string
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Status(code int) {
	c.Resp.WriteHeader(code)
}

func (c *Context) Fail(code int, errInfo string) {
	c.Status(code)
	if _, err := c.Resp.Write([]byte(errInfo)); err != nil {
		Warn.Panicln("resp write err:", err)
	}
}

func (c *Context) Query(key string, defaultValue string) string {
	s := c.Req.URL.Query().Get(key)
	if s != "" {
		return s
	}
	return defaultValue
}

func (c *Context) GetHeader(key string, defaultValue string) string {
	s := c.Req.Header.Get(key)
	if s != "" {
		return s
	}
	return defaultValue
}

func (c *Context) SetHeader(key string, value string) {
	c.Resp.Header().Set(key, value)
}

//第一个参数name 为 cookie 名
//第二个参数value 为 cookie 值
//第三个参数path 为 cookie 所在的目录
//第四个domain 为所在域，表示我们的 cookie 作用范围，里面可以是localhost也可以是你的域名，看自己情况
//第五个参数maxAge 为 cookie 有效时长，当 cookie 存在的时间超过设定时间时，cookie 就会失效，它就不再是我们有效的 cookie，他的时间单位是秒second
//第六个secure 表示是否只能通过 https 访问，为true只能是https
//第七个httpOnly 表示 cookie 是否可以通过 js代码进行操作，为true时不能被js获取

func (c *Context) RequestAddCookie(name, value string) {
	c.Req.AddCookie(&http.Cookie{
		Name:  name,
		Value: value,
	})
}

func (c *Context) SetCookie(name, value, path, domain string, maxAge int, secure, httpOnly bool) {
	http.SetCookie(c.Resp, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *Context) CookieValue(name string) (string, error) {
	if cookie, err := c.Req.Cookie(name); err != nil {
		return "", err
	} else {
		return cookie.Value, nil
	}
}

func (c *Context) Uri() string {
	return c.Req.RequestURI
}

func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.Req.Cookie(name)
}

func (c *Context) DelCookie(name string) {
	if cookie, _ := c.Req.Cookie(name); cookie != nil {
		cookie.MaxAge = -1
	}
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	objData, err := json.Marshal(obj)
	if err != nil {
		Warn.Panicln("json marshal err:", err)
		c.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := c.Resp.Write(objData); err != nil {
		Warn.Panicln("resp write err:", err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if _, err := c.Resp.Write(data); err != nil {
		Warn.Panicln("resp write err:", err)
	}
}

func (c *Context) HTML(code int, html string, data interface{}) {
	//c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	t := templates.Lookup(html)
	if t != nil {
		if err := t.Execute(c.Resp, data); err != nil {
			Warn.Panicln("execute html err:", err)
			c.Resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		c.Resp.WriteHeader(http.StatusNotFound)
	}
}

func (c *Context) XML(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/xml")
	c.Status(code)
	objData, err := xml.Marshal(obj)
	if err != nil {
		Warn.Panicln("xml marshal err:", err)
		c.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := c.Resp.Write(objData); err != nil {
		Warn.Panicln("resp write err:", err)
	}
}

func (c *Context) execute() {
	c.fs[0](c)
}

func (c *Context) Next() {
	if len(c.fs) != 0 {
		c.fs = c.fs[1:]
		if len(c.fs) != 0 {
			c.fs[0](c)
		}
	}
}

func (c *Context) Abort() {
	c.fs = c.fs[:0]
}
