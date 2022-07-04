package api

type NewsFullDetailed struct {
	ID       int       // номер записи
	Title    string    // заголовок публикации
	PubTime  int64     // время публикации
	Link     string    // ссылка на источник
	Content  string    // содержание публикации
	Comments []Comment // комментарии к публикации
}

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
