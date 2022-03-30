# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

ENV TZ=Asia/Shanghai
# 设置固定的项目路径
ENV BUILDDIR /app/build

###############################################################################
#                                Building
###############################################################################
WORKDIR $BUILDDIR
COPY . .

# 编译项目
RUN go build -o app .

###############################################################################
#                                   START
###############################################################################
FROM alpine:3.15.3
ENV WORKDIR /app
ENV TZ=Asia/Shanghai

WORKDIR $WORKDIR

COPY --from=0 /app/build/app ./

ADD config.yaml $WORKDIR/config.yaml

CMD ["./app"]
