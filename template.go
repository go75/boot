package boot

import "html/template"

// 全局的模板变量
var templates *template.Template

// 加载模板
func LoadTemplates(pattern string) {
	templates = template.New("templates")
	template.Must(templates.ParseGlob(pattern))
}
