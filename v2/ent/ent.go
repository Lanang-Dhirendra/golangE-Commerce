package ent

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type UserData struct {
	// ID -> @U.... (@U + 4 rand chars; 14m perms)

	Username string   `json:"nam"`
	Money    int      `json:"mon"`
	Logs     []string `json:"lgs"`
	Ratings  []string `json:"rts"` // score -> n; 0 >= n >= 100
	Shop     ShopData `json:"sop"` // no id -> shop isn't setupped; id -> shop setupped
}
type LogData struct {
	// ID -> @L........ (@L + 8 rand chars; 218t perms)

	Owner string `json:"own"`
	Time  int64  `json:"tym"`
	State int    `json:"stt"`
	Data  string `json:"dat"`
}
type RatingData struct {
	// ID -> @R........ (R + 8 rand chars; 218t perms)

	UserID string `json:"uid"`
	ShopID string `json:"sid"`
	Value  uint8  `json:"val"`
	//Review string `json:"rvw"`
}
type ShopData struct {
	// ID is tied to user's ID

	ShopName string              `json:"nam"`
	State    uint8               `json:"stt"` // 0 -> not created; 1 -> closed; 2 -> open
	Ratings  []string            `json:"rts"`
	Items    map[string]ItemData `json:"tms"`
}
type ItemData struct {
	// ID -> @I....-... (item id is shop id + 3 chars; 238k perms/shop)

	Name   string `json:"nam"`
	Desc   string `json:"dsc"`
	Cost   uint   `json:"cst"`
	Amount uint   `json:"amo"`
}
type ShopState struct {
	Str string
	Num int
}

// symbols for IDs
var syms string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// banned characters when naming account/item/desc
var BannedChars string = "@[]{}\\\""

// users, logs & ratings data

var Users map[string]UserData
var Logs map[string]LogData
var Ratings map[string]RatingData

// file paths

var goFilePath string
var usrDatPath string
var lgsDatPath string
var rtingsPath string

func init() {
	var err error
	goFilePath, err = os.Getwd()
	if err != nil {
		log.Fatal("can't read root dir, err:", err)
	}
	usrDatPath = filepath.Join(goFilePath, "json", "userData.json")
	rtingsPath = filepath.Join(goFilePath, "json", "ratings.json")
	lgsDatPath = filepath.Join(goFilePath, "json", "logs.json")
}

// internal entity functions

// check for error, returns if given data is an error
func IsError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
		return true
	}
	return false
}

func catchPanicJsonRead() { // (recover catch isn't perfect yet)
	if r := recover(); r != nil {
		filePath := ""
		switch r {
		case "userData":
			filePath = usrDatPath
		case "logData":
			filePath = lgsDatPath
		case "ratingData":
			filePath = rtingsPath
		default:
			failRecover(r)
		}
		file, test2 := os.Stat(filePath)
		if IsError(test2) {
			if os.IsNotExist(test2) {
				createFile(filePath)
			}
			file, test2 = os.Stat(filePath)
			if IsError(test2) {
				failRecover(r)
			}
		}
		if test1 := file.Size(); test1 == 0 {
			writeFile(filePath, []byte("{}"), 0)
		}
		if test3 := UpdateData(1); !test3 {
			failRecover(r)
		}
	}
}
func failRecover(err any) {
	fmt.Println("\033[31mUNFIXABLE ERROR, PLEASE CHECK AGAIN. ERROR:", err, "\033[0m")
	os.Exit(0)
}
func createFile(filePath string) {
	file, err := os.Create(filePath)
	if IsError(err) {
		return
	}
	defer file.Close()
}

func readFile(filePath string) []byte {
	var file, err = os.OpenFile(filePath, os.O_RDONLY, 0644)
	if IsError(err) {
		return []byte{}
	}
	defer file.Close()

	var text = make([]byte, 65536)
	var fileLen int
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if IsError(err) {
				return []byte{}
			}
		}
		if n == 0 {
			break
		}
		fileLen = n
	}
	if fileLen == 0 {
		panic("fwtvvyuybdhijn")
	}
	return text[:fileLen]
}

