#!/bin/bash

# Запуск тестов для всех пакетов с выводом подробной информации
go test ./... -v

# Генерация отчета о покрытии кода тестами
go test ./... -coverprofile=coverage.out

# Вывод процента покрытия кода
go tool cover -func=coverage.out

# Опционально - генерация HTML-отчета о покрытии кода
# Раскомментируйте следующую строку, чтобы создать HTML-отчет
# go tool cover -html=coverage.out -o coverage.html 