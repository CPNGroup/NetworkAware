# 基于go基础镜像构建
FROM golang:latest as builder

# 维护者信息
LABEL maintainer="Yang Li"

# 设置容器内部的工作目录 
WORKDIR /app

# 拷贝  go mod 和 sum 文件
COPY go.mod go.sum ./

# 根据go mod下载依赖
RUN go mod download

# 从当前目录拷贝到容器工作目录
COPY . .

# build go应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## alpine镜像 #######
FROM alpine:latest  

# 运行包管理
RUN apk --no-cache add ca-certificates

# 工作目录
WORKDIR /root/

# 从前面的stage拷贝预编译的二进制文件 
COPY --from=builder /app/main .


# 执行命令
CMD ["./main"] 