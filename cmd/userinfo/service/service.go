package service

import (
	"chat/cmd/userinfo/routing/api"
	"chat/cmd/userinfo/service/dao"
	"chat/cmd/userinfo/service/model"
	"chat/pkg/log"
)

func AdduserInfo(param api.AddUserInfoParam) bool {
	err := dao.AddUserInfo(model.UserInfo{
		AccountName: param.AccountName,
		NickName:    param.NickName,
		Birthday:    param.Birthday,
		Gender:      param.Gender,
	})
	if err != nil {
		log.Errorf("account[%v] add userInfo failed, err: %v", param.AccountName, err)
		return false
	}
	return true
}

func DelUserInfo(accountName string) bool {
	err := dao.DelUserInfo(accountName)
	if err != nil {
		log.Errorf("account[%v] del failed, err: %v", accountName, err)
		return false
	}
	return true
}

func QueryUserInfo(accountName string) (api.UserInfoResp, error) {
	userInfo, err := dao.GetUserInfo(accountName)
	if err != nil {
		log.Errorf("account[%v] query userInfo failed, err: %v", accountName, err)
		return api.UserInfoResp{}, err
	}
	return api.UserInfoResp{
		NickName: userInfo.NickName,
		Birthday: userInfo.Birthday,
		Gender:   userInfo.Gender,
	}, nil
}

func UpdateUserInfo(param api.UpdateUserInfoParam) bool {
	err := dao.UpdateUserInfo(model.UserInfo{
		AccountName: param.AccountName,
		NickName:    param.NickName,
		Birthday:    param.Birthday,
		Gender:      param.Gender,
	})
	if err != nil {
		log.Errorf("account[%v] update info failed, err: %v", param.AccountName, err)
		return false
	}
	return true
}

func UpdateNickName(accountName, nickname string) bool {
	err := dao.UpdateNickName(accountName, nickname)
	if err != nil {
		log.Errorf("account[%v] nickname update failed, err: %v", accountName, err)
		return false
	}
	return true
}

func UpdateGender(accountName string, gender int) bool { // TODO support in the future
	err := dao.UpdateGender(accountName, gender)
	if err != nil {
		log.Errorf("account[%v] gender update failed, err: %v", accountName, err)
		return false
	}
	return true
}

func UpdateBirthday(accountName string, birthday string) bool { // TODO support in the future
	err := dao.UpdateBirthday(accountName, birthday)
	if err != nil {
		log.Errorf("account[%v] birthday update failed, err: %v", accountName, err)
		return false
	}
	return true
}

func UpdateAvatarPath(accountName string, avatarPath string) bool {
	err := dao.UpdateAvatarPath(accountName, avatarPath)
	if err != nil {
		log.Errorf("account[%v] birthday update failed, err: %v", accountName, err)
		return false
	}
	return true
}
