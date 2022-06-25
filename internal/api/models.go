package api

type NewsFullDetailed struct {
	ID       int       // номер записи
	Title    string    // заголовок публикации
	PubTime  int64     // время публикации
	Link     string    // ссылка на источник
	Content  string    // содержание публикации
	Comments []Comment // комментарии к публикации
}

type NewsShortDetailed struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
	//  	ShortContent string // первый абзац публикации
}

type Comment struct {
	ID       int    // id комментария
	ParentID int    // id родительского комментария
	NewsID   int    // id новости
	Text     string // тело комментария
	PubTime  int64  // время публикации
}
