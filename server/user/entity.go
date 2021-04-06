package user

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	Token string `json:"token"`
}
