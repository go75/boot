package boot

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// 找出错误信息
func trace(message string) string {
	var pcs [32]uintptr
	//跳过前3个Callers
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			Warn.Printf("%s\n\n", trace(message))
			c.Fail(http.StatusInternalServerError, "Internal Server Error")
		}
	}()
	c.Next()
}
