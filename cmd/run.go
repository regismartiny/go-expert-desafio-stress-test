/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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

func runStressTest(url string, requests int, concurrency int) {
	fmt.Printf("Running stresst test for url: %s, with %d total requests and %d concurrent calls\n", url, requests, concurrency)
}
