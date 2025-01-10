package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	dataSizeMB      = 10 // データサイズ（MB）
	numMeasurements = 5  // 測定回数
	threads         = 4  // 並列ダウンロード数
)

// 並列ダウンロード速度測定
func parallelDownload(url string, threads int) float64 {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()
			io.Copy(io.Discard, resp.Body)
		}()
	}

	wg.Wait()
	duration := time.Since(start).Seconds()
	totalData := float64(dataSizeMB*threads) * 8 // データ量（ビット）
	return totalData / (duration * 1024 * 1024)  // Mbpsで返す
}

// 測定結果を分析
func analyzeSpeeds(speeds []float64) (float64, float64) {
	sort.Float64s(speeds)
	return calculateAverage(speeds), calculateMedian(speeds)
}

func calculateAverage(speeds []float64) float64 {
	var total float64
	for _, speed := range speeds {
		total += speed
	}
	return total / float64(len(speeds))
}

func calculateMedian(speeds []float64) float64 {
	mid := len(speeds) / 2
	if len(speeds)%2 == 0 {
		return (speeds[mid-1] + speeds[mid]) / 2
	}
	return speeds[mid]
}

// 結果の表示
func displayResults(speeds []float64, avg, median float64) {
	fmt.Println("\n===== LAN Speed Tester Results =====")
	for i, speed := range speeds {
		fmt.Printf("Measurement %d: %.2f Mbps\n", i+1, speed)
	}
	fmt.Printf("\nAverage Speed: %.2f Mbps\n", avg)
	fmt.Printf("Median Speed: %.2f Mbps\n", median)
	fmt.Println("=====================================")
}

func main() {
	downloadURL := "http://localhost:8080/download"

	// 測定開始
	fmt.Println("Measuring download speed...")
	speeds := make([]float64, numMeasurements)

	for i := 0; i < numMeasurements; i++ {
		speeds[i] = parallelDownload(downloadURL, threads)
		fmt.Printf("Measurement %d complete.\n", i+1)
	}

	// 測定結果の分析
	avg, median := analyzeSpeeds(speeds)

	// 結果の表示
	displayResults(speeds, avg, median)
}
