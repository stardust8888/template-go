package model

type KeyCenter struct {
	GormModel
	Secret string `gorm:"column:secret;size:256;comment:密钥密文(hex);default:'';not null" json:"-"`
	Desc   string `gorm:"column:description;size:256;comment:描述;default:'';not null" json:"desc"`
}

func (KeyCenter) TableName() string {
	return "key_center"
}
