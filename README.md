# Сервис управления баннерами Avito
- Сервис реализован на языке **Go**
- Логгирование на основе **log/slog**
- Для работы с web используется **chi** + **http/net**
- Для хранения in-memory K-V значений был выбран **go-cache**
- База данных: **PostgreSQL** (driver **pgx/v5**)
- Миграции реализованы с помощью **golang-migrate/v4**
- Для выполнения фоновых задач используется **gocron**
- Реализованы интеграционные тесты с помощью **testify/suite**
- Сервис и тесты поднимаются с помощью компоновки **Docker**-контейнеров (сервис, БД, БД для тестов)
- Набор инструкций по взаимодействию с проектом описан в **Makefile**
- Проведено load-тестирование для различных сценариев получения баннеров и обновления данных с помощью **k6**. Результаты тестирования находятся в папке *k6-loadtest/results*

## Доп. информация по взаимодействию с сервисом:
- Для запуска сервиса используется команда __make up__, для завершения сервиса - __make down__ для запуска тестов - __make test__
- Токены для аунтефикации хранятся в таблице __users__. В нее помещены 2 токена: пользовательский ('user_token') и админский ('admin_token')
- Взаимодействовать с сервисом можно, например, с помощью __curl__-запросов, или используя сервисы тестирования API (Postman,Insomnia и т.д.)

## Примеры запросов

### Получение баннера для юзера
```bash
curl -X GET "http://0.0.0.0:8080/user_banner?tag_id=123&feature_id=456&use_last_revision=true" -H "token: user_token"
```

### Получение баннеров по фильтру:
```bash
curl -X GET "http://0.0.0.0:8080/banner?tag_id=123&feature_id=456&limit=10&offset=0" -H "token: admin_token"
```

### Сохранение баннера:
```bash
curl -X POST "http://0.0.0.0:8080/banner" -H "token: admin_token" -H "Content-Type: application/json" -d '{
  "tag_ids": [123, 456],
  "feature_id": 789,
  "content": {"title": "test title", "text": "test text", "url": "test url"},
  "is_active": true
}'
```

### Обновление баннера по айди:
```bash
curl -X PATCH "http://0.0.0.0:8080/banner/1" -H "token: admin_token" -H "Content-Type: application/json" -d '{
  "tag_ids": [2, 3],
  "feature_id": 4,
  "content": {"title": "updated title", "text": "updated text", "url": "updated url"},
  "is_active": true
}'
```

### Удаление баннера по айди:
```bash
curl -X DELETE "http://0.0.0.0:8080/banner/1" -H "token: admin_token"
```

### Удаление баннеров по фильтру:
```bash
curl -X DELETE "http://0.0.0.0:8080/banner?feature_id=5&tag_id=8" -H "token: admin_token"
```