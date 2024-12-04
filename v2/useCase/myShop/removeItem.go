package myShop

import (
	"ecommerce/ent"
	"fmt"
)

func RemoveItem(uID string) {
	fmt.Println(" === Remove Item")
	ent.UpdateData(0)

	// get user index & failsafe if shop is empty
	UserDatum, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}
	if UserDatum.Shop.State == 0 {
		fmt.Println("Error: shop for called uID is not made, uID:", uID)
		return
	}

	itemDatas := ent.Users[uID].Shop.Items
	if len(itemDatas) == 0 {
		fmt.Println("No item on list.")
		return
	}
	for key, val := range itemDatas {
		fmt.Printf("- %s : %s\n", key[6:], val.Name)
	}

	// asks for item to remove & failsafe if index is out of range
	itemID := ""
	for scs := false; !scs; {
		itemID, scs = ent.InputText("Item ID (0 to cancel): ", true)
		if itemID == "0" {
			fmt.Println("Cancelled removing item.")
			return
		}
		for key := range itemDatas {
			if key[6:] == itemID {
				itemID = key
				scs = false
				break
			}
		}
		scs = !scs
		if !scs {
			fmt.Println("Error: Item ID doesn't exist.")
		}
	}

	// remove prompt
	inputTxt, _ := ent.InputText("Type 'Remove.' exactly to remove item: ", true)
	if inputTxt == "Remove." {
		ent.RemoveItem(uID, itemID)
		fmt.Println("Item removed successfully.")
		return
	} else {
		fmt.Println("Remove failed, typed incorrectly.")
		return
	}
}
