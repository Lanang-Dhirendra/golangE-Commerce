package main

import (
	"ecommerce/ent"
	"ecommerce/useCase/home"
	"ecommerce/useCase/start"
	"fmt"
	"os"
)

// ---++===+ Start Page +===++---
func main() {
	help, title := true, true
	for {
		if title { // program title
			fmt.Println()
			fmt.Println("===============================")
			fmt.Println("===  E-Commerce Program v2  ===")
			fmt.Println("===============================")
			fmt.Println("  ++   Made by IAmAGuest   ++  ")
			fmt.Println()
			title = false
		}

		if help { // actions list
			fmt.Println("---- Start Page ----")
			fmt.Println("1. Login")
			fmt.Println("2. Register")
			fmt.Println("0. Exit")
			help = false
		}

		// input action
		action, _ := ent.InputText("\nAction: ", true)
		switch action { // actions
		case "0": // exit program
			fmt.Println("Exiting program.")
			ent.WaitEnter()
			os.Exit(0)
		case "1": // login
			if check, uID := start.Login(); check {
				home.Home(uID)
			}
			help = true
		case "2": // register
			start.Regis()
			ent.WaitEnter()
			help = true
		default: // failsafe
			fmt.Println("Error: Unknown action.")
		}
	}
}
