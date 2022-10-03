# Kubernetes Intro

Для локальной работы с Kubernetes необходимо `minikube` и `kubectl`.
Инструкции по установке можно найти [тут](https://kubernetes.io/docs/tasks/tools/).

## Запуск контейнеров в Kubernetes

Запускаем локальный Kubernetes в виртуалке

```
minikube delete # если до этого игрались с minikube
minikube start # если не работает, попробуйте --vm-driver=hyperkit для мака на x86 или --vm-driver=virtualbox для Linux 
```

- запускаем под с echo сервером

```
kubectl run echo-server --image=ealen/echo-server
kubectl get pods -o yaml # подробное описание пода, в конце есть его IP
```


- запускаем еще один под с убунтой

```
kubectl run -it ubuntu --image=ubuntu:20.04
# внутри пода
apt update && apt install -y curl
curl 127.0.0.3/hello # здесь 127.0.0.3 --- IP пода, которое мы получили на предыдущем шаге, может отличаться в вашем случае
```

## Добавление репликасета

- Добавляем репликасет

```
kubectl create -f echo-rs.yaml
```

- изменяем количество реплик пода

```
kubectl scale --replicas=20 rs echo-service
```

## Добавление сервиса

- добавляем сервис, ссылающиеся на эхо-поды

```
kubectl create -f echo-svc.yaml
```

- проверяем доступность подов по DNS имени с пода ubuntu

```
kubectl exec -it ubuntu -- bash
# внутри пода ubuntu
apt update && apt install -y curl
curl echo-service:8080/foo
```

## NodePort для доступа извне кластера

```
kubectl create -f echo-nodeport.yaml
curl `minikube ip`:30030/pupa/lupa # на Mac с --driver=docker не работает
minikube ssh curl localhost:300300 # альтернатива предыдущей команде
```

## Плавная раскатка новой версии сервиса

По очереди накатываем deployments:

- [./echo-deployment-00-initial.yaml](./echo-deployment-00-initial.yaml)
- [./echo-deployment-01-recreate.yaml](./echo-deployment-01-recreate.yaml)
- [./echo-deployment-02-rollout.yaml](./echo-deployment-02-rollout.yaml)

## Рестарт упавших инстансов

- Собираем образ ``lfyuomrgylo/probe-server`` в папке [probe-server](./probe-server) и пушим его в Docker Hub.
- Раскатываем сервис k8s
  ```
  kubectl apply -f probe-server.yaml
  ```
- В отдельном терминале смотрим на статусы подов
  ```
  kubectl get pods -w
  ```
- Отключаем readiness одного из подов
  ```
  kubectl exec POD_ID_HERE -- curl 'localhost:8080/set-probe?ready=f'
  ```
- Отключаем liveness одного из подов
  ```
  kubectl exec POD_ID_HERE -- curl 'localhost:8080/set-probe?alive=f'
  ```

## Полезные ссылки

- интерактивный туториал [Learn Kubernetes Basics](https://kubernetes.io/docs/tutorials/kubernetes-basics/)