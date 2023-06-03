package dao

import (
	"errors"

	"LogAnalyse/app/service/user/rpc/model"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("no such user")
	ErrUserExist    = errors.New("user already exist")
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	m := db.Migrator()
	if !m.HasTable(&model.User{}) {
		err := m.CreateTable(&model.User{})
		if err != nil {
			panic(err)
		}
	}
	return &User{db: db}
}

func (u *User) CreateUser(user *model.User) error {
	err := u.db.Model(&model.User{}).
		Where(&model.User{Username: user.Username}).First(&model.User{}).Error
	if err == nil {
		return ErrUserExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return u.db.Model(&model.User{}).Create(user).Error
}

func (u *User) GetUserByUsername(username string) (user *model.User, err error) {
	user = new(model.User)
	err = u.db.Model(user).Where(&model.User{Username: username}).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return
}
