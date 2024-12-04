package myShop

import (
	"ecommerce/ent"
	"fmt"
)

func ChangeShopName(uID string) {
	fmt.Println(" === Change Shop Name")
	ent.UpdateData(0)
	UserDatum, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}
	if UserDatum.Shop.State == 0 {
		fmt.Println("Error: shop for called uID is not made, uID:", uID)
		return
	}

	shopName := UserDatum.Shop.ShopName
	shopNameNew := ""
	for inputScs := false; !inputScs; {
		shopNameNew, inputScs = ent.InputText("New shop name (0 to cancel): ", false)
		if !inputScs {
			continue
		}
		if shopNameNew == "0" {
			fmt.Println("Cancelled changing shop name.")
			return
		}
		if shopNameNew == shopName {
			fmt.Println("Error: name inputted is the same name.")
			inputScs = false
			continue
		}
		for _, val := range ent.Users {
			if shopNameNew == val.Shop.ShopName {
				fmt.Println("Error: Name already exists.")
				inputScs = false
				break
			}
		}
	}

	ent.WriteShopData(uID, "", shopNameNew, 1)
	ent.CreateLog("change shop name", uID, shopName, shopNameNew)
	fmt.Println("Shop name successfully changed.")
}
