FROM golang:1.16-alpine as builder

WORKDIR /www/job
# 准备工作
#RUN export 
COPY go.mod ./
COPY go.sum ./
COPY . .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk add --no-cache ca-certificates tzdata

# 编译
RUN CGO_ENABLED=0 GOPROXY="https://goproxy.cn" GO111MODULE=on go build -ldflags "-s -w" -o job
RUN chmod +x /www/job/job

FROM alpine as runner
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /www/job/job /www/job/
COPY --from=builder /www/job/config/config_dev.toml /www/job/config/config.toml
COPY --from=builder /www/job/static /www/job/static/
COPY --from=builder /www/job/views /www/job/views/
WORKDIR /www/job

# 执行编译生成的二进制文件
CMD ["/www/job/job", "-c", "/www/job/config/config.toml"]

# 暴露端口
EXPOSE 8001