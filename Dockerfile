FROM golang:latest 
# 指定构建过程中的工作目录
WORKDIR /opt/application

COPY . .

# 执行代码编译命令。操作系统参数为linux，编译后的二进制产物命名为main，并存放在当前目录下。
RUN GOOS=linux go build -o main .

RUN apt update && \
    apt upgrade && \
    apt-get install bash && \
    apt-get install curl 
USER root