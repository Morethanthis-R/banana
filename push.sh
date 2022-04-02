#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go wire_gen.go
#docker build -t registry.cn-shenzhen.aliyuncs.com/satsun/china_labor_union_backend:backend .
#docker push registry.cn-shenzhen.aliyuncs.com/satsun/china_labor_union_backend:backend

docker build -t registry.cn-shenzhen.aliyuncs.com/morethanthis/peach:common-v1 .
docker push registry.cn-shenzhen.aliyuncs.com/morethanthis/peach:common-v1

docker build -t registry.cn-shenzhen.aliyuncs.com/morethanthis/peach:transfer-v1 .
docker push registry.cn-shenzhen.aliyuncs.com/morethanthis/peach:transfer-v1

docker build -t registry.cn-shenzhen.aliyuncs.com/morethanthis/peach:account_center-v1 .
docker push registry.cn-shenzhen.aliyuncs.com/morethanthis/peach:account_center-v1

docker build -t registry.cn-shenzhen.aliyuncs.com/morethanthis/banana:account_center-v1 .
docker push registry.cn-shenzhen.aliyuncs.com/morethanthis/banana:accoundocker push registry.cn-shenzhen.aliyuncs.com/morethanthis/peach:transfer-v1t_center-v1