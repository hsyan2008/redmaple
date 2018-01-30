package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Task struct {
	base
}

func (this *Task) Index() {
	this.Data["title"] = "任务列表"

	where := models.Cond{
		"user_id":    this.User.Id,
		"is_deleted": "N",
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

func (this *Task) Detail() {
	id, _ := strconv.Atoi(this.Request.FormValue("id"))
	task, _ := taskModel.GetById(id)
	if task == nil {
		this.Throw(99400, "参数错误")
	}

	this.Data["title"] = "任务详情"
	this.Data["task"] = task
	this.Data["taskStatus"] = taskStatus

	this.Data["messages"], _ = messageModel.Search(models.Cond{"task_id": id, "orderby": "id desc"})
}

//只有status为0才可以删除，表示放弃本任务
func (this *Task) Del() {
	id, err := strconv.Atoi(this.Request.PostFormValue("id"))
	task, _ := taskModel.GetById(id)
	if task != nil && task.Status == 0 {

		err = taskModel.Update(models.Cond{"is_deleted": "Y", "status": 9}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		tmpTaskProjectes, _ := taskProjectesModel.Search(models.Cond{
			"task_id": task.Id,
		})
		GitTools.Lock()
		defer GitTools.Unlock()
		for _, val := range tmpTaskProjectes {
			project, _ := projectModel.GetById(val.ProjectId)
			GitTools.DelBranch(project.Wwwroot, task.Branch)
		}
		err = taskProjectesModel.Update(models.Cond{"is_deleted": "Y"},
			models.Cond{"task_id": id})
		hfw.CheckErr(err)

		err = taskReviewModel.Update(models.Cond{"is_deleted": "Y"},
			models.Cond{"task_id": id})
		hfw.CheckErr(err)

		this.saveMessage(task.Id, 9, "删除")
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Task) Add() {
	this.TemplateFile = "task/edit.html"
	this.Data["title"] = "创建任务"
	this.Data["task"] = taskModel

	users, _ := userModel.Search(models.Cond{
		"orderby":    "id asc",
		"is_deleted": "N",
	})
	this.Data["users"] = users
	this.Data["user_id"] = this.User.Id

	projectes, _ := projectModel.Search(models.Cond{
		"orderby":    "id asc",
		"is_deleted": "N",
	})
	this.Data["projectes"] = projectes
}

func (this *Task) Save() {
	if this.Request.Method == "POST" {
		task := &models.Tasks{}

		Code := this.Request.PostFormValue("Code")
		Name := this.Request.PostFormValue("Name")
		ReviewUserId, _ := strconv.Atoi(this.Request.PostFormValue("ReviewUserId"))
		projectIds := this.Request.Form["projectIds"]
		Comment := this.Request.PostFormValue("Comment")

		if Code == "" || Name == "" || ReviewUserId <= 0 || len(projectIds) == 0 || Comment == "" {
			this.Throw(99400, "请完善信息")
		}

		// task.Code = fmt.Sprintf("%06d%s", this.User.Id, time.Now().Format("20060102150405"))
		task.Code = Code
		task.Branch = fmt.Sprintf("%s_%d", task.Code, time.Now().Unix())
		task.Name = Name
		task.UserId = this.User.Id
		task.Comment = Comment
		task.Status = 0
		task.IsDeleted = "N"

		user, _ := userModel.GetById(ReviewUserId)
		if user == nil {
			this.Throw(99400, "错误的用户id")
		}

		Projectes := make([]*models.Projectes, 0)
		for _, id := range projectIds {
			project, _ := projectModel.GetById(id)
			if project == nil {
				this.Throw(99400, "错误的项目编号")
			}
			Projectes = append(Projectes, project)
		}

		hfw.CheckErr(taskModel.Save(task))

		taskProjectes := &models.TaskProjectes{}
		GitTools.Lock()
		defer GitTools.Unlock()
		for _, v := range Projectes {
			GitTools.NewBranch(v.Wwwroot, task.Branch)
			taskProjectes.StartCommit, _, _ = GitTools.GetCommitId(v.Wwwroot, task.Branch)
			taskProjectes.Id = 0
			taskProjectes.IsPatch = "N"
			taskProjectes.IsMerge = "N"
			taskProjectes.IsFinish = "N"
			taskProjectes.IsDeleted = "N"
			taskProjectes.TaskId = task.Id
			taskProjectes.ProjectId = v.Id
			taskProjectes.Status = task.Status
			hfw.CheckErr(taskProjectesModel.Save(taskProjectes))
		}

		taskReview := &models.TaskReviews{}
		taskReview.TaskId = task.Id
		taskReview.UserId = ReviewUserId
		taskReview.Status = task.Status
		taskReview.IsDeleted = "N"
		hfw.CheckErr(taskReviewModel.Save(taskReview))

		this.saveMessage(task.Id, task.Status, "创建")
	} else {
		this.Throw(99400, "非法请求")
	}
}

func (this *Task) ToReview() {
	id, err := strconv.Atoi(this.Request.PostFormValue("id"))
	task, _ := taskModel.GetById(id)
	if task != nil && task.Id > 0 && task.Status < 10 {

		taskProjectes, _ := taskProjectesModel.Search(models.Cond{"task_id": id})
		var isChange bool //代码是否有变动
		GitTools.Lock()
		defer GitTools.Unlock()
		for _, v := range taskProjectes {
			endCommit, _, _ := GitTools.GetCommitId(v.Project.Wwwroot, task.Branch)
			if endCommit != v.StartCommit {
				isChange = true
			}
			err = taskProjectesModel.Update(models.Cond{"status": "20", "end_commit": endCommit},
				models.Cond{"id": v.Id})
			hfw.CheckErr(err)
			defer func(id int) {
				if isChange == false {
					err = taskProjectesModel.Update(models.Cond{"status": "0"}, models.Cond{"id": id})
				}
			}(v.Id)
		}
		if isChange == false {
			this.Throw(99400, "你并没有提交任何变更")
		}

		err = taskModel.Update(models.Cond{"status": "20"}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		err = taskReviewModel.Update(models.Cond{"status": "20"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)

		this.saveMessage(task.Id, 20, "提交Review")

		users := make([]*models.Users, 0)
		for _, v := range task.TaskReviews {
			users = append(users, v.User)
		}
		this.sendMail(task, "提交Reveiw了，请尽快进行Review", users...)

	} else {
		this.Throw(99400, "参数错误")
	}
}
