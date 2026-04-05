package transformers

import "nexus/backend/models"

type UserItem struct {
	ID        uint              `json:"id"`
	UUID      string            `json:"uuid"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Role      string            `json:"role"`
	RootAdmin bool              `json:"root_admin"`
	Coins     int               `json:"coins"`
	Suspended bool              `json:"suspended"`
	CreatedAt string            `json:"created_at"`
	Servers   int              `json:"servers_count"`
}

type UserDetail struct {
	ID         uint              `json:"id"`
	UUID       string            `json:"uuid"`
	Username   string            `json:"username"`
	Email      string            `json:"email"`
	NameFirst  string            `json:"name_first"`
	NameLast   string            `json:"name_last"`
	Role       string            `json:"role"`
	RootAdmin  bool              `json:"root_admin"`
	Coins      int               `json:"coins"`
	Suspended  bool              `json:"suspended"`
	CreatedAt  string            `json:"created_at"`
	UpdatedAt  string            `json:"updated_at"`
	ServerIDs  []uint            `json:"server_ids,omitempty"`
}

func TransformUser(user models.User) UserItem {
	item := UserItem{
		ID:        user.ID,
		UUID:      user.UUID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		RootAdmin: user.RootAdmin,
		Coins:     user.Coins,
		Suspended: user.Suspended,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if len(user.Servers) > 0 {
		item.Servers = len(user.Servers)
	}
	return item
}

func TransformUserDetail(user models.User) UserDetail {
	detail := UserDetail{
		ID:        user.ID,
		UUID:      user.UUID,
		Username:  user.Username,
		Email:     user.Email,
		NameFirst: user.NameFirst,
		NameLast:  user.NameLast,
		Role:      user.Role,
		RootAdmin: user.RootAdmin,
		Coins:     user.Coins,
		Suspended: user.Suspended,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	for _, s := range user.Servers {
		detail.ServerIDs = append(detail.ServerIDs, s.ID)
	}

	return detail
}

func TransformUsers(users []models.User) []UserItem {
	items := make([]UserItem, len(users))
	for i, u := range users {
		items[i] = TransformUser(u)
	}
	return items
}
