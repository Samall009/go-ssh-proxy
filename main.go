package main

// 从本地端口9000转发到远程端口9999
import (
	config2 "github/lucky/ssh_proxy/config"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

// ssh代理全局链接指针
var GableClient *ssh.Client

// 入口方法
func main() {
	// 设置监听 监听本地端口数据
	localListener, err := net.Listen("tcp", config2.LocalAddrString)
	if err != nil {
		log.Fatalf("net.Listen failed: %v", err)
	}

	// 兜底关闭
	defer localListener.Close()

	log.Println("开始连接代理服务器")

	// 链接ssh服务器
	sshClient()

	log.Println("开始接受链接")

	// 兜底执行 关闭ssh链接
	defer GableClient.Close()

	// 执行任务
	task(localListener)
}

// 链接ssh服务器
func sshClient() {
	// 错误信息类型定义
	var err error

	// 链接ssh
	GableClient, err = ssh.Dial("tcp", config2.ServerAddrString, config2.Config)

	// 判断错误
	if err != nil {
		log.Fatalf("代理服务器链接错误信息: %v", err)
		os.Exit(-1)
	}
}

// 任务监听
func task(localListener net.Listener) {
	for {
		// 接收本地数据
		client, err := localListener.Accept()
		if err != nil {
			log.Fatalf("本地数据hook错误信息: %v", err)
			return
		}

		// 为空关闭
		if client == nil {
			return
		}

		// 携程处理器执行代理请求
		go proxy(client)
	}
}

// 获取代理链接
func getHostPort(client net.Conn) string {
	var b [1024]byte
	// 从链接中读取数据
	n, err := client.Read(b[:])
	if err != nil {
		log.Fatalf("代理数据读取错误信息: %v", err)
		return ""
	}

	//只处理Socks5协议
	if b[0] == 0x05 {
		// 客户端回应：Socks服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})

		// 从链接中读取数据
		n, err = client.Read(b[:])

		// 初始化变量
		var host, port string

		// 获取host
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}

		// 获取端口
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		// 返回链接字符串
		return net.JoinHostPort(host, port)
	}

	return ""
}

// 执行代理
func proxy(client net.Conn) {
	// 兜底执行 关闭链接
	defer client.Close()

	// 获取链接地址
	hostPort := getHostPort(client)

	// 判断存在地址链接
	if hostPort != "" {
		// 转发地址
		ssh, err := GableClient.Dial("tcp", hostPort)
		if err != nil {
			log.Fatalf("代理数据请求错误信息: %v", err)
			return
		}
		// 兜底关闭
		defer ssh.Close()

		// 响应客户端连接成功
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		//进行转发
		go io.Copy(ssh, client)
		io.Copy(client, ssh)
	}
}
