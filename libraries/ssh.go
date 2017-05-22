package libraries

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hsyan2008/hfw"

	"golang.org/x/crypto/ssh"
)

type SshConfig struct {
	Username string
	Ip       string
	Port     string
	Auth     string
}

type Ssh struct {
	sc *ssh.Client
	rc *ssh.Client
	c  *ssh.Client
	s  SshConfig
	r  SshConfig
}

// s 外网服务器
// r 内网服务器，可以为空
func NewSsh(s, r SshConfig) *Ssh {

	ssh := &Ssh{
		s: s,
		r: r,
	}

	ssh.Dial()

	return ssh
}

func (this *Ssh) Dial() {
	config := &ssh.ClientConfig{
		User: this.s.Username,
		Auth: []ssh.AuthMethod{
			this.getAuth(this.s.Auth),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var err error
	this.sc, err = ssh.Dial("tcp", this.s.Ip+":"+this.s.Port, config)
	hfw.CheckErr(err)

	this.c = this.sc

	//如果是跳板，继续连接
	if this.r.Ip != "" && this.r.Port != "" {
		rc, err := this.sc.Dial("tcp", this.r.Ip+":"+this.r.Port)
		hfw.CheckErr(err)

		conn, nc, req, err := ssh.NewClientConn(rc, "", &ssh.ClientConfig{
			User: this.r.Username,
			Auth: []ssh.AuthMethod{
				this.getAuth(this.r.Auth),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		hfw.CheckErr(err)

		this.rc = ssh.NewClient(conn, nc, req)
		this.c = this.rc
	}

	return
}

func (this *Ssh) getAuth(auth string) ssh.AuthMethod {
	//是文件
	var key []byte

	if _, err := os.Stat(auth); err == nil {
		key, _ = ioutil.ReadFile(auth)
	}

	//密码
	if len(key) == 0 {
		if len(auth) < 50 {
			return ssh.Password(auth)
		} else {
			key = []byte(auth)
		}
	}

	signer, _ := ssh.ParsePrivateKey(key)
	return ssh.PublicKeys(signer)
}

//一个Session只能执行一次
func (this *Ssh) Exec(cmd string) ([]byte, error) {

	sess, err := this.c.NewSession()
	hfw.CheckErr(err)
	defer func() {
		_ = sess.Close()
	}()

	return sess.CombinedOutput(cmd)
}

func (this *Ssh) Close() {
	defer func() {
		if this.sc != nil {
			_ = this.sc.Close()
		}
	}()
	if this.rc != nil {
		_ = this.rc.Close()
	}
}

//实现了目录的上传，未限速
//可以实现过滤，不支持正则过滤，不过滤最外层
func (this *Ssh) Scp(src, des, exclude string) (err error) {

	exclude = strings.Replace(exclude, "\r", "", -1)
	tmp := strings.Split(exclude, "\n")
	excludes := make(map[string]string)
	for _, v := range tmp {
		excludes[v] = v
	}

	file, err := os.Open(src)
	if err != nil {
		return
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return
	}
	if fileinfo.Mode().IsDir() {
		//如果srcDir是目录，走这个
		return this.scpDir(src, des, excludes, 0755, true)
	} else if fileinfo.Mode().IsRegular() {
		//如果srcDir是文件，则执行ssh.Run的时候，不用mkdir
		return this.scpDir(src, des, excludes, 0755, false)
	}

	return nil
}

func (this *Ssh) scpDir(src, des string, excludes map[string]string, fm os.FileMode, isDir bool) (err error) {

	file, err := os.Open(src)
	if err != nil {
		return
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return
	}
	if fileinfo.Mode().IsDir() {
		des := des + "/" + filepath.Base(src)
		for {
			files, err := file.Readdir(3)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			for _, v := range files {
				if _, ok := excludes[v.Name()]; ok {
					continue
				}
				err = this.scpDir(src+"/"+v.Name(), des, excludes, fileinfo.Mode().Perm(), isDir)
				if err != nil {
					return err
				}
			}
		}
	} else if fileinfo.Mode().IsRegular() {
		return this.scpFile(src, des, fm, isDir)
	}

	return nil
}
func (this *Ssh) scpFile(src, des string, fm os.FileMode, isDir bool) (err error) {
	sess, err := this.c.NewSession()
	if err != nil {
		return
	}
	defer func() {
		_ = sess.Close()
	}()

	go func() {
		w, err := sess.StdinPipe()
		if err != nil {
			return
		}
		defer func() {
			_ = w.Close()
		}()
		File, err := os.Open(src)
		if err != nil {
			return
		}
		info, err := File.Stat()
		if err != nil {
			return
		}
		// fmt.Fprintln(w, "C0755", info.Size(), info.Name())
		// fmt.Fprintf(w, "C%#o %d %s\n", info.Mode().Perm(), info.Size(), info.Name())
		//发布代码，文件默认是644
		fmt.Fprintf(w, "C0644 %d %s\n", info.Size(), info.Name())
		_, err = io.Copy(w, File)
		if err != nil {
			return
		}

		fmt.Fprint(w, "\x00")
	}()

	var b bytes.Buffer
	sess.Stdout = &b
	var cmd string
	if isDir {
		// cmd = fmt.Sprintf("mkdir -m %#o -p %s; /usr/bin/scp -qrt %s", fm, des, des)
		//发布代码，目录默认是755
		cmd = fmt.Sprintf("mkdir -m 0755 -p %s; /usr/bin/scp -qrt %s", des, des)
	} else {
		cmd = fmt.Sprintf("/usr/bin/scp -qrt %s", des)
	}
	if err := sess.Run(cmd); err != nil {
		// hfw.Warn(err)
		if err.Error() != "Process exited with status 1" {
			return err
		}
	}

	return nil
}
