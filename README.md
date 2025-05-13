# HTTP Load Balancer & Rate Limiter.


## 1. Пакет `loadbalancer`
[pkg/loadbalancer](https://github.com/andreyxaxa/http-load_balancer-rate_limiter/tree/main/pkg/loadbalancer)

Простой балансировщик нагрузки, который принимает входящие HTTP-запросы и распределяет их по пулу бэкенд-серверов.

- Реализован алгоритм `Robin Round`.

- Список бэкендов и порт для прослушивания получает из [config.json](https://github.com/andreyxaxa/http-load_balancer-rate_limiter/blob/main/config.json)

- Логирование входящих запросов, ошибок и смены бэкенда.

- Реализован HealthCheck.

## 2. Пакет `ratelimiter`
[pkg/ratelimiter](https://github.com/andreyxaxa/http-load_balancer-rate_limiter/tree/main/pkg/ratelimiter)

Модуль для ограничения частоты запросов (rate-limiting).

- Реализован алгоритм `Token Bucket`.

- Отслеживает состояние каждого клиента по IP.

- Поддерживает возможность настройки разных лимитов для разных клиентов.

- Настройки для разных клиентов сохраняются в базе данных.

- По умолчанию у каждого клиента 5 токенов, 0.5 токена в секунду.

## Детали 

- Конфиг - [config/config.go](https://github.com/andreyxaxa/http-load_balancer-rate_limiter/blob/main/config/config.go); Читается из `.json` файла.
- Логгер - [pkg/logger](https://github.com/andreyxaxa/http-load_balancer-rate_limiter/tree/main/pkg/logger); Интерфейс позволяет подменить логгер.
- Реализован graceful shutdown - [internal/app/app.go](https://github.com/andreyxaxa/http-load_balancer-rate_limiter/blob/main/internal/app/app.go).
- Удобная и гибкая конфигурация HTTP сервера - [pkg/httpserver/options.go](https://github.com/andreyxaxa/http-load_balancer-rate_limiter/blob/main/pkg/httpserver/options.go).
  Позволяет конфигурировать сервер в конструкторе таким образом:
```go
httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
```

## Запуск

### Docker:

1. Клонируем репозиторий.

2. `make compose-up`

Начинаем выполнять `curl localhost:8080`.

Можем нагрузить `ab -n 5000 -c 1000 http://localhost:8080/`

Поскольку `ab` выполнит все запросы почти мгновенно, получим всего 5 успешных запросов. Остальные не пройдут из-за rate-limiting'а. (По дефолту на клиента 5 токенов, 0.5 токена в секунду).

![image](https://github.com/user-attachments/assets/9eb35bb8-e804-44cb-9f66-e6e1e06fa6a0)



## Прочие `make` команды
- `make deps`:
```
go mod tidy && go mod verify
```
- `make compose-down`:
```
docker compose -f docker-compose.yml down
```
