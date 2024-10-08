package test

import (
	"github.com/bertoxic/graphqlChat/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	log.Print("hello starting main tests now")
	config.SetPasswordCost(bcrypt.MinCost)
}

func tearDown() {
	log.Print("tearing down after tests")
}
