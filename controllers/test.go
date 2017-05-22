package controllers

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Test struct {
	base
}

func (this *Test) Index() {
	this.Data["title"] = "测试列表"

	where := models.Cond{
		"is_deleted": "N",
		"where":      "status >= 40 AND status < 50",
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

func (this *Test) TestSuccess() {
	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	task, _ := taskModel.GetById(id)
	if task != nil && task.Status == 43 && task.TestUserId == this.User.Id {
		err = taskProjectesModel.Update(models.Cond{"status": "60"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)
		err = taskReviewModel.Update(models.Cond{"status": "60"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)
		err = taskModel.Update(models.Cond{"status": "60"}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		saveMessage(id, this.User.Id, 60, "测试通过")
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Test) TestFail() {

	GitTools.Lock()
	defer GitTools.Unlock()

	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	task, _ := taskModel.GetById(id)
	if task != nil && task.Status == 43 && task.TestUserId == this.User.Id {
		err = taskModel.Update(models.Cond{"status": "2"}, models.Cond{"Id": id})
		_ = taskReviewModel.Update(models.Cond{"status": "2"}, models.Cond{"task_id": id})
		_ = taskProjectesModel.Update(models.Cond{"status": "2",
			"is_patch": "N"}, models.Cond{"task_id": id})

		var wg = &sync.WaitGroup{}
		//排除本任务，重新部署和本任务相关的分支的test
		for _, val := range task.TaskProjectes {
			//本分支没有改动，不需要重新建test
			if val.EndCommit == val.StartCommit {
				continue
			}
			GitTools.ReBranch(val.Project.Wwwroot, "test")

			taskProjectes, _ := taskProjectesModel.Search(models.Cond{
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

		saveMessage(id, this.User.Id, 2, this.Request.PostFormValue("msg"))
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Test) StartTest() {

	//防止并发
	GitTools.Lock()
	defer GitTools.Unlock()

	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	task, _ := taskModel.GetById(id)
	if task != nil && (task.Status == 40 || task.Status == 42) {

		err = taskModel.Update(models.Cond{"status": "41", "test_user_id": this.User.Id}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		defer func() {
			if err := recover(); err != nil {
				_ = taskModel.Update(models.Cond{"status": "42"}, models.Cond{"Id": id})
				_ = taskReviewModel.Update(models.Cond{"status": "42"}, models.Cond{"task_id": id})
				_ = taskProjectesModel.Update(models.Cond{"status": "42"}, models.Cond{"task_id": id})
				panic(err)
			}
		}()

		var wg = &sync.WaitGroup{}
		for _, val := range task.TaskProjectes {
			//把代码合并到test
			if val.StartCommit != val.EndCommit {
				if val.IsPatch == "N" {
					err = GitTools.Merge(val.Project.Wwwroot, task.Code, "test", val.StartCommit, val.EndCommit, fmt.Sprintf("taskCode:%s\n%s", task.Code, task.Comment), fmt.Sprintf("%s <%s>", task.User.Realname, task.User.Email))
					hfw.CheckErr(err)
					err = taskProjectesModel.Update(models.Cond{"is_patch": "Y"}, models.Cond{"id": val.Id})
				}
				//对于每台机器，并发发布
				for _, v := range val.Project.TestMachines {
					wg.Add(1)
					go func() {
						_ = release("test", wg, val.Project, v)
					}()
				}
			}
		}

		wg.Wait()

		err = taskModel.Update(models.Cond{"status": "43"}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		_ = taskReviewModel.Update(models.Cond{"status": "43"}, models.Cond{"task_id": id})
		_ = taskProjectesModel.Update(models.Cond{"status": "43"}, models.Cond{"task_id": id})

		saveMessage(id, this.User.Id, 43, "开始测试")
	} else {
		this.Throw(99400, "参数错误")
	}
}
