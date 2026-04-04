package transformers

import (
	"nexus/backend/models"
)

type UserTransformed struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	RootAdmin bool   `json:"root_admin"`
	Coins     int    `json:"coins"`
	NameFirst string `json:"name_first"`
	NameLast  string `json:"name_last"`
	Language  string `json:"language"`
	CreatedAt string `json:"created_at"`
}

func TransformUser(user models.User) UserTransformed {
	return UserTransformed{
		ID:        user.ID,
		UUID:      user.UUID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		RootAdmin: user.RootAdmin,
		Coins:     user.Coins,
		NameFirst: user.NameFirst,
		NameLast:  user.NameLast,
		Language:  user.Language,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func TransformUsers(users []models.User) []UserTransformed {
	result := make([]UserTransformed, len(users))
	for i, u := range users {
		result[i] = TransformUser(u)
	}
	return result
}
