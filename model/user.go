package model

import "template-go/utils/x/timex"

type User struct {
	GormModel
	Username    string `gorm:"column:username;size:64;comment:用户名;default:'';not null" json:"username"`
	Phone       string `gorm:"column:phone;size:64;comment:手机号;default:'';not null" json:"-"`
	PhoneSuffix string `gorm:"column:phone_suffix;size:32;comment:手机号;default:'';not null" json:"-"`
	Wx          string `gorm:"column:wx;size:64;comment:微信;default:'';not null" json:"wx"`
	WxOpenID    string `gorm:"column:wx_open_id;size:64;uniqueIndex;comment:微信openid;default:'';not null" json:"wxOpenID"`
	WxUnionID   string `gorm:"column:wx_union_id;size:64;comment:微信unionid;default:'';not null" json:"-"`
	SessionKey  string `gorm:"column:session_key;size:64;comment:微信sessionkey;default:'';not null" json:"-"`
	Avatar      string `gorm:"column:avatar;size:256;comment:头像;default:'';not null" json:"avatar"`
	Role        int    `gorm:"column:role;type:tinyint;comment:角色 0:未记录1:父亲2:母亲;default:0;not null" json:"role"`
	IsVip       int    `gorm:"column:is_vip;type:tinyint;comment:是否会员 0:否1:是;default:0;not null" json:"isVip"`
	Status      int    `gorm:"column:status;type:tinyint;comment:状态 0:正常1:禁用;default:0;not null" json:"status"`
	RegisterIP  string `gorm:"column:register_ip;size:64;comment:注册ip;default:'';not null" json:"registerIP"`
	Promoter    string `gorm:"column:promoter;size:64;comment:推广人;default:'';not null" json:"promoter"`
}

func (*User) TableName() string {
	return "user"
}

type UserDeleteApply struct {
	GormModel
	UserID    int        `gorm:"column:user_id;index;comment:用户ID;default:0;not null" json:"userID"`
	ApplyAt   timex.Time `gorm:"column:apply_at;type:timestamp;comment:申请时间;not null" json:"applyAt"`
	ProcessAt timex.Time `gorm:"column:process_at;type:timestamp;comment:处理时间(到期删除或取消时间)" json:"processAt"`
	Status    int        `gorm:"column:status;type:tinyint;comment:状态 0:待处理 1:已处理 2:已取消;default:0;not null" json:"status"`
}

func (*UserDeleteApply) TableName() string {
	return "user_delete_apply"
}
