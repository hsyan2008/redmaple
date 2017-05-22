package models

import (
	"time"
)

type TaskReviews struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	TaskId    int       `xorm:"not null default 0 unique(uidx_task_id_user_id) INT(10)"`
	UserId    int       `xorm:"not null default 0 unique(uidx_task_id_user_id) index(idx_user_id_status) INT(10)"`
	Status    int       `xorm:"not null default 0 index(idx_user_id_status) TINYINT(3)"`
	IsDeleted string    `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt time.Time `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
	User      *Users    `xorm:"-"`
}

func (this *TaskReviews) TableName() string {

	return "task_reviews"
}

func (this *TaskReviews) Save(t *TaskReviews) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *TaskReviews) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *TaskReviews) SearchOne(cond Cond) (t *TaskReviews, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *TaskReviews) Search(cond Cond) (t []*TaskReviews, err error) {

	err = dao.Search(&t, cond)

	for k, v := range t {
		t[k].User, _ = userModel.SearchOne(Cond{"id": v.UserId})
	}

	return
}

func (this *TaskReviews) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *TaskReviews) GetMulti(ids ...interface{}) (t []*TaskReviews, err error) {
	err = dao.GetMulti(&t, ids...)

	for k, v := range t {
		t[k].User, _ = userModel.GetById(v.UserId)
	}

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *TaskReviews) GetById(id interface{}) (t *TaskReviews, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}
