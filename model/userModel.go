package model

type User struct {
	UserID     string `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	HakAkses   string `json:"hak_akses"`
	IdHakAkses int    `json:"hak_id"`
}
