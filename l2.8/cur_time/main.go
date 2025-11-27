package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	response, err := ntp.Time("pool.ntp.org")
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println(response.Format("2006-01-02 15:04:05 MST"))
}
