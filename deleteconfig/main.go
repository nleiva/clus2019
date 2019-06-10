package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	xr "github.com/nleiva/xrgrpc"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("This process took %s\n", elapsed)

}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	// YANG config; defaults to "yangconfig.json"
	ypath := flag.String("ypath", "../input/yangdelocconfig.json", "YANG path arguments")
	flag.Parse()

	// ID for the transaction.
	var id int64 = 1

	// Manually specify target parameters.
	router, err := xr.BuildRouter(
		xr.WithUsername("cisco"),
		xr.WithPassword("cisco"),
		//xr.WithHost("[2001:420:2cff:1204::5502:1]:57344"),
		//xr.WithCert("../input/certificate/router1.pem"),
		xr.WithHost("[2001:420:2cff:1204::5502:2]:57344"),
		xr.WithCert("../input/certificate/router2.pem"),
		xr.WithTimeout(5),
	)
	if err != nil {
		log.Fatalf("could not build a router, %v", err)
	}

	// Setup a connection to the target.
	conn, ctx, err := xr.Connect(*router)
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", router.Host, err)
	}
	defer conn.Close()

	// Get YANG config file to delete
	js, err := ioutil.ReadFile(*ypath)
	if err != nil {
		log.Fatalf("could not read file: %v: %v\n", *ypath, err)
	}

	// Delete 'js' config on target
	ri, err := xr.DeleteConfig(ctx, conn, string(js), id)
	if err != nil {
		log.Fatalf("failed to delete config from %s, %v", router.Host, err)
	} else {
		fmt.Printf("\nconfig deleted on %s -> Request ID: %v, Response ID: %v\n\n", router.Host, id, ri)
	}
}
