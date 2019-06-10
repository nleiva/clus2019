package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	proto "github.com/golang/protobuf/proto"
	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
)

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func main() {
	// Subs options; LLDP, we will add some more
	p := flag.String("subs", "LLDP", "Telemetry Subscription")
	// Encoding option; defaults to GPBKV (only one supported in this example)
	enc := flag.String("enc", "gpbkv", "Encoding: 'json', 'gpb' or 'gpbkv'")
	flag.Parse()

	mape := map[string]int64{
		"gpb":   2,
		"gpbkv": 3,
		"json":  4,
	}
	e, ok := mape[*enc]
	if !ok {
		log.Fatalf("encoding option '%v' not supported", *enc)
	}

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
		xr.WithTimeout(60),
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, ech, err := xr.GetSubscription(ctx, conn, *p, id, e)
	if err != nil {
		log.Fatalf("could not setup Telemetry Subscription: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			fmt.Printf("\nmanually cancelled the session to %v\n\n", router.Host)
			cancel()
			return
		case <-ctx.Done():
			// Timeout: "context deadline exceeded"
			err = ctx.Err()
			fmt.Printf("\ngRPC session timed out after %v seconds: %v\n\n", router.Timeout, err.Error())
			return
		case err = <-ech:
			// Session canceled: "context canceled"
			fmt.Printf("\ngRPC session to %v failed: %v\n\n", router.Host, err.Error())
			return
		}
	}()

	line := strings.Repeat("*", 90)
	for tele := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(tele, message)
		if err != nil {
			log.Fatalf("could not unmarshall the message: %v\n", err)
		}
		ts := message.GetMsgTimestamp()
		ts64 := int64(ts * 1000000)
		fmt.Println(line)
		fmt.Printf("Time %v, Path: %v\n", time.Unix(0, ts64).Format("03:04:05PM"), message.GetEncodingPath())
		fmt.Println(line)
		exploreFields(message.GetDataGpbkv(), "")
	}

}

func exploreFields(f []*telemetry.TelemetryField, indent string) {
	for _, field := range f {
		switch field.GetFields() {
		case nil:
			decodeKV(field, indent)
		default:
			exploreFields(field.GetFields(), indent+" ")
		}
	}
}

func decodeKV(f *telemetry.TelemetryField, indent string) {
	// This is incomplete, but covers most of the cases I've seen so far.
	switch f.GetValueByType().(type) {
	case *telemetry.TelemetryField_StringValue:
		fmt.Printf("%s%s: %s\n", indent, f.GetName(), f.GetStringValue())
	case *telemetry.TelemetryField_BoolValue:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetBoolValue())
	case *telemetry.TelemetryField_Uint32Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetUint32Value())
	case *telemetry.TelemetryField_Uint64Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetUint64Value())
	case *telemetry.TelemetryField_BytesValue:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetBytesValue())
	case *telemetry.TelemetryField_Sint32Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetSint32Value())
	case *telemetry.TelemetryField_Sint64Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetSint64Value())
	default:
	}
}
