package model

type Cart struct{
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	UserID int64 `gorm: "not_null" json:"user_id"`
	ProductID int64 `json:"product_id"`
	SizeID int64 `gorm: "not_null" json:"size_id"`
	Num int64 `gorm: "not_null" json:"num"`
}

