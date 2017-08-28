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
	"strconv"
)

func main() {
	name := flag.String("n", "", "get domain for vm name")
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
		data, err := a.GetCategory(id, "C__CATG__IP")
		if err != nil {
			log.Fatalf("GetCategory: %s\n", err)
		}
		fmt.Printf("%+v\n", data.Result[0]["domain"])
	} else {
		log.Fatalf("object not found")
	}

	a.Logout()
}
