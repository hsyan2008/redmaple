package controllers

import (
	"math"
	"strconv"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/models"
)

type Review struct {
	base
}

func (this *Review) Index() {
	this.Data["title"] = "Review列表"

	where := models.Cond{
		"user_id":    this.User.Id,
		"is_deleted": "N",
		"where":      "status >= 20 AND status < 30",
		"orderby":    "status asc, task_id desc",
	}
	total, _ := taskReviewModel.Count(where)
	this.Data["total"] = int(math.Ceil(float64(total / pageSize)))
	page, _ := strconv.Atoi(this.Request.FormValue("page"))
	page = hfw.Min(hfw.Max(1, page), int(total))
	where["page"] = page
	where["pagesize"] = int(pageSize)
	tmp, _ := taskReviewModel.Search(where)

	taskIds := make([]int, 0)
	for _, v := range tmp {
		taskIds = append(taskIds, v.TaskId)
	}

	tasks, _ := taskModel.GetMulti(taskIds)
	this.Data["tasks"] = tasks
	this.Data["taskStatus"] = taskStatus
	this.Data["prePage"] = page - 1
	this.Data["page"] = page
	this.Data["nextPage"] = page + 1
	this.Data["pageSize"] = pageSize
}

func (this *Review) StartReview() {
	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	taskReview, _ := taskReviewModel.SearchOne(models.Cond{
		"task_id": id,
		"user_id": this.User.Id,
	})

	if taskReview == nil {
		this.Throw(99400, "参数错误")
	}

	task, _ := taskModel.GetById(id)
	if task != nil && task.Status == 20 {

		err = taskModel.Update(models.Cond{"status": "21", "review_user_id": this.User.Id}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		err = taskReviewModel.Update(models.Cond{"status": "21"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)

		err = taskProjectesModel.Update(models.Cond{"status": "21"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)

		this.saveMessage(task.Id, 21, "开始Review")
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Review) ReviewSuccess() {
	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	taskReview, _ := taskReviewModel.SearchOne(models.Cond{
		"task_id": id,
		"user_id": this.User.Id,
	})

	if taskReview == nil {
		this.Throw(99400, "参数错误")
	}

	task, _ := taskModel.GetById(id)
	if task != nil && task.ReviewUserId == this.User.Id && task.Status == 21 {

		err = taskModel.Update(models.Cond{"status": "40"}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		err = taskReviewModel.Update(models.Cond{"status": "40"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)

		err = taskProjectesModel.Update(models.Cond{"status": "40"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)

		this.saveMessage(task.Id, 40, "Review通过")
	} else {
		this.Throw(99400, "参数错误")
	}
}

func (this *Review) ReviewFail() {
	id, err := strconv.Atoi(this.Request.PostFormValue("id"))

	taskReview, _ := taskReviewModel.SearchOne(models.Cond{
		"task_id": id,
		"user_id": this.User.Id,
	})

	if taskReview == nil {
		this.Throw(99400, "参数错误")
	}

	task, _ := taskModel.GetById(id)
	if task != nil && task.ReviewUserId == this.User.Id && task.Status == 21 {

		err = taskModel.Update(models.Cond{"status": "1", "review_user_id": 0}, models.Cond{"Id": id})
		hfw.CheckErr(err)

		err = taskReviewModel.Update(models.Cond{"status": "1"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)

		err = taskProjectesModel.Update(models.Cond{"status": "1"}, models.Cond{"task_id": id})
		hfw.CheckErr(err)

		this.saveMessage(task.Id, 1, this.Request.PostFormValue("msg"))
	} else {
		this.Throw(99400, "参数错误")
	}
}
