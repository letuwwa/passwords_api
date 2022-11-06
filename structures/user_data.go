package structures

type UserData struct {
	UserName     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	DatabaseName string `json:"database_name"`
}
