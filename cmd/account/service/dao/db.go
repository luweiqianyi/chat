package dao

import (
	"chat/cmd/account/global"
	"chat/cmd/account/service/model"
	"chat/pkg/log"
	"errors"
	"fmt"
)

func CreateTable() error {
	db := global.MySqlDB()
	if db == nil {
		log.Errorf("db required")
		return errors.New("db required")
	}

	err := db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&model.UserAccount{})
	if err != nil {
		log.Errorf("db table:TbUserAccount  autoMigrate failed, err: %v", err)
	}
	return err
}

func FindAccount(accountName string) bool {
	db := global.MySqlDB()
	if db == nil {
		log.Errorf("db required")
		return false
	}

	var userAccount model.UserAccount
	rowsAffected := db.Select([]string{model.AccountNameColumn}).Find(
		&userAccount,
		fmt.Sprintf("%s=?", model.AccountNameColumn),
		accountName).RowsAffected
	if rowsAffected == 0 {
		log.Errorf("account[%v] find none record", accountName)
		return false
	}
	return true
}

func AddAccount(accountName, password string) bool {
	db := global.MySqlDB()
	if db == nil {
		log.Errorf("db required")
		return false
	}

	err := db.Create(&model.UserAccount{
		AccountName: accountName,
		Password:    password,
	}).Error
	if err != nil {
		log.Errorf("account[%v] add failed, err: %v", accountName, err)
		return false
	}
	return true
}

func DelAccount(accountName string) bool {
	db := global.MySqlDB()
	if db == nil {
		log.Errorf("db required")
		return false
	}

	var userAccount model.UserAccount
	err := db.Delete(&userAccount, fmt.Sprintf("%s=?", model.AccountNameColumn), accountName).Error
	if err != nil {
		log.Errorf("account[%v] delete failed, err: %v\n", accountName, err)
		return false
	}
	return true
}
