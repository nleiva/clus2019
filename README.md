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

4. Merge config

```bash
$ cd mergeconfig
$ go build
$ ./mergeconfig

config merged on [2001:420:2cff:1204::5502:2]:57344 -> Request ID: 1, Response ID: 1

2019/06/10 16:41:16 This process took 2.08032828s
```

5. Delete config

```bash
$ cd deleteconfig
$ go build
$ ./deleteconfig

config deleted on [2001:420:2cff:1204::5502:2]:57344 -> Request ID: 1, Response ID: 1
```

## gRPC config on router

```
grpc
 port 57344
 address-family ipv6
!
```

## Telemetry subscription ID 

It has to be preconfigured on the device <sup>[1](#myfootnote1)</sup>.

```
telemetry model-driven
 sensor-group LLDPNeighbor
  sensor-path Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail
 !
 subscription LLDP
  sensor-group-id LLDPNeighbor sample-interval 15000
 !
!
```

<a name="myfootnote1">[1]</a>: [gNMI](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi.proto) defines a variant where you do not need this config.

## Certificate files

You need to retrive the `ems.pem` file from the IOS XR device (after enabling gRPC/TLS) and put it in the [input](example/input) folder. You can find the file in the router on either `/misc/config/grpc/` or `/var/xr/config/grpc`.

- /var/xr/config/grpc

```console
$ ls -la
total 20
drwxr-xr-x  3 root root 4096 Jul  5 17:47 .
drwxr-xr-x 10 root root 4096 Jul  3 12:50 ..
drwx------  2 root root 4096 Jul  3 12:50 dialout
-rw-------  1 root root 1675 Jul  5 17:47 ems.key
-rw-rw-rw-  1 root root 1513 Jul  5 17:47 ems.pem
```
