package controllers

import (
	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Index struct {
	base
}

func (this *Index) Index() {

	if this.Session.IsExist("userinfo") {
		// this.Redirect("/admin", http.StatusFound)
		this.Data["title"] = "首页"
		this.TemplateFile = "index/index.html"

		this.Data["messages"], _ = messageModel.Search(models.Cond{"orderby": "id desc", "pageSize": 10})

		this.Data["tasks"], _ = taskModel.Search(models.Cond{"orderby": "id desc", "pageSize": 10})
		this.Data["taskStatus"] = taskStatus
		this.FuncMap["getProcess"] = func(status int) int {
			return status * 100 / 81
		}

		this.Data["users"], _ = userModel.Search(models.Cond{"is_deleted": "N", "orderby": "id desc"})
		tmpGroups, _ := groupModel.Search(models.Cond{"is_deleted": "N"})
		groups := make(map[int]string)
		for _, v := range tmpGroups {
			groups[v.Id] = v.Name
		}

		this.FuncMap["getGroupName"] = func(groupId int) string {
			return groups[groupId]
		}
	} else {
		this.Data["title"] = "登录"
		this.TemplateFile = "index/login.html"
	}
}

func (this *Index) Login() {
	username := this.Request.PostFormValue("username")
	password := this.Request.PostFormValue("password")

	user, _ := userModel.SearchOne(models.Cond{"Name": username})

	if user == nil {
		this.Throw(99400, "用户名或密码错误")
	}

	if user.IsDeleted == "Y" {
		this.Throw(99400, "该用户已被禁用")
	}

	if user.Password == hfw.Md5(password+user.Salt) {
		this.Session.Set("userinfo", user)
		this.Results = "登陆成功"
	} else {
		this.Throw(99400, "用户名或密码错误")
	}
}

func (this *Index) Logout() {
	this.Session.Destroy()
	this.Redirect("/")
}
