package start

import (
	"ecommerce/ent"
	"fmt"
)

// asks name, makes an account if it doesn't exist
func Regis() bool {
	fmt.Println(" === Register")
	name, scs := ent.InputText("(case-sensitive; 0 to cancel)\nEnter name: ", false)
	if scs {
		if name == "0" {
			fmt.Println("Register cancelled.")
			return false
		}
		if found, bc := ent.CheckBannedChars(name); found {
			fmt.Printf("Register failed, contains banned characters %s.\n(all banned characters: %s)", bc, ent.BannedChars)
		}
		ent.UpdateData(0)
		for _, user := range ent.Users {
			if name == user.Username {
				fmt.Println("Register failed, name already exists.")
				return false
			}
		}
		userID := ent.CreateUser(name)
		fmt.Printf("Register successful.\nYour ID: %s\n", userID[2:])
		return true
	}
	return false
}
