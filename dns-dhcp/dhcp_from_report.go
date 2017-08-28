/*
	i-doit-go-tools

	Copyright (C) 2017 Carsten Seeger

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.

	@author Carsten Seeger
	@copyright Copyright (C) 2017 Carsten Seeger
	@license http://www.gnu.org/licenses/gpl-3.0 GNU General Public License 3
	@link https://github.com/cseeger-epages/i-doit-go-tools
*/

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
