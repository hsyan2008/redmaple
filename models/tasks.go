package models

import (
	"time"
)

type Tasks struct {
	Id            int              `xorm:"not null pk autoincr INT(10)"`
	Code          string           `xorm:"not null default '' unique CHAR(20)"`
	Name          string           `xorm:"not null default '' VARCHAR(32)"`
	UserId        int              `xorm:"not null default 0 index INT(10)"`
	ReviewUserId  int              `xorm:"not null default 0 INT(10)"`
	TestUserId    int              `xorm:"not null default 0 INT(10)"`
	Comment       string           `xorm:"not null TEXT"`
	Status        int              `xorm:"not null default 0 TINYINT(3)"`
	IsDeleted     string           `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt     time.Time        `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt     time.Time        `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
	User          *Users           `xorm:"-"`
	ReviewUser    *Users           `xorm:"-"`
	TestUser      *Users           `xorm:"-"`
	TaskProjectes []*TaskProjectes `xorm:"-"`
	TaskReviews   []*TaskReviews   `xorm:"-"`
}

func (this *Tasks) TableName() string {

	return "tasks"
}

func (this *Tasks) Save(t *Tasks) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *Tasks) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *Tasks) SearchOne(cond Cond) (t *Tasks, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *Tasks) Search(cond Cond) (t []*Tasks, err error) {

	err = dao.Search(&t, cond)

	for k, v := range t {
		t[k].TaskReviews, _ = taskReviewModel.Search(Cond{
			"task_id": v.Id,
			"orderby": "id asc",
		})
		t[k].TaskProjectes, _ = taskProjectesModel.Search(Cond{
			"task_id": v.Id,
			"orderby": "id asc",
		})
		t[k].User, _ = userModel.GetById(v.UserId)
		t[k].ReviewUser, _ = userModel.GetById(v.ReviewUserId)
		t[k].TestUser, _ = userModel.GetById(v.TestUserId)
	}

	return
}

func (this *Tasks) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *Tasks) GetMulti(ids ...interface{}) (t []*Tasks, err error) {
	err = dao.GetMulti(&t, ids...)

	for k, v := range t {
		t[k].TaskReviews, _ = taskReviewModel.Search(Cond{
			"task_id": v.Id,
			"orderby": "id asc",
		})
		t[k].TaskProjectes, _ = taskProjectesModel.Search(Cond{
			"task_id": v.Id,
			"orderby": "id asc",
		})
		t[k].User, _ = userModel.GetById(v.UserId)
		t[k].ReviewUser, _ = userModel.GetById(v.ReviewUserId)
		t[k].TestUser, _ = userModel.GetById(v.TestUserId)
	}

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *Tasks) GetById(id interface{}) (t *Tasks, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}
