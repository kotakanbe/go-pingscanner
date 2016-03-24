package main

import (
	"fmt"

	ps "github.com/kotakanbe/go-pingscanner"
)

func main() {
	scanner := ps.PingScanner{
		//  CIDR: "192.168.11.0/24",
		CIDR: "192.168.11.2/32",
		PingOptions: []string{
			"-c1",
			"-t1",
		},
		NumOfConcurrency: 100,
	}
	if aliveIPs, err := scanner.Scan(); err != nil {
		fmt.Println(err)
	} else {
		if len(aliveIPs) < 1 {
			fmt.Println("no alive hosts")
		}
		for _, ip := range aliveIPs {
			fmt.Println(ip)
		}
	}
}
