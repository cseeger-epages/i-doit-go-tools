# i-doit-go-tools

are a set of tools used for automating vm creation process using vmware and i-doit.

You can find more information in my blogposts:
- [creating-systems-with-pipelining](https://developer.epages.com/blog/2015/08/03/creating-systems-with-pipelining.html)
- [going-forward](https://developer.epages.com/blog/2017/08/01/going-forward.html)

The software can be used under the GPLv3 see LICENCE.

## Folders

### dns-dhcp
Here you can find 2 scripts for generating dhcp and dns configurations for default linux dhcp (isc-dhcp-server) and a typical linux dns server (bind9).

### get-data
A set of scripts used for getting specific data from i-doit using [i-doit-go](https://github.com/cseeger-epages/i-doit-go-api).

### vm-data
create and delete scripts for vm data in i-doit

every folder has an own README.md describing the scripts and there use cases and usage in detail.

## Attention !

Before you start using these scripts you have to adjust some lines in the script. The easiest way is to search for the following line in every go script:

```
goidoit.NewLogin
```
and adjust these lines by adding your i-doit-url, api-key, username and password to it. Also there are some scripts having multiple lines feel free to add a global var with you data
and change the NewLogin using these variables.

After changing your credentials the scripts should work properly and can by compiled using the `make.sh` script. The script creates a bin folder where all binaries are in.
