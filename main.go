package main

import (
	_ "MYProject/PC_V1.0/routers"
	"github.com/astaxie/beego"
)

func main() {
	//beego.EnableGzip = false
	beego.Run()
}
