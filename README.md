# Заглушки и моки для нагрузочного тестирования

[[_TOC_]]

## Wiremock

### Ссылки

- [Документация](https://wiremock.org/docs/)
- [GitHub](https://github.com/wiremock/wiremock)
- [Docker образ](https://hub.docker.com/r/wiremock/wiremock)

### Шаги

#### 1. Скачивание и запуск

```bash
cd wiremock
```

**jar:**

Скачать: https://repo1.maven.org/maven2/org/wiremock/wiremock-standalone/3.13.2/wiremock-standalone-3.13.2.jar

или через wget:
```bash
wget https://repo1.maven.org/maven2/org/wiremock/wiremock-standalone/3.13.2/wiremock-standalone-3.13.2.jar
```

```bash
java -jar wiremock-standalone-3.13.2.jar --help
java -jar wiremock-standalone-3.13.2.jar --port 8080
```

**docker:**
```bash
docker run -it --rm -p 8080:8080 --name wiremock wiremock/wiremock:3.13.2 --help
docker run -it --rm -p 8080:8080 --name wiremock wiremock/wiremock:3.13.2
```

Админка: http://localhost:8080/__admin/mappings

#### 2. Запуск с локальными конфигами

Папки для маппингов и статических файлов уже созданы:

```
wiremock/
├── __files/
├── mappings/
└── wiremock-standalone-3.13.2.jar
```

**jar:**
```bash
java -jar wiremock-standalone-3.13.2.jar --port 8080 --root-dir .
```

**docker:**
```bash
docker run -it --rm -p 8080:8080 --name wiremock \
  -v $PWD:/home/wiremock \
  wiremock/wiremock:3.13.2
```

#### 3. Создание файла маппинга

Создаём `mappings/stub.json`:

```json
{
  "request": {
    "method": "GET",
    "url": "/api/user"
  },
  "response": {
    "status": 200,
    "body": "{\"id\": 1, \"name\": \"test\"}",
    "headers": {
      "Content-Type": "application/json"
    }
  }
}
```

Готовый пример: [stub.json](wiremock/stub.json)

#### 4. Создание маппингов через API

```bash
curl -X POST -d @stub.json http://localhost:8080/__admin/mappings
```

Swagger UI: http://localhost:8080/__admin/swagger-ui/

#### 5. Включение HTTPS

**jar:**
```bash
java -jar wiremock-standalone-3.13.2.jar --https-port 8443 --verbose
```

**docker:**
```bash
docker run -it --rm -p 8443:8443 --name wiremock \
  wiremock/wiremock:3.13.2 --https-port 8443 --verbose
```

#### 6. Response templating

**jar:**
```bash
java -jar wiremock-standalone-3.13.2.jar --port 8080 --global-response-templating
```

**docker:**
```bash
docker run -it --rm -p 8080:8080 --name wiremock \
  -v $PWD:/home/wiremock \
  wiremock/wiremock:3.13.2 --global-response-templating
```

Примеры шаблонов:

**Path parameter:**
```json
{
  "request": {
    "urlPattern": "/api/user/[0-9]+"
  },
  "response": {
    "body": "user id: {{request.path.[2]}}",
    "transformers": ["response-template"]
  }
}
```

**Query параметры:**
```json
{
  "request": {
    "urlPattern": "/api/user.*",
    "queryParameters": {
      "id": { "matches": "[0-9]+" },
      "name": { "matches": "\\w+" }
    }
  },
  "response": {
    "body": "user id: {{request.query.id}}, user name: {{request.query.name}}",
    "transformers": ["response-template"]
  }
}
```

**Математические операции:**
```json
{
  "request": {
    "urlPattern": "/api/calc.*",
    "queryParameters": {
      "a": { "matches": "[0-9]+" },
      "b": { "matches": "[0-9]+" }
    }
  },
  "response": {
    "body": "{{math request.query.a '+' request.query.b}}",
    "transformers": ["response-template"]
  }
}
```

#### 7. Задержки ответов (важно для нагрузочного тестирования!)

**Фиксированная задержка:**
```json
{
  "request": {
    "url": "/api/slow"
  },
  "response": {
    "status": 200,
    "body": "slow response",
    "fixedDelayMilliseconds": 500
  }
}
```

Пример: [stub-with-delay.json](wiremock/stub-with-delay.json)

**Случайная задержка (lognormal распределение):**
```json
{
  "request": {
    "url": "/api/random-delay"
  },
  "response": {
    "status": 200,
    "body": "random delay",
    "delayDistribution": {
      "type": "lognormal",
      "median": 200,
      "sigma": 0.4
    }
  }
}
```

Пример: [stub-random-delay.json](wiremock/stub-random-delay.json)

**Uniform distribution:**
```json
{
  "response": {
    "delayDistribution": {
      "type": "uniform",
      "lower": 100,
      "upper": 500
    }
  }
}
```

#### 8. Extensions

Для генерации случайных данных можно использовать [wiremock-faker-extension](https://github.com/wiremock/wiremock-faker-extension):

```bash
mkdir extensions
wget -O extensions/wiremock-faker-extension-standalone-0.2.0.jar \
  https://repo1.maven.org/maven2/org/wiremock/extensions/wiremock-faker-extension-standalone/0.2.0/wiremock-faker-extension-standalone-0.2.0.jar
```

**jar:**
```bash
java -cp "wiremock-standalone-3.13.2.jar:extensions/*" wiremock.Run --port 8080 \
  --extensions org.wiremock.RandomExtension
```

**docker:**
```bash
docker run -it --rm -p 8080:8080 --name wiremock \
  -v $PWD/extensions:/var/wiremock/extensions \
  wiremock/wiremock:3.13.2 --extensions org.wiremock.RandomExtension
```

Пример использования в маппинге:

```json
{
  "request": {
    "url": "/api/random-user"
  },
  "response": {
    "body": "{\"name\": \"{{random 'Name.first_name'}}\", \"email\": \"{{random 'Internet.email_address'}}\"}",
    "transformers": ["response-template"]
  }
}
```

#### 9. Proxy & record

Записывает все запросы к реальному сервису для последующего воспроизведения:

**jar:**
```bash
java -jar wiremock-standalone-3.13.2.jar \
  --proxy-all="https://jsonplaceholder.typicode.com" \
  --record-mappings \
  --root-dir .
```

**docker:**
```bash
docker run -it --rm -p 8080:8080 --name wiremock \
  -v $PWD:/home/wiremock \
  wiremock/wiremock:3.13.2 \
  --proxy-all="https://jsonplaceholder.typicode.com" \
  --record-mappings
```

#### 10. Performance tuning

Параметры для повышения производительности при нагрузочном тестировании:

**jar:**
```bash
java -jar wiremock-standalone-3.13.2.jar --port 8080 \
  --no-request-journal \
  --disable-request-logging \
  --async-response-enabled true \
  --async-response-threads 10
```

**docker:**
```bash
docker run -it --rm -p 8080:8080 --name wiremock \
  -v $PWD:/home/wiremock \
  wiremock/wiremock:3.13.2 \
  --no-request-journal \
  --disable-request-logging \
  --async-response-enabled true \
  --async-response-threads 10
```

| Параметр | Описание |
|----------|----------|
| `--no-request-journal` | Отключает сохранение истории запросов (экономит память) |
| `--disable-request-logging` | Отключает логирование запросов в консоль |
| `--async-response-enabled` | Включает асинхронную обработку ответов |
| `--async-response-threads` | Количество потоков для асинхронной обработки |

---

## Mockintosh

Уникальная особенность — **performanceProfiles** для симуляции реалистичного поведения сервиса с ошибками и задержками.

### Ссылки

- [Документация](https://mockintosh.io/)
- [GitHub](https://github.com/up9inc/mockintosh)

### Установка

**pip (рекомендуется в venv):**
```bash
python3 -m venv venv
source venv/bin/activate
pip install mockintosh
```

**docker (рекомендуется):**
```bash
docker pull up9inc/mockintosh
```

### Шаги

#### 1. Запуск с конфигом

```bash
cd mockintosh
```

**standalone (если установлен через pip):**
```bash
mockintosh config.yaml
```

**docker:**
```bash
docker run -it --rm -p 8001:8001 \
  -v $PWD:/config \
  up9inc/mockintosh /config/config.yaml
```

#### 2. Performance Profiles

Главная фича для нагрузочного тестирования — возможность задать распределение ответов:

```yaml
performanceProfiles:
  realistic:
    ratio: 1.0          # вероятность применения (1.0 = 100%)
    delay: 0.05         # задержка в секундах
    faults:
      '200': 0.90       # 90% - успешный ответ
      '500': 0.05       # 5% - internal server error
      '503': 0.03       # 3% - service unavailable
      RST: 0.01         # 1% - сброс соединения (TCP RST)
      FIN: 0.01         # 1% - закрытие соединения (TCP FIN)

services:
  - name: Mock API
    performanceProfile: realistic
    port: 8001
    endpoints:
      - path: "/api/user"
        method: GET
        response:
          body: '{"id": 1}'
```

Полный пример: [config.yaml](mockintosh/config.yaml)

**Уровни применения профиля:**
- `globals` — для всех сервисов
- `services[].performanceProfile` — для конкретного сервиса
- `endpoints[].performanceProfile` — для конкретного endpoint

#### 3. Проверка

```bash
# несколько запросов - увидите разные статус коды
for i in {1..20}; do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8001/api/user; done
```

---

## Заглушка на Go (Golang Stub)

Простой и производительный сервер на Go (net/http) для демонстрации скорости нативного кода.

### Шаги

#### 1. Локальный запуск

```bash
cd golang-stub
go run main.go -port 8080
```

#### 2. Сборка и запуск в Docker

```bash
docker build -t golang-stub .
docker run -it --rm -p 8080:8080 golang-stub
```

#### 3. Проверка

```bash
curl http://localhost:8080/api/user
curl http://localhost:8080/api/health
```

Код: [main.go](golang-stub/main.go)

---

## Бенчмарк

Сравнение производительности заглушек с помощью `wrk`:

```bash
# установка wrk
# macOS: brew install wrk
# Linux: apt install wrk

# запуск бенчмарка
wrk -c 40 -d 30s -t 4 --latency http://localhost:8080/api/user
```

| Параметр | Описание |
|----------|----------|
| `-c 40` | 40 одновременных соединений |
| `-d 30s` | длительность теста 30 секунд |
| `-t 4` | 4 потока |
| `--latency` | показать распределение latency |

### Результаты бенчмарка

**Тестовый стенд:** MacBook Pro, Apple M1 Pro, 16 GB RAM

**Параметры теста:** `wrk -c 40 -d 10s -t 4` (3 прогона)

| Инструмент | Requests/sec | Примечание |
|------------|-------------|------------|
| Golang Stub | 78,000 - 80,000 | Go (net/http) |
| Mockserver | 78,000 - 80,000 | JVM (Netty) |
| Wiremock (inline body) | 49,000 | JVM (Jetty), с `--no-request-journal` |
| Wiremock (bodyFileName) | 18,000 | чтение тела из файла |
| Wiremock (static __files) | 16,000 | раздача файла напрямую |
| Mockintosh | ~50* | Docker эмуляция amd64 на arm64 |

> **Вывод:** Для нагрузочного тестирования с Wiremock используйте inline `body`, а не `bodyFileName`.

> На локальном тестировании (loopback) Golang и Mockserver показывают схожие результаты ~80k RPS — это лимит тестового стенда.

> *Mockintosh Docker-образ доступен только для `linux/amd64`.
> На Apple Silicon (M1/M2/M3) работает через эмуляцию, что снижает производительность.
> На x86_64 Linux ожидается ~1,000-5,000 RPS.

---

## Mockserver (deprecated, но быстрый)

> **Проект не развивается с 2024 года**: https://github.com/mock-server/mockserver/issues/1935
>
> Часто встречается на legacy-проектах. Несмотря на отсутствие развития, работает стабильно и показывает высокую производительность.

### Ссылки

- [Документация](https://www.mock-server.com/mock_server/getting_started.html)
- [GitHub](https://github.com/mock-server/mockserver)
- [Docker образ](https://hub.docker.com/r/mockserver/mockserver)

### Шаги

#### 1. Скачивание и запуск

```bash
cd mockserver
```

**jar:**
```bash
wget -O mockserver-netty.jar https://repo1.maven.org/maven2/org/mock-server/mockserver-netty-no-dependencies/5.15.0/mockserver-netty-no-dependencies-5.15.0.jar
java -jar mockserver-netty.jar -serverPort 1080
```

**docker:**
```bash
docker run -it --rm -p 1080:1080 \
  --env MOCKSERVER_LOG_LEVEL=ERROR \
  mockserver/mockserver:5.15.0
```

Админка: http://localhost:1080/mockserver/dashboard

#### 2. Создание expectation

```bash
curl -X PUT "http://localhost:1080/mockserver/expectation" -d '{
  "httpRequest": {
    "method": "GET",
    "path": "/api/user"
  },
  "httpResponse": {
    "statusCode": 200,
    "headers": {"Content-Type": ["application/json"]},
    "body": "{\"id\": 1, \"name\": \"test\"}"
  },
  "times": {"unlimited": true}
}'
```

#### 3. Проверка

```bash
curl http://localhost:1080/api/user
```

