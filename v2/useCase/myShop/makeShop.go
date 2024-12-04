package myShop

import (
	"ecommerce/ent"
	"fmt"
)

func MakeShop(uID string) {
	fmt.Println(" === Make Shop")
	UserDatum, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}
	if UserDatum.Shop.State != 0 {
		fmt.Println("Error: shop for called uID is already made, uID:", uID)
		return
	}

	shopName := ""
	for scs := false; !scs; {
		shopName, scs = ent.InputText("Shop name (0 to cancel): ", false)
	}
	if shopName == "0" {
		fmt.Println("Action cancelled.")
		return
	}
	ent.AddShop(uID, shopName)
	fmt.Println("Shop successfully made.")
}
