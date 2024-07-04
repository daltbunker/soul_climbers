package types

type Blog struct {
	Body string `json:"body"`
	Title string `json:"title"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	Id int32 `json:"id"`
}