func writeFile(filePath string, addition []byte, at int) {
	var file, err = os.OpenFile(filePath, os.O_RDWR, 0644)
	if IsError(err) {
		return
	}
	defer file.Close()

	_, err = file.WriteAt(addition, int64(at))
	if IsError(err) {
		return
	}
	err = file.Sync()
	if IsError(err) {
		return
	}
}

func clearFile(filePath string) {
	var file, err = os.OpenFile(filePath, os.O_RDWR, 0644)
	if IsError(err) {
		return
	}
	defer file.Close()
	file.Truncate(0)
	err = file.Sync()
	if IsError(err) {
		return
	}
}

func findByte(bytes []byte, from int, bit ...byte) int {
	currBit := 0
	for i, b := range bytes {
		if i < from {
			continue
		}
		if b != bit[currBit] {
			continue
		}
		currBit++
		if currBit == len(bit) {
			return i
		}
	}
	return -1
}

func findString(bytes []byte, from int, str string) int {
	searchState := 0
	for i, b := range bytes {
		if i < from {
			continue
		}
		if searchState == 0 {
			if b != str[0] {
				continue
			}
			searchState++
		} else {
			if searchState == len(str) {
				return i - searchState
			}
			if b != str[searchState] {
				searchState = 0
				continue
			}
			searchState++
		}
	}
	return -1
}

func findBrackets(bytes []byte, from int, pair ...[2]byte) int {
	pairing, currPair, foundOpen := 0, 0, false
	for i, b := range bytes {
		if i < from {
			continue
		}
		if b != pair[currPair][0] && b != pair[currPair][1] {
			continue
		}
		if b == pair[currPair][0] {
			pairing += 1
			foundOpen = true
		} else {
			pairing -= 1
		}
		if foundOpen && pairing == 0 {
			currPair++
			foundOpen = false
		}
		if currPair == len(pair) {
			return i
		}
	}
	return -1
}

func writeLog(logID, logOwner, logData string, logTime, logState int) {
	jsonBytes := readFile(lgsDatPath)
	UpdateData(0)

	_, found := Logs[logID]
	if found {
		fmt.Println("Error: log entry already exists")
		return
	}

	add := fmt.Sprintf("\n\"%s\":{\"own\":\"%s\", \"tym\":%d, \"stt\":%d, \"dat\":\"%s\"}\n}", logID, logOwner, logTime, logState, logData)
	at := 0
	if len(Logs) > 0 {
		add = "," + add
		at = findBrackets(jsonBytes, at, [2]byte{'{', '}'}) - 1
	} else {
		add = "{" + add
		clearFile(lgsDatPath)
	}
	writeFile(lgsDatPath, []byte(add), at)
}
func writeRating(ratingID, user, shop string, ratingVal int) {
	jsonBytes := readFile(rtingsPath)
	UpdateData(0)

	_, found := Ratings[ratingID]
	var addition []byte
	at, tempStr := 0, ""
	tempStr = fmt.Sprintf("\"%s\":{\"uid\":\"%s\", \"sid\":\"%s\", \"val\":%d}", ratingID, user, shop, ratingVal)
	if found {
		at = findString(jsonBytes, at, ratingID)
		tempStr = string(jsonBytes[:at-1]) + tempStr
		at = findByte(jsonBytes, at, '}') + 1
		tempStr = tempStr + string(jsonBytes[at:])
		at = 0
		clearFile(rtingsPath)
	} else {
		tempStr = "\n" + tempStr + "\n}"
		at = findBrackets(jsonBytes, at, [2]byte{'{', '}'})
		if len(Ratings) > 0 {
			tempStr = "," + tempStr
			at -= 1
		}
	}
	addition = []byte(tempStr)
	writeFile(rtingsPath, addition, at)
}
func WriteUserData(uID, newData string, editData int) {
	jsonBytes := readFile(usrDatPath)
	UpdateData(0)

	if _, found := Users[uID]; !found {
		fmt.Println("Error: user w/ id given not found")
		return
	}

	editDataStr := ""
	switch editData {
	case 1: // name
		editDataStr = "\"nam\""
		newData = "\"" + newData + "\""
	case 2: // money
		editDataStr = "\"mon\""
	case 3: // logs
		editDataStr = "\"lgs\""
		newData = "\"" + newData + "\"]"
		if len(Users[uID].Logs) > 0 {
			newData = "," + newData
		}
	case 4: // ratings
		editDataStr = "\"rts\""
		newData = "\"" + newData + "\"]"
		if len(Users[uID].Ratings) > 0 {
			newData = "," + newData
		}
	default: // failsafe
		fmt.Println("Error: unknown data to edit")
		return
	}

	at := 0
	at = findString(jsonBytes, at, uID)
	at = findString(jsonBytes, at, editDataStr+":")
	at += len(editDataStr) + 1
	if editData == 3 || editData == 4 {
		at = findByte(jsonBytes, at, ']')
	}
	bytesBefore := jsonBytes[0:at]
	at = findByte(jsonBytes, at, ',')
	bytesAfter := jsonBytes[at:]
	addition := []byte(strings.Join([]string{string(bytesBefore), newData, string(bytesAfter)}, ""))

	clearFile(usrDatPath)
	writeFile(usrDatPath, addition, 0)
}

