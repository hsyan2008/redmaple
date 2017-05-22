package main

import (
	"os"
	"os/user"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/controllers"
)

func main() {

	hfw.Debug("Pid:", os.Getpid(), "Starting ...")
	defer hfw.Debug("Pid:", os.Getpid(), "Shutdown complete!")

	user, _ := user.Current()
	if user.Uid == "0" {
		panic("不能以root用户运行")
	}

	//两端/可以都不写
	hfw.RegisterRoute("/", &controllers.Index{})
	hfw.RegisterRoute("/machine", &controllers.Machine{})
	hfw.RegisterRoute("/task", &controllers.Task{})
	hfw.RegisterRoute("/user", &controllers.User{})
	hfw.RegisterRoute("/project", &controllers.Project{})
	hfw.RegisterRoute("/review", &controllers.Review{})
	hfw.RegisterRoute("/test", &controllers.Test{})
	hfw.RegisterRoute("/release", &controllers.Release{})

	hfw.RegisterStatic("/js/", hfw.Config.Template.Static)
	hfw.RegisterStatic("/css/", hfw.Config.Template.Static)
	hfw.RegisterStatic("/images/", hfw.Config.Template.Static)
	hfw.RegisterStatic("/fonts/", hfw.Config.Template.Static)
	hfw.RegisterFile("/favicon.ico", hfw.Config.Template.Static+"/images/")

	hfw.Run()
}
