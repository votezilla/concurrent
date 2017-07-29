// concurrent.go
package main

import (
	"fmt"
	"os"
	"net/http"
	"strconv"
	"time"
)

func MakeRequest(url string, ch chan<-time.Duration) {
	start := time.Now()
	http.Get(url)
	
	ch <- time.Since(start)
}

func timeSite(url string, numRequests int) {
	var min, max, sum time.Duration
	
	ch := make(chan time.Duration, numRequests)
	for i := 0; i < numRequests; i++ {
		go MakeRequest(url, ch)
	}

	for i := 0; i < numRequests; i++ {
		duration := <-ch
		
		if i == 0 {
			min = duration
			max = duration
			sum = duration
		} else {
			if duration < min { min = duration }
			if duration > max { max = duration }
			sum += duration
		}
	}
	
	fmt.Printf("Timing for site: %s\n", url)
	fmt.Printf("Average request: %.2fs\n", sum.Seconds() / float64(numRequests))
	fmt.Printf("Min request: %.2fs\n", min.Seconds())
	fmt.Printf("Max request: %.2fs\n\n", max.Seconds())
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: concurrent.exe [num_requests] [website_url1] [website2...?]")
		return
	}
	
	numRequests, _ := strconv.Atoi(os.Args[1])
	
	for _, url := range os.Args[2:] {
		timeSite(url, numRequests)
	}
}