// func deleteUser(uID string) { ...add??
//
//		fmt.Println("test")
//	}

func createCode(dataType int, itemID ...string) string {
	switch dataType {
	case 1:
		return createRandChars(4, "@U")
	case 2:
		return createRandChars(8, "@L")
	case 3:
		return createRandChars(8, "@R")
	case 4:
		return createRandChars(3, "@I"+itemID[0][2:]+"-")
	}
	return ""
}
func createRandChars(length int, prefix string) string {
	runtime.GOMAXPROCS(2)
	var codeStr = make(chan string, 4)
	var validPrefix = []string{"@U", "@L", "@R"}
	UpdateData(0)

	// find code by randomly generating
	go func(strCh chan<- string) {
		for {
			var1 := ""
			for j := 0; j < length; j++ {
				var1 += string(syms[rand.Intn(len(syms))])
			}
			strCh <- var1
		}
	}(codeStr)

	// find code alphabetically, only starts after 10s pause
	if prefix != validPrefix[0] || prefix != validPrefix[1] || prefix != validPrefix[2] {
		go func(strCh chan<- string) {
			time.Sleep(time.Second * 10)
			for i := 0; true; i++ {
				var1 := ""
				for j := 0; j < length; j++ {
					var2 := int(i / int(math.Pow(float64(len(syms)), float64(j))))
					var1 = string(syms[var2%len(syms)]) + var1
				}
				strCh <- var1
			}
		}(codeStr)
	}

	// retrieve code not available already
	return func(prefix string, strCh <-chan string) string {
		tempCodeStr := ""
		for {
			tempCodeStr = <-strCh
			tempCodeStr = prefix + tempCodeStr
			found := false
			switch prefix {
			case validPrefix[0]:
				_, found = Users[tempCodeStr]
			case validPrefix[1]:
				_, found = Logs[tempCodeStr]
			case validPrefix[2]:
				_, found = Ratings[tempCodeStr]
			}
			if found {
				tempCodeStr = ""
			} else {
				break
			}
		}
		return tempCodeStr
	}(prefix, codeStr)
}

// check for banned chars in given string
func CheckBannedChars(str string) (bool, string) {
	banned := ""
	for _, c := range str {
		for _, b := range BannedChars {
			if c == b {
				banned += string(c)
				break
			}
		}
	}
	return len(banned) > 0, banned
}

