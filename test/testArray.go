package main

import (
	//"encoding/json"
	"encoding/base64"
	"fmt"
	//"u9/tool"
)

func main() {

	plant := `银子1两`

	cpStr := base64.URLEncoding.EncodeToString([]byte(plant))

	fmt.Println(cpStr)

	ppBy, err := base64.URLEncoding.DecodeString("cpStr")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(ppBy))
	//test, err := tool.DecodeStdEncodeing("successsuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccesssuccess")
	//fmt.Print(test)
	//fmt.Print(err)
	// params := make(map[string]string)
	// params["1"] = "2"
	// params["2"] = "3"
	// b, _ := json.Marshal(params)
	// fmt.Print(string(b))
}
