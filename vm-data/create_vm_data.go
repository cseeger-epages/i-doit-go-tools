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
	"github.com/cseeger-epages/i-doit-go-api"
	"log"
	"strconv"
)

func main() {
	name := flag.String("n", "", "vm name")
	desc := flag.String("t", "", "description text")
	domain := flag.String("d", "", "domain")
	net3 := flag.String("l3", "", "net name for layer 3 net assignment")
	net2 := flag.String("l2", "", "net name for layer 2 net assignment")
	mac := flag.String("m", "", "mac address")
	ip := flag.String("i", "", "ip address")
	flag.Parse()

	// create api object
	goidoit.SkipTLSVerify(true)

	// enable Debug
	//goidoit.Debug(true)

	a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")
	//goidoit.Debug(true)

	if len(*name) > 0 && len(*desc) > 0 && len(*domain) > 0 && len(*ip) > 0 && len(*net3) > 0 && len(*net2) > 0 && len(*mac) > 0 {

		l3, err := a.GetObjectByType("C__OBJTYPE__LAYER3_NET", *net3)
		if err != nil {
			log.Fatal(err)
		}

		l3Id, _ := strconv.Atoi(l3.Result[0]["id"].(string))
		l2, err := a.GetObjectByType("C__OBJTYPE__LAYER2_NET", *net2)
		if err != nil {
			log.Fatal(err)
		}

		l2Id, _ := strconv.Atoi(l2.Result[0]["id"].(string))

		// object data
		data := struct {
			Type  string `json:"type"`
			Title string `json:"title"`
			Desc  string `json:"description"`
		}{"VIRTUALMACHINE(CLUSTER)", *name, *desc}

		result, err := a.CreateObj(data)
		if err != nil {
			log.Fatal(err)
		}
		id, _ := strconv.Atoi(result.Result[0]["id"].(string))

		// Host address

		NetData := struct {
			Hostname       string `json:"hostname"`
			Ip             string `json:"ipv4_address"`
			Ipv4Assingment int    `json:"ipv4_assignment"`
			NetType        int    `json:"net_type"`
			Net            int    `json:"net"`
			Domain         string `json:"domain"`
		}{*name, *ip, 1, 1, l3Id, *domain}

		a.CreateCat(id, "C__CATG__IP", NetData)

		// interfaces
		InterfData := struct {
			Title string `json:"title"`
		}{"eth0"}

		interf, _ := a.CreateCat(id, "C__CMDB__SUBCAT__NETWORK_INTERFACE_P", InterfData)

		interfId, _ := strconv.Atoi(interf.Result[0]["id"].(string))

		PortData := struct {
			Title       string   `json:"title"`
			Interface   int      `json:"interface"`
			PortType    string   `json:"port_type"`
			PortMode    int      `json:"port_mode"`
			PlugType    string   `json:"plug_type"`
			Negotiation int      `json:"negotiation"`
			Duplex      int      `json:"dublex"`
			Speed       int      `json:"speed"`
			SpeedType   int      `json:"speed_type"`
			Standard    int      `json:"standard"`
			Mac         string   `json:"mac"`
			Active      int      `json:"active"`
			Addresses   []string `json:"addresses"`
			L2Assign    int      `json:"layer2_assignment"`
			Desc        string   `json:"description"`
		}{"vm-network 172.20.0.0", interfId, "Ethernet", 1, "RJ-45", 1, 2, 1000, 3, 1, *mac, 1, []string{*ip}, l2Id, "autogenerated entry"}

		a.CreateCat(id, "C__CMDB__SUBCAT__NETWORK_PORT", PortData)

	} else {
		log.Fatal("missing parameter please set all parameters")
	}

	a.Logout()
}
