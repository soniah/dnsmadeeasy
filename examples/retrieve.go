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
	recordID := os.Getenv("recordid")

	fmt.Println("Using these values:")
	fmt.Println("akey:", akey)
	fmt.Println("skey:", skey)
	fmt.Println("domainid:", domainID)
	fmt.Println("recordid:", recordID)

	if len(akey) == 0 || len(skey) == 0 || len(domainID) == 0 || len(recordID) == 0 {
		log.Fatalf("Environment variable(s) not set\n")
	}

	client, err := dme.NewClient(akey, skey)
	client.URL = "http://api.sandbox.dnsmadeeasy.com/V2.0"
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	req, err2 := client.RetrieveRecord(domainID, recordID)
	if err2 != nil {
		log.Fatalf("RetrieveRecord result: %v error %v", req, err2)
	}

	log.Print("Result: ", *req)
}
