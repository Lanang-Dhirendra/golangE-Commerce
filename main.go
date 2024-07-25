package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type users struct {
	Name string
	Shop []itemData
}
type itemData struct {
	Name   string
	Desc   string
	Amount uint
}

// e-commerce user data
var userData []users = []users{
	{
		Name: "guest",
		Shop: []itemData{
			{Name: "mouse", Desc: "computer component", Amount: 2},
			{Name: "mice", Desc: "the living rodent", Amount: 0},
		},
	},
	{
		Name: "lanang",
		Shop: []itemData{},
	},
}

// alternate yes answer
var yesAns []string = []string{"1", "y", "yes", "ya", "iya"}

func inputText(question string, allowEmpty bool) (string, bool) {
	// asks for input
	txtInput := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	txt, err := txtInput.ReadString('\n')
	txt = strings.TrimSpace(txt)

	// error check
	if err != nil {
		fmt.Println("Error: ", err)
		return "", false
	}
	if !allowEmpty && txt == "" {
		fmt.Println("Error: Empty text.")
		return "", false
	}

	return txt, true // returns inputted text & if input success
}
func inputNum(question string, allowNeg bool) (int, bool) {
	// asks for input
	var numInput float64
	fmt.Print(question)
	_, err := fmt.Scanln(&numInput)

	// error check
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, false
	}
	if math.Floor(numInput) != numInput {
		fmt.Println("Error: Value cannot be a decimal number.")
		return 0, false
	}
	if !allowNeg && numInput < 0 {
		fmt.Println("Error: Value cannot be negative.")
		return 0, false
	}

	return int(numInput), true // returns inputted number & if input success
}
func waitEnter() { // waits for enter input
	fmt.Print("Press Enter to proceed.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
func findUser(user string) (bool, int) {
	// finding user & saves user index
	userIndex := 0
	for userIndex < len(userData) {
		tempName := userData[userIndex].Name
		if user == tempName {
			break
		}
		userIndex++
	}

	if userIndex == len(userData) { // failsafe if the user is not found
		return false, -1
	}

	return true, userIndex // returns if the user is found & the user index
}

func showShop(user string, showType uint8) (int, bool) {
	// title & check if shop is empty
	fmt.Printf("-- %s's Item List --\n", user)
	scs, userIndex := findUser(user)
	if !scs {
		fmt.Println("Error: User not found.")
		return -1, true
	}
	itemDatas := userData[userIndex].Shop
	if len(itemDatas) == 0 {
		fmt.Println("No item on list.")
		return userIndex, true
	}

	// determine settings for output
	for j := 0; j < len(itemDatas); j++ { // output
		itemName := itemDatas[j].Name
		itemDesc := itemDatas[j].Desc
		itemAmount := itemDatas[j].Amount
		if itemDesc == "" {
			itemDesc = "- (no description)"
		}

		switch showType { // output
		case 1: // w/ desc, w/ amount, no stock label
			fmt.Printf("%d. %s (Amount: %d)\n%s\n", j+1, itemName, itemAmount, itemDesc)
			break
		case 2: // w/ desc, w/ amount, w/ stock label
			if itemAmount == 0 {
				fmt.Printf("%d. %s (Out of stock)\n%s\n", j+1, itemName, itemDesc)
			} else {
				fmt.Printf("%d. %s (Amount: %d)\n%s\n", j+1, itemName, itemAmount, itemDesc)
			}
			break
		case 3: // no desc, no amount
			fmt.Printf("%d. %s\n", j+1, itemName)
			break
		default:
			fmt.Println("Error: Unknown show shop print type.")
			break
		}
	}

	return userIndex, false // returns index of user & if the shop is empty
}
func redoAction(question string) bool { // redos action from called function if answers yes
	addAgain, _ := inputText(question, true)
	for i := 0; i < len(yesAns); i++ {
		if strings.ToLower(addAgain) == yesAns[i] {
			return true
		}
	}
	return false
}

// ---++===+ Start Page +===++---
func main() {
	var loop, help, title bool = true, true, true
	for loop {
		if title { // program title
			fmt.Println()
			fmt.Println("===============================")
			fmt.Println("== Simple E-Commerce Program ==")
			fmt.Println("===============================")
			fmt.Println("   ++  Made by IAmAGuest  ++   ")
			title = false
		}

		if help { // actions list
			fmt.Println()
			fmt.Println("---- Start Page ----")
			fmt.Println("1. Login")
			fmt.Println("2. Register")
			fmt.Println("0. Exit")
			help = false
		}

		// input action
		action, _ := inputText("\nAction: ", true)
		switch action { // actions
		case "0": // exit program
			fmt.Println("Exiting program.\n ")
			loop = false
			os.Exit(0)
			break
		case "1": // login
			if check, name := p0a_Login(); check {
				p1_Home(name)
			}
			help = true
			break
		case "2": // register
			p0b_Regis()
			waitEnter()
			help = true
			break
		default: // failsafe
			fmt.Println("Error: Unknown action.")
			break
		}
	}

}
func p0a_Login() (bool, string) { // asks name, logs in w/ that name if it exists
	fmt.Println(" == Login")
	name, scs := inputText("(case-sensitive; 0 to cancel)\nEnter name: ", false)
	if scs {
		if userFound, _ := findUser(name); userFound {
			fmt.Printf("Login successful as %s.\n", name)
			return true, name
		}
		if name == "0" {
			fmt.Println("Login cancelled.")
			return false, ""
		}
		fmt.Println("Login failed, unregistered name.")
		waitEnter()
	}
	return false, ""
}
func p0b_Regis() bool { // asks name, makes an account if it doesn't exist
	fmt.Println(" == Register")
	name, scs := inputText("(case-sensitive; 0 to cancel)\nEnter name: ", false)
	if scs {
		if userFound, _ := findUser(name); !userFound {
			if name == "0" {
				fmt.Println("Register cancelled.")
				return false
			}
			userData = append(userData, []users{{Name: name, Shop: []itemData{}}}...)
			fmt.Println("Register successful.")
			return true
		}
		fmt.Println("Register failed, name already exists.")
	}
	return false
}

// ---++===+ Home Page +===++---
func p1_Home(user string) {
	var loop, help bool = true, true
	for loop {
		if help { // actions list
			fmt.Println()
			fmt.Println("---- Home Page ----")
			fmt.Printf("Logged in as %s.\n", user)
			fmt.Println("1. Browse")
			fmt.Println("2. My Shop")
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := inputText("\nAction: ", true)
		switch action { // actions
		case "0": // logout
			fmt.Println("Logging out.\n ")
			loop = false
			break
		case "1": // browse menu
			p2_Browse(user)
			help = true
			break
		case "2": // my shop menu
			p3_MyShop(user)
			help = true
			break
		default: // failsafe
			fmt.Println("Error: Unknown action.")
			break
		}
	}
}

// ---++===+ Browse Page +===++---
func p2_Browse(user string) {
	fmt.Println("Browse.")
	var loop, help bool = true, true
	for loop {
		if help { // actions list
			fmt.Println("\n---- Browse ----")
			fmt.Printf("Logged in as %s.\n", user)
			fmt.Println("1. Buy Item(s)")
			fmt.Println("2. Show Sellers")
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := inputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			loop = false
			break
		case "1": // buy item(s)
			for temp := true; temp; {
				temp = p2a_BuyItems(user)
			}
			waitEnter()
			help = true
			break
		case "2": // show sellers
			p2b_ShowSellers(user)
			waitEnter()
			help = true
			break
		default: // failsafe
			fmt.Println("Error: Unknown action.")
			break
		}
	}
}
func p2a_BuyItems(user string) bool {
	fmt.Println(" == Buy Items")

	// find all sellers that isn't the current user
	sellers, userSells, itemDatas := []users{}, false, []itemData{}
	for i := 0; i < len(userData); i++ {
		tempUser := userData[i]
		if tempUser.Name == user {
			userSells = true
			continue
		}
		itemDatas = tempUser.Shop
		if len(itemDatas) == 0 {
			continue
		}
		sellers = append(sellers, tempUser)
		fmt.Printf("%d. %s", len(sellers), tempUser.Name)
		fmt.Println()
	}
	if len(sellers) == 0 {
		fmt.Print("There are currently no sellers")
		if userSells {
			fmt.Print(" (that isn't you)")
		}
		fmt.Println(".")
		return false
	}

	// input seller
	var sellerIndex int
	for scs := false; !scs; {
		// input
		sellerIndex, scs = inputNum("Enter seller index (0 to cancel): ", true)

		// incorrect input check
		if sellerIndex < 0 || sellerIndex > len(sellers) {
			fmt.Println("Error: Seller index out of reach.")
			scs = false
			continue
		}
		sellerIndex -= 1
	}
	if sellerIndex == -1 {
		fmt.Println("Action cancelled.")
		return false
	}
	fmt.Println()

	// loop for buying multiple items from same seller
	for buyLoop := true; buyLoop; {
		showShop(sellers[sellerIndex].Name, 2)
		itemDatas = sellers[sellerIndex].Shop
		itemIndex := 0

		// get item index
		for scs := false; !scs; {
			itemIndex, scs = inputNum("Item index (0 to cancel): ", true)
			if itemIndex < 0 || itemIndex > len(itemDatas) {
				fmt.Println("Error: Item index out of reach.")
				scs = false
				continue
			}
			itemIndex -= 1
		}
		if itemIndex == -1 {
			fmt.Println("Action cancelled.")
			return redoAction("Buy again? (1 for yes): ")
		}
		if itemDatas[itemIndex].Amount == 0 {
		}

		// output item data
		itemName := itemDatas[itemIndex].Name
		itemDesc := itemDatas[itemIndex].Desc
		itemAmount := itemDatas[itemIndex].Amount
		if itemDesc == "" {
			itemDesc = "- (no desc)"
		}
		fmt.Printf("Item selected: %s\n", itemName)
		itemBought := 0
		for inputScs := false; !inputScs; {
			itemBought, inputScs = inputNum("Item amount (0 to cancel): ", false)
			if int(itemAmount)-itemBought < 0 {
				fmt.Println("Insufficient stock.")
				inputScs = false
			}
		}
		if itemBought == 0 {
			fmt.Println("Action cancelled.")
		} else {
			userData[sellerIndex].Shop[itemIndex].Amount -= uint(itemBought)
			fmt.Println("Item successfully bought.")
		}

		if !redoAction("Buy again from this seller? (1 for yes): ") {
			buyLoop = false
		}
	}

	return redoAction("Buy again? (1 for yes): ")
}
func p2b_ShowSellers(user string) {
	fmt.Println(" == Show Sellers")

	sellers := 0
	for i := 0; i < len(userData); i++ {
		tempUser := userData[i]
		itemDatas := tempUser.Shop
		if len(itemDatas) == 0 {
			continue
		}
		sellers++
		fmt.Printf("%d. %s", sellers, tempUser.Name)
		if tempUser.Name == user {
			fmt.Print(" (You)")
		}
		fmt.Println()
	}
	if sellers == 0 {
		fmt.Println("There are currently no sellers.")
	}
}

// ---++===+ My Shop Page +===++---
func p3_MyShop(user string) {
	var loop, help bool = true, true
	for loop {
		if help { // actions list
			fmt.Printf("\n---- %s's Shop Page ----\n", user)
			fmt.Printf("Logged in as %s.\n", user)
			fmt.Println("1. Show Item List")
			fmt.Println("2. Add New Item")
			fmt.Println("3. Edit Item")
			fmt.Println("4. Remove Item")
			fmt.Println("0. Return")
			help = false
		}

		// input action
		action, _ := inputText("\nAction: ", true)
		switch action { // actions
		case "0": // return
			fmt.Println("Returning to home page.\n ")
			loop = false
			break
		case "1": // show item list
			p3a_ShowItemList(user)
			waitEnter()
			help = true
			break
		case "2": // add a new item
			for temp := true; temp; {
				temp = p3b_AddItem(user)
			}
			waitEnter()
			help = true
			break
		case "3": // edit item data
			for temp := true; temp; {
				temp = p3c_EditItem(user)
			}
			waitEnter()
			help = true
			break
		case "4": // remove an item
			p3d_RemoveItem(user)
			waitEnter()
			help = true
			break
		default: // failsafe
			fmt.Println("Error: Unknown action.")
			break
		}
	}
}
func p3a_ShowItemList(user string) { // self explanatory
	fmt.Println(" == Show Item List")
	showShop(user, 1)
}
func p3b_AddItem(user string) bool {
	fmt.Println("\n == Add New Item")
	fmt.Println("(type '-1' at any point to cancel adding item)")

	// find index for current user & failsafe if not found
	scs, userIndex := findUser(user)
	if !scs {
		fmt.Println("Error: User not found.")
		return false
	}

	// input item data
	var itemName, itemDesc string
	var itemAmount int
	itemDatas := userData[userIndex].Shop
	for inputScs := false; !inputScs; { // input item name
		itemName, inputScs = inputText("Item name: ", false)
		for i := 0; i < len(itemDatas); i++ {
			if itemName == itemDatas[i].Name {
				fmt.Println("Error: Item already exists.")
				inputScs = false
				break
			}
		}
	}
	if itemName == "-1" {
		fmt.Println("Cancelled adding item.")
		return false
	}

	// input item description
	itemDesc, _ = inputText("Item description (optional): ", true)
	if itemDesc == "-1" {
		fmt.Println("Cancelled adding item.")
		return false
	}

	for inputScs := false; !inputScs; { // input item amount
		itemAmount, inputScs = inputNum("Item amount: ", true)
		if itemAmount == -1 {
			fmt.Println("Cancelled adding item.")
			return false
		}
		if itemAmount < 0 {
			fmt.Println("Error: Value cannot be negative (except -1).")
			inputScs = false
		}
	}

	// append item to slice
	userData[userIndex].Shop = append(
		userData[userIndex].Shop,
		itemData{Name: itemName, Desc: itemDesc, Amount: uint(itemAmount)},
	)

	// output & ask to re-add new item(s)
	fmt.Printf("Item %s successfully added.\n", itemName)
	return redoAction("Add again? (1 for yes): ")
}
func p3c_EditItem(user string) bool {
	fmt.Println("\n == Edit Item")

	// get user index & failsafe if shop is empty
	userIndex, emptyShop := showShop(user, 3)
	if emptyShop {
		return false
	}

	// asks for item to edit
	itemDatas := userData[userIndex].Shop
	itemIndex := 0
	for scs := false; !scs; {
		itemIndex, scs = inputNum("Item index (0 to cancel): ", true)
		if itemIndex < 0 || itemIndex > len(itemDatas) {
			fmt.Println("Error: Item index out of reach.")
			scs = false
			continue
		}
		itemIndex -= 1
	}
	if itemIndex == -1 {
		fmt.Println("Cancelled editing item.")
		return false
	}

	// outputs item data
	itemName := itemDatas[itemIndex].Name
	itemDesc := itemDatas[itemIndex].Desc
	itemAmount := itemDatas[itemIndex].Amount
	if itemDesc == "" {
		itemDesc = "- (no desc)"
	}
	fmt.Printf("\nItem selected:\n%s (Amount: %d)\n%s\n", itemName, itemAmount, itemDesc)

	// asks for item data to edit
	var element string
	for scs := false; !scs; {
		element, _ = inputText("Select item element to edit.\n1 -> Name\n2 -> Desc\n3 -> Amount\n0 -> Cancel\nInput: ", true)
		for _, c := range []string{"0", "1", "2", "3"} {
			if element == c {
				scs = true
				break
			}
		}
		if !scs {
			fmt.Println("Error: Element index out of reach.")
		}
	}

	// edit process
	switch element {
	case "0": // cancel
		fmt.Println("Cancelled editing item.")
		return false
	case "1": // edit item name
		itemNameNew := ""
		for inputScs := false; !inputScs; {
			itemNameNew, inputScs = inputText("New item name (0 to cancel): ", false)
			if !inputScs {
				break
			}
			for i := 0; i < len(itemDatas); i++ {
				if itemNameNew == itemDatas[i].Name {
					fmt.Println("Error: Item already exists.")
					inputScs = false
					break
				}
			}
		}
		if itemNameNew == "0" {
			fmt.Println("Cancelled editing item.")
			break
		}
		userData[userIndex].Shop[itemIndex].Name = itemNameNew
		fmt.Println("Item name successfully changed.")
		break
	case "2": // edit item description
		itemDescNew, _ := inputText("New item description (0 to cancel): ", false)
		if itemDescNew == "0" {
			fmt.Println("Cancelled editing item.")
			break
		}
		userData[userIndex].Shop[itemIndex].Desc = itemDescNew
		fmt.Println("Item description successfully changed.")
		break
	case "3": // edit item amount
		itemAmountNew := 0
		for inputScs := false; !inputScs; {
			itemAmountNew, inputScs = inputNum("Item amount (-1 to cancel): ", true)
			if itemAmountNew == -1 {
				fmt.Println("Cancelled adding item.")
				return false
			}
			if itemAmountNew < 0 {
				fmt.Println("Error: Value cannot be negative (except -1).")
				inputScs = false
			}
		}
		userData[userIndex].Shop[itemIndex].Amount = uint(itemAmountNew)
		fmt.Println("Item amount successfully changed.")
		break
	default: // failsafe
		break
	}
	return redoAction("Edit again? (1 for yes): ")
}
func p3d_RemoveItem(user string) {
	fmt.Println(" == Remove Item")

	// get user index & failsafe if shop is empty
	userIndex, emptyShop := showShop(user, 3)
	if emptyShop {
		return
	}

	// asks for item to remove & failsafe if index is out of range
	itemDatas := userData[userIndex].Shop
	itemIndex := 0
	for scs := false; !scs; {
		itemIndex, scs = inputNum("Item index (0 to cancel): ", true)
		if itemIndex < 0 || itemIndex > len(itemDatas) {
			fmt.Println("Error: Item index out of reach.")
			scs = false
			continue
		}
		itemIndex -= 1
	}
	if itemIndex == -1 {
		fmt.Println("Cancelled removing item.")
		return
	}

	// remove prompt
	inputTxt, _ := inputText("Type 'Remove.' exactly to remove item: ", true)
	if inputTxt == "Remove." {
		var userNewShop []itemData
		for i := 0; i < len(itemDatas); i++ {
			if i == itemIndex {
				i++
			}
			if i < len(itemDatas) {
				userNewShop = append(userNewShop, itemDatas[i])
			}
		}
		userData[userIndex].Shop = userNewShop
		fmt.Println("Item removed successfully.")
		return
	} else {
		fmt.Println("Remove failed, typed incorrectly.")
		return
	}
}
