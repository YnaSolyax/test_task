# Инициализация и загрузка зависимостей
go mod tidy

# Прямой запуск
go run main.go

# Или сборка бинарного файла и его запуск
go build -o server main.go
./server

# Запуск тестов с подробным выводом
go test -v ./internal/server

#POST
curl -X POST http://localhost:8080/todos -d '{"title":"Тест"}'

#GET all
curl http://localhost:8080/todos

#GET id
curl http://localhost:8080/todos/1

#PUT id
curl -X PUT http://localhost:8080/todos/1 -d '{"title":"Новое","status":"finished"}'

#Delete Id
curl -X DELETE http://localhost:8080/todos/1
