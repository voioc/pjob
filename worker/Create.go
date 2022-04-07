package worker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/astaxie/beego"
	gote "github.com/linxiaozhi/go-telnet"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/utils"
	"golang.org/x/crypto/ssh"
)

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		beego.Warn(fmt.Sprintf("任务执行时间超过%d秒，进程将被强制杀掉: %d", int(timeout/time.Second), cmd.Process.Pid))
		go func() {
			<-done // 读出上面的goroutine数据，避免阻塞导致无法退出
		}()
		if err = cmd.Process.Kill(); err != nil {
			beego.Error(fmt.Sprintf("进程无法杀掉: %d, 错误信息: %s", cmd.Process.Pid, err))
		}
		return err, true
	case err = <-done:
		return err, false
	}
}

func NewCommandJob(TaskID int, serverID int, name string, command string) *Job {
	job := &Job{
		TaskID:   TaskID,
		Name:     name,
		ServerID: serverID,
	}

	// job.JobKey = libs.JobKey(id, serverID)
	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		//cmd := exec.Command("/bin/bash", "-c", command)
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("CMD", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}

		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, timeout)
		jobresult = new(JobResult)
		jobresult.OutMsg = bufOut.String()
		jobresult.ErrMsg = bufErr.String()

		jobresult.IsOk = true
		if err != nil {
			jobresult.IsOk = false
		}

		jobresult.IsTimeout = isTimeout

		return jobresult
	}
	return job
}

//远程执行任务 密钥验证
func RemoteCommandJob(TaskID, serverID int, name, command string, servers *model.TaskServer) *Job {
	job := &Job{
		TaskID:   TaskID,
		Name:     name,
		ServerID: serverID,
	}

	// job.JobKey = libs.JobKey(id, serverID)

	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		jobresult = new(JobResult)
		jobresult.OutMsg = ""
		jobresult.ErrMsg = ""
		jobresult.IsTimeout = false
		jobresult.IsOk = true

		key, err := ioutil.ReadFile(servers.PrivateKeySrc)
		if err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("读取私钥失败，%v", err.Error())
			return
		}
		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("创建签名失败，%v", err.Error())
			return
		}
		addr := fmt.Sprintf("%s:%d", servers.ServerIP, servers.Port)
		config := &ssh.ClientConfig{
			User: servers.ServerAccount,
			Auth: []ssh.AuthMethod{
				// Use the PublicKeys method for remote authentication.
				ssh.PublicKeys(signer),
			},
			//HostKeyCallback: ssh.FixedHostKey(hostKey),
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		// Connect to the remote server and perform the SSH handshake.47.93.220.5
		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败，%v", err.Error())
			return
		}

		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败，%v", err.Error())
			return
		}
		defer session.Close()

		// Once a Session is created, you can execute a single command on
		// the remote side using the Run method.

		var b bytes.Buffer
		var c bytes.Buffer
		session.Stdout = &b
		session.Stderr = &c
		jobresult.IsTimeout = false
		//session.Output(command)
		if err := session.Run(command); err != nil {
			jobresult.ErrMsg = c.String()
			jobresult.OutMsg = b.String()
			jobresult.IsOk = false
			return
		}
		jobresult.OutMsg = b.String()
		jobresult.ErrMsg = c.String()
		jobresult.IsOk = true

		return
	}
	return job
}

