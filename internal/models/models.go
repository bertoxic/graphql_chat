package models

type InputDetails interface {
	Sanitize()
	Validate() error
}
type UserDetails struct {
	ID       string `bson:"id" json:"id"`
	UserName string `bson:"user_name" json:"user_name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}
type RegistrationInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

func (in *RegistrationInput) Sanitize() {

}
func (in *RegistrationInput) Validate() error {
	return nil
}
