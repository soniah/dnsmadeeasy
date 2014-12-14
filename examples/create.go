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

	ip := os.Getenv("ip")

	fmt.Println("Using these values:")
	fmt.Println("akey:", akey)
	fmt.Println("skey:", skey)
	fmt.Println("domainid:", domainID)
	fmt.Println("ip:", ip)

	if len(akey) == 0 || len(skey) == 0 || domainID == 0 || len(ip) == 0 {
		log.Fatalf("Environment variable(s) not set\n")
	}

	client, err := dme.NewClient(akey, skey)
	client.URL = "http://api.sandbox.dnsmadeeasy.com/V2.0"
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	cr := map[string]interface{}{
		"name":  "test",
		"type":  "A",
		"value": ip,
		"ttl":   86400,
	}

	result, err2 := client.CreateRecord(domainID, cr)
	if err2 != nil {
		log.Fatalf("Result: '%d' Error: %s", result, err2)
	}

	log.Printf("Result: '%d'", result)
}
