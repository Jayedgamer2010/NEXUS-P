package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UUID          string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	Username      string    `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email         string    `gorm:"size:191;uniqueIndex;not null" json:"email"`
	Password      string    `gorm:"size:255;not null" json:"-"`
	Role          string    `gorm:"size:20;index;default:'client'" json:"role"`
	RootAdmin     bool      `gorm:"default:false" json:"root_admin"`
	Coins         int       `gorm:"default:0;index" json:"coins"`
	NameFirst     string    `gorm:"size:255" json:"name_first"`
	NameLast      string    `gorm:"size:255" json:"name_last"`
	Language      string    `gorm:"size:10;default:'en'" json:"language"`
	UseTotp       bool      `gorm:"default:false" json:"use_totp"`
	TotpSecret    string    `gorm:"size:255" json:"-"`
	RememberToken string    `gorm:"size:255" json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Servers []Server `gorm:"foreignKey:UserID" json:"servers,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	if u.Language == "" {
		u.Language = "en"
	}
	if len(u.Password) > 0 && u.Password != "changeme" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") && len(u.Password) > 0 && u.Password != "changeme" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
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

func (u *User) Sanitize() User {
	copy := *u
	copy.Password = ""
	copy.TotpSecret = ""
	copy.RememberToken = ""
	return copy
}

func (u *User) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"uuid":       u.UUID,
		"username":   u.Username,
		"email":      u.Email,
		"role":       u.Role,
		"root_admin": u.RootAdmin,
		"coins":      u.Coins,
		"name_first": u.NameFirst,
		"name_last":  u.NameLast,
		"language":   u.Language,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
