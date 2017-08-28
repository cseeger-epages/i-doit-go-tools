# create and delete vm's

## Step 1: build
Create your binarys using the `make.sh` script

## Step 2: required data

### create
- vm-name
- domain
- mac address
- ip address
- (optional) description
- layer 3 net (i-doit)
- layer 2 net (i-doit)

### delete
- vm-name

## Step 3: usage

### create
```
#!/bin/bash

bin/create_vm_data \
 -n <vm-name> \
 -d <domain> \
 -m <mac address> \
 -i <ip> \
 -t <description> \
 -l3 <i-doit layer-3 net> \
 -l2 <i-doit layer-2 net> \
```

### delete
```
bin/delete_vm -n <vm-name>
```
