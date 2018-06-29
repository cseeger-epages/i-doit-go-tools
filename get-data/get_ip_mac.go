package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cseeger-epages/i-doit-go-api"
)

func main() {
	generate := flag.Bool("g", false, "generate next free ip/mac combination flags\n\t-n <network name> must be specified\n\t-a <subnet in CIDR> optional to define a specific subnet for searching free addresses eg: -a 172.20.60.0/24 by having a full /16 net in -n defined\n\t-r <Number> to reserve the first <Number> ip addresses these will not used as free addresses")
	network := flag.String("n", "", "define i-doit network used for -g and -ip")
	macPrefix := flag.String("mac", "00:50:57:3F:00:00", "define mac prefix used for mac generation algorithm")
	cidr := flag.String("a", "", "CIDR address subnet used for -g or -ip ")
	reserve := flag.Int("r", 0, "reserve N first ip addresses in combination with -g")
	vm := flag.String("vm", "", "get mac for virtual machine (cannot be used with -g)")
	ip := flag.String("ip", "", "get mac for specified ip, requires -n <network>, optional -a")

	flag.Parse()

	if len(*ip) > 0 {
		if len(*network) > 0 {
			if len(*macPrefix) > 0 {
				if validateIp(*ip, *network) {
					var cidrAddr string
					if len(*cidr) == 0 {
						goidoit.SkipTLSVerify(true)
						a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")
						data, _ := a.GetObjectByType("C__OBJTYPE__LAYER3_NET", *network)

						id, _ := strconv.Atoi(data.Result[0]["id"].(string))
						netData, _ := a.GetCat(id, "C__CATS__NET")
						cidrAddr = fmt.Sprintf("%v/%v", netData.Result[0]["address"], netData.Result[0]["cidr_suffix"])
					} else {
						cidrAddr = *cidr
					}
					fmt.Println(calcMac(*ip, *macPrefix, cidrAddr))
				} else {
					log.Fatal("ip is not valid or network does not exists")
				}
			} else {
				log.Fatal("no valid mac prefix")
			}
		} else {
			log.Fatal("no network defined please use -n <i-doit network>")
		}
		os.Exit(0)
	}

	if *generate {
		if len(*network) > 0 {
			var ipaddr string
			ipaddr = getNextFreeIpMac(*network, *cidr, *reserve)
			if len(*macPrefix) > 0 {
				var cidrAddr string
				if len(*cidr) == 0 {
					goidoit.SkipTLSVerify(true)
					a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")
					data, _ := a.GetObjectByType("C__OBJTYPE__LAYER3_NET", *network)

					id, _ := strconv.Atoi(data.Result[0]["id"].(string))
					netData, _ := a.GetCat(id, "C__CATS__NET")
					cidrAddr = fmt.Sprintf("%v/%v", netData.Result[0]["address"], netData.Result[0]["cidr_suffix"])
				} else {
					cidrAddr = *cidr
				}
				mac := calcMac(ipaddr, *macPrefix, cidrAddr)
				fmt.Printf("%v,%v\n", ipaddr, mac)
			} else {
				log.Fatal("no valid mac prefix")
			}
		} else {
			fmt.Println("-n flag is needed to define Network")
			os.Exit(0)
		}
		os.Exit(0)
	}

	if len(*vm) > 0 {
		macs := getMac(*vm)
		for _, v := range macs {
			fmt.Printf("%v\n", v)
		}
		os.Exit(0)
	}

	flag.PrintDefaults()
}

func incHex(hexStr string) string {
	hexCounter, _ := hex.DecodeString(hexStr)
	hexCounter[0]++
	return hex.EncodeToString(hexCounter)
}

func incHexByVal(hexStr string, val int) string {
	hexCounter, _ := hex.DecodeString(hexStr)
	hexCounter[0] += byte(val)
	return hex.EncodeToString(hexCounter)
}

