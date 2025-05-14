package model

type user struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	hakAkses int    `json:"hak_akses"`
}
