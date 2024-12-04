package myRatings

import (
	"ecommerce/ent"
	"fmt"
)

func ViewUserRts(uID string) {
	fmt.Println(" === View User Ratings")
	ent.UpdateData(0)

	UserDat, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}

	if len(UserDat.Ratings) == 0 {
		fmt.Println("User has not rated any shop.")
		return
	}
	for i, id := range UserDat.Ratings {
		rating := ent.Ratings[id]
		fmt.Printf("%d. %s; rated %d\n", i+1, ent.Users[rating.ShopID].Shop.ShopName, rating.Value)
	}
}
