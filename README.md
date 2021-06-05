<p align="center">
    <a href="https://github.com/jkuup/go-xops" target="_blank">
        <img src="https://raw.githubusercontent.com/jkuup/go-xops/master/img/gopher.png?v=0.2.2" width="180" />
    </a>
    <h3 align="center">go-xops</h3>
    <p align="center">Golang 自动化运维平台</p>
    <p align="center">
        <a href="https://travis-ci.com/jkuup/go-xops"><img src="https://travis-ci.com/jkuup/go-xops.svg?branch=master"></a>
        <a href="https://github.com/jkuup/go-xops/releases"><img src="https://img.shields.io/badge/Version-v1.0.0-red.svg"></a>
        <a href="https://goreportcard.com/report/github.com/jkuup/go-xops"><img src="https://goreportcard.com/badge/github.com/jkuup/go-xops?v=1.0.0"></a>
        <a href="https://hub.docker.com/r/jkuup/go-xops"><img src="https://img.shields.io/badge/Docker-Latest-orange"></a>
        <a href="https://github.com/jkuup/go-xops/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-Apache License-orange.svg"></a>
    </p>
</p>
<br/>
功能开发展示和进度(更新慢请包涵)

- 已完成
- [X] 基于Casbin的用户权限管理
- [X] Linux主机管理
- [X] Linux服务器的命令执行，文件下发，web界面终端
- [X] 主机监控,数据库监控基于prometheus和grafana
- [X] 基于docker和k8s的镜像构建和打包
- 待开发
- [ ] CI/CD持续构建和集成

安装部署

- 先把xops.sql给导入导xops数据库中
- 先克隆git clone https://github.com/jkuup/go-xops.git
- cd go-xops 在conf目录下修改config-dev.yml文件对内容
![avatar](https://github.com/jkuup/go-xops/master/config-dev-1.png)
![avatar](https://github.com/jkuup/go-xops/master/config-dev-2.png)
注意上面需要修改自己对数据库地址，关于prometheus和k8s可以不用修改项目也可以跑起来
- 执行命令 go run main.go

