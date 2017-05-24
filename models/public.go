package models

import (
	"encoding/gob"

	"github.com/hsyan2008/hfw"
)

var dao = hfw.NewNoCacheDao()

type Cond map[string]interface{}

var groupModel = &Groups{}
var machineModel = &Machines{}
var messageModel = &Messages{}
var projectModel = &Projectes{}
var settingModel = &Settings{}
var taskModel = &Tasks{}
var taskReviewModel = &TaskReviews{}
var taskProjectesModel = &TaskProjectes{}
var userModel = &Users{}

func init() {
	//gob: type not registered for interface: models.Users
	gob.Register(groupModel)
	gob.Register(machineModel)
	gob.Register(messageModel)
	gob.Register(projectModel)
	gob.Register(settingModel)
	gob.Register(taskProjectesModel)
	gob.Register(taskReviewModel)
	gob.Register(taskModel)
	gob.Register(userModel)
}
