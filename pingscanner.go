/*
Package pingscanner scan alive IPs of the given CIDR range in parallel.

Example usage:

package main

import (
	"fmt"

	ps "github.com/kotakanbe/go-pingscanner"
)

func main() {
	scanner := ps.PingScanner{
		CIDR: "192.168.11.0/24",
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
*/
package pingscanner

import (
	"net"
	"os/exec"
	"sort"
	"strings"
)

// PingScanner has information of Scanning.
type PingScanner struct {
	// CIDR (ex. 192.168.0.0/24)
	CIDR string

	// Number of concurrency ping process. (ex. 100)
	NumOfConcurrency int

	// ping command options. (ex. []string{"-c1", "-t1"})
	PingOptions []string
}

type pong struct {
	IP    string
	Alive bool
}

// Scan ping to hosts in CIDR range.
func (d PingScanner) Scan() (aliveIPs []string, err error) {
	var hostsInCidr []string
	if hostsInCidr, err = expandCidrIntoIPs(d.CIDR); err != nil {
		return nil, err
	}
	pingChan := make(chan string, d.NumOfConcurrency)
	pongChan := make(chan pong, len(hostsInCidr))
	doneChan := make(chan []pong)

	for i := 0; i < d.NumOfConcurrency; i++ {
		go ping(pingChan, pongChan, d.PingOptions...)
	}

	go receivePong(len(hostsInCidr), pongChan, doneChan)

	for _, ip := range hostsInCidr {
		pingChan <- ip
	}

	alives := <-doneChan
	for _, a := range alives {
		aliveIPs = append(aliveIPs, a.IP)
	}
	sort.Strings(aliveIPs)
	return
}

func expandCidrIntoIPs(cidr string) ([]string, error) {
	splitted := strings.Split(cidr, "/")
	if splitted[1] == "32" {
		return []string{splitted[0]}, nil
	}
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func ping(pingChan <-chan string, pongChan chan<- pong, pingOptions ...string) {
	for ip := range pingChan {
		pingOptions = append(pingOptions, ip)
		_, err := exec.Command("ping", pingOptions...).Output()
		var alive bool
		if err != nil {
			alive = false
		} else {
			alive = true
		}
		pongChan <- pong{IP: ip, Alive: alive}
	}
}

func receivePong(pongNum int, pongChan <-chan pong, doneChan chan<- []pong) {
	var alives []pong
	for i := 0; i < pongNum; i++ {
		pong := <-pongChan
		if pong.Alive {
			alives = append(alives, pong)
		}
	}
	doneChan <- alives
}
