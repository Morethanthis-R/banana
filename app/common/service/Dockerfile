FROM alpine:3.6

RUN sed -i 's/dl-cdn\.alpinelinux\.org/mirrors\.aliyun\.com/g' /etc/apk/repositories
RUN apk update --no-cache
ENV TZ=Asia/Shanghai
RUN apk update \
    && apk add tzdata \
    && echo "${TZ}" > /etc/timezone \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && rm /var/cache/apk/*
RUN mkdir /main
RUN mkdir /main/configs
COPY cmd/server/main /main/
COPY configs/config.yaml /main/configs/
COPY configs/registry.yaml /main/configs/
WORKDIR /main
ENTRYPOINT ["./main"]