**Функционал API:**

Авторизация и получение JWT-токена (новый пользователь получает 1000 монет)

Просмотр баланса монет, инвентаря и истории транзакций

Покупка товаров из внутреннего магазина

Перевод монет между пользователями


**Запуск проекта:**

`git clone https://github.com/Pihtaw/merchstore.git && cd merchstore`

`docker-compose build --no-cache`

`docker-compose up -d`

**Запуск тестов:**

`go test -v ./tests`

**Проведенное нагрузочное тестирование**
```
wrk -t50 -c1000 -d30s http://localhost:8080/auth 
Running 30s test @ http://localhost:8080/auth
  50 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.41ms   12.73ms 188.82ms   85.88%
    Req/Sec     2.29k   518.47    17.27k    89.68%
  3422653 requests in 30.10s, 414.54MB read
  Socket errors: connect 29, read 0, write 0, timeout 0
  Non-2xx or 3xx responses: 3422653
Requests/sec: 113723.35
Transfer/sec:     13.77MB
```
