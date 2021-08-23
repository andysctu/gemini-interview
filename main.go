package main

import (
	"fmt"
    "net/http"

    "github.com/andysctu/jobcoin/mixer"
)

func main() {
	http.HandleFunc("/getDepositAddressForMixing", mixer.GetDepositAddressForMixing)
	http.ListenAndServe(":9001", nil)

	fmt.Println("Ready for mixing on port 9001")
}
