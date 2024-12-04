package myLogs

import (
	"ecommerce/ent"
	"fmt"
)

func GetLogs(uID string) {
	fmt.Println(" === Get Logs")

	ent.UpdateData(0)
	_, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}

	ent.CreateLog("get logs", uID)
	ent.WriteLogTxtFile(uID)
}
