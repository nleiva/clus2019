package main

import (
	"flag"
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

	// IPv6 prefix to setup; defaults to "2001:db8::/32"
	pfx := flag.String("pfx", "2001:db8::/32", "IPv6 prefix to setup")
	// IPv6 next-hop to setup; defaults to "2001:db8:cafe::1"
	nh := flag.String("nh", "2001:db8:cafe::1", "IPv6 next-hop to setup")
	flag.Parse()

	// Admin Distance
	var admdis uint32 = 2

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
	conn, _, err := xr.Connect(*router)
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", router.Host, err)
	}
	defer conn.Close()

	// CSCva95005: Return SL_NOT_CONNECTED when the init session is killed from the Client.
	err = xr.ClientInit(conn)
	if err != nil {
		log.Fatalf("Failed to initialize connection to %s, %v", router.Host, err)
	}

	// VRF Register Operation (= 1),
	err = xr.VRFOperation(conn, 1, admdis)
	if err != nil {
		log.Fatalf("Failed to register the VRF Operation on %s, %v", router.Host, err)
	}
	// VRF EOF Operation (= 3),
	err = xr.VRFOperation(conn, 3, admdis)
	if err != nil {
		log.Fatalf("Failed to send VRF Operation EOF to %s, %v", router.Host, err)
	}
	// Route Add Operation (= 1),
	err = xr.SetRoute(conn, 1, *pfx, admdis, *nh)
	if err != nil {
		log.Fatalf("Failed to set Route on %s, %v", router.Host, err)
	}

}
