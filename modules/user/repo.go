package user

import (
	"github.com/google/uuid"
	"github.com/sammyhass/web-ide/server/modules/crypto"
	"github.com/sammyhass/web-ide/server/modules/db"
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
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ur *UserRepository) Create(userDto CreateUserDto) (User, error) {
	hashedPassword, err := crypto.HashPassword(userDto.Password)

	if err != nil {
		return User{}, err
	}

	newUser := User{
		ID:       uuid.New().String(),
		Username: userDto.Username,
		Password: hashedPassword,
	}

	res := ur.db.Create(&newUser)

	if res.Error != nil {
		return User{}, res.Error
	}

	return newUser, nil
}

func (ur *UserRepository) FindById(
	id string,
) (User, error) {
	var user User

	res := ur.db.First(&user, id)

	if res.Error != nil {
		return User{}, res.Error
	}

	return user, nil
}

func (ur *UserRepository) FindByUsername(
	username string,
) (User, error) {
	var user User

	res := ur.db.Where("username = ?", username).First(&user)

	if res.Error != nil {
		return User{}, res.Error
	}

	return user, nil
}
