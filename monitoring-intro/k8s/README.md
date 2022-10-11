# Запускаем сетап в k8s

1. Запускаем minikube и открываем Kubernetes Dashboard для обзора
   ресурсов кластера.
   * запускаем minikube
     ```bash
     minikube delete && minikube start \
        --cpus=4 \
        --memory=3g \
        --addons=metrics-server \
        --ports=8001:8001 \
        --ports=30000:30000 \
        --ports=30300:30300
     ```
   * запускаем Kubernetes Dashboard
     ```bash
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
   После этого интерфейс Prometheus должен быть доступен на http://localhost:30000

1. Деплоим сервис в Prometheus и выставляем его по NodePort 30300
   ```bash
   kubectl apply -f app-deployment.yaml
   ```
   Дергаем руками сервис:
   ```bash
   curl -XPOST 'localhost:30300/generate-pair?zeroBits=20' -D-
   curl 'localhost:30300/metrics'
   ```

1. Смотрим в Prometheus UI по адресу http://localhost:30000
   на потребление CPU подами
   ```
   sum by (kubernetes_pod_name) (rate(process_cpu_seconds_total{kubernetes_pod_name=~"demo-miner-.*"}[1m]))
   ```
   а также смотрим на latency запросов
   ```
   histogram_quantile(0.5, sum by (le, kubernetes_pod_name) (rate(myapp_http_duration_seconds_bucket{kubernetes_pod_name=~"demo-miner-.*",path="/generate-pair"}[1m])))
   ```

1. Запускаем нагрузочный тест с помощью Apache Benchmark.
   ```bash
   seq 1 200 | xargs -t -n1 -J%% ab \
       -c %% \
       -t 7 \
       -m POST \
       'http://localhost:30300/generate-pair?zeroBits=12'
   ```
   
1. Видно, что как только поды уперлись в ограничение по CPU, latency идет вверх.
   Решение -- добавим HPA:
   ```bash
   kubectl apply -f hpa.yaml
   ```
   
1. Перезапускаем load test и наслаждаемся масштабированием подом.