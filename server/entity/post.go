package entity

type Post struct {
	Id      uint32 `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body"`
	PicUrl  string `json:"picUrl"`
	//Bookmark bool   `json:"bookmark"` 每個人 bookmark 不同，所以應由別的db 關連它們
}
