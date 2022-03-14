# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

ENV TZ=Asia/Shanghai
# 设置固定的项目路径
ENV WORKDIR /app
ENV BUILDDIR /app/build

###############################################################################
#                                INSTALLATION
###############################################################################
WORKDIR $BUILDDIR
COPY . .

# 编译项目
RUN go build .

RUN cp ./qqbot $WORKDIR

RUN rm -rf $BUILDDIR

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
ADD config.yaml $WORKDIR/config.yaml
CMD ./qqbot
