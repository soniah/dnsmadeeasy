package main

import (
	"fmt"
	dme "github.com/soniah/dnsmadeeasy"
	"log"
	"os"
)

func main() {
	akey := os.Getenv("akey")
	skey := os.Getenv("skey")
	domainID := os.Getenv("domainid")
	ip := os.Getenv("ip")

	fmt.Println("Using these values:")
	fmt.Println("akey:", akey)
	fmt.Println("skey:", skey)
	fmt.Println("domainid:", domainID)
	fmt.Println("ip:", ip)

	if len(akey) == 0 || len(skey) == 0 || len(domainID) == 0 || len(ip) == 0 {
		log.Fatalf("Environment variable(s) not set\n")
	}

	client, err := dme.NewClient(akey, skey)
	client.URL = "http://api.sandbox.dnsmadeeasy.com/V2.0"
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	opts := dme.ChangeRecord{
		Name:  "test",
		Type:  "A",
		Value: ip,
		TTL:   86400,
	}

	result, err2 := client.CreateRecord(domainID, &opts)
	if err2 != nil {
		log.Fatalf("Result: '%s' Error: %s", result, err2)
	}

	log.Printf("Result: '%s'", result)
}
