package myShop

import (
	"ecommerce/ent"
	"fmt"
)

func ShowItemList(uID string) { // self explanatory
	fmt.Println(" === Show Item List")
	ent.ShowShop(uID, 1, 0, 0)
}
