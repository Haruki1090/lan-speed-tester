package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

const dataSizeMB = 10     // データサイズ（MB）
const numMeasurements = 5 // 測定回数

// 並列ダウンロード速度測定
func parallelDownload(url string, threads int) float64 {
	var wg sync.WaitGroup // ゴルーチンの完了を待つためのWaitGroup
	start := time.Now()   // 開始時刻

	for i := 0; i < threads; i++ {
		wg.Add(1) // ゴルーチンの数だけWaitGroupに追加
		go func() {
			defer wg.Done() // ゴルーチン終了時にWaitGroupをデクリメント
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close() // 関数終了時にリソースを解放
			io.Copy(io.Discard, resp.Body)
		}()
	}

	wg.Wait()
	duration := time.Since(start).Seconds()
	totalData := float64(dataSizeMB*threads) * 8 // ダウンロードデータ量（ビット）
	return totalData / (duration * 1024 * 1024)  // ダウンロード速度（Mbps）
}

// 測定結果を平均値・中央値で分析
func analyzeSpeed(speeds []float64) (float64, float64) {
	sort.Float64s(speeds) // 昇順ソート
	avg := calculateAverage(speeds)
	median := calculateMedian(speeds)
	return avg, median
}

// 平均値を計算
func calculateAverage(speeds []float64) float64 {
	var total float64
	for _, speed := range speeds {
		total += speed
	}
	return total / float64(len(speeds))
}

// 中央値を計算
func calculateMedian(speeds []float64) float64 {
	mid := len(speeds) / 2
	if len(speeds)%2 == 0 {
		return (speeds[mid-1] + speeds[mid]) / 2
	}
	return speeds[mid]
}

func main() {
	// サーバーURL
	downloadURL := "http://localhost:8080/download"

	// 測定開始
	fmt.Println("Measuring download speed...")
	speeds := make([]float64, numMeasurements)

	// 測定
	for i := 0; i < numMeasurements; i++ {
		speeds[i] = parallelDownload(downloadURL, 4) // 並列処理でダウンロード速度を測定
		fmt.Printf("Measurement %d: %.2f Mbps\n", i+1, speeds[i])
	}

	// 測定結果の分析
	avg, median := analyzeSpeed(speeds)
	fmt.Printf("\nAvarage Speed: %.2f Mbps\n", avg)
	fmt.Printf("Median Speed: %.2f Mbps\n", median)
}
