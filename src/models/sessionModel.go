package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"slashbase.com/backend/src/config"
	"slashbase.com/backend/src/db"
)

type UserSession struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    string `gorm:"not null"`
	IsActive  bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignkey:user_id"`
}

func NewUserSession(userID string) (*UserSession, error) {
	var err error = nil
	if userID == "" {
		return nil, errors.New("user id cannot be empty")
	}
	return &UserSession{
		UserID:   userID,
		IsActive: true,
	}, err
}

func (session UserSession) GetAuthToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionID": session.ID,
	})
	tokenString, err := token.SignedString(config.GetAuthTokenSecret())
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (session UserSession) SetInActive() error {
	session.IsActive = false
	return session.Save()
}

func (session UserSession) Save() error {
	return db.GetDB().Save(&session).Error
}
