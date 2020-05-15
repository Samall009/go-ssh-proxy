package config

import (
	"encoding/json"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
)

// 配置文件机构提
type globalConfig struct {
	UserName         string // 服务器用户名
	Password         string // 服务器密码
	LocalAddrString  string // 本地端口监听	// 需要在代理设置中填写的端口号
	ServerAddrString string // ssh地址
}

// 配置文件指针
var GlobalConfig *globalConfig

// ssh 配置指针
var Config *ssh.ClientConfig

// 包加载执行
func init() {
	// 加载全局配置文件
	loadGlobalConfig()
	// 加载ssh文件
	loadSSHConfig()
}

// 读取配置文件 加载全局参数
func loadGlobalConfig() {
	// 配置文件读取
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	GlobalConfig = &globalConfig{}
	err = decoder.Decode(GlobalConfig)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// 初始化ssh链接配置文件
func loadSSHConfig() {
	// 配置赋值
	Config = &ssh.ClientConfig{
		User: GlobalConfig.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(GlobalConfig.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
}
