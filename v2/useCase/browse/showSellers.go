package browse

import (
	"ecommerce/ent"
	"fmt"
)

func ShowSellers(uID string) { // +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=
	fmt.Println(" === Show Sellers")
	ent.UpdateData(0)

	sellers := 0
	for id, seller := range ent.Users {
		if seller.Shop.State == 0 {
			continue
		}
		if len(seller.Shop.Items) == 0 {
			continue
		}
		sellers++
		shopState := ent.GetShopState(id).Str
		if ent.GetShopRating(id) == -1 {
			fmt.Printf("%d. %s %s (Rating: not rated)", sellers, seller.Shop.ShopName, shopState)
		} else {
			fmt.Printf("%d. %s %s (Rating: %.2f)", sellers, seller.Shop.ShopName, shopState, ent.GetShopRating(id))
		}
		if seller.Username == uID {
			fmt.Print(" (You)")
		}
		fmt.Println()
	}
	if sellers == 0 {
		fmt.Println("There are currently no sellers.")
	}
}
