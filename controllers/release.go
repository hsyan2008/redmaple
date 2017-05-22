package controllers

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Release struct {
	base
}

func (this *Release) Index() {
	this.Data["title"] = "上线列表"

	where := models.Cond{
		"is_deleted": "N",
		"where":      "status >= 60",
	}
	total, _ := taskModel.Count(where)
	this.Data["total"] = int(math.Ceil(float64(total / pageSize)))
	page, _ := strconv.Atoi(this.Request.FormValue("page"))
	page = hfw.Min(hfw.Max(1, page), int(total))
	where["page"] = page
	where["pagesize"] = int(pageSize)
	tasks, _ := taskModel.Search(where)

	this.Data["tasks"] = tasks
	this.Data["taskStatus"] = taskStatus
	this.Data["prePage"] = page - 1
	this.Data["page"] = page
	this.Data["nextPage"] = page + 1
	this.Data["pageSize"] = pageSize
}

//是否需要重建test和pre_release TODO
func (this *Release) ReleaseSuccess() {

	GitTools.Lock()
	defer GitTools.Unlock()

	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	task, _ := taskModel.GetById(id)
	if task != nil && task.Status == 80 && this.User.GroupId == 1 {

		for _, val := range task.TaskProjectes {
			if val.IsFinish == "N" && val.EndCommit != val.StartCommit {
				err = GitTools.Merge(val.Project.Wwwroot, task.Code, "master", val.StartCommit, val.EndCommit, fmt.Sprintf("taskCode:%s\n%s", task.Code, task.Comment), fmt.Sprintf("%s <%s>", task.User.Realname, task.User.Email))
				hfw.CheckErr(err)
			}

			err = taskProjectesModel.Update(models.Cond{"status": "81", "is_finish": "Y"}, models.Cond{"id": val.Id})
			GitTools.DelBranch(val.Project.Wwwroot, task.Code)
		}

		err = taskReviewModel.Update(models.Cond{"status": "81"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)
		err = taskModel.Update(models.Cond{"status": "81"}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		saveMessage(id, this.User.Id, 81, "上线成功完成")
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Release) ReleaseFail() {

	GitTools.Lock()
	defer GitTools.Unlock()

	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	task, _ := taskModel.GetById(id)
	if task != nil && task.Status == 80 && this.User.GroupId == 1 {
		err = taskModel.Update(models.Cond{"status": "3"}, models.Cond{"Id": id})
		_ = taskReviewModel.Update(models.Cond{"status": "3"}, models.Cond{"task_id": id})
		_ = taskProjectesModel.Update(models.Cond{"status": "3",
			"is_patch": "N", "is_merge": "N"}, models.Cond{"task_id": id})

		var wg = &sync.WaitGroup{}
		//排除本任务，重新部署和本任务相关的分支的test和pre_release
		for _, val := range task.TaskProjectes {
			//本分支没有改动，不需要重新建test和pre_release
			if val.EndCommit == val.StartCommit {
				continue
			}

			//重建pre_release
			GitTools.ReBranch(val.Project.Wwwroot, "pre_release")
			taskProjectes, _ := taskProjectesModel.Search(models.Cond{
				"project_id": val.ProjectId,
				"where":      "status = 80",
			})
			for _, v := range taskProjectes {
				tmpTask, _ := taskModel.GetById(v.TaskId)
				//没有更改代码的，不会重新合并，免得报错
				err = GitTools.Merge(v.Project.Wwwroot, tmpTask.Code, "pre_release", v.StartCommit, v.EndCommit, fmt.Sprintf("taskCode:%s\n%s", tmpTask.Code, tmpTask.Comment), fmt.Sprintf("%s <%s>", tmpTask.User.Realname, tmpTask.User.Email))
				hfw.CheckErr(err)
			}
			for _, v := range val.Project.ProdMachines {
				wg.Add(1)
				go func() {
					_ = release("pre_release", wg, val.Project, v)
				}()
			}

			//重建test
			GitTools.ReBranch(val.Project.Wwwroot, "test")
			taskProjectes, _ = taskProjectesModel.Search(models.Cond{
				"project_id": val.ProjectId,
				"where":      "status >= 43 AND status <= 80",
				"orderby":    "status desc", //优先合并测试通过的分支
			})
			for _, v := range taskProjectes {
				tmpTask, _ := taskModel.GetById(v.TaskId)
				//没有更改代码的，不会重新合并，免得报错
				err = GitTools.Merge(v.Project.Wwwroot, tmpTask.Code, "test", v.StartCommit, v.EndCommit, fmt.Sprintf("taskCode:%s\n%s", tmpTask.Code, tmpTask.Comment), fmt.Sprintf("%s <%s>", tmpTask.User.Realname, tmpTask.User.Email))
				hfw.CheckErr(err)
			}
			for _, v := range val.Project.TestMachines {
				wg.Add(1)
				go func() {
					_ = release("test", wg, val.Project, v)
				}()
			}

		}
		wg.Wait()

		saveMessage(id, this.User.Id, 3, this.Request.PostFormValue("msg"))
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Release) ToRelease() {

	GitTools.Lock()
	defer GitTools.Unlock()

	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	task, _ := taskModel.GetById(id)
	if task != nil && (task.Status == 60 || task.Status == 62) && this.User.GroupId == 1 {

		err = taskModel.Update(models.Cond{"status": "61"}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		defer func() {
			if err := recover(); err != nil {
				_ = taskModel.Update(models.Cond{"status": "62"}, models.Cond{"Id": id})
				_ = taskReviewModel.Update(models.Cond{"status": "62"}, models.Cond{"task_id": id})
				_ = taskProjectesModel.Update(models.Cond{"status": "62"}, models.Cond{"task_id": id})
				panic(err)
			}
		}()

		var wg = &sync.WaitGroup{}
		for _, val := range task.TaskProjectes {
			//把代码合并到test
			if val.StartCommit != val.EndCommit {
				if val.IsMerge == "N" {
					err = GitTools.Merge(val.Project.Wwwroot, task.Code, "pre_release", val.StartCommit, val.EndCommit, fmt.Sprintf("taskCode:%s\n%s", task.Code, task.Comment), fmt.Sprintf("%s <%s>", task.User.Realname, task.User.Email))
					hfw.CheckErr(err)
					err = taskProjectesModel.Update(models.Cond{"is_merge": "Y"}, models.Cond{"id": val.Id})
				}
				//对于每台机器，并发发布
				for _, v := range val.Project.ProdMachines {
					wg.Add(1)
					go func() {
						_ = release("pre_release", wg, val.Project, v)
					}()
				}
			}
		}

		wg.Wait()

		err = taskModel.Update(models.Cond{"status": "80"}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		_ = taskReviewModel.Update(models.Cond{"status": "80"}, models.Cond{"task_id": id})
		_ = taskProjectesModel.Update(models.Cond{"status": "80"}, models.Cond{"task_id": id})

		saveMessage(id, this.User.Id, 80, "上线等待验证")
	} else {
		this.Throw(99400, "参数错误")
	}
}
