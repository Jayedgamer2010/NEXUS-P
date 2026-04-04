package requests

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Role     string `json:"role" validate:"required,oneof=admin client"`
	Coins    int    `json:"coins"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=100"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=8,max=72"`
	Role     *string `json:"role" validate:"omitempty,oneof=admin client"`
	Coins    *int    `json:"coins"`
	RootAdmin *bool  `json:"root_admin"`
}

type UpdateAccountRequest struct {
	Email       *string `json:"email" validate:"omitempty,email"`
	Password    *string `json:"password" validate:"omitempty,min=8,max=72"`
	NameFirst   *string `json:"name_first"`
	NameLast    *string `json:"name_last"`
	Language    *string `json:"language"`
}
