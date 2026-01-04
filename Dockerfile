# 本文件描述了构建 GoAdmin 开发环境镜像的标准方式，以及如何使用容器
#
# 使用说明:
#
# # 组装代码开发环境，从 docker-compose.yml 中获取相关数据库工具，首次运行会比较慢
# docker build -t goadmin:1.0 .
#
# # 将源代码挂载到容器中以便快速开发:
# docker run -v `pwd`:/home/goadmin --name -d goadmin:1.0
# docker exec -it goadmin /bin/bash
# # 如果本地代码已更改，可以重启容器以使更改生效
# docker restart goadmin
#  

# 使用最新的 Golang 官方镜像作为基础镜像
FROM golang:latest
# 设置维护者信息
MAINTAINER josingcjx

# 将当前目录下的所有文件复制到容器的 /home/goadmin 目录
COPY . /home/goadmin

# 设置环境变量：
# GOPATH: Go 语言工作路径，包含 /home/goadmin 目录
# GOPROXY: Go 模块代理，使用阿里云镜像和 goproxy.cn 加速依赖下载
ENV GOPATH=$GOPATH:/home/goadmin/ GOPROXY=https://mirrors.aliyun.com/goproxy,https://goproxy.cn,direct

# 运行安装命令：
# 1. 更新 apt 包索引并修复缺失的包
# 2. 安装 zip（压缩工具）、vim（文本编辑器）、postgresql（PostgreSQL 数据库）、mysql-common（MySQL 公共文件）、default-mysql-server（MySQL 服务器）
# 3. 解压 Go 依赖工具包到根目录
RUN apt-get update --fix-missing && \
    apt-get install -y zip vim postgresql mysql-common default-mysql-server && \
    tar -C / -xvf /home/goadmin/tools/godependacy.tgz 
    # 如果安装依赖工具失败，可以将本地的工具复制到远程
    #mkdir -p /go/bin  && \
    #mv /home/goadmin/tools/{gotest,goimports,golint,golangci-lint,adm} /go/bin
    #go get golang.org/x/tools/cmd/goimports && \
    #go get github.com/rakyll/gotest && \
    #go get -u golang.org/x/lint/golint && \
    #go install github.com/purpose168/GoAdmin-adm@latest && \
    #go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

# 设置工作目录为 /home/goadmin
WORKDIR /home/goadmin