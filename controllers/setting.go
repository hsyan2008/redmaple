package controllers

import (
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Setting struct {
	base
}

func (this *Setting) Index() {
	this.Data["title"] = "修改系统设置"
	this.TemplateFile = "setting/edit.html"
	tmp, _ := settingModel.SearchOne(models.Cond{})
	if tmp == nil {
		this.Data["setting"] = new(models.Settings)
	} else {
		this.Data["setting"] = tmp
	}
}

func (this *Setting) Save() {
	if this.Request.Method == "POST" {
		setting := new(models.Settings)
		Id := this.Request.PostFormValue("Id")
		if Id != "" {
			setting, _ = settingModel.GetById(Id)
			if setting == nil {
				this.Throw(99400, "参数错误")
			}
		}

		setting.SmtpAddr = this.Request.PostFormValue("SmtpAddr")
		setting.SmtpPort, _ = strconv.Atoi(this.Request.PostFormValue("SmtpPort"))
		setting.SmtpSsl = this.Request.PostFormValue("SmtpSsl")
		if setting.SmtpSsl != "Y" {
			setting.SmtpSsl = "N"
		}
		setting.SmtpName = this.Request.PostFormValue("SmtpName")
		setting.SmtpUser = this.Request.PostFormValue("SmtpUser")
		setting.SmtpPass = this.Request.PostFormValue("SmtpPass")

		setting.IsDeleted = "N"

		d := gomail.NewDialer(setting.SmtpAddr, setting.SmtpPort, setting.SmtpUser, setting.SmtpPass)
		s, err := d.Dial()
		if err != nil {
			this.Throw(99400, "邮箱配置不正确:"+err.Error())
		}
		defer func() {
			_ = s.Close()
		}()

		// user, _ := userModel.GetById(1)
		// m := gomail.NewMessage()
		// m.SetHeader("From", setting.SmtpUser)
		// m.SetAddressHeader("To", user.Email, "test")
		// m.SetHeader("Subject", "test")
		// m.SetBody("text/html", "test")
		// err = gomail.Send(s, m)
		// if err != nil {
		// 	this.Throw(99400, "测试邮件发送失败"+err.Error())
		// }

		hfw.CheckErr(settingModel.Save(setting))

	} else {
		this.Throw(99400, "非法请求")
	}
}
