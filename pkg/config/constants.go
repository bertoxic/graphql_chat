package config

import "golang.org/x/crypto/bcrypt"

var PasswordCost = bcrypt.DefaultCost

func SetPasswordCost(cost int) {
	PasswordCost = cost
}
