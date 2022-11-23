package auth

import (
	"errors"

	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/user"
	"golang.org/x/crypto/bcrypt"
)

const JWT_CLAIM_USER_ID = "user_id"

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
		JWT_CLAIM_USER_ID: u.ID,
	})
}

func (as *AuthService) Login(dto loginDto) (model.User, string, error) {
	found, err := as.userRepo.FindByEmail(dto.Email)

	if err != nil {
		return model.User{}, "", errors.New("username or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(found.Password),
		[]byte(dto.Password),
	); err != nil {
		return model.User{}, "", errors.New("username or password is incorrect")
	}

	token, err := as.generateJWTFomUser(found)
	if err != nil {
		return model.User{}, "", err
	}

	return found, token, nil
}

// Register creates a new user
func (as *AuthService) Register(registerDto loginDto) (model.User, string, error) {

	_, err := as.userRepo.FindByEmail(registerDto.Email)

	if err == nil {
		return model.User{}, "", errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDto.Password), bcrypt.DefaultCost)

	if err != nil {
		return model.User{}, "", err
	}

	u, err := as.userRepo.Create(user.CreateUserDto{
		Email:    registerDto.Email,
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
