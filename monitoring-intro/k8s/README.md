# Запускаем сетап в k8s

1. Запускаем minikube и открываем Kubernetes Dashboard для обзора
   ресурсов кластера.
   ```bash
   minikube start # --ports=8001:8001 --ports=30000:30000 --ports=30300:30300 for docker driver
   minikube dashboard --port=8001
   ```

1. Собираем docker-образ ``lfyuomrgylo/demo-miner`` в папке [app](../app)
   и пушим его в docker registry

1. Создаем неймспейс `monitoring` в k8s, чтобы запустить в нем Prometheus:
   ```bash
   kubectl create namespace monitoring
   ```
   
1. Создаем ``ClusterRole`` с необходимыми для мониторинга правами и делаем
   её дефолтом для неймспейса ``monitoring``:
   ```bash
   kubectl apply -f role.yaml
   ```
   
1. Создаем в неймспейсе `monitoring` ресурс ``ConfigMap``, 
   в который сохарним конфигурацию для нашего Prometheus:
   ```bash
   kubectl apply -f confmap.yaml
   ```
   
1. Деплоим инстанс Prometheus и делаем его доступным по NodePort 30000
   ```bash
   kubectl apply -f prom-deployment.yaml
   ```
   После этого интерфейс Prometheus должен быть доступен на ``http://$(minikube ip):30000``
   либо на ``http://localhost:30000`` (если используется драйвер `docker` для `minikube`).

1. Деплоим сервис в Prometheus и выставляем его по NodePort 30300
   ```bash
   kubectl apply -f app-deployment.yaml
   ```
   Дергаем руками сервис:
   ```bash
   curl -XPOST 'localhost:30300/generate-pair?zeroBits=20' -D-
   curl 'localhost:30300/metrics'
   ```
   Смотрим на метрики в Prometheus UI на ``http://localhost:30000``