func WriteShopData(uID, itemID, newData string, editData int) {
	jsonBytes := readFile(usrDatPath)
	UpdateData(0)

	user, found := Users[uID]
	if !found {
		fmt.Println("Error: user w/ id given not found")
		return
	}
	if editData > 3 {
		found2 := false
		for k := range user.Shop.Items {
			if k == itemID {
				found2 = true
				break
			}
		}
		if !found2 {
			fmt.Println("Error: item w/ id given not found")
			return
		}
	}

	editDataStr := ""
	switch editData {
	case 1: // name
		editDataStr = "\"nam\""
		newData = "\"" + newData + "\""
	case 2: // state
		editDataStr = "\"stt\""
	case 3: // ratingID
		editDataStr = "\"rts\""
		newData = "\"" + newData + "\"]"
		if len(user.Shop.Ratings) > 0 {
			newData = "," + newData
		}
	case 4: // itemName
		editDataStr = "\"nam\""
		newData = "\"" + newData + "\""
	case 5: // itemDesc
		editDataStr = "\"dsc\""
		newData = "\"" + newData + "\""
	case 6: // itemCost
		editDataStr = "\"cst\""
	case 7: // itemAmount
		editDataStr = "\"amo\""
		newData += "\n\t\t"
	default: // failsafe
		fmt.Println("Error: unknown data to edit")
		return
	}

	at := 0
	at = findString(jsonBytes, at, uID)
	at = findString(jsonBytes, at, "\"sop\"")
	if editData > 3 {
		at = findString(jsonBytes, at, itemID)
	}
	at = findString(jsonBytes, at, editDataStr+":")
	at += len(editDataStr) + 1
	if editData == 3 {
		at = findByte(jsonBytes, at, ']')
	}
	bytesBefore := jsonBytes[0:at]
	if editData != 7 {
		at = findByte(jsonBytes, at, ',')
	} else {
		at = findByte(jsonBytes, at, '}')
	}
	bytesAfter := jsonBytes[at:]
	addition := []byte(strings.Join([]string{string(bytesBefore), newData, string(bytesAfter)}, ""))

	clearFile(usrDatPath)
	writeFile(usrDatPath, addition, 0)
}

// updates current user, log, & rating data
func UpdateData(from int) bool {
	runtime.GOMAXPROCS(3)
	var upDataScs = make(chan bool)
	Users = map[string]UserData{}
	Logs = map[string]LogData{}
	Ratings = map[string]RatingData{}
	temp := [3]bool{}
	for i := 0; i < 3; i++ {
		go func(data int, ch chan<- bool) {
			defer catchPanicJsonRead()
			var err error
			panTxt := ""
			switch data {
			case 1:
				err = json.Unmarshal(readFile(usrDatPath), &Users)
				panTxt = "userData"
			case 2:
				err = json.Unmarshal(readFile(lgsDatPath), &Logs)
				panTxt = "logData"
			case 3:
				err = json.Unmarshal(readFile(rtingsPath), &Ratings)
				panTxt = "ratingData"
			default:
				ch <- false
				return
			}
			if IsError(err) {
				if from == 0 {
					panic(panTxt)
				} else {
					ch <- false
				}
			}
			ch <- true
		}(i+1, upDataScs)
	}
	for i := 0; i < 3; i++ {
		temp[i] = <-upDataScs
	}
	for i := 0; i < 3; i++ {
		if !temp[i] {
			return false
		}
	}
	return true
}

// logType: "user regis", "user login", "make shop", "item buy", "money transfer", "rate shop", "change shop state", "add item", "edit item desc",
// "edit item x desc", "stock item", "remove item", "deposit money", "view logs", "get logs", "change user name", "change shop name"
func CreateLog(logType, logOwner string, datas ...any) {
	logCode := createCode(2)
	logTime := int(time.Now().Unix())
	logData, logState := "", 0
	switch logType {
	case "user regis":
		logData, logState = fmt.Sprintf("account %s created with name %v", logOwner, datas[0]), 1
	case "user login":
		logData, logState = fmt.Sprintf("account %s login", logOwner), 2
	case "make shop":
		logData, logState = fmt.Sprintf("account %s created shop %v", logOwner, datas[0]), 1
	case "item buy":
		logData, logState = fmt.Sprintf("account %s bought %v item %v from shop %v", datas[0], datas[1], datas[2], datas[3]), 1
	case "money transfer":
		logData, logState = fmt.Sprintf("account %s pays %v money to %v", datas[0], datas[1], datas[2]), 1
	case "rate shop":
		logData = fmt.Sprintf("account %v rates shop %v with score %v", datas[0], datas[1], datas[2])
		if datas[3] == 1 || datas[3] == 2 {
			logState, _ = datas[3].(int)
		}
	case "change shop state":
		logData, logState = fmt.Sprintf("shop %s changes shop state from %v to %v", logOwner, datas[0], datas[1]), 2
	case "add item":
		logData, logState = fmt.Sprintf("shop %s added item %v with name %v, cost %v, and amount %v", logOwner, datas[0], datas[1], datas[2], datas[3]), 1
	case "edit item desc":
		logData, logState = fmt.Sprintf("shop %s edited item %v's description", logOwner, datas[0]), 1
	case "edit item x desc":
		logData, logState = fmt.Sprintf("shop %s edited item %v's %v to %v", logOwner, datas[0], datas[1], datas[2]), 1
	case "stock item":
		logData, logState = fmt.Sprintf("shop %s stocked item %v by %v (now: %v)", logOwner, datas[0], datas[1], datas[2]), 1
	case "remove item":
		logData, logState = fmt.Sprintf("shop %s removed item %v", logOwner, datas[0]), 1
	case "change shop name":
		logData, logState = fmt.Sprintf("shop %s changes name from %v to %v", logOwner, datas[0], datas[1]), 1
	case "deposit money":
		logData, logState = fmt.Sprintf("account %s deposited %v money (now: %v)", logOwner, datas[0], datas[1]), 1
	case "view logs":
		logData, logState = fmt.Sprintf("account %s view logs", logOwner), 2
	case "get logs":
		logData, logState = fmt.Sprintf("account %s get log file", logOwner), 1
	case "change user name":
		logData, logState = fmt.Sprintf("account %s changes name from %v to %v", logOwner, datas[0], datas[1]), 1
	default:
		return
	}
	writeLog(logCode, logOwner, logData, logTime, logState)
	WriteUserData(logOwner, logCode, 3)
}

