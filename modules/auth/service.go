package auth

import (
	"errors"

	"github.com/sammyhass/web-ide/server/modules/model"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo *userRepository
}

func newService(
	ur *userRepository,
) *Service {
	return &Service{
		userRepo: ur,
	}
}

func (as *Service) generateJWTFomUser(u model.User) (string, error) {
	return generateJWTFromUser(u.ID)
}

var (
	errIncorrectEmailOrPassword = errors.New("email or password is incorrect")
	errUserAlreadyExists        = errors.New("email already in use")
)

func (as *Service) login(dto loginDto) (model.User, string, error) {
	found, err := as.userRepo.findByEmail(dto.Email)

	if err != nil {
		return model.User{}, "", errIncorrectEmailOrPassword
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(found.Password),
		[]byte(dto.Password),
	); err != nil {
		return model.User{}, "", errIncorrectEmailOrPassword
	}

	token, err := as.generateJWTFomUser(found)
	if err != nil {
		return model.User{}, "", err
	}

	return found, token, nil
}

// Register creates a new user
func (as *Service) Register(registerDto loginDto) (model.User, string, error) {

	_, err := as.userRepo.findByEmail(registerDto.Email)

	if err == nil {
		return model.User{}, "", errUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDto.Password), bcrypt.DefaultCost)

	if err != nil {
		return model.User{}, "", err
	}

	u, err := as.userRepo.create(CreateUserDto{
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
