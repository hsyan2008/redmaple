package models

import (
	"time"
)

type TaskProjectes struct {
	Id          int        `xorm:"not null pk autoincr INT(10)"`
	TaskId      int        `xorm:"not null default 0 index INT(10)"`
	ProjectId   int        `xorm:"not null default 0 INT(10)"`
	Status      int        `xorm:"not null default 0 TINYINT(3)"`
	StartCommit string     `xorm:"not null default '' CHAR(7)"`
	EndCommit   string     `xorm:"not null default '' CHAR(7)"`
	IsPatch     string     `xorm:"not null default 'N' ENUM('Y','N')"`
	IsMerge     string     `xorm:"not null default 'N' ENUM('Y','N')"`
	IsFinish    string     `xorm:"not null default 'N' ENUM('Y','N')"`
	IsDeleted   string     `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt   time.Time  `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt   time.Time  `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
	Project     *Projectes `xorm:"-"`
}

func (this *TaskProjectes) TableName() string {

	return "task_projectes"
}

func (this *TaskProjectes) Save(t *TaskProjectes) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *TaskProjectes) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *TaskProjectes) SearchOne(cond Cond) (t *TaskProjectes, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *TaskProjectes) Search(cond Cond) (t []*TaskProjectes, err error) {

	if user_id, ok := cond["user_id"]; ok {
		if user_id == 1 {
			delete(cond, "user_id")
		}
	}

	err = dao.Search(&t, cond)

	for k, v := range t {
		t[k].Project, _ = projectModel.SearchOne(Cond{"id": v.ProjectId})
	}

	return
}

func (this *TaskProjectes) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *TaskProjectes) GetMulti(ids ...interface{}) (t []*TaskProjectes, err error) {
	err = dao.GetMulti(&t, ids...)

	for k, v := range t {
		t[k].Project, _ = projectModel.GetById(v.ProjectId)
	}

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *TaskProjectes) GetById(id interface{}) (t *TaskProjectes, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}
