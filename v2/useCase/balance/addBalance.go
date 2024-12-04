package balance

import (
	"ecommerce/ent"
	"fmt"
)

func AddBalance(uID string) {
	fmt.Println(" === Add Balance")
	ent.UpdateData(0)

	UserDat, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}

	fmt.Printf("Your current balance: %d\n", UserDat.Money)

	plusBal := 0
	for inputScs := false; !inputScs; {
		plusBal, inputScs = ent.InputNum("Add amount (0 to cancel): ", false)
		if !inputScs {
			continue
		}
		if plusBal == 0 {
			fmt.Println("Cancelled adding balance.")
			return
		}
	}

	balanceNew := UserDat.Money + plusBal
	ent.WriteUserData(uID, ent.StrConvIToA(balanceNew), 2)
	ent.CreateLog("deposit money", uID, plusBal, balanceNew)
	fmt.Printf("Successfully added %d to your account (now: %d)\n", plusBal, balanceNew)
}
