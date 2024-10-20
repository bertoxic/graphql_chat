package models

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	Floatmap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

type InputDetails interface {
	Sanitize()
	Validate() error
}

//	type UserDetails struct {
//		ID       string `bson:"id" json:"id"`
//		UserName string `bson:"user_name" json:"user_name"`
//		Email    string `bson:"email" json:"email"`
//		Password string `bson:"password" json:"password"`
//	}
type RegistrationInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type LoginInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type UserPostStats struct {
	TotalPosts   int `json:"total_posts"`
	TotalLikes   int `json:"total_likes"`
	TotalReposts int `json:"total_reposts"`
}
