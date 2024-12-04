package myShop

import (
	"ecommerce/ent"
	"fmt"
)

func StockItem(uID string) bool {
	fmt.Println("\n === Stock Item")
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
		fmt.Printf("- %s : %s (%d)\n", key[6:], val.Name, val.Amount)
	}

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

	stockType := ""
	for inputScs := false; !inputScs; {
		stockType, _ = ent.InputText("Stock type (add/remove): ", true)
		switch ent.ToLowerTxt(stockType) {
		case "add", "+":
			stockType = "+"
			inputScs = true
		case "remove", "-":
			stockType = "-"
			inputScs = true
		case "cancel", "0":
			fmt.Println("Cancelled stocking item.")
			return false
		default:
			fmt.Println("Unknown input.")
		}
	}
	stockDif := 0
	for inputScs := false; !inputScs; {
		stockDif, inputScs = ent.InputNum("Amount (-1 to cancel): ", true)
		if stockDif == -1 {
			fmt.Println("Cancelled stocking item.")
			return false
		}
		if stockDif < 0 {
			fmt.Println("Error: Value cannot be negative (except -1).")
			inputScs = false
		}
	}
	if stockType == "-" {
		stockDif *= -1
	}
	itemAmountNew := int(itemDatas[itemID].Amount) + stockDif
	ent.WriteShopData(uID, itemID, ent.StrConvIToA(itemAmountNew), 7)
	ent.CreateLog("stock item", uID, itemID, stockDif, itemAmountNew)
	fmt.Printf("Stock for item %s is changed to %d.\n", itemDatas[itemID].Name, itemAmountNew)
	return ent.RedoAction("Again? ")
}
