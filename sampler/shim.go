package sampler

import "fmt"

var reverse *cache

func Init() {
	reverse = newcache()
	fmt.Println("Reverse DNS cache initialized")
}

func GetHostnameByByteSlice(ip []byte) string {
	if reverse, ok := reverse.get(ip); ok {
		return reverse
	}

	// Latency will be low at the cost of missing items during the discovery phase
	go reverse.set(ip)
	return ""
}