func CreateUser(name string) string {
	jsonBytes := readFile(usrDatPath)
	UpdateData(0)
	uID := createCode(1)

	add := fmt.Sprintf("\n\"%s\":{\n\"nam\":\"\",\n\"mon\":0,\n\"lgs\":[],\n\"rts\":[],\n\"sop\":{}\n}\n}", uID)
	at := 0
	if len(Users) > 0 {
		add = "," + add
		at = findBrackets(jsonBytes, 0, [2]byte{'{', '}'}) - 1
	} else {
		add = "{" + add
		clearFile(usrDatPath)
	}

	writeFile(usrDatPath, []byte(add), at)
	WriteUserData(uID, name, 1)
	CreateLog("user regis", uID, name)
	return uID
}

func AddShop(sID, shopName string) {
	jsonBytes := readFile(usrDatPath)
	UpdateData(0)

	add := "{\n\"nam\":\"\",\n\"stt\":0,\n\"rts\":[],\n\"tms\":{}\n}"
	at := findString(jsonBytes, 0, sID)
	at = findString(jsonBytes, at, "\"sop\"") + 6
	add += string(jsonBytes)[at+2:]
	writeFile(usrDatPath, []byte(add), at)

	CreateLog("make shop", sID, shopName)
	WriteShopData(sID, "", shopName, 1)
	WriteShopData(sID, "", "1", 2)
}

func AddItem(sID, name, desc string, cost, amount int) string {
	UpdateData(0)
	_, scs := Users[sID]
	if !scs {
		fmt.Println("Error: User not found.")
		return ""
	}
	jsonBytes := readFile(usrDatPath)
	itemID := createCode(4, sID)

	add := fmt.Sprintf("\n\"%s\":{\n\"nam\":\"\",\n\"dsc\":\"\",\n\"cst\":0,\n\"amo\":0\n}\n", itemID)
	at := findString(jsonBytes, 0, "\""+sID+"\"")
	at = findString(jsonBytes, at, "\"sop\"") + 7
	at = findBrackets(jsonBytes, at, [2]byte{'{', '}'}) - 1
	if len(Users[sID].Shop.Items) > 0 {
		add = "," + add
	} else {
		add = "{" + add
	}
	add += string(jsonBytes[at+1:])
	writeFile(usrDatPath, []byte(add), at)
	WriteShopData(sID, itemID, name, 4)
	WriteShopData(sID, itemID, desc, 5)
	WriteShopData(sID, itemID, strconv.Itoa(cost), 6)
	WriteShopData(sID, itemID, strconv.Itoa(amount), 7)
	return itemID
}

