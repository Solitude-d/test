#基础镜像
FROM ubuntu:20.04
#把编译后的打包进来这个镜像，放到工作目录  /app  这个是可以更换的
COPY webook /app/webook
WORKDIR /app

ENTRYPOINT ["/app/webook"]