func calcMac(ip string, prefix string, cidr string) string {

	ipaddr := net.ParseIP(ip)
	ipRange, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal("no valid CIDR address")
	}

	if !ipNet.Contains(ipaddr) {
		log.Fatal("ip address not in range of your CIDR address")
	}

	netMask := ipNet.String()
	netMaskSplit := strings.Split(netMask, ".")

	calc := strings.Split(ip, ".")
	ipSplit := strings.Split(ipRange.String(), ".")
	prefixSplit := strings.Split(prefix, ":")

	i := 0
	mac := make([]string, 4)
	inc := make([]int, 4)
	offset := make([]int, 4)
	for k, v := range netMaskSplit {
		val, _ := strconv.Atoi(v)
		if val == 0 {
			mac[k] = "00"
			inc[k] = 0
			inc[k], _ = strconv.Atoi(calc[k])
			offset[k], _ = strconv.Atoi(ipSplit[k])
			mac[k] = incHexByVal(mac[k], inc[k]-offset[k])
			i++
		}
	}
	plen := len(prefixSplit) - (i + 1)
	retval := ""
	for j := 0; j <= plen; j++ {
		if j == 0 {
			retval = prefixSplit[j]
		} else {
			retval = fmt.Sprintf("%v:%v", retval, prefixSplit[j])
		}
	}
	for j := 0; j < 4; j++ {
		if mac[j] != "" {
			retval = fmt.Sprintf("%v:%v", retval, mac[j])
		}
	}
	return retval
}

func validateIp(ip string, network string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		fmt.Printf("%v is not a valid IPv4 address\n", trial)
		return false
	}
	goidoit.SkipTLSVerify(true)
	a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")

	obj, _ := a.GetObjectByType("C__OBJTYPE__LAYER3_NET", network)

	if obj.Result == nil {
		return false
	}

	id, _ := strconv.Atoi(obj.Result[0]["id"].(string))
	netw, _ := a.GetCat(id, "C__CATS__NET")

	cidr := fmt.Sprintf("%v/%v", netw.Result[0]["address"], netw.Result[0]["cidr_suffix"])

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal("invalid CIDR address")
	}

	if ipNet.Contains(trial) {
		return true
	}
	return false

}

func getMac(obj string) []string {
	goidoit.SkipTLSVerify(true)
	a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")

	objdata, _ := a.GetObject(obj)
	id, _ := strconv.Atoi(objdata.Result[0]["id"].(string))
	objcat, _ := a.GetCat(id, "C__CMDB__SUBCAT__NETWORK_PORT")
	var ret []string
	for _, v := range objcat.Result {
		ret = append(ret, fmt.Sprintf("%v,%v", v["title"], v["mac"]))
	}
	a.Logout()
	return ret
}

func getNextFreeIpMac(network string, cidr string, reserve int) string {
	goidoit.SkipTLSVerify(true)
	a, _ := goidoit.NewLogin("<i-doit-url>/src/jsonrpc.php", "<apikey>", "<username>", "<password>")

	data, err := a.GetObjectByType("C__OBJTYPE__LAYER3_NET", network)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.Atoi(data.Result[0]["id"].(string))
	netData, _ := a.GetCat(id, "C__CATS__NET")

	if len(cidr) == 0 {
		cidr = fmt.Sprintf("%v/%v", netData.Result[0]["address"], netData.Result[0]["cidr_suffix"])
	}

	cat, _ := a.GetCat(id, "C__CATS__NET_IP_ADDRESSES")

	var ipList []net.IP
	for _, v := range cat.Result {
		if v["ipv4_assignment"] != nil {
			ip := net.ParseIP(v["title"].(string))
			if ip != nil {
				ipList = append(ipList, ip)
			}
		}
	}

	sort.Sort(ByAddr(ipList))
	a.Logout()
	ip := getNextFreeIP(ipList, cidr, reserve)
	test := net.ParseIP(netData.Result[0]["range_to"].(string))
	trial := net.ParseIP(ip)
	if bytes.Compare(trial, test) >= 0 {
		log.Fatalf("no free ip addresses in %v\n", network)
	}
	return ip
}

func getNextFreeIP(ipList []net.IP, cidr string, reserve int) string {
	var ip net.IP
	addr, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatalf("CIDR ERROR: %s\n", err)
	}
	count := 0
	for search := addr; ipNet.Contains(search); inc(search) {
		count++
		// reserve first N ip addresses
		if count <= reserve {
			continue
		}
		if !contains(ipList, search) {
			ip = search
			break
		}
	}

	return ip.String()
}

// slice contains
func contains(s []net.IP, e net.IP) bool {
	for _, a := range s {
		if a.Equal(e) {
			return true
		}
	}
	return false
}

// increment ip helper function
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// sort functions for IP addresses
type ByAddr []net.IP

func (ip ByAddr) Len() int {
	return len(ip)
}

func (ip ByAddr) Swap(i, j int) {
	ip[i], ip[j] = ip[j], ip[i]
}

func (ip ByAddr) Less(i, j int) bool {
	ipi := ip[i].To4()
	ipj := ip[j].To4()

	for k := 0; k <= 3; k++ {
		if ipi[k] < ipj[k] {
			return true
		} else if ipi[k] > ipj[k] {
			return false
		}
	}
	return false
}
