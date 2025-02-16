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
