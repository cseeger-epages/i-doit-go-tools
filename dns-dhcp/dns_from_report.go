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
	"flag"
	"fmt"
	"github.com/cseeger-epages/i-doit-go-api"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	reverse := flag.Bool("r", false, "create reverse entrys for /16 zone files")
	splitSubnet := flag.Int("s", 0, "split reverse entrys to subnet and create only /24 reverse entrys, e.g: -s 154 splits the x.x.x.x/16 subnet to only create entrys for the subnet for x.x.154.x/24")
	overwriteDomain := flag.String("d", "", "overwrite domain with a specific value")
	id := flag.String("id", "", "set id")

	flag.Parse()

	var reportId = 0
	reportId, _ = strconv.Atoi(*id)

	goidoit.SkipTLSVerify(true)
	a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")

	goidoit.SkipTLSVerify(true)

	data, err := a.GetReport(reportId)
	if err != nil {
		log.Fatal("Error while requesting report ", err)
	}

	for _, v := range data.Result {
		if len(*overwriteDomain) > 0 {
			v["Domain"] = *overwriteDomain
		}
		if v["Hostname"] == nil || v["Domain"] == nil || v["Host address"] == nil {
			if len(v["Hostname"].(string)) == 0 || len(v["Domain"].(string)) == 0 || len(v["Host address"].(string)) == 0 {
				fmt.Fprintf(os.Stderr, "ERROR: some variable is empty Hostname: \"%s\"; Host address: \"%s\"; Domain: \"%s\"\n", v["Hostname"], v["Host address"], v["Domain"])
			}
		} else {
			if *reverse {
				var ip_sep = strings.Split(v["Host address"].(string), ".")
				if len(v["Domain"].(string)) > 0 {
					if *splitSubnet != 0 {
						test, _ := strconv.Atoi(ip_sep[2])
						if test == *splitSubnet {
							fmt.Printf("%s\tIN\tPTR\t%s.%s.\n", ip_sep[3], v["Hostname"], v["Domain"])
						}
					} else {
						fmt.Printf("%s.%s\tIN\tPTR\t%s.%s.\n", ip_sep[3], ip_sep[2], v["Hostname"], v["Domain"])
					}
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: Domain is empty for Hostname %s with address %s\n", v["Hostname"], v["Host address"])
				}
			} else {
				fmt.Printf("%s\tIN\tA\t %s\n", v["Hostname"], v["Host address"])
			}
		}
	}
	a.Logout()
}
