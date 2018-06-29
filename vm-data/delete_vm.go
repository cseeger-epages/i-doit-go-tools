package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/cseeger-epages/i-doit-go-api"
)

func main() {
	name := flag.String("n", "", "vm name")
	flag.Parse()

	// create api object
	goidoit.SkipTLSVerify(true)
	a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")

	// enable Debug
	//goidoit.Debug(true)

	idData, err := a.GetObject(*name)
	if err != nil {
		log.Fatalf("GetObject: %s\n", err)
	}
	if idData.Result != nil {
		id, err := strconv.Atoi(idData.Result[0]["id"].(string))
		if err != nil {
			log.Fatalf("strconv: %s\n", err)
		}
		_, err = a.Quickpurge(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v with id: %v successfully deleted", *name, id)
	} else {
		log.Fatalf("vm not found")
	}

	a.Logout()
}