func RemoveItem(sID, itemID string) {
	UpdateData(0)
	_, scs := Users[sID]
	if !scs {
		fmt.Println("Error: User not found.")
		return
	}
	jsonBytes := readFile(usrDatPath)

	at1 := findString(jsonBytes, 0, "\""+itemID+"\"") - 1
	at2 := findBrackets(jsonBytes, at1, [2]byte{'{', '}'}) + 2
	remType := 0
	if jsonBytes[at2] == '}' {
		remType += 2
	}
	if jsonBytes[at1-1] == '{' {
		remType++
	}
	switch remType {
	case 1: // itemIndex: first
		at1--
		at2--
	case 2: // itemIndex: last
		at1++
		at2++
	}

	bytes := []byte(string(jsonBytes[:at1]) + string(jsonBytes[at2:]))
	clearFile(usrDatPath)
	writeFile(usrDatPath, bytes, 0)
	CreateLog("remove item", sID, itemID)
}

func DoRating(uID, sID string, ratVal int) {
	rID := func() string {
		for ratingID, val := range Ratings {
			if val.UserID == uID && val.ShopID == sID {
				return ratingID
			}
		}
		return ""
	}()
	if rID == "" {
		rID = createCode(3)
		WriteUserData(uID, rID, 4)
		WriteShopData(sID, "", rID, 3)
	}
	writeRating(rID, uID, sID, ratVal)
	CreateLog("rate shop", uID, uID, sID, ratVal, 1)
	CreateLog("rate shop", sID, uID, sID, ratVal, 2)
}

func GetShopState(uID string) ShopState {
	UpdateData(0)
	UserDatum, scs := Users[uID]
	if !scs {
		fmt.Println("Error: User not found.")
		return ShopState{}
	}
	shopState := UserDatum.Shop.State
	switch shopState {
	case 2: // open
		return ShopState{Str: "\033[32m(Open)\033[0m", Num: 2}
	case 1: //closed
		return ShopState{Str: "\033[31m(Closed)\033[0m", Num: 1}
	}
	return ShopState{}
}

func GetShopRating(sID string) float32 {
	UpdateData(0)
	UserDatum, scs := Users[sID]
	if !scs {
		fmt.Println("Error: User not found.")
		return -1
	}

	ratingIDs := UserDatum.Shop.Ratings
	if len(ratingIDs) == 0 {
		return -1
	}

	rateVal := float32(0)
	for _, rID := range ratingIDs {
		rateVal += float32(Ratings[rID].Value)
	}
	return rateVal / float32(len(ratingIDs))
}

func ChangeShopState(sID string) {
	shopStt := GetShopState(sID).Num
	shopSttNew := shopStt%2 + 1
	shopSttTxt1, shopSttTxt2 := "", ""
	if shopSttNew == 2 {
		shopSttTxt1 = "closed"
		shopSttTxt2 = "open"
	} else {
		shopSttTxt1 = "open"
		shopSttTxt2 = "closed"
	}
	WriteShopData(sID, "", strconv.Itoa(shopSttNew), 2)
	CreateLog("change shop state", sID, shopSttTxt1, shopSttTxt2)
}

