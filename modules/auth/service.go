package auth

import (
	"errors"

	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/user"
	"golang.org/x/crypto/bcrypt"
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

func (as *AuthService) generateJWTFomUser(u model.User) (string, error) {
	return generateJWTFromClaims(map[string]interface{}{
		"user_id": u.ID,
	})
}

func (as *AuthService) Login(username string, password string) (model.User, string, error) {
	found, err := as.userRepo.FindByUsername(username)

	if err != nil {
		return model.User{}, "", errors.New("username or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(found.Password),
		[]byte(password),
	); err != nil {
		return model.User{}, "", errors.New("username or password is incorrect")
	}

	token, err := as.generateJWTFomUser(found)
	if err != nil {
		return model.User{}, "", err
	}

	return found, token, nil
}

// Register creates a new user with the given username and password and returns the created user
func (as *AuthService) Register(username string, password string) (model.User, string, error) {

	_, err := as.userRepo.FindByUsername(username)

	if err == nil {
		return model.User{}, "", errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return model.User{}, "", err
	}

	u, err := as.userRepo.Create(user.CreateUserDto{
		Username: username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return model.User{}, "", err
	}

	jwt, err := as.generateJWTFomUser(u)

	if err != nil {
		return model.User{}, "", err
	}

	return u, jwt, nil

}
