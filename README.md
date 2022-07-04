# API GATEWAY 

* **GET /news** - возвращает страницу со списком новостей. Поддерживает фильтрацию по названию новости (параметр filter) и запрашивемый номер страницы (параметр page).
* **GET /news/{id}** - возвращает объект полной новости по ее id и комментарии к ней.
* **POST /news/{id}/comment** - добавляет комментарий к новости по ее id.

### Микросервис является единой точкой входа для:
* сервиса - парсера новостей: https://github.com/MarySmirnova/news_reader
* сервиса комментариев: https://github.com/MarySmirnova/comments_service
* сервиса модерации комментариев: https://github.com/MarySmirnova/moderation_service

### .env example:

    GATEWAY_LISTEN=:8080
	GATEWAY_READ_TIMEOUT=30s
	GATEWAY_WRITE_TIMEOUT=30s

    NEWS_ADDRESS=localhost:8081
	COMMENTS_ADDRESS=localhost:8082
	MODERATE_ADDRESS=localhost:8083


#### Структуры:

	//ответ на GET /news/{id}
	type NewsFullDetailed struct {
		ID       int       // номер записи
		Title    string    // заголовок публикации
		PubTime  int64     // время публикации
		Link     string    // ссылка на источник
		Content  string    // содержание публикации
		Comments []Comment // комментарии к публикации
	}

	//ответ на GET /news
	type NewsList struct {
		Posts []NewsShortDetailed // список новостей в сокращенном формате
		Page  Page                // объект пагинации для новостей
	}

	type NewsShortDetailed struct {
		ID      int    // номер записи
		Title   string // заголовок публикации
		PubTime int64  // время публикации
		Link    string // ссылка на источник
	}

	type Comment struct {
		ID       int    // id комментария
		ParentID int    // id родительского комментария
		NewsID   int    // id новости
		Text     string // тело комментария
		PubTime  int64  // время публикации
	}

	type Page struct {
		TotalPages   int // общее количество страниц по запросу
		NumberOfPage int // номер страницы
		ItemsPerPage int // количество элементов на одной странице
	}
