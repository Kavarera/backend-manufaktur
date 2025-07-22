package model

type Satuan struct {
	IDSatuan   int    `json:"idSatuan"`
	NamaSatuan string `json:"satuanNama"`
}

type SatuanTurunan struct {
	IDTurunan         int    `json:"idSatuanTurunan"`
	NamaSatuanTurunan string `json:"namaTurunan"`
	IDSatuan          int    `json:"idSatuan"`
}
