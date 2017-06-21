package controllers

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Project struct {
	base
}

func (this *Project) Index() {
	this.Data["title"] = "项目列表"

	total, _ := projectModel.Count(models.Cond{})
	this.Data["total"] = int(math.Ceil(float64(total / pageSize)))
	page, _ := strconv.Atoi(this.Request.FormValue("page"))
	page = hfw.Min(hfw.Max(1, page), int(total))
	where := models.Cond{
		"page":     page,
		"pagesize": int(pageSize),
	}
	this.Data["projectes"], _ = projectModel.Search(where)
	this.Data["prePage"] = page - 1
	this.Data["page"] = page
	this.Data["nextPage"] = page + 1
	this.Data["pageSize"] = pageSize
}

//暂不清理git代码 TODO
func (this *Project) Del() {
	id, _ := strconv.Atoi(this.Request.PostFormValue("id"))
	if id > 0 {
		err := projectModel.Update(models.Cond{"is_deleted": "Y"}, models.Cond{"Id": id})
		hfw.CheckErr(err)
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Project) Add() {
	this.TemplateFile = "project/edit.html"
	this.Data["title"] = "添加项目"
	this.Data["projectes"] = projectModel
	this.Data["devMachines"], _ = machineModel.Search(models.Cond{"env": "1", "is_deleted": "N"})
	this.Data["testMachines"], _ = machineModel.Search(models.Cond{"env": "2", "is_deleted": "N"})
	this.Data["prodMachines"], _ = machineModel.Search(models.Cond{"env": "3", "is_deleted": "N"})

	this.FuncMap["InArray"] = func(id int, machines []*models.Machines) (f bool) {
		for _, v := range machines {
			if v.Id == id {
				return true
			}
		}

		return
	}
}

func (this *Project) Edit() {
	this.Data["title"] = "修改项目信息"
	id, _ := strconv.Atoi(this.Request.FormValue("id"))
	if id <= 0 {
		this.Throw(99400, "参数错误")
	}
	project, _ := projectModel.GetById(id)
	if project == nil {
		this.Throw(99400, "参数错误")
	}
	this.Data["projectes"] = project
	this.Data["devMachines"], _ = machineModel.Search(models.Cond{"env": "1", "is_deleted": "N"})
	this.Data["testMachines"], _ = machineModel.Search(models.Cond{"env": "2", "is_deleted": "N"})
	this.Data["prodMachines"], _ = machineModel.Search(models.Cond{"env": "3", "is_deleted": "N"})

	this.FuncMap["InArray"] = func(id int, machines []*models.Machines) (f bool) {
		for _, v := range machines {
			if v.Id == id {
				return true
			}
		}

		return
	}
}

func (this *Project) Save() {
	if this.Request.Method == "POST" {
		project := &models.Projectes{}
		Id, _ := strconv.Atoi(this.Request.PostFormValue("Id"))
		IsLock := "N"
		if Id > 0 {
			project, _ = projectModel.GetById(Id)
			if project == nil {
				this.Throw(99400, "参数错误")
			}
			IsLock = project.IsLock
		} else {
			Name := this.Request.PostFormValue("Name")
			tmp, _ := projectModel.SearchOne(models.Cond{"Name": Name})
			if tmp != nil {
				this.Throw(99400, "项目已存在")
			}
			project.Name = Name
		}

		Git := this.Request.PostFormValue("Git")
		Wwwroot := this.Request.PostFormValue("Wwwroot")

		if IsLock == "Y" && (Git != project.Git || Wwwroot != project.Wwwroot) {
			this.Throw(99400, "项目在开发中，请勿修改git和wwwroot")
		}

		DevWwwroot := this.Request.PostFormValue("DevWwwroot")
		devMachineIds := this.Request.Form["devMachineIds"]
		if len(devMachineIds) > 0 {
			tmp, _ := machineModel.GetMulti(devMachineIds)
			if len(tmp) != len(devMachineIds) {
				this.Throw(99400, "参数错误")
			}
		}
		DevMachineIds := strings.Join(devMachineIds, ",")
		DevAfterRelease := this.Request.PostFormValue("DevAfterRelease")

		TestWwwroot := this.Request.PostFormValue("TestWwwroot")
		testMachineIds := this.Request.Form["testMachineIds"]
		if len(testMachineIds) > 0 {
			tmp, _ := machineModel.GetMulti(testMachineIds)
			if len(tmp) != len(testMachineIds) {
				this.Throw(99400, "参数错误")
			}
		} else {
			this.Throw(99400, "请选择测试服务器")
		}
		TestMachineIds := strings.Join(testMachineIds, ",")
		TestAfterRelease := this.Request.PostFormValue("TestAfterRelease")

		ProdWwwroot := this.Request.PostFormValue("ProdWwwroot")
		prodMachineIds := this.Request.Form["prodMachineIds"]
		if len(prodMachineIds) > 0 {
			tmp, _ := machineModel.GetMulti(prodMachineIds)
			if len(tmp) != len(prodMachineIds) {
				this.Throw(99400, "参数错误")
			}
		} else {
			this.Throw(99400, "请选择生产服务器")
		}
		ProdMachineIds := strings.Join(prodMachineIds, ",")
		ProdAfterRelease := this.Request.PostFormValue("ProdAfterRelease")

		//开发环境没有的话，将不部署
		if Wwwroot == "" && (TestWwwroot == "" || ProdWwwroot == "") {
			this.Throw(99400, "缺少代码部署路径，请检查所有环境设置")
		}

		//wwwroot有变化，删除原来的
		if project.Wwwroot != "" && project.Wwwroot != Wwwroot {
			// hfw.CheckErr(os.RemoveAll(project.Wwwroot))
			hfw.CheckErr(os.Rename(project.Wwwroot, Wwwroot))
		}

		GitTools.Lock()
		defer GitTools.Unlock()

		//如果更改git地址
		if project.Git != "" && Git != project.Git {
			//初始化
			hfw.CheckErr(os.RemoveAll(Wwwroot))
			GitTools.Clone(Git, Wwwroot)
		}

		if project.Id == 0 {
			hfw.CheckErr(os.MkdirAll(filepath.Dir(Wwwroot), 0755))

			//初始化
			GitTools.Clone(Git, Wwwroot)
			//建立测试和开发分支
			GitTools.ReBranch(Wwwroot, "pre_release", "test")
		}

		project.Git = Git
		project.Wwwroot = Wwwroot
		project.DevWwwroot = DevWwwroot
		project.DevMachineIds = DevMachineIds
		project.DevAfterRelease = DevAfterRelease
		project.TestWwwroot = TestWwwroot
		project.TestMachineIds = TestMachineIds
		project.TestAfterRelease = TestAfterRelease
		project.ProdWwwroot = ProdWwwroot
		project.ProdMachineIds = ProdMachineIds
		project.ProdAfterRelease = ProdAfterRelease
		project.IsLock = IsLock
		project.IsDeleted = "N"

		hfw.CheckErr(projectModel.Save(project))
	} else {
		this.Throw(99400, "非法请求")
	}
}
