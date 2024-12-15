package types

type BlogImg struct {
	ImgName string `json:"imageName"`
	Img []byte `json:"img"`
	BlogId int32 `json:"blogId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}