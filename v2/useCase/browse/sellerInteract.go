package browse

import (
	"ecommerce/ent"
	"fmt"
)

func SellerInteract(uID, sID string) {
	fmt.Println(" === Interact")

	ent.UpdateData(0)
	if _, scs := ent.Users[uID]; !scs {
		fmt.Println("Error: user not found")
		return
	}
	if _, scs := ent.Users[sID]; !scs {
		fmt.Println("Error: shop not found")
		return
	}

	help := true
	for {
		if help {
			fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
			if ent.GetShopRating(sID) == -1 {
				fmt.Printf("Shop: %s (Rating: not rated)\n", ent.Users[sID].Shop.ShopName)
			} else {
				fmt.Printf("Shop: %s (Rating: %.2f)\n", ent.Users[sID].Shop.ShopName, ent.GetShopRating(sID))
			}
			fmt.Println("1. Buy Items")
			fmt.Println("2. Rate Shop")
			fmt.Println("0. Return")
			help = false
		}

		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Printf("Exiting shop %s.\n", ent.Users[sID].Shop.ShopName)
			return
		case "1":
			for redo := true; redo; {
				redo = buyItems(uID, sID)
			}
			ent.WaitEnter()
			help = true
		case "2":
			rateShop(uID, sID)
			ent.WaitEnter()
			help = true
		default:
			fmt.Println("Unknown action.")
		}
	}
}

func buyItems(uID, sID string) bool {
	fmt.Println(" === Buy Items")

	ent.UpdateData(0)
	if _, scs := ent.Users[uID]; !scs {
		fmt.Println("Error: user not found")
		return false
	}
	if _, scs := ent.Users[sID]; !scs {
		fmt.Println("Error: shop not found")
		return false
	}

	ent.ShowShop(sID, 2, 0, 0)
	itemDatas := ent.Users[sID].Shop.Items
	itemID := ""

	// get item index
	fmt.Println("Your money:", ent.Users[uID].Money)
	for scs := false; !scs; {
		itemID, scs = ent.InputText("Item ID (0 to cancel): ", true)
		if itemID == "0" {
			fmt.Println("Cancelled buying item.")
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

	// output item data
	itemName := itemDatas[itemID].Name
	itemAmount := itemDatas[itemID].Amount
	fmt.Printf("\nItem selected: %s\nCost: %d; Amount: %d\n", itemName, itemDatas[itemID].Cost, itemAmount)
	itemBought, purchaseCost := 0, 0
	fmt.Println("Your money:", ent.Users[uID].Money)
	for inputScs := false; !inputScs; {
		itemBought, inputScs = ent.InputNum("Item amount (0 to cancel): ", false)
		if int(itemAmount)-itemBought < 0 {
			fmt.Println("Insufficient stock.")
			inputScs = false
		}
		purchaseCost = itemBought * int(itemDatas[itemID].Cost)
		if ent.Users[uID].Money-purchaseCost < 0 {
			fmt.Println("Not enough money.")
			inputScs = false
		}
	}
	if itemBought == 0 {
		fmt.Println("Action cancelled.")
	} else {
		ent.WriteShopData(sID, itemID, ent.StrConvIToA(int(itemDatas[itemID].Amount)-itemBought), 7)
		ent.CreateLog("item buy", uID, uID, itemBought, itemID, sID)
		ent.CreateLog("item buy", sID, uID, itemBought, itemID, sID)

		ent.WriteUserData(uID, ent.StrConvIToA(ent.Users[uID].Money-purchaseCost), 2)
		ent.CreateLog("money transfer", uID, uID, purchaseCost, sID)

		ent.WriteUserData(sID, ent.StrConvIToA(ent.Users[sID].Money+purchaseCost), 2)
		ent.CreateLog("money transfer", sID, uID, purchaseCost, sID)

		fmt.Println("Item successfully bought.")
	}
	return ent.RedoAction("Buy again? (1 for yes): ")
}

func rateShop(uID, sID string) {
	fmt.Println(" === Rate Shop")

	ent.UpdateData(0)
	if _, scs := ent.Users[uID]; !scs {
		fmt.Println("Error: user not found")
		return
	}
	if _, scs := ent.Users[sID]; !scs {
		fmt.Println("Error: shop not found")
		return
	}

	rateVal := 0
	for inputScs := false; !inputScs; {
		rateVal, inputScs = ent.InputNum("Rate value (0-100; -1 to cancel): ", true)
		if !inputScs {
			continue
		}
		if rateVal == -1 {
			fmt.Println("Cancelled rating shop.")
			return
		}
		if 0 <= rateVal && rateVal <= 100 {
			break
		}
		fmt.Println("Error: Value cannot be negative (except -1).")
		inputScs = false
	}
	ent.DoRating(uID, sID, rateVal)
	fmt.Printf("Successfully rate shop %s with value %d.\n", ent.Users[sID].Shop.ShopName, rateVal)
}
