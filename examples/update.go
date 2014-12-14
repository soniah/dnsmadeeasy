package main

import (
	"fmt"
	dme "github.com/soniah/dnsmadeeasy"
	"log"
	"os"
	"strconv"
)

func main() {
	akey := os.Getenv("akey")
	skey := os.Getenv("skey")
	domainID, err := strconv.ParseInt(os.Getenv("domainid"), 10, 64)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	recordID, err := strconv.ParseInt(os.Getenv("recordid"), 10, 64)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	fmt.Println("Using these values:")
	fmt.Println("akey:", akey)
	fmt.Println("skey:", skey)
	fmt.Println("domainid:", domainID)
	fmt.Println("recordid:", recordID)

	if len(akey) == 0 || len(skey) == 0 || domainID == 0 || recordID == 0 {
		log.Fatalf("Environment variable(s) not set\n")
	}

	client, err := dme.NewClient(akey, skey)
	client.URL = "http://api.sandbox.dnsmadeeasy.com/V2.0"
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	cr := map[string]interface{}{
		"name": "test-update",
	}

	result, err2 := client.UpdateRecord(domainID, recordID, cr)
	if err2 != nil {
		log.Fatalf("UpdateRecord result: %v error %v", result, err2)
	}
	log.Print("Result: ", result)
}
