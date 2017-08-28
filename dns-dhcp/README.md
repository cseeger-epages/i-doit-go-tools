# DHCP generation

## Step 1: i-doit report
The fasted way i found to create these data is to make use of the report manager in i-doit.

The used Attributes shown in the report are
- Hostname (Hostaddress)
- MAC-Address (Port)
- Hostaddress (Hostaddress)
- Domain (Hostaddress)

the Conditions are the following:
```
Category: Port (Network)
Attribute: MAC-Address LIKE %...%  _ AND

Category: Hostaddress
Attribute: Hostname LIKE %...%  _ AND

Category: Hostaddress
Attribute: Net =  <select your layer 3 net where you want to create DHCP data for> AND/OR

Category: Hostaddress
Attribute: Net =  <select another layer 3 net where you want to create DHCP data for> AND/OR

Category: Hostaddress
Attribute: Addresstype IPv4 = DHCP  AND

Category: Global
Attribute: Status = Normal
```

The first two conditions check if MAC and Hostname are not empty.
The other should be self explaining :P.

Create this report and remember the report id!

## Step 2: general templates for configuration

Under templates-dhcp you will find a simple template example for a 172.20.0.0/16 Net.
The `range` command is outcommented depending on the fact if you want to allow a dynamic ip range in your net or only known mac/ip combinations.
The `include` includes the dhcpd.lease file wich will be created using the `dhcp_from_report` binary.

Change these files to fit your needs.

## Step 3: create the dhcpd.lease file

To test if your configuration is working you can use:
`go run dhcp_from_report.go <report id>`
you should see some output like
```
host vm-name { hardware ethernet XX:XX:XX:XX:XX:XX; fixed-address <ipv4address>; option host-name "vm-name"; option domain-name "your.domain"; }
```
these are the leases entrys. To automate the process its better to create a binary file so there is no need to install go to the system used.
```
go build dhcp_from_report.go -o dhcp_from_report
```
or use the `make.sh` in the root directory of this repository, this creates all go binarys and puts them under `bin`

## Step 4: deployment

The script itself can be used in the following way:

```
bin/dhcp_from_report 22 1>dhcpd.leases.vm 2>dhcpd.error
```
The stdout only shows the correct entrys while the stderr will display any errors.

You can test if your dhcpd configuration works using
```
dhcpd -t -cf dhcpd.conf
```
if yes you can deploy it to your server. If the dhcpd server has any problems reading these files try to perm the files `644` and the dhcpd.d directory to `755`.

# DNS generation

## Step 1: i-doit report

The used Attributes shown in the report are
- Hostaddress (Hostaddress)
- Hostname (Hostaddress)
- Domain (Hostaddress)

the Conditions are the following:
```
Category: Hostaddress
Attribute: Domain LIKE %...%  _ AND

Category: Hostaddress
Attribute: Hostname LIKE %...%  _ AND

Category: Hostaddress
Attribute: Net =  <select your layer 3 net where you want to create DHCP data for> AND

Category: Hostaddress
Attribute: Net =  <select another layer 3 net where you want to create DHCP data for> AND

Category: Global
Attribute: Status = Normal
```

## Step 2: general templates for configuration

Under templates-dns you will find an example configuration `named.conf.local` for the domain `your.domain` using the 172.20.0.0/16 network.
This file includes the `db.172.20` and `db.your.domain` files wich will be generated using the `db.172.20.template` and `db.your.domain.template`.

Change these files to fit your needs. As you can see the file names use a specific naming convention by changing the net and domain please change the filenames too.
The `notify` option should be used in case you want to notify any other e.g. internal domain controler dns server.

## Step 3: create the db. files

creating the binary is the same procedure as by the dhcp go files.
The `dns_from_report` has some more options since you have to define if you want to create forward or reverse entrys.
To create the forward entrys:
```
bin/dns_from_report -id <report id> 1>db.your.domain.data 2>dns.error
```
for reverse entrys:
```
bin/dns_from_report -id <report id> -r 1>db.20.172.data 2>dns.error
```

## Step 4: deployment

if you have these files just cut them together with your templates and deploy the named.conf.local and your db. files to your bind9 server and test the configuration:
```
# create temporary directory
cd /tmp
mkdir etc/bind

# copy all files
cp -v name.conf.local etc/bind/

cp -v db.20.172.template etc/bind/db.20.172
cat db.20.172.data >> etc/bind/db.20.172

cp -v db.your.domain.template db.your.domain
cat db.your.domain.data >> db.your.domain

# set proper SERIAL
sed -i "s/\$SERIAL/$(date +%Y%m%d)$BUILD_NUMBER/" etc/bind/db.20.172
sed -i "s/\$SERIAL/$(date +%Y%m%d)$BUILD_NUMBER/" etc/bind/db.your.domain

# check the configuration
named-checkconf etc/bind/named.conf.local
named-checkzone your.domain etc/bind/db.20.172
named-checkzone your.domain etc/bind/db.your.domain

cp -v etc/bind/* /etc/bind
service bind9 reload
```


