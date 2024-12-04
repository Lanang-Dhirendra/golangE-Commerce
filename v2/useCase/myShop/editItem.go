package myShop

import (
	"ecommerce/ent"
	"fmt"
)

func EditItem(uID string) bool {
	fmt.Println("\n === Edit Item")
	ent.UpdateData(0)

	// get user index & failsafe if shop is empty
	UserDatum, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return false
	}
	if UserDatum.Shop.State == 0 {
		fmt.Println("Error: shop for called uID is not made, uID:", uID)
		return false
	}
	itemDatas := ent.Users[uID].Shop.Items
	if len(itemDatas) == 0 {
		fmt.Println("No item on list.")
		return false
	}
	for key, val := range itemDatas {
		fmt.Printf("- %s : %s\n", key[6:], val.Name)
	}

	// asks for item to remove & failsafe if index is out of range
	itemID := ""
	for scs := false; !scs; {
		itemID, scs = ent.InputText("Item ID (0 to cancel): ", true)
		if itemID == "0" {
			fmt.Println("Cancelled editing item.")
			return false
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

	// outputs item data
	itemDat := itemDatas[itemID]
	itemDesc := itemDat.Desc
	itemAmount := ""
	if itemDesc == "" {
		itemDesc = "- (no desc)"
	}
	if itemDat.Amount == 0 {
		itemAmount = "Empty"
	} else {
		itemAmount = "Amount: " + ent.StrConvIToA(int(itemDat.Amount))
	}
	fmt.Printf("\nItem selected\n%s (%s)\nCost: %d\n%s\n", itemDat.Name, itemAmount, itemDat.Cost, itemDesc)

	// asks for item data to edit
	var element string
	for scs := false; !scs; {
		element, _ = ent.InputText("Select item element to edit.\n1 -> Name\n2 -> Desc\n3 -> Cost\n0 -> Cancel\nInput: ", true)
		for _, c := range []string{"0", "1", "2", "3"} {
			if element == c {
				scs = true
				break
			}
		}
		if !scs {
			fmt.Println("Error: Element index out of reach.")
		}
	}

	// edit process
	switch element {
	case "0": // cancel
		fmt.Println("Cancelled editing item.")
		return false
	case "1": // edit item name
		itemNameNew := ""
		for inputScs := false; !inputScs; {
			itemNameNew, inputScs = ent.InputText("New item name (0 to cancel): ", false)
			if !inputScs {
				continue
			}
			for _, val := range itemDatas {
				if itemNameNew == val.Name {
					fmt.Println("Error: Item already exists.")
					inputScs = false
					break
				}
			}
		}
		if itemNameNew == "0" {
			fmt.Println("Cancelled editing item.")
			break
		}
		ent.WriteShopData(uID, itemID, itemNameNew, 4)
		ent.CreateLog("edit item x desc", uID, itemID, "name", itemNameNew)
		fmt.Println("Item name successfully changed.")
	case "2": // edit item description
		itemDescNew, _ := ent.InputText("New item description (0 to cancel): ", false)
		if itemDescNew == "0" {
			fmt.Println("Cancelled editing item.")
			break
		}
		ent.WriteShopData(uID, itemID, itemDescNew, 5)
		ent.CreateLog("edit item desc", uID, itemID)
		fmt.Println("Item description successfully changed.")
	case "3": // edit item cost
		itemCostNew := 0
		for inputScs := false; !inputScs; {
			itemCostNew, inputScs = ent.InputNum("Item cost (-1 to cancel): ", true)
			if itemCostNew == -1 {
				fmt.Println("Cancelled adding item.")
				return false
			}
			if itemCostNew < 0 {
				fmt.Println("Error: Value cannot be negative (except -1).")
				inputScs = false
			}
		}
		ent.WriteShopData(uID, itemID, ent.StrConvIToA(itemCostNew), 6)
		ent.CreateLog("edit item x desc", uID, itemID, "cost", itemCostNew)
		fmt.Println("Item cost successfully changed.")
	}
	return ent.RedoAction("Edit again? (1 for yes): ")
}
