package main

import (
	"flag"
	"fmt"
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

	// CLI config to apply; defaults to "interface lo1 desc test"
	cli := flag.String("cli", "interface lo1 desc test", "Config to apply")
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

	// Apply 'cli' config to target
	err = xr.CLIConfig(ctx, conn, *cli, id)
	if err != nil {
		log.Fatalf("failed to config %s, %v", router.Host, err)
	}
	fmt.Printf("\nconfig applied to %s\n\n", router.Host)
}
