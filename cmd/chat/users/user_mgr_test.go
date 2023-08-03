package users

import (
	"chat/pkg/log"
	"testing"
)

func TestUserManager(t *testing.T) {
	userManager := UserManager{}
	user, found := userManager.GetUser("zhangSan")
	if found {
		log.Infof("found: %v", user)
	} else {
		log.Infof("zhangSan not found")
	}

	userManager.AddUser(User{
		AccountName: "zhangSan",
	})
	user, found = userManager.GetUser("zhangSan")
	if found {
		log.Infof("found: %v", user)
	} else {
		log.Infof("zhangSan not found")
	}

	userManager.DeleteUser("zhangSan")
	user, found = userManager.GetUser("zhangSan")
	if found {
		log.Infof("found: %v", user)
	} else {
		log.Infof("zhangSan not found")
	}
}
