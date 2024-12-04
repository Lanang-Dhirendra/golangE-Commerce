package home

import (
	"ecommerce/ent"
	"ecommerce/useCase/balance"
	"ecommerce/useCase/browse"
	"ecommerce/useCase/myLogs"
	"ecommerce/useCase/myRatings"
	"ecommerce/useCase/myShop"
	"ecommerce/useCase/settings"
	"fmt"
)

// ---++===+ Home Page +===++---
func Home(uID string) {
	help := true
	for {
		ent.UpdateData(0)
		if help { // actions list
			fmt.Println()
			fmt.Println("---- Home Page ----")
			fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
			fmt.Println("1. Online Shops")
			fmt.Println("2. My Shop")
			fmt.Println("3. Balance")
			fmt.Println("4. My Ratings")
			fmt.Println("5. My Logs")
			fmt.Println("6. Settings")
			fmt.Println("0. Logout")
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // logout
			fmt.Println("Logging out.\n ")
			return
		case "1": // browse menu
			browseMenu(uID)
			help = true
		case "2": // my shop menu
			myShopMenu(uID)
			help = true
		case "3": //
			balanceMenu(uID)
			help = true
		case "4": //
			myRatingsMenu(uID)
			help = true
		case "5": //
			myLogsMenu(uID)
			help = true
		case "6": //
			settingsMenu(uID)
			help = true
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}

// ---++===+ Browse Page +===++---
func browseMenu(uID string) {
	help := true
	for {
		if help { // actions list
			fmt.Println("\n---- Online Shops ----")
			fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
			fmt.Println("1. Search")
			fmt.Println("2. Show Sellers")
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			return
		case "1": // buy item(s)
			for temp := true; temp; {
				temp = browse.SearchShop(uID)
			}
			ent.WaitEnter()
			help = true
		case "2": // show sellers
			browse.ShowSellers(uID)
			ent.WaitEnter()
			help = true
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}

// ---++===+ My Shop Page +===++---
func myShopMenu(uID string) {
	help := true
	for {
		shopStt := ent.GetShopState(uID)
		if help { // actions list
			if shopStt.Num == 0 {
				fmt.Println("\n---- Shop Page ----")
				fmt.Println("(shop isn't made yet)")
				fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
				fmt.Println("1. Make Shop")
				fmt.Println("0. Return")
			} else {
				fmt.Printf("\n---- %s's Shop Page %s ----\n", ent.Users[uID].Shop.ShopName, shopStt.Str)
				fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
				fmt.Println("1. Change State")
				fmt.Println("2. Show Item List")
				fmt.Println("3. Add New Item")
				fmt.Println("4. Edit Item")
				fmt.Println("5. Stock Item")
				fmt.Println("6. Remove Item")
				fmt.Println("7. Change Shop Name")
				fmt.Println("0. Return")
			}
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			return
		case "1": // make shop/change shop state
			if shopStt.Num != 0 {
				ent.ChangeShopState(uID)
				fmt.Println("State changed.")
			} else {
				myShop.MakeShop(uID)
			}
			ent.WaitEnter()
			help = true
		case "2": // show item list
			if shopStt.Num != 0 {
				myShop.ShowItemList(uID)
				ent.WaitEnter()
				help = true
			} else {
				fmt.Println("Error: Unknown action.")
			}
		case "3": // add a new item
			if shopStt.Num != 0 {
				for temp := true; temp; {
					temp = myShop.AddItem(uID)
				}
				ent.WaitEnter()
				help = true
			} else {
				fmt.Println("Error: Unknown action.")
			}
		case "4": // edit item data
			if shopStt.Num != 0 {
				for temp := true; temp; {
					temp = myShop.EditItem(uID)
				}
				ent.WaitEnter()
				help = true
			} else {
				fmt.Println("Error: Unknown action.")
			}
		case "5": // stock items
			if shopStt.Num != 0 {
				for temp := true; temp; {
					temp = myShop.StockItem(uID)
				}
				ent.WaitEnter()
				help = true
			} else {
				fmt.Println("Error: Unknown action.")
			}
		case "6": // remove an item
			if shopStt.Num != 0 {
				myShop.RemoveItem(uID)
				ent.WaitEnter()
				help = true
			} else {
				fmt.Println("Error: Unknown action.")
			}
		case "7": // change shop name
			if shopStt.Num != 0 {
				myShop.ChangeShopName(uID)
				ent.WaitEnter()
				help = true
			} else {
				fmt.Println("Error: Unknown action.")
			}
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}

func balanceMenu(uID string) {
	help := true
	for {
		if help { // actions list
			fmt.Println("\n---- Balance ----")
			fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
			fmt.Printf("Your current balance: %d\n", ent.Users[uID].Money)
			fmt.Println("1. Add Balance")
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			return
		case "1": // add balance
			balance.AddBalance(uID)
			ent.WaitEnter()
			help = true
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}

func myRatingsMenu(uID string) {
	help := true
	for {
		if help { // actions list
			fmt.Println("\n---- My Ratings ----")
			fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
			fmt.Println("1. View User Ratings")
			fmt.Println("2. Change Rating")
			if ent.Users[uID].Shop.State != 0 {
				fmt.Printf("3. View Shop Ratings (%s)\n", ent.Users[uID].Shop.ShopName)
			}
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			return
		case "1": // view user ratings
			myRatings.ViewUserRts(uID)
			ent.WaitEnter()
			help = true
		case "2": // change rating
			for temp := true; temp; {
				temp = myRatings.ChangeRating(uID)
			}
			ent.WaitEnter()
			help = true
		case "3": // view shop ratings
			if ent.Users[uID].Shop.State != 0 {
				myRatings.ViewShopRts(uID)
				ent.WaitEnter()
				help = true
			} else {
				fmt.Println("Error: Unknown action.")
			}
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}

func myLogsMenu(uID string) {
	help := true
	for {
		if help { // actions list
			fmt.Println("\n---- My Logs ----")
			fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
			fmt.Println("1. View Logs")
			fmt.Println("2. Get Logs")
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			return
		case "1": // view logs
			myLogs.ViewLogs(uID)
			ent.WaitEnter()
			help = true
		case "2": // get logs
			myLogs.GetLogs(uID)
			ent.WaitEnter()
			help = true
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}

func settingsMenu(uID string) {
	help := true
	for {
		if help { // actions list
			fmt.Println("\n---- Settings ----")
			fmt.Printf("Logged in as %s.\n", ent.Users[uID].Username)
			fmt.Println("1. Change Name")
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			return
		case "1": // change name
			settings.ChangeUserName(uID)
			ent.WaitEnter()
			help = true
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}
