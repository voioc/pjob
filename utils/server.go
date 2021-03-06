package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"time"

	gote "github.com/linxiaozhi/go-telnet"
	"github.com/pkg/errors"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
	"golang.org/x/crypto/ssh"
)

func RemoteCommandByTelnetPassword(servers *model.TaskServer) error {

	addr := fmt.Sprintf("%s:%d", servers.ServerIP, servers.Port)
	conn, err := gote.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		return err
	}

	defer conn.Close()

	buf := make([]byte, 4096)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(servers.ServerAccount + "\r\n"))
	if err != nil {
		return err
	}

	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(servers.Password + "\r\n"))
	if err != nil {
		return err
	}

	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	str := GbkAsUtf8(string(buf[:]))

	if strings.Contains(str, ">") {
		return nil
	}

	return errors.Errorf("连接失败!")
}

func RemoteCommandByPassword(servers *model.TaskServer) error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(servers.Password))

	clientConfig = &ssh.ClientConfig{
		User: servers.ServerAccount,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	addr = fmt.Sprintf("%s:%d", servers.ServerIP, servers.Port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err == nil {
		defer client.Close()
	}
	return err
}

func RemoteCommandByKey(servers *model.TaskServer) error {
	key, err := ioutil.ReadFile(servers.PrivateKeySrc)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
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
		Timeout: 5 * time.Second,
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err == nil {
		client.Close()
	}
	return err
}

func RemoteAgent(servers *model.TaskServer) error {

	conn, err := net.Dial("tcp", servers.ServerIP+":"+strconv.Itoa(servers.Port))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	var reply *common.RpcResult
	defer client.Close()

	ping := "ping"
	err = client.Call("RpcTask.HeartBeat", ping, &reply)
	if err != nil {
		return err
	}
	if reply.Status == 200 {
		return nil
	} else {
		return fmt.Errorf("链接错误：%v", reply.Message)
	}
}
