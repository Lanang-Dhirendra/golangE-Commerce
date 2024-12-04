package myShop

import (
	"ecommerce/ent"
	"fmt"
)

func AddItem(uID string) bool {
	fmt.Println("\n === Add New Item")
	fmt.Println("(type '-1' at any point to cancel adding item)")
	ent.UpdateData(0)

	// find key for current user & failsafe if not found
	UserDatum, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return false
	}
	if UserDatum.Shop.State == 0 {
		fmt.Println("Error: shop for called uID is not made, uID:", uID)
		return false
	}

	// input item data
	var itemName, itemDesc string
	var itemAmount, itemCost int
	itemDatas := UserDatum.Shop.Items
	for inputScs := false; !inputScs; { // input item name
		itemName, inputScs = ent.InputText("Item name: ", false)
		for _, itemDat := range itemDatas {
			if itemName == itemDat.Name {
				fmt.Println("Error: Item already exists.")
				inputScs = false
				break
			}
		}
	}
	if itemName == "-1" {
		fmt.Println("Cancelled adding item.")
		return false
	}

	// input item description
	itemDesc, _ = ent.InputText("Item description (optional): ", true)
	if itemDesc == "-1" {
		fmt.Println("Cancelled adding item.")
		return false
	}

	for inputScs := false; !inputScs; { // input item amount
		itemCost, inputScs = ent.InputNum("Item cost: ", true)
		if itemCost == -1 {
			fmt.Println("Cancelled adding item.")
			return false
		}
		if itemCost < 0 {
			fmt.Println("Error: Value cannot be negative (except -1).")
			inputScs = false
		}
	}

	for inputScs := false; !inputScs; { // input item amount
		itemAmount, inputScs = ent.InputNum("Item amount: ", true)
		if itemAmount == -1 {
			fmt.Println("Cancelled adding item.")
			return false
		}
		if itemAmount < 0 {
			fmt.Println("Error: Value cannot be negative (except -1).")
			inputScs = false
		}
	}

	itemID := ent.AddItem(uID, itemName, itemDesc, itemCost, itemAmount)
	ent.CreateLog("add item", uID, itemID, itemName, itemCost, itemAmount)

	// output & ask to re-add new item(s)
	fmt.Printf("Item %s successfully added.\n", itemName)
	return ent.RedoAction("Add again? (1 for yes): ")
}
