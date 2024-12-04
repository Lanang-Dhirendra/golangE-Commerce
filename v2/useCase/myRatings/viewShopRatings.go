package myRatings

import (
	"ecommerce/ent"
	"fmt"
)

func ViewShopRts(uID string) {
	fmt.Println(" === View Shop Ratings")
	ent.UpdateData(0)

	UserDat, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}
	if UserDat.Shop.State == 0 {
		fmt.Println("Error: shop for called uID is not made, uID:", uID)
		return
	}

	if len(UserDat.Shop.Ratings) == 0 {
		fmt.Println("Shop has not been rated yet.")
		return
	}

	for i, id := range UserDat.Shop.Ratings {
		rating := ent.Ratings[id]
		fmt.Printf("%d. %s; rated %d\n", i+1, ent.Users[rating.UserID].Username, rating.Value)
	}
}
