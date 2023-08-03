package model

const (
	AccountNameColumn = "accountName"
)

type UserAccount struct {
	AccountName string `gorm:"primarykey;column:accountName;not null"` // user customize column name
	Password    string `gorm:"column:password"`
}

// TableName customize table name
func (u *UserAccount) TableName() string {
	return "TbUserAccount"
}
