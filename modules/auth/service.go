package auth

import (
	"errors"

	"github.com/sammyhass/web-ide/server/modules/crypto"
	"github.com/sammyhass/web-ide/server/modules/user"
)

type AuthService struct {
	userRepo *user.UserRepository
}

func NewService(
	ur *user.UserRepository,
) *AuthService {
	return &AuthService{
		userRepo: ur,
	}
}

func (as *AuthService) Login(username string, password string) (user.User, error) {
	found, err := as.userRepo.FindByUsername(username)

	if err != nil {
		return user.User{}, err
	}

	if ok := crypto.Compare(found.Password, password); !ok {
		return user.User{}, errors.New("login failed, check your credentials")
	}

	return found, nil
}

// Register creates a new user with the given username and password and returns the created user
func (as *AuthService) Register(username string, password string) (user.User, error) {

	_, err := as.userRepo.FindByUsername(username)

	if err == nil {
		return user.User{}, errors.New("username already exists")
	}

	hashedPassword, err := crypto.HashPassword(password)

	if err != nil {
		return user.User{}, err
	}

	newUser := user.User{
		ID:       user.NewID(),
		Username: username,
		Password: hashedPassword,
	}

	u, err := as.userRepo.Create(user.CreateUserDto{
		Username: newUser.Username,
		Password: newUser.Password,
	})

	if err != nil {
		return user.User{}, err
	}

	return u, nil

}
