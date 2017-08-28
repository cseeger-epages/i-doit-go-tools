# get scripts

## create

using the `make.sh` script.

## get\_domain

Is as simple as the name sais.
```
bin/get_domain -n <vm-name>
```

## get\_ip\_mac

Here the hard magic happens !

The `get_ip_mac` script has a bunch of options, `bin/get_ip_mac -h`
```
Usage of ./get_ip_mac:
  -a string
    	CIDR address subnet used for -g or -ip 
  -g	generate next free ip/mac combination flags
	-n <network name> must be specified
	-a <subnet in CIDR> optional to define a specific subnet for searching free addresses eg: -a 172.20.60.0/24 by having a full /16 net in -n defined
	-r <Number> to reserve the first <Number> ip addresses these will not used as free addresses
  -ip string
    	get mac for specified ip, requires -n <network>, optional -a
  -mac string
    	define mac prefix used for mac generation algorithm (default "00:50:57:3F:00:00")
  -n string
    	define i-doit network used for -g and -ip
  -r int
    	reserve N first ip addresses in combination with -g
  -vm string
    	get mac for virtual machine (cannot be used with -g)
```

gives you some hints. This script is used to get the next free ip/mac combination in an i-doit network. It also supports the option to have e.g. /16 network in i-doit where you only use 
a /24 subnet for autogeneration. So if you define the `-a` flag you can define a start address by `<somestartaddress>/24`. Also you have the option to reserver some ip addresses. Lets say
you have a /24 network and want to have the first 10 addresses free for maybe network environment stuff you can use the `-r` flag to reserve these ip addresses so they are not used as a free
address.

Here is an example usage:
```
bin/get_ip_mac -g -r 3 -n '<some idoit network with size and range 172.20.0.0/16>' -a 172.20.80.5/16
```
will use the defined i-doit network with the range 172.20.0.0/16 but only looks at the subnet 172.20.80.5/16 starting with address 172.20.80.5 and reserving the first 3 addresses. This
will give you the next free address in this net between 172.20.80.8 and 172.20.80.256.

Also it contains a mac generation algorithm using a mac prefix (by default `00:50:57:3F:00:00`) to generate unique mac addresses.
