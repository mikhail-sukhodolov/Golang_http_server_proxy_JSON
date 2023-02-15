package structs


type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

type Storage struct {
	users      map[int]*User
	lastUserId int
}

type MakeFriendsRequest struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type CreateUserResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

type MakeFriendsResponse struct {
	Message string `json:"msg"`
}

type DeleteUserRequest struct {
	TargetId int `json:"target_id"`
}

type UserAge struct {
	Age int `json:"new_age"`
}
