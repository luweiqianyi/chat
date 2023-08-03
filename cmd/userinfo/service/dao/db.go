package dao

import (
	"chat/cmd/userinfo/global"
	"chat/cmd/userinfo/service/model"
	"chat/pkg/log"
	"chat/pkg/mysql"
	"errors"
	"fmt"
)

func CreateTable() error {
	db := global.GetMySqlDB()
	if db == nil {
		log.Errorf("db required")
		return errors.New("db required")
	}

	err := db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&model.UserInfo{})
	if err != nil {
		log.Errorf("db table:TbUserInfo  autoMigrate failed, err: %v", err)
	}
	return err
}

func AddUserInfo(info model.UserInfo) error {
	return mysql.AddDTO(global.GetMySqlDB(), &info)
}

func DelUserInfo(accountName string) error {
	var userInfo model.UserInfo
	return mysql.DelDTO(
		global.GetMySqlDB(),
		&userInfo,
		fmt.Sprintf("%s=?", model.AccountNameColumn),
		accountName)
}

func GetUserInfo(accountName string) (model.UserInfo, error) {
	var userInfo model.UserInfo
	err := mysql.QueryDTO(
		global.GetMySqlDB(),
		&userInfo,
		fmt.Sprintf("%v=?", model.AccountNameColumn),
		accountName)
	if err != nil {
		return model.UserInfo{}, err
	}
	return userInfo, nil
}

func UpdateUserInfo(info model.UserInfo) error {
	return mysql.UpdateDTO(global.GetMySqlDB(), info)
}

func UpdateNickName(accountName, nickname string) error {
	return mysql.UpdateOneColumn(
		global.GetMySqlDB(),
		&model.UserInfo{},
		model.NickNameColumn,
		nickname,
		fmt.Sprintf("%v=?", model.AccountNameColumn),
		accountName)
}

func UpdateGender(accountName string, gender int) error {
	return mysql.UpdateOneColumn(
		global.GetMySqlDB(),
		&model.UserInfo{},
		model.GenderColumn,
		gender,
		fmt.Sprintf("%v=?", model.AccountNameColumn),
		accountName)
}

func UpdateBirthday(accountName string, birthday string) error {
	return mysql.UpdateOneColumn(
		global.GetMySqlDB(),
		&model.UserInfo{},
		model.BirthdayColumn,
		birthday,
		fmt.Sprintf("%v=?", model.AccountNameColumn),
		accountName)
}

func UpdateAvatarPath(accountName string, avatarPath string) error {
	return mysql.UpdateOneColumn(
		global.GetMySqlDB(),
		&model.UserInfo{},
		model.AvatarPathColumn,
		avatarPath,
		fmt.Sprintf("%v=?", model.AccountNameColumn),
		accountName)
}
