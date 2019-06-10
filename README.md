# clus2019

CLUS 2019 Demo examples

1. GetConfig

```bash
$ cd getconfig
$ go run main.go

config from [2001:420:2cff:1204::5502:1]:57344
 {
 "data": {
  "openconfig-interfaces:interfaces": {
   "interface": [
    {
     "name": "Loopback60",
     "config": {
      "name": "Loopback60",
      "type": "iana-if-type:softwareLoopback",
      "enabled": true
     },
...
```

2. Show Command output (text)

```bash
$ cd showcmd
$ go build
$ ./showcmd -cli "show isis database" -enc text

output from [2001:420:2cff:1204::5502:2]:57344

----------------------------- show isis database ------------------------------

IS-IS BB2 (Level-2) Link State Database
LSPID                 LSP Seq Num  LSP Checksum  LSP Holdtime/Rcvd  ATT/P/OL
mrstn-5502-1.cisco.com.00-00  0x0000033a   0xe092        3126 /4000         0/0/0
mrstn-5502-2.cisco.com.00-00* 0x00000118   0xbf45        2023 /*            0/0/0
...
```

3. Set config (text)

```bash
$ cd setconfig
$ go build
$ ./setconfig -cli "interface Lo11 ipv6 address 2001:db8::/128"

config applied to [2001:420:2cff:1204::5502:2]:57344

2019/06/10 14:43:28 This process took 2.17090586s
```
