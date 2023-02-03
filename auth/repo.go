package auth

import (
	"github.com/google/uuid"
	"github.com/sammyhass/web-ide/server/db"
	"github.com/sammyhass/web-ide/server/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func newRepository() *userRepository {
	return &userRepository{
		db: db.GetConnection(),
	}
}

type CreateUserDto struct {
	Email    string `json:"username"`
	Password string `json:"password"` // accepts the hashed password
}

func (ur *userRepository) create(userDto CreateUserDto) (model.User, error) {

	newUser := model.User{
		ID:       uuid.New().String(),
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	res := ur.db.Create(&newUser)

	if res.Error != nil {
		return model.User{}, res.Error
	}

	return newUser, nil
}

func (ur *userRepository) findById(
	id string,
) (model.User, error) {
	var user model.User

	res := ur.db.Where("id = ?", id).First(&user)

	if res.Error != nil {
		return model.User{}, res.Error
	}

	return user, nil
}

func (ur *userRepository) findByEmail(
	email string,
) (model.User, error) {
	var user model.User

	res := ur.db.Where("email = ?", email).First(&user)

	if res.Error != nil {
		return model.User{}, res.Error
	}

	return user, nil
}
