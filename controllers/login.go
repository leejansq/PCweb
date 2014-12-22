package controllers

import (
	"MYProject/PC_V1.0/models"
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"time"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Post() {
	user := this.GetString("USERNAME")
	pwd := this.GetString("Password")
	ret, err := models.YILogin(user, pwd)
	if err != nil {
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	ret_R, err := models.GetList(ret.Data.Token, ret.Data.Token_secret, ret.Data.Userid)
	if err != nil {
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	fmt.Printf("%#v\n", ret_R)
	this.SetSession("list", ret_R)
	for _, v := range ret_R.Data {
		go http.Get("http://127.0.0.1:8685/creat?uid=" + v.Uid + "&secret=" + models.AesDecrpto(v.Password, v.Uid))
	}

	//this.Data["json"] = ret_R
	//this.ServeJson()
	this.Data["Data"] = ret_R.Data
	this.TplNames = "Vplayer.html"
}

func (this *MainController) Get() {
	time.Sleep(time.Second * 5)
	//this.Data["Website"] = "beego.me"
	//this.Data["Email"] = "astaxie@gmail.com"
	this.TplNames = "login.html"
}
