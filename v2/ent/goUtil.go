package ent

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func ToLowerTxt(str string) string {
	return strings.ToLower(str)
}

func StrConvIToA(num int) string {
	return strconv.Itoa(num)
}

// get input text from question, can check for empty string input
// returns inputted text & if input is a success
func InputText(question string, allowEmpty bool) (string, bool) {
	// asks for input
	txtInput := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	txt, err := txtInput.ReadString('\n')
	txt = strings.TrimSpace(txt)

	// error check
	if IsError(err) {
		return "", false
	}
	if !allowEmpty && txt == "" {
		fmt.Println("Error: empty text")
		return "", false
	}

	return txt, true // returns inputted text & if input success
}

// get input number from question, can check for negative value
// returns inputted number & if input is a success
func InputNum(question string, allowNeg bool) (int, bool) {
	// asks for input
	var numInput float64
	fmt.Print(question)
	_, err := fmt.Scanln(&numInput)

	// error check
	if IsError(err) {
		return 0, false
	}
	if math.Floor(numInput) != numInput {
		fmt.Println("Error: value cannot be a decimal number")
		return 0, false
	}
	if !allowNeg && numInput < 0 {
		fmt.Println("Error: value cannot be negative")
		return 0, false
	}

	return int(numInput), true // returns inputted number & if input success
}

// waits for enter input
func WaitEnter() {
	fmt.Println("Press Enter to proceed.")
	fmt.Scanln()
}

// redos action from called function if answers yes
func RedoAction(question string) bool {
	var yesAns = []string{"1", "y", "yes", "ya", "iya"}
	addAgain, _ := InputText(question, true)
	for i := 0; i < len(yesAns); i++ {
		if strings.ToLower(addAgain) == yesAns[i] {
			return true
		}
	}
	return false
}
