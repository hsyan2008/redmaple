package controllers

import (
	"math"
	"strconv"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"

	"github.com/pborman/uuid"
)

type User struct {
	base
}

func (this *User) Index() {
	this.Data["title"] = "用户列表"

	tmpGroups, _ := groupModel.Search(models.Cond{"orderby": "id asc"})
	groups := make(map[int]string)
	for _, v := range tmpGroups {
		groups[v.Id] = v.Name
	}
	this.Data["groups"] = groups

	total, _ := userModel.Count(models.Cond{})
	this.Data["total"] = int(math.Ceil(float64(total / pageSize)))
	page, _ := strconv.Atoi(this.Request.FormValue("page"))
	page = hfw.Min(hfw.Max(1, page), int(total))
	where := models.Cond{
		"page":     page,
		"pagesize": int(pageSize),
	}
	users, _ := userModel.Search(where)
	this.Data["users"] = users
	this.Data["prePage"] = page - 1
	this.Data["page"] = page
	this.Data["nextPage"] = page + 1
	this.Data["pageSize"] = pageSize
}

func (this *User) Del() {
	id, _ := strconv.Atoi(this.Request.PostFormValue("id"))
	if id > 0 {
		err := userModel.Update(models.Cond{"is_deleted": "Y"}, models.Cond{"Id": id})
		hfw.CheckErr(err)
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *User) Restore() {
	id, _ := strconv.Atoi(this.Request.PostFormValue("id"))
	if id > 0 {
		err := userModel.Update(models.Cond{"is_deleted": "N"}, models.Cond{"Id": id})
		hfw.CheckErr(err)
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *User) Add() {
	this.TemplateFile = "user/edit.html"
	this.Data["title"] = "添加用户"
	this.Data["user"] = userModel
	groups, _ := groupModel.Search(models.Cond{"orderby": "id asc"})
	this.Data["groups"] = groups
}

func (this *User) Edit() {
	this.Data["title"] = "修改用户信息"
	id, _ := strconv.Atoi(this.Request.FormValue("id"))
	if id <= 0 {
		this.Throw(99400, "参数错误")
	}
	user, _ := userModel.GetById(id)
	if user == nil {
		this.Throw(99400, "参数错误")
	}
	this.Data["user"] = user
	groups, _ := groupModel.Search(models.Cond{"orderby": "id asc"})
	this.Data["groups"] = groups
}

func (this *User) Save() {
	if this.Request.Method == "POST" {
		user := &models.Users{}
		password := this.Request.PostFormValue("password")
		id, _ := strconv.Atoi(this.Request.PostFormValue("id"))
		if id > 0 {
			user, _ = userModel.GetById(id)
			if user == nil {
				this.Throw(99400, "参数错误")
			}
		} else {
			if password == "" {
				this.Throw(99400, "请填写密码")
			}
			groupId, _ := strconv.Atoi(this.Request.PostFormValue("group_id"))
			if groupId == 0 {
				this.Throw(99400, "请选择分组")
			}

			name := this.Request.PostFormValue("name")
			tmp, _ := userModel.SearchOne(models.Cond{"Name": name})
			if tmp != nil {
				this.Throw(99400, "用户已存在")
			}
			user.Name = name

			user.GroupId = groupId

			user.IsDeleted = "N"
		}
		if password != "" {
			user.Salt = hfw.Md5(uuid.New())
			user.Password = hfw.Md5(password + user.Salt)
		}
		user.Realname = this.Request.PostFormValue("realname")
		user.Email = this.Request.PostFormValue("email")

		hfw.CheckErr(userModel.Save(user))
	} else {
		this.Throw(99400, "非法请求")
	}
}

func (this *User) Profile() {
	this.Data["title"] = "修改个人信息"
	if this.Request.Method == "POST" {
		user := this.Session.Get("userinfo").(models.Users)
		oldpassword := this.Request.PostFormValue("oldpassword")
		password := this.Request.PostFormValue("password")
		if password != "" {
			if user.Password != hfw.Md5(oldpassword+user.Salt) {
				this.Throw(99400, "原密码错误")
			}
			user.Salt = hfw.Md5(uuid.New())
			user.Password = hfw.Md5(password + user.Salt)
		}
		user.Realname = this.Request.PostFormValue("realname")
		user.Email = this.Request.PostFormValue("email")

		hfw.CheckErr(userModel.Save(&user))
		this.Session.Set("userinfo", user)
	} else {
		this.Data["user"] = this.Session.Get("userinfo").(*models.Users)
		groups, _ := groupModel.Search(models.Cond{"orderby": "id asc"})
		this.Data["groups"] = groups
	}
}
