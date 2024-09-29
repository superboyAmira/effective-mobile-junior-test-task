# effective-mobile-junior-test-task
test assignment for the junior position. 2024 Go.


# Описание задачи

# Стек

* go 1.22.3
* gin
* GORM
* Testify & linters
* Swagger/SwaggerUI
* Docker-compose

в проекте старался имитировать подход GitFlow

# QuickStart

1. Запуск контейнера с бд

```make container_up``` на порту 9900. При изменении порта в docker-compose так же поменяйте его и в .env

2. Запуск АПИ на порту 8080

```cd cmd/online-song-library```
```go run main.go```

3. http://localhost:8080/api/swagger/index.html#/

по этому адресу вы сможете нати документацию Swagger, по которой можно вручную потестировать АПИ

4. Завершение

```^Z``` в работающем процессе - graceful shutdown сервера.

Чтобы удалить контейнер с бд пропишите ```make container_rm```

# Настройки 

Все переменные окружения задаются через ```config/.env```. Список необходимых переменных представлен в репозитории.


* Путь до внешнего АПИ для обогащения данных можно задать ```PATH_EXTERNAL_API_HTTPTEST_SERVER=""```, и поменять ```EXTERNAL_API_HTTPTEST_SERVER="true"``` на false
* По дефолту используется обрезанная версия апи с захардкоженными ответами

Логгер используется slog/log. Уровень логгирования так же представлен в .env ```LOG_LEVEL="PROD"```
Код покрыт INFO и DEBUG сообщениями.

# Тестирование

Компоненты апи покрыты юнит-тестами с использованием моков. Чтобы запустить все тесты пропишите ```make test```. Запустятся тесты и линтер golangci-lint.
