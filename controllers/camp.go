package controllers

import (
	//"MYProject/PC_V1.0/models"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type CampController struct {
	beego.Controller
}

func (this *CampController) Get() {
	uid := this.GetString("uid")

	resp, err := http.Get("http://127.0.0.1:8685/start?uid=" + uid)
	if err != nil {
		return
	}
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		this.Data["Url"] = string(body)
		this.TplNames = "video.html"
	}
}