func ShowShop(uID string, showType uint8, show, from int) bool {
	// finds user from shop id
	UpdateData(0)
	UserDatum, scs := Users[uID]
	if !scs {
		fmt.Println("Error: User not found.")
		return true
	}
	shopState := UserDatum.Shop.State
	if shopState == 0 {
		fmt.Println("Error: shop for this user is not made")
		return true
	}

	// title & check if shop is empty
	fmt.Printf("\n-- %s's Item List %s --\n", UserDatum.Shop.ShopName, GetShopState(uID).Str)
	itemDatas := UserDatum.Shop.Items
	if len(itemDatas) == 0 {
		fmt.Println("No item on list.")
		return true
	}

	// determine the lengthiest item to remove offset
	idxLen := 0
	for {
		if len(itemDatas) >= int(math.Pow(10, float64(idxLen))) {
			idxLen++
			continue
		}
		break
	}
	idxLen += 2
	nameLen := 0
	for _, data := range itemDatas {
		itemName := data.Name
		if len(itemName) > nameLen {
			nameLen = len(itemName)
		}
	}
	nameLen += 1

	// determine settings for output
	space := "\033[36m.\033[0m"
	printRes := ""
	i := 0
	if show == 0 {
		show -= 1
	}
	for key, data := range itemDatas {
		if i < from {
			i++
			continue
		}
		// remove offset from item indexs & names
		itemIndex := strconv.Itoa(i+1) + "."
		for j := len(itemIndex); j < idxLen; j++ {
			itemIndex = itemIndex + space
		}
		itemName := data.Name
		for j := len(itemName); j < nameLen; j++ {
			itemName = itemName + space
		}
		itemDesc := data.Desc
		if itemDesc == "" {
			itemDesc = "-"
		}
		itemAmount := data.Amount
		itemCost := data.Cost

		switch showType { // output
		case 1: // w/ desc, w/ labelled stock, w/ cost
			if itemAmount == 0 {
				printRes += fmt.Sprintf("%s%s(\033[31mEmpty\033[0m)\nCost: %d\n%s\n", itemIndex, itemName, itemCost, itemDesc)
			} else {
				printRes += fmt.Sprintf("%s%s(Stock: %d)\nCost: %d\n%s\n", itemIndex, itemName, itemAmount, itemCost, itemDesc)
			}
		case 2: // no desc, w/ labelled stock, w/ cost
			if itemAmount == 0 {
				printRes += fmt.Sprintf("%s%s(\033[31mEmpty\033[0m)\nCost: %d; ID: %s\n", itemIndex, itemName, itemCost, key[6:])
			} else {
				printRes += fmt.Sprintf("%s%s(Stock: %d)\nCost: %d; ID: %s\n", itemIndex, itemName, itemAmount, itemCost, key[6:])
			}
		case 3: // no desc, no stock, no cost
			printRes = fmt.Sprintf("%s%s\n", itemIndex, itemName)
		default: // failsafe
			printRes = "Error: unknown shop print type."
		}
		i++
		if uint(i-from) >= uint(show) {
			break
		}
	}
	fmt.Print(printRes, "\n")

	return false // returns if the shop is empty
}

func ShowLogs(uID string, show, from, state int) {
	UpdateData(0)

	UserDatum, scs := Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}

	fmt.Printf("\n-- %s's Logs --\n", UserDatum.Username)
	userLogs := UserDatum.Logs
	if len(userLogs) == 0 {
		fmt.Println("No logs.")
		return
	}

	timezone, _ := time.Now().Zone()
	fmt.Printf("Time shown in UTC%s.\n", timezone)
	i := 0
	if show == 0 {
		show -= 1
	}
	for _, logID := range userLogs {
		if i < from {
			i++
			continue
		}
		logData := Logs[logID]
		if !(state >= logData.State) {
			continue
		}
		timeNow := time.Unix(logData.Time, 0)
		test1, test2 := [3]int{}, [3]string{}
		test1[0], test1[1], test1[2] = timeNow.Clock()
		for j := 0; j < 3; j++ {
			test2[j] = strconv.Itoa(test1[j])
			if len(test2[j]) < 2 {
				test2[j] = "0" + test2[j]
			}
		}
		timetime := fmt.Sprintf("%d-%d-%d %s:%s:%s", timeNow.Year(), int(timeNow.Month()), timeNow.Day(), test2[0], test2[1], test2[2])
		fmt.Printf("ID: %s; Time: %v; Data: %s\n", logID, timetime, logData.Data)
		i++
		if uint(i-from) >= uint(show) {
			break
		}
	}
}

// this program still prints all logs from all user, not from this user only =+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=
func WriteLogTxtFile(uID string) {
	UpdateData(0)
	_, scs := Users[uID]
	if !scs {
		fmt.Println("Error: User not found.")
		return
	}

	path := ""
	i := 1
	for {
		path = filepath.Join(goFilePath, "txtLogs", "logtxt-"+Users[uID].Username+"-"+strconv.Itoa(i)+".txt")
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			break
		}
		if err == nil {
			i++
			continue
		}
		if IsError(err) {
			return
		}
	}
	createFile(path)

	dataToWrite := fmt.Sprintf(" ===| %s 's Logs |=== \n", uID)
	for logID, val := range Logs {
		if val.Owner == uID {
			dataToWrite += fmt.Sprintf("\n%s\n\ttime: %v\n\tdata: %s\n", logID, time.Unix(val.Time, 0).UTC(), val.Data)
		}
	}
	CreateLog("get logs", uID)
	writeFile(path, []byte(dataToWrite), 0) // =+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
}
