### go-gateway-practice

- 框架gin_scaffold: https://github.com/e421083458/gin_scaffold

- learn-from: https://github.com/e421083458/go_gateway

纯属学习



### k8s部署

- 创建docker文件 vim dockerfile_dashboard
- 创建docker镜像：
```
docker build -f dockerfile_dashboard -t dockerfile_dashboard .
```
- 运行测试docker镜像: 
```
docker run -it --rm --name go_gateteway_dashboard go_gateteway_dashboard
```
- 创建交叉编译脚本，解决build太慢问题  vim docker_build.sh
- 编写服务编排文件，vim k8s_dashboard.yaml
- 启动服务
```
kubectl apply -f k8s_dashboard.yaml
kubectl apply -f k8s_server.yaml
```
- 查看所有部署
```
kubectl get all
```

## Kubernetes安装
通过Minikube快速搭建一个本地的Kubernetes单节点环境
https://m.imooc.com/article/23785