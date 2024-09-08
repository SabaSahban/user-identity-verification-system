package model

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type User struct {
	Email      string
	LastName   string
	NationalID string
	IP         string
	UUID       string
	State      string
}

type UserRepo interface {
	Save(user User) error
	FindByNationalID(nationalID string) (*User, error)
	UpdateStateByUserID(userID string, newState string) error
	FindUserByUserID(userID string) (*User, error)
}

func NewSQLUserRepo(db *gorm.DB) SQLUserRepo {
	return SQLUserRepo{DB: db}
}

type SQLUserRepo struct {
	DB *gorm.DB
}

func (s SQLUserRepo) Save(user User) (finalErr error) {
	err := s.DB.Save(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (s SQLUserRepo) FindByNationalID(nationalID string) (*User, error) {
	user := &User{}
	err := s.DB.Where("national_id = ?", nationalID).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			logrus.Errorf("national id doesn't exist")

			return nil, err
		}
		return nil, err
	}
	return user, nil
}

func (s SQLUserRepo) UpdateStateByUserID(userID string, newState string) error {
	err := s.DB.Model(&User{}).
		Where("uuid = ?", userID).
		Update("state", newState).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s SQLUserRepo) FindUserByUserID(userID string) (*User, error) {
	user := &User{}
	err := s.DB.Where("uuid = ?", userID).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
