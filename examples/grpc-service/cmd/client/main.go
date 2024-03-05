package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/555f/gg/examples/grpc-service/pkg/client"
	"github.com/555f/gg/examples/grpc-service/pkg/dto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:9001", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	cc := client.NewProfileControllerCient(conn)

	profileCh := make(chan *dto.Profile)

	go func() {
		for i := 1; i < 5; i++ {
			profileCh <- &dto.Profile{ID: i, FistName: "Doe"}
		}
		close(profileCh)
	}()

	outCh, err := cc.Stream(profileCh)
	if err != nil {
		log.Fatalf("fail create stream: %v", err)
	}

	for s := range outCh {
		fmt.Printf("%v\n", s)
	}

	profileCh = make(chan *dto.Profile, 1)

	profileCh <- &dto.Profile{ID: 4, FistName: "Doe"}

	outCh, err = cc.Stream(profileCh)
	if err != nil {
		log.Fatalf("fail create stream: %v", err)
	}

	for s := range outCh {
		fmt.Printf("%v\n", s)
	}
}
