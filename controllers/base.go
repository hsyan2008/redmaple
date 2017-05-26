package controllers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hsyan2008/hfw"
	"github.com/hsyan2008/redmaple/libraries"
	"github.com/hsyan2008/redmaple/models"
	gomail "gopkg.in/gomail.v2"
)

var GitTools = libraries.NewGitTools()

var pageSize int64 = 10

var groupModel = &models.Groups{}
var machineModel = &models.Machines{}
var messageModel = &models.Messages{}
var projectModel = &models.Projectes{}
var settingModel = &models.Settings{}
var taskModel = &models.Tasks{}
var taskReviewModel = &models.TaskReviews{}
var taskProjectesModel = &models.TaskProjectes{}
var userModel = &models.Users{}

var taskStatus = map[int]string{
	0: "开发中",
	1: "Review未通过，重新开发中",
	2: "测试未通过，重新开发中",
	3: "上线验证未通过，重新开发中",
	4: "部署失败，重新开发中",
	9: "已删除",

	20: "等待Review",
	21: "Review中",

	40: "等待测试",
	41: "部署中",
	42: "部署失败",
	43: "测试中",

	60: "测试通过，等待上线",
	61: "上线中",
	62: "发布失败",

	80: "上线完成，等待验证",
	81: "恭喜，上线完成",
}

type base struct {
	hfw.Controller
	User *models.Users
}

func (this *base) Before() {
	// hfw.Debug("base Before")
	this.Layout = "layout.html"
	this.Data["siteName"] = "-RedMaple开发上线管理系统"
	if this.Controll != "index" {
		if this.Session.IsExist("userinfo") == false {
			this.Redirect("/")
		}
	}
	if this.Session.IsExist("userinfo") {
		this.User = this.Session.Get("userinfo").(*models.Users)
	}
	this.Data["userinfo"] = this.User
}

//在服务器上执行，后期改成scp TODO
func (this *base) release(env string, wg *sync.WaitGroup, project *models.Projectes, machine *models.Machines) (err error) {
	defer wg.Done()
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("代码发布失败")
		}
	}()

	ssh, err := libraries.NewSsh(libraries.SshConfig{
		Username: machine.User,
		Auth:     machine.Auth,
		Ip:       machine.Ip,
		Port:     machine.Port,
	}, libraries.SshConfig{
		Username: machine.InnerUser,
		Auth:     machine.InnerAuth,
		Ip:       machine.InnerIp,
		Port:     machine.InnerPort,
	})
	if err != nil {
		return
	}
	defer ssh.Close()

	var wwwroot string
	if env == "test" {
		wwwroot = project.TestWwwroot
	} else if env == "pre_release" {
		wwwroot = project.ProdWwwroot
	} else if env == "dev" {
		wwwroot = project.DevWwwroot
	}
	if wwwroot == "" {
		wwwroot = project.Wwwroot
	}

	dir := filepath.Dir(wwwroot)
	tmpWwwroot := fmt.Sprintf("%s_%s", wwwroot, time.Now().Format("20060102150405"))

	err = GitTools.Cp(project.Wwwroot, tmpWwwroot)
	if err != nil {
		return
	}
	err = ssh.Scp(tmpWwwroot, dir, ".git")
	if err != nil {
		return
	}
	err = os.RemoveAll(tmpWwwroot)
	if err != nil {
		return
	}
	_, err = ssh.Exec(fmt.Sprintf("ln -sfT %s %s && ls -dt %s_* | awk 'NR>5{print $0}' | xargs rm -rf", tmpWwwroot, wwwroot, wwwroot))
	if err != nil {
		return
	}

	var afterRelease string
	if env == "test" {
		afterRelease = strings.Replace(project.TestAfterRelease, "\r", "", -1)
	} else if env == "pre_release" {
		afterRelease = strings.Replace(project.ProdAfterRelease, "\r", "", -1)
	} else if env == "dev" {
		afterRelease = strings.Replace(project.DevAfterRelease, "\r", "", -1)
	}
	if afterRelease != "" {
		afterReleases := strings.Split(afterRelease, "\n")
		for _, v := range afterReleases {
			if v != "" {
				_, err = ssh.Exec(v)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (this *base) saveMessage(taskId, status int, msg string) {
	message := &models.Messages{}

	message.TaskId = taskId
	message.UserId = this.User.Id
	message.Status = status
	message.IsDeleted = "N"
	message.Msg = msg

	_ = messageModel.Save(message)
}

var settingInstance, _ = settingModel.SearchOne(models.Cond{})

func (this *base) sendMail(task *models.Tasks, msg string, users ...*models.Users) {
	if len(users) == 0 || settingInstance == nil {
		return
	}
	d := gomail.NewDialer(settingInstance.SmtpAddr, settingInstance.SmtpPort, settingInstance.SmtpUser, settingInstance.SmtpPass)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = s.Close()
	}()

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(settingInstance.SmtpUser, settingInstance.SmtpName))
	to := make([]string, len(users))
	for k, user := range users {
		to[k] = m.FormatAddress(user.Email, user.Realname)
	}
	m.SetHeader("To", to...)
	msg = fmt.Sprintf("项目 %s %s", task.Name, msg)
	m.SetHeader("Subject", msg)
	m.SetBody("text/html", msg)
	err = gomail.Send(s, m)
	if err != nil {
		this.Throw(99400, "邮件发送失败"+err.Error())
	}
}
