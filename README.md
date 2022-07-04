# API GATEWAY 

* **GET /news** - возвращает страницу со списком новостей. Поддерживает фильтрацию по названию новости (параметр filter) и запрашивемый номер страницы (параметр page).
* **GET /news/{id}** - возвращает объект полной новости по ее id и комментарии к ней.
* **POST /news/{id}/comment** - добавляет комментарий к новости по ее id.

## Микросервис является единой точкой входа для:
* сервиса - парсера новостей: https://github.com/MarySmirnova/news_reader
* сервиса комментариев: https://github.com/MarySmirnova/comments_service
* сервиса модерации комментариев: https://github.com/MarySmirnova/moderation_service

## .env example:

    GATEWAY_LISTEN=:8080
	GATEWAY_READ_TIMEOUT=30s
	GATEWAY_WRITE_TIMEOUT=30s

    NEWS_ADDRESS=localhost:8081
	COMMENTS_ADDRESS=localhost:8082
	MODERATE_ADDRESS=localhost:8083