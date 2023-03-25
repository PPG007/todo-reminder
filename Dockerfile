FROM phusion/baseimage:master

CMD ["/sbin/my_init"]

RUN apt-get update && \
    apt-get install -y wget && \
    wget -q https://mirrors.ustc.edu.cn/golang/go1.19.3.linux-amd64.tar.gz -O golang.tar.gz && \
    tar -C /usr/local -xf golang.tar.gz && \
    rm golang.tar.gz

ENV GOPATH /app
ENV GOROOT /usr/local/go
ENV PATH ${GOPATH}/bin:${GOROOT}/bin:$PATH
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
ENV GOMOD=/root/go.mod

WORKDIR /root

COPY . /root

RUN go mod tidy && \
    go build

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
