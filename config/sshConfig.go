package config

import (
	"golang.org/x/crypto/ssh"
	"net"
)

// 基本信息
var (
	UserName         = "UserName"       // 服务器用户名
	Password         = "Password"       // 服务器密码
	LocalAddrString  = "localhost:3090" // 本地端口监听	// 需要在代理设置中填写的端口号
	ServerAddrString = "serve_ip:port"  // ssh地址
)

// 设置config参数
var Config = &ssh.ClientConfig{
	User: UserName,
	Auth: []ssh.AuthMethod{
		ssh.Password(Password),
	},
	HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	},
}
