# go-ssh-proxy通过ssh隧道访问梯子使用请修改配置文件中的配置信息# 执行```shell scriptgo run main.go```利用软件或浏览器 设置好socks5代理即可正常使用了监听地址: 127.0.0.1 或者 localhost监听端口: 1090 默认 其他任意未被占用端口代理软件中的地址需要与本脚本中的代理端口一致否则无法链接网络golang新手第一个程序吧# 使用方式复制config目录下的config.go.example 为 config.go更改 func loadGlobalConfig 中的配置信息打包 go build 生成可执行文件启动可执行文件 浏览器设置代理即可