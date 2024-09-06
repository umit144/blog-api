package types

type Post struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Slug    string `json:"-"`
	Content string `json:"content"`
	Author  User   `json:"author"`
}
