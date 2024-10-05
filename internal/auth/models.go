package auth

// models.go

type InputDetails interface {
	Sanitize()
	Validate() error
}

type RegistrationInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type LoginInput struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type UserDetails struct {
	ID       string `bson:"id" json:"id"`
	UserName string `bson:"user_name" json:"user_name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type AuthResponse struct {
	AccessToken string      `bson:"access_token" json:"access_token"`
	User        UserDetails `bson:"user" json:"user"`
}
