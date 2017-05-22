package models

import (
	"strings"
	"time"
)

type Projectes struct {
	Id               int         `xorm:"not null pk autoincr INT(11)"`
	Name             string      `xorm:"not null default '' unique VARCHAR(32)"`
	Git              string      `xorm:"not null default '' VARCHAR(255)"`
	Wwwroot          string      `xorm:"not null default '' VARCHAR(255)"`
	DevWwwroot       string      `xorm:"not null default '' VARCHAR(255)"`
	DevMachineIds    string      `xorm:"not null default '' VARCHAR(64)"`
	DevAfterRelease  string      `xorm:"not null TEXT"`
	TestWwwroot      string      `xorm:"not null default '' VARCHAR(255)"`
	TestMachineIds   string      `xorm:"not null default '' VARCHAR(64)"`
	TestAfterRelease string      `xorm:"not null TEXT"`
	ProdWwwroot      string      `xorm:"not null default '' VARCHAR(255)"`
	ProdMachineIds   string      `xorm:"not null default '' VARCHAR(64)"`
	ProdAfterRelease string      `xorm:"not null TEXT"`
	IsLock           string      `xorm:"not null default 'N' ENUM('Y','N')"`
	IsDeleted        string      `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt        time.Time   `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt        time.Time   `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
	DevMachines      []*Machines `xorm:"-"`
	TestMachines     []*Machines `xorm:"-"`
	ProdMachines     []*Machines `xorm:"-"`
}

func (this *Projectes) TableName() string {

	return "projectes"
}

func (this *Projectes) Save(t *Projectes) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *Projectes) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *Projectes) SearchOne(cond Cond) (t *Projectes, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *Projectes) Search(cond Cond) (t []*Projectes, err error) {

	err = dao.Search(&t, cond)

	for key, val := range t {
		if val.DevMachineIds != "" {
			machine_ids := strings.Split(val.DevMachineIds, ",")
			t[key].DevMachines, _ = machineModel.GetMulti(machine_ids)
		}
		if val.TestMachineIds != "" {
			machine_ids := strings.Split(val.TestMachineIds, ",")
			t[key].TestMachines, _ = machineModel.GetMulti(machine_ids)
		}
		if val.ProdMachineIds != "" {
			machine_ids := strings.Split(val.ProdMachineIds, ",")
			t[key].ProdMachines, _ = machineModel.GetMulti(machine_ids)
		}
	}

	return
}

func (this *Projectes) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *Projectes) GetMulti(ids ...interface{}) (t []*Projectes, err error) {

	err = dao.GetMulti(&t, ids...)

	for key, val := range t {
		if val.DevMachineIds != "" {
			machine_ids := strings.Split(val.DevMachineIds, ",")
			t[key].DevMachines, _ = machineModel.GetMulti(machine_ids)
		}
		if val.TestMachineIds != "" {
			machine_ids := strings.Split(val.TestMachineIds, ",")
			t[key].TestMachines, _ = machineModel.GetMulti(machine_ids)
		}
		if val.ProdMachineIds != "" {
			machine_ids := strings.Split(val.ProdMachineIds, ",")
			t[key].ProdMachines, _ = machineModel.GetMulti(machine_ids)
		}
	}

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *Projectes) GetById(id interface{}) (t *Projectes, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}
