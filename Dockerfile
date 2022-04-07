FROM public-env-mirror-service-registry.cn-beijing.cr.aliyuncs.com/dist/golang:1.16

WORKDIR /tmp/job
# 准备工作
#RUN export 
COPY go.mod .
COPY go.sum .
COPY . .

#加入git访问权限
# COPY .netrc /root/.netrc

# 编译
# RUN go env -w GOPRIVATE=codeup.aliyun.com
RUN go env
RUN GOPROXY="https://goproxy.cn" GO111MODULE=on go build -o /www/job/job .
#RUN go build -o ./out/go-sample-app .
RUN chmod +x /www/job/job

# ARG envType=test
# COPY config/config_${envType}.toml config/env.toml
COPY config /www/job/config/
COPY static /www/job/static/
COPY views /www/job/views/

# 执行编译生成的二进制文件
CMD ["./www/job/job"]
# 暴露端口
EXPOSE 8001