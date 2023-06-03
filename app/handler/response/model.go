package response

type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type Group struct {
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
	Users   []User `json:"users"`
}
