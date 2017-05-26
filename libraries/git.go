package libraries

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/hsyan2008/go-logger/logger"
	"github.com/hsyan2008/hfw"
)

type GitTools struct {
	lock *sync.Mutex
}

const gitCmd = "/usr/bin/git"

func NewGitTools() *GitTools {

	return &GitTools{
		lock: &sync.Mutex{},
	}
}

func (this *GitTools) exec(cmd string, throw bool) (rs []byte, err error) {

	cmd = strings.Replace(cmd, "gitCmd", gitCmd, -1)
	logger.Debug(cmd)
	rs, err = exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	logger.Debug(string(rs), err)
	if throw {
		hfw.CheckErr(err)
	}

	return
}

func (this *GitTools) compBranchName(branch string) string {
	if branch != "master" && branch != "test" && branch != "pre_release" {
		branch = fmt.Sprintf("branch_%s", branch)
	}

	return branch
}

func (this *GitTools) Clone(git, wwwroot string) {

	cmd := fmt.Sprintf("gitCmd clone %s %s", git, wwwroot)
	_, _ = this.exec(cmd, true)
}

func (this *GitTools) NewBranch(wwwroot, branch string) {
	logger.Debug("NewBranch", wwwroot, branch)

	branch = this.compBranchName(branch)

	cmd := fmt.Sprintf("cd %s && gitCmd pull -q && gitCmd checkout -q -b %s origin/master && gitCmd push -q --set-upstream origin %s", wwwroot, branch, branch)
	_, _ = this.exec(cmd, true)
}

func (this *GitTools) DelBranch(wwwroot, branch string) {
	logger.Debug("DelBranch", wwwroot, branch)

	branch = this.compBranchName(branch)

	cmd := fmt.Sprintf("cd %s && gitCmd pull -q && gitCmd checkout -q master && gitCmd branch -q -D %s && gitCmd push -q origin --delete %s", wwwroot, branch, branch)
	_, _ = this.exec(cmd, false)
}

//重建分支
func (this *GitTools) ReBranch(wwwroot string, branches ...string) {
	for _, branch := range branches {
		this.DelBranch(wwwroot, branch)
		this.NewBranch(wwwroot, branch)
	}
}

func (this *GitTools) GetCommitId(wwwroot string, branch string) (commitid, msg string, err error) {
	logger.Debug("GetCommitId", wwwroot, branch)

	branch = this.compBranchName(branch)

	cmd := fmt.Sprintf("cd %s && gitCmd pull -q && gitCmd log origin/%s --pretty=\"%%h %%s\" -1 | cat", wwwroot, branch)
	//注意如果分支不存在，err也是nil
	rs, err := this.exec(cmd, true)
	sp := strings.SplitN(strings.TrimSpace(string(rs)), " ", 2)

	return sp[0], sp[1], nil
}

//合并src 代码到 des
//预发布分支合并到master，不需要精简commit
func (this *GitTools) Patch(wwwroot, src, des, start, end string) (err error) {
	logger.Debug("Patch", wwwroot, src, des, start, end)

	if start == end {
		return
	}

	src = this.compBranchName(src)
	des = this.compBranchName(des)

	defer func() {
		cmd := fmt.Sprintf("cd %s && rm -rf *.patch", wwwroot)
		_, _ = this.exec(cmd, false)
	}()

	cmd := fmt.Sprintf("cd %s && gitCmd pull -q && gitCmd checkout -q %s && gitCmd format-patch -q %s..%s && gitCmd checkout -q %s && gitCmd apply --check *.patch && gitCmd am -q *.patch && gitCmd push -q", wwwroot, src, start, end, des)

	_, err = this.exec(cmd, true)

	return
}

//开发分支合并到测试分支或开发分支合并到预发布分支，需要精简commit
func (this *GitTools) Merge(wwwroot, src, des, start, end, msg, author string) (err error) {
	logger.Debug("Merge", wwwroot, src, des, start, end)

	if start == end {
		return
	}

	src = this.compBranchName(src)
	des = this.compBranchName(des)

	defer func() {
		cmd := fmt.Sprintf("cd %s && rm -rf *.patch", wwwroot)
		_, _ = this.exec(cmd, false)
	}()
	defer func() {
		cmd := fmt.Sprintf("cd %s && gitCmd checkout -q %s && gitCmd branch -q -D %s_tmp", wwwroot, des, des)
		_, _ = this.exec(cmd, false)
	}()

	//format-patch到临时分支
	//merge到目标分支
	cmd := fmt.Sprintf("cd %s && gitCmd pull -q && gitCmd checkout -q %s && gitCmd format-patch -q %s..%s && gitCmd checkout -q -b %s_tmp origin/%s && gitCmd apply --check *.patch && gitCmd am -q *.patch && gitCmd checkout -q %s && gitCmd merge -q --squash %s_tmp && gitCmd commit -q --author='%s' -m '%s' && gitCmd push -q", wwwroot, src, start, end, des, des, des, des, author, msg)

	_, err = this.exec(cmd, true)

	return
}

func (this *GitTools) Cp(from, to string) (err error) {

	cmd := fmt.Sprintf("cp -r %s %s", from, to)

	_, err = this.exec(cmd, true)

	return
}

func (this *GitTools) Lock() {
	this.lock.Lock()
}

func (this *GitTools) Unlock() {
	this.lock.Unlock()
}
