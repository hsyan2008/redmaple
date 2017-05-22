package controllers

import (
	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type base struct {
	hfw.Controller
	User *models.Users
}

func (this *base) Before() {
	// hfw.Debug("base Before")
	this.Layout = "layout.html"
	this.Data["siteName"] = "-RedMaple开发上线管理系统"
	if this.Controll != "index" {
		if this.Session.IsExist("userinfo") == false {
			this.Redirect("/")
		}
	}
	if this.Session.IsExist("userinfo") {
		this.User = this.Session.Get("userinfo").(*models.Users)
	}
	this.Data["userinfo"] = this.User
}
