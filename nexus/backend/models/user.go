package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"uniqueIndex;size:36" json:"uuid"`
	Username  string    `gorm:"uniqueIndex;size:255" json:"username"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Password  string    `gorm:"size:255" json:"-"`
	NameFirst string    `gorm:"size:100" json:"name_first"`
	NameLast  string    `gorm:"size:100" json:"name_last"`
	Role      string    `gorm:"size:20;default:client" json:"role"`
	RootAdmin bool      `gorm:"default:false" json:"root_admin"`
	Coins     int       `gorm:"default:0" json:"coins"`
	Suspended bool      `gorm:"default:false" json:"suspended"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Servers   []Server  `gorm:"foreignKey:UserID" json:"servers,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	return nil
}

func (u *User) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin" || u.RootAdmin
}
