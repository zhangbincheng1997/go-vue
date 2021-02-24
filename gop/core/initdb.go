package core

import (
	"fmt"
	"main/model"
	"main/utils"
	"os"
)

func main() {
	db := MySQL()

	err := db.AutoMigrate(
		model.User{},
		model.Role{},
		model.UserRole{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	users := []model.User{
		{Username: "admin", Password: utils.MD5("admin"), Role: "admin"},
		{Username: "123", Password: utils.MD5("123"), Role: "user"},
	}
	if err := db.Create(&users); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
