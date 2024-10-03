package packages

import "context"

type UserRepo interface {
	GetUserByName(ctx context.Context, username string) (User, error)
}

type User struct {
	ID        string `bson:"i_d" json:"i_d"`
	Username  string `bson:"user_name" json:"user_name"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password"`
	CreatedAt string `bson:"created_at" json:"created_at"`
	UpdatedAt string `bson:"updated_at" json:"updated_at"`
}

func (us *User) GetUserByName(ctx context.Context, username string) (*User, error) {

	return nil, nil
}
