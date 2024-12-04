package myLogs

import (
	"ecommerce/ent"
	"fmt"
)

func ViewLogs(uID string) {
	fmt.Println(" === View Logs")

	ent.UpdateData(0)
	_, scs := ent.Users[uID]
	if !scs {
		fmt.Println("Error: user not found")
		return
	}

	ent.CreateLog("view logs", uID)
	ent.ShowLogs(uID, 0, 0, 1)
}
