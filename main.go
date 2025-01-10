package main

import (
	"bytes"
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
	threads         = 4  // 並列ダウンロード/アップロード数
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

// 並列アップロード速度測定
func parallelUpload(url string, threads int) float64 {
	var wg sync.WaitGroup
	start := time.Now()

	data := bytes.Repeat([]byte("A"), dataSizeMB*1024*1024) // アップロード用のダミーデータ

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(url, "application/octet-stream", bytes.NewReader(data))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()
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
func displayResults(testType string, speeds []float64, avg, median float64) {
	fmt.Printf("\n===== %s Speed Test Results =====\n", testType)
	for i, speed := range speeds {
		fmt.Printf("Measurement %d: %.2f Mbps\n", i+1, speed)
	}
	fmt.Printf("\nAverage Speed: %.2f Mbps\n", avg)
	fmt.Printf("Median Speed: %.2f Mbps\n", median)
	fmt.Println("===================================")
}

func main() {
	downloadURL := "http://localhost:8080/download"
	uploadURL := "http://localhost:8080/upload"

	// ダウンロード速度測定
	fmt.Println("Measuring download speed...")
	downloadSpeeds := make([]float64, numMeasurements)
	for i := 0; i < numMeasurements; i++ {
		downloadSpeeds[i] = parallelDownload(downloadURL, threads)
		fmt.Printf("Download Measurement %d complete.\n", i+1)
	}
	downloadAvg, downloadMedian := analyzeSpeeds(downloadSpeeds)
	displayResults("Download", downloadSpeeds, downloadAvg, downloadMedian)

	// アップロード速度測定
	fmt.Println("\nMeasuring upload speed...")
	uploadSpeeds := make([]float64, numMeasurements)
	for i := 0; i < numMeasurements; i++ {
		uploadSpeeds[i] = parallelUpload(uploadURL, threads)
		fmt.Printf("Upload Measurement %d complete.\n", i+1)
	}
	uploadAvg, uploadMedian := analyzeSpeeds(uploadSpeeds)
	displayResults("Upload", uploadSpeeds, uploadAvg, uploadMedian)
}
