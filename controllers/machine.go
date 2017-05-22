package controllers

import (
	"math"
	"strconv"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Machine struct {
	base
}

func (this *Machine) Before() {
	this.base.Before()
	this.Data["envs"] = map[int]string{
		1: "开发环境",
		2: "测试环境",
		3: "生产环境",
	}
}

func (this *Machine) Index() {
	this.Data["title"] = "服务器列表"

	total, _ := machineModel.Count(models.Cond{})
	this.Data["total"] = int(math.Ceil(float64(total / pageSize)))
	page, _ := strconv.Atoi(this.Request.FormValue("page"))
	page = hfw.Min(hfw.Max(1, page), int(total))
	where := models.Cond{
		"page":     page,
		"pagesize": int(pageSize),
	}
	machines, _ := machineModel.Search(where)
	this.Data["machines"] = machines
	this.Data["prePage"] = page - 1
	this.Data["page"] = page
	this.Data["nextPage"] = page + 1
	this.Data["pageSize"] = pageSize
}

func (this *Machine) Del() {
	id, _ := strconv.Atoi(this.Request.PostFormValue("id"))
	if id > 0 {
		err := machineModel.Update(models.Cond{"is_deleted": "Y"}, models.Cond{"Id": id})
		hfw.CheckErr(err)
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Machine) Add() {
	this.TemplateFile = "machine/edit.html"
	this.Data["title"] = "添加直连服务器"
	this.Data["machine"] = machineModel
}

func (this *Machine) AddInner() {
	this.TemplateFile = "machine/inneredit.html"
	this.Data["title"] = "添加内网服务器"
	this.Data["machine"] = machineModel
}

func (this *Machine) Edit() {
	this.Data["title"] = "修改服务器信息"
	id, _ := strconv.Atoi(this.Request.FormValue("id"))
	if id <= 0 {
		this.Throw(99400, "参数错误")
	}
	machine, _ := machineModel.GetById(id)
	if machine == nil {
		this.Throw(99400, "参数错误")
	}
	this.Data["machine"] = machine
	if machine.InnerIp == "" {
		this.TemplateFile = "machine/edit.html"
	} else {
		this.TemplateFile = "machine/inneredit.html"
	}
}

func (this *Machine) Save() {
	if this.Request.Method == "POST" {
		machine := &models.Machines{}
		Id, _ := strconv.Atoi(this.Request.PostFormValue("Id"))
		if Id > 0 {
			machine, _ = machineModel.GetById(Id)
			if machine == nil {
				this.Throw(99400, "参数错误")
			}
		}
		Env, _ := strconv.Atoi(this.Request.PostFormValue("Env"))
		if Env == 0 {
			this.Throw(99400, "请选择服务器环境")
		}
		machine.Env = Env
		machine.Name = this.Request.PostFormValue("Name")
		machine.Ip = this.Request.PostFormValue("Ip")
		machine.Port = this.Request.PostFormValue("Port")
		machine.User = this.Request.PostFormValue("User")
		machine.Auth = this.Request.PostFormValue("Auth")
		machine.InnerIp = this.Request.PostFormValue("InnerIp")
		machine.InnerPort = this.Request.PostFormValue("InnerPort")
		machine.InnerUser = this.Request.PostFormValue("InnerUser")
		machine.InnerAuth = this.Request.PostFormValue("InnerAuth")
		machine.IsDeleted = "N"

		hfw.CheckErr(machineModel.Save(machine))
	} else {
		this.Throw(99400, "非法请求")
	}
}
