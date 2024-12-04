package settings

import (
	"ecommerce/ent"
	"fmt"
)

func ChangeUserName(uID string) {
	fmt.Println(" === Change User Name")
	UserDatum, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}

	name := UserDatum.Username
	nameNew := ""
	for inputScs := false; !inputScs; {
		nameNew, inputScs = ent.InputText("New name (0 to cancel): ", false)
		if !inputScs {
			continue
		}
		if nameNew == "0" {
			fmt.Println("Cancelled changing shop name.")
			return
		}
		if nameNew == name {
			fmt.Println("Error: name inputted is the same name.")
			inputScs = false
			continue
		}
		for _, val := range ent.Users {
			if nameNew == val.Username {
				fmt.Println("Error: Name already exists.")
				inputScs = false
				break
			}
		}
	}

	ent.WriteUserData(uID, nameNew, 1)
	ent.CreateLog("change user name", uID, name, nameNew)
	fmt.Println("User name successfully changed.")
}
