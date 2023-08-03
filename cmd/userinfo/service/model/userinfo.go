package model

const (
	AccountNameColumn = "accountName"
	NickNameColumn    = "nickName"
	BirthdayColumn    = "birthday"
	GenderColumn      = "gender"
	AvatarPathColumn  = "avatarPath"
)

type UserInfo struct {
	AccountName string `gorm:"primarykey;column:accountName;not null"`
	NickName    string `gorm:"column:nickName"`
	Birthday    string `gorm:"column:birthday"`
	Gender      int    `gorm:"column:gender"`
	AvatarPath  string `gorm:"column:avatarPath"`
}

func (u *UserInfo) TableName() string {
	return "TbUserInfo"
}
