package boot

import (
	"log"
	"os"
)

var (
	//debug信息
	Debug *log.Logger = log.New(os.Stdout, "\u001B[1;36m[Debug]:\u001B[0m", log.Ltime|log.Llongfile)
	//重要信息
	Info *log.Logger = log.New(os.Stdout, "\u001B[1;34m[Info]:\u001B[0m", log.Ldate|log.Ltime|log.Lshortfile)
	//警告
	Warn *log.Logger = log.New(os.Stdout, "\u001B[1;33m[Warn]:\u001B[0m", log.Ldate|log.Ltime|log.Lshortfile)
)