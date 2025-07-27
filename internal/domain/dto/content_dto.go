package dto

type ContentResponse struct {
	Id        int64  `json:"id"`
	ChapterId int64  `json:"chapter_id"`
	Filename  string `json:"filename"`
	Url       string `json:"url"`
}