func RemoteCommandJobByPassword(id int, serverID int, name string, command string, servers *model.TaskServer) *Job {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	job := &Job{
		// ID:         id,
		Name:       name,
		ServerID:   serverID,
		ServerType: servers.ConnectionType,
	}

	// job.JobKey = libs.JobKey(id, serverID)
	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		jobresult = new(JobResult)
		jobresult.OutMsg = ""
		jobresult.ErrMsg = ""
		jobresult.IsTimeout = false
		jobresult.IsOk = true

		// get auth method
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(servers.Password))

		clientConfig = &ssh.ClientConfig{
			User: servers.ServerAccount,
			Auth: auth,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			//Timeout: 1000 * time.Second,
		}

		// connet to ssh
		addr = fmt.Sprintf("%s:%d", servers.ServerIP, servers.Port)

		if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("连接服务器失败，%v", err.Error())
			return
		}

		defer client.Close()

		// create session
		if session, err = client.NewSession(); err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("连接服务器失败，%v", err.Error())
			return
		}

		var b bytes.Buffer
		var c bytes.Buffer
		session.Stdout = &b
		session.Stderr = &c
		//session.Output(command)
		if err := session.Run(command); err != nil {
			jobresult.IsOk = false
		}
		jobresult.OutMsg = b.String()
		jobresult.ErrMsg = c.String()
		return
	}

	return job
}

func RemoteCommandJobByTelnetPassword(TaskID, serverID int, name, command string, servers *model.TaskServer) *Job {

	job := &Job{
		TaskID:   TaskID,
		Name:     name,
		ServerID: serverID,
	}

	// job.JobKey = libs.JobKey(id, serverID)
	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		jobresult = new(JobResult)
		jobresult.OutMsg = ""
		jobresult.ErrMsg = ""
		jobresult.IsTimeout = false
		jobresult.IsOk = true

		addr := fmt.Sprintf("%s:%d", servers.ServerIP, servers.Port)
		conn, err := gote.DialTimeout("tcp", addr, timeout)
		if err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败0，%v", err.Error())
			return
		}

		defer conn.Close()

		buf := make([]byte, 4096)

		if _, err = conn.Read(buf); err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败-1，%v", err.Error())
			return
		}

		if _, err = conn.Write([]byte(servers.ServerAccount + "\r\n")); err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败-2，%v", err.Error())
			return
		}

		if _, err = conn.Read(buf); err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败-3，%v", err.Error())
			return
		}

		if _, err = conn.Write([]byte(servers.Password + "\r\n")); err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败-4，%v", err.Error())
			return
		}

		if _, err = conn.Read(buf); err != nil {
			jobresult.IsOk = false
			jobresult.ErrMsg = fmt.Sprintf("服务器连接失败-5，%v", err.Error())
			return
		}

		loginStr := utils.GbkAsUtf8(string(buf[:]))
		if !strings.Contains(loginStr, ">") {
			jobresult.ErrMsg = jobresult.ErrMsg + "Login failed!"
			jobresult.IsOk = false
			return
		}

		commandArr := strings.Split(command, "\n")

		out, n := "", 0
		for _, c := range commandArr {
			_, err = conn.Write([]byte(c + "\r\n"))
			if err != nil {
				jobresult.ErrMsg = fmt.Sprintf("服务器连接失败-6，%v", err.Error())
				jobresult.IsOk = false
				return
			}

			n, err = conn.Read(buf)

			out = out + utils.GbkAsUtf8(string(buf[0:n]))
			if err != nil ||
				strings.Contains(out, "'"+c+"' is not recognized as an internal or external command") ||
				strings.Contains(out, "'"+c+"' 不是内部或外部命令，也不是可运行的程序") {
				jobresult.ErrMsg = jobresult.ErrMsg + " " + utils.GbkAsUtf8(string(buf[0:n]))
				jobresult.IsOk = false
				jobresult.OutMsg = out
				return
			}
		}
		jobresult.IsOk = true
		jobresult.OutMsg = out
		return
	}

	return job
}

func RemoteCommandJobByAgentPassword(TaskID int, serverID int, name string, command string, servers *model.TaskServer) *Job {

	job := &Job{
		TaskID:     TaskID,
		Name:       name,
		ServerType: servers.ConnectionType,
		ServerID:   serverID,
	}

	// job.JobKey = libs.JobKey(id, serverID)
	job.RunFunc = func(timeout time.Duration) *JobResult {
		return new(JobResult)
	}
	return job
}
