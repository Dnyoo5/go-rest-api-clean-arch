package models

type Product struct {
	Id    int    `json:"id"`
	Nama  string `json:"nama" validate:"required"`
	Harga int    `json:"harga" validate:"required"`
	Stok  int    `json:"stok" validate:"required,gte=0"`
}