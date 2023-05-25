package boot

import "net/http"

// 加载静态资源
func LoadStatic(dir string) {
	http.Handle("/"+dir+"/", http.StripPrefix("/"+dir+"/", http.FileServer(http.Dir(dir))))
}
