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

	// Encoding option; defaults to JSON
	enc := flag.String("enc", "json", "Encoding: 'json' or 'text'")
	// CLI to issue; defaults to "show grpc status"
	cli := flag.String("cli", "show grpc status", "Command to execute")
	flag.Parse()

	// ID for the transaction.
	id := 1
	var output string

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

	// Return show command output based on encoding selected
	switch *enc {
	case "text":
		output, err = xr.ShowCmdTextOutput(ctx, conn, *cli, int64(id))
	case "json":
		output, err = xr.ShowCmdJSONOutput(ctx, conn, *cli, int64(id))
	default:
		log.Fatalf("don't recognize encoding: %v\n", *enc)
	}
	if err != nil {
		log.Fatalf("couldn't get the cli output: %v\n", err)
	}
	fmt.Printf("\noutput from %s\n %s\n", router.Host, output)
}
