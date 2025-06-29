package main

import (
	"fmt"
	"log"

	"github.com/PeterCullenBurbery/go_functions_002/date_time_functions"
)

func main() {
	stamp, err := date_time_functions.Date_time_stamp()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("🧾 Java timestamp:", date_time_functions.Safe_time_stamp(stamp, 1))
}