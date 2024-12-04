package start

import (
	"ecommerce/ent"
	"fmt"
)

// asks name, logs in w/ that name if it exists
func Login() (bool, string) {
	fmt.Println(" === Login")
	id, scs := ent.InputText("(case-sensitive; 0 to cancel)\nEnter ID: ", false)
	if scs {
		if id == "0" {
			fmt.Println("Login cancelled.")
			return false, ""
		}
		ent.UpdateData(0)
		if id[0] != '@' {
			id = "@U" + id
		}
		if data, userFound := ent.Users[id]; userFound {
			ent.CreateLog("user login", id)
			fmt.Printf("Login successful as %s.\n", data.Username)
			return true, id
		}
		fmt.Println("Login failed, unregistered ID.")
		ent.WaitEnter()
	}
	return false, ""
}
