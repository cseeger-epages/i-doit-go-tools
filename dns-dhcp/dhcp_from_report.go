package main

import (
	"fmt"
	"github.com/cseeger-epages/i-doit-go-api"
	"log"
	"os"
	"strconv"
)

func help() {
	fmt.Println("usage: ./dhcp_from_report <report_id>")
	os.Exit(1)
}

func main() {
	var reportId = 0
	if len(os.Args) > 1 {
		reportId, _ = strconv.Atoi(os.Args[1])
	} else {
		help()
	}

	goidoit.SkipTLSVerify(true)
	a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")

	data, err := a.GetReport(reportId)
	if err != nil {
		log.Fatal("Error while requesting report ", err)
	}

	for _, v := range data.Result {
		if len(v["Hostname"].(string)) == 0 || len(v["MAC-address"].(string)) == 0 || len(v["Domain"].(string)) == 0 || len(v["Host address"].(string)) == 0 {
			fmt.Fprintf(os.Stderr, "ERROR: some variable is empty Hostname: \"%s\"; MAC-address: \"%s\"; Host address: \"%s\"; Domain: \"%s\"\n", v["Hostname"], v["MAC-address"], v["Host address"], v["Domain"])
		} else {
			fmt.Printf("host %s { hardware ethernet %s; fixed-address %s; option host-name \"%s\"; option domain-name \"%s\"; }\n", v["Hostname"], v["MAC-address"], v["Host address"], v["Hostname"], v["Domain"])
		}
	}
	a.Logout()
}
