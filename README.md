# pingscanner

Scanning alive hosts of the given CIDR range in parallel.

## Usage

```go
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
```

## Author

[kotakanbe](https://github.com/kotakanbe)

## License

Please see [LICENSE](https://github.com/kotakanbe/go-pingscanner/blob/master/LICENSE).
