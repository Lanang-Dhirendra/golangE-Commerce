package myRatings

import (
	"ecommerce/ent"
	"fmt"
)

func ChangeRating(uID string) bool {
	fmt.Println(" === Change Rating")
	ent.UpdateData(0)

	UserDat, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return false
	}

	if len(UserDat.Ratings) == 0 {
		fmt.Println("User has not rated any shop.")
		return false
	}
	for i, id := range UserDat.Ratings {
		rating := ent.Ratings[id]
		fmt.Printf("%d. %s; rated %d\n", i+1, ent.Users[rating.ShopID].Shop.ShopName, rating.Value)
	}

	sellerID := ""
	for sellerID == "" {
		// input
		sellerName, _ := ent.InputText("Enter shop name (0 to cancel): ", true)

		// incorrect input check
		if sellerName == "0" {
			fmt.Println("Action cancelled.")
			return false
		}
		for _, id := range UserDat.Ratings {
			shop := ent.Ratings[id].ShopID
			if ent.Users[shop].Shop.ShopName == sellerName {
				sellerID = shop
				break
			}
		}
		if sellerID == "" {
			fmt.Println("Error: seller name not found")
		}
	}

	rateVal := 0
	for inputScs := false; !inputScs; {
		rateVal, inputScs = ent.InputNum("Rate value (0-100; -1 to cancel): ", true)
		if !inputScs {
			continue
		}
		if rateVal == -1 {
			fmt.Println("Cancelled rating shop.")
			return false
		}
		if 0 <= rateVal && rateVal <= 100 {
			break
		}
		fmt.Println("Error: Value cannot be negative (except -1).")
		inputScs = false
	}
	ent.DoRating(uID, sellerID, rateVal)
	fmt.Printf("Successfully rate shop %s with value %d.\n", ent.Users[sellerID].Shop.ShopName, rateVal)
	return ent.RedoAction("Again? ")
}
