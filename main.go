package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	URL               = "http://srv.msk01.gigacorp.local/_stats"
	loadAverageLimit  = 30.0
	ramUsageLimit     = 0.8
	diskSpaceLimit    = 0.9
	networkUsageLimit = 0.9
	maxErrorCount     = 3
)

func main() {
	errorCount := 0

	for {
		responce, err := http.Get(URL)
		if err != nil || responce.StatusCode != http.StatusOK {
			errorCount++
			if errorCount >= maxErrorCount {
				fmt.Println("Unable to fetch server statistic")
				break
			}
			time.Sleep(10 * time.Second)
			continue
		}

		body, err := io.ReadAll(responce.Body)
		responce.Body.Close()
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		stats := strings.Split(string(body), ",")
		if len(stats) < 7 {
			fmt.Println("Incorrect data format")
			errorCount++
			continue
		}

		errorCount = 0

		loadAvg, _ := strconv.ParseFloat(stats[0], 64)
		totalRAM, _ := strconv.ParseFloat(stats[1], 64)
		usedRAM, _ := strconv.ParseFloat(stats[2], 64)
		totalDisk, _ := strconv.ParseFloat(stats[3], 64)
		usedDisk, _ := strconv.ParseFloat(stats[4], 64)
		totalNet, _ := strconv.ParseFloat(stats[5], 64)
		usedNet, _ := strconv.ParseFloat(stats[6], 64)

		if loadAvg > loadAverageLimit {
			fmt.Printf("Load Average is too high: %d\n", int(loadAvg))
		}

		ramUsage := usedRAM / totalRAM
		if ramUsage > ramUsageLimit {
			fmt.Printf("Memory usage too high: %d%%\n", int(ramUsage*100))
		}

		diskFree := (totalDisk - usedDisk) / (1024 * 1024)
		if (totalDisk - usedDisk) < (1-diskSpaceLimit)*totalDisk {
			fmt.Printf("Free disk space is too low: %d Mb left\n", int(diskFree))
		}

		netFree := (totalNet - usedNet) / (1000 * 1000)
		if usedNet > networkUsageLimit*totalNet {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(netFree))
		}

		time.Sleep(10 * time.Second)
	}
}
