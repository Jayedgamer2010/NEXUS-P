package requests

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	Role     string `json:"role" validate:"oneof=client admin"`
	RootAdmin bool   `json:"root_admin"`
	Coins    int    `json:"coins"`
}

type UpdateUserRequest struct {
	Username  string `json:"username" validate:"omitempty,min=2,max=50"`
	Email     string `json:"email" validate:"omitempty,email,max=255"`
	Role      string `json:"role" validate:"omitempty,oneof=client admin"`
	RootAdmin *bool  `json:"root_admin"`
	Coins     *int   `json:"coins"`
	Password  string `json:"password" validate:"omitempty,min=8,max=255"`
	Suspended *bool  `json:"suspended"`
}

type UpdateAccountRequest struct {
	Email    string `json:"email" validate:"omitempty,email,max=255"`
	Password string `json:"password" validate:"omitempty,min=8,max=255"`
}

type CreateEggRequest struct {
	Author      string `json:"author" validate:"required,max=255"`
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
	DockerImage string `json:"docker_image" validate:"required,max=255"`
	Startup     string `json:"startup"`
	ConfigStop  string `json:"config_stop"`
}

type UpdateEggRequest struct {
	Author      string `json:"author" validate:"omitempty,max=255"`
	Name        string `json:"name" validate:"omitempty,min=2,max=255"`
	Description string `json:"description"`
	DockerImage string `json:"docker_image" validate:"omitempty,max=255"`
	Startup     string `json:"startup"`
	ConfigStop  string `json:"config_stop"`
}
