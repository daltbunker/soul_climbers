package types

type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role 	  string `json:"role"`
	SoulScore int32  `json:"soulScore"`
}