package user

import (
	"github.com/google/uuid"
	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewRepository() *UserRepository {
	return &UserRepository{
		db: db.GetConnection(),
	}
}

type CreateUserDto struct {
	Email    string `json:"username"`
	Password string `json:"password"` // accepts the hashed password
}

func (ur *UserRepository) Create(userDto CreateUserDto) (model.User, error) {

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

func (ur *UserRepository) FindById(
	id string,
) (model.User, error) {
	var user model.User

	res := ur.db.Where("id = ?", id).First(&user)

	if res.Error != nil {
		return model.User{}, res.Error
	}

	return user, nil
}

func (ur *UserRepository) FindByEmail(
	email string,
) (model.User, error) {
	var user model.User

	res := ur.db.Where("email = ?", email).First(&user)

	if res.Error != nil {
		return model.User{}, res.Error
	}

	return user, nil
}
