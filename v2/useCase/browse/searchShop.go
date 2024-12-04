package browse

import (
	"ecommerce/ent"
	"fmt"
)

func SearchShop(uID string) bool {
	fmt.Println(" === Search Shop")
	ent.UpdateData(0)

	// find all sellers that isn't the current user
	sellers, userSells := map[string]ent.UserData{}, false
	for key, sellerData := range ent.Users {
		if key == uID && sellerData.Shop.State != 0 {
			userSells = true
			continue
		}
		if sellerData.Shop.State != 2 {
			continue
		}
		sellers[key] = sellerData
		fmt.Printf("%d. %s", len(sellers), sellerData.Shop.ShopName)
		fmt.Println()
	}
	if len(sellers) == 0 {
		fmt.Print("There are currently no shops")
		if userSells {
			fmt.Print(" (that isn't yours)")
		}
		fmt.Println(".")
		return false
	}

	// input seller
	sellerID := ""
	for sellerID == "" {
		// input
		sellerName, _ := ent.InputText("Enter shop name (0 to cancel): ", true)

		// incorrect input check
		if sellerName == "0" {
			fmt.Println("Action cancelled.")
			return false
		}
		for key, temp := range sellers {
			if temp.Shop.ShopName == sellerName {
				sellerID = key
				break
			}
		}
		if sellerID == "" {
			if userSells && sellerName == ent.Users[uID].Username {
				fmt.Println("Can't buy from self.")
			} else {
				fmt.Println("Error: seller name not found")
			}
		}
	}
	fmt.Println()
	SellerInteract(uID, sellerID)
	return ent.RedoAction("Search again? (1 for yes): ")
}
