package structures

type UserData struct {
	UserName       string `json:"username"`
	PasswordHash   string `json:"password_hash"`
	CollectionName string `json:"collection_name"`
}
