package types

type Blog struct {
	Body string `json:"body"`
	Title string `json:"title"`
	Excerpt string `json:"excerpt"`
	IsPublished bool `json:"isPublished"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	Id int32 `json:"id"`
}