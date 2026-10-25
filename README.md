Проект создан для удобства отслеживания работоспособности запрашиваемых сервисов с заданным интервалом времени.

"repository.go является тестовым файлом для сохранение записей в память"

"По мере разработки информация по запуску и сборке образа будет дополнена"

Команды для поднятия контейнера: 
`docker run --name gowatch-postgres \
  -e POSTGRES_USER=gowatch_user \
  -e POSTGRES_PASSWORD=gowatch_pass \
  -e POSTGRES_DB=gowatch_db \
  -p 5432:5432 \
  -d postgres:16-alpine`


  Команды для запуска миграций:
  `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` //установка утилиты golang-migrate

  `migrate -path migrations -database "postgres://gowatch_user:gowatch_pass@localhost:5432/gowatch_db?sslmode=disable" up` // команда для "поднятия" миграций

  команды для goose:

  `go install github.com/pressly/goose/v3/cmd/goose@latest` установка goose

  `goose -dir migrations postgres "postgres://gowatch_user:gowatch_pass@localhost:5432/gowatch_db?sslmode=disable" up` поднятие миграций

  `goose -dir migrations postgres "postgres://gowatch_user:gowatch_pass@localhost:5432/gowatch_db?sslmode=disable" down` откат миграций