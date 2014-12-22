package routers

import (
	"MYProject/PC_V1.0/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/login", &controllers.MainController{})
	beego.Router("/camp", &controllers.CampController{})
}
