package main

import (
	"os"

	"github.com/hsyan2008/go-logger/logger"
	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/controllers"
)

func main() {

	logger.Debug("Pid:", os.Getpid(), "Starting ...")
	defer logger.Debug("Pid:", os.Getpid(), "Shutdown complete!")

	//两端/可以都不写
	hfw.RegisterRoute("/", &controllers.Index{})
	hfw.RegisterRoute("/machine", &controllers.Machine{})
	hfw.RegisterRoute("/task", &controllers.Task{})
	hfw.RegisterRoute("/user", &controllers.User{})
	hfw.RegisterRoute("/project", &controllers.Project{})
	hfw.RegisterRoute("/review", &controllers.Review{})
	hfw.RegisterRoute("/test", &controllers.Test{})
	hfw.RegisterRoute("/release", &controllers.Release{})
	hfw.RegisterRoute("/setting", &controllers.Setting{})

	hfw.RegisterStatic("/js/", hfw.Config.Template.Static)
	hfw.RegisterStatic("/css/", hfw.Config.Template.Static)
	hfw.RegisterStatic("/images/", hfw.Config.Template.Static)
	hfw.RegisterStatic("/fonts/", hfw.Config.Template.Static)

	hfw.RegisterFile("/favicon.ico", hfw.Config.Template.Static+"/images/")

	hfw.Run()
}
