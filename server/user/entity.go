package user

import "ca-tech-dojo/model/user"

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	Token string `json:"token"`
}

type GetUserResponse struct {
	Name  string `json:"name"`
	Point uint   `json:"point"`
}

func newGetUserResponse(userEntity user.User) GetUserResponse {
	resp := GetUserResponse{
		Name:  userEntity.Name,
		Point: userEntity.Point,
	}
	return resp
}
