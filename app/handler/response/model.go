package response

type User struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type Group struct {
	GroupID string `json:"groupId"`
	Name    string `json:"name"`
	Users   []User `json:"users"`
}
