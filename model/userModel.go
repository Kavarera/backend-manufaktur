package model

type User struct {
	UserID     string `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	HakAkses   int    `json:"hak_akses"`
	IdHakAkses int    `json:"hak_id"`
}

type GetUser struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	HakAkses int    `json:"hak_akses"`
}
