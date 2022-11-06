package structures

type UserData struct {
	UserName       string `json:"user_name"`
	PasswordHash   string `json:"password_hash"`
	CollectionName string `json:"collection_name"`
}
