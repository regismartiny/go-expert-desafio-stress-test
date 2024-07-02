/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the stress test",
	Long:  `Run the stress test providing the web service URL, the total number of requests and the number of concurrent calls`,
	Run: func(cmd *cobra.Command, args []string) {

		flagUrl, _ := cmd.Flags().GetString("url")
		flagRequests, _ := cmd.Flags().GetInt("requests")
		flagConcurrency, _ := cmd.Flags().GetInt("concurrency")

		runStressTest(flagUrl, flagRequests, flagConcurrency)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	runCmd.PersistentFlags().StringP("url", "u", "", "URL of the service")
	runCmd.PersistentFlags().IntP("requests", "r", 10, "Total number of requests")
	runCmd.PersistentFlags().IntP("concurrency", "c", 10, "Number of simultaneous calls")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type RequestInfo struct {
	startTime  time.Time
	endTime    time.Time
	statusCode int
}

func runStressTest(url string, requests int, concurrency int) {
	fmt.Printf("Running stress test for url: %s, with %d total requests and %d concurrent calls\n", url, requests, concurrency)

	requestsInfo := make([]RequestInfo, 0, requests)

	nConcurrentCalls := concurrency
	if nConcurrentCalls > requests {
		nConcurrentCalls = requests
	}

	nNonConcurrentCalls := requests - nConcurrentCalls

	wg := sync.WaitGroup{}
	wg.Add(nConcurrentCalls + nNonConcurrentCalls)
	startTime := time.Now()

	fmt.Printf("\nStart time: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Println("----------------------------------")

	// concurrent calls
	for range nConcurrentCalls {
		go func() {
			defer wg.Done()
			requestInfo := doRequest(url)
			requestsInfo = append(requestsInfo, requestInfo)
		}()
	}

	// non concurrent calls
	go func() {
		for range nNonConcurrentCalls {
			defer wg.Done()
			requestInfo := doRequest(url)
			requestsInfo = append(requestsInfo, requestInfo)
		}
	}()

	wg.Wait()
	endTime := time.Now()

	totalRequests := len(requestsInfo)

	fmt.Println("----------------------------------")
	fmt.Printf("End time: %s\n", endTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total time: %s\n", endTime.Sub(startTime))
	fmt.Printf("Total number of requests: %d\n", totalRequests)
	requestWithStatus200 := getNumberRequestsWithStatus200(requestsInfo)
	fmt.Printf("Average number of requests per second: %.2f\n", float64(totalRequests)/endTime.Sub(startTime).Seconds())
	fmt.Printf("Number of requests with status 200: %d (%d%%)\n", requestWithStatus200, requestWithStatus200/totalRequests*100)
	fmt.Println("Distribution of other status codes: \n", getFormattedStatusCodeDistribution(getStatusCodeDistribution(requestsInfo), totalRequests))
}

func doRequest(url string) RequestInfo {
	requestInfo := RequestInfo{startTime: time.Now()}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	requestInfo.statusCode = resp.StatusCode

	requestInfo.endTime = time.Now()

	fmt.Printf("Request completed -- %s\n", requestInfo.endTime.Sub(requestInfo.startTime))

	return requestInfo
}

func getNumberRequestsWithStatus200(requestsInfo []RequestInfo) int {
	var requestsWithStatus200 int
	for _, requestInfo := range requestsInfo {
		if requestInfo.statusCode == 200 {
			requestsWithStatus200++
		}
	}
	return requestsWithStatus200
}

func getFormattedStatusCodeDistribution(statusCodeDistribution map[int]int, totalNumberOfRequests int) string {
	formattedStatusCodeDistribution := ""
	for statusCode, count := range statusCodeDistribution {
		formattedStatusCodeDistribution += fmt.Sprintf("HTTP %d: %d (%d%%)\n", statusCode, count, count/totalNumberOfRequests*100)
	}
	return formattedStatusCodeDistribution
}

func getStatusCodeDistribution(requestsInfo []RequestInfo) map[int]int {
	statusCodeDistribution := make(map[int]int)
	for _, requestInfo := range requestsInfo {
		statusCodeDistribution[requestInfo.statusCode]++
	}
	return statusCodeDistribution
}
