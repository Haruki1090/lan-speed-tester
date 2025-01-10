package main

import (
	"fmt"
	"io"
	"net/http"
)

const dataSizeMB = 10 // ダウンロード用データサイズ（MB）

// ヘルスチェックエンドポイント
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

// ダウンロード用エンドポイント
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, dataSizeMB*1024*1024) // ダミーデータ
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// アップロード用エンドポイント
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	uploadedBytes, err := io.Copy(io.Discard, r.Body) // アップロードデータを捨てつつサイズを測定
	if err != nil {
		http.Error(w, "Failed to read upload data", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Received %d bytes\n", uploadedBytes)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/upload", uploadHandler)

	port := ":8080"
	fmt.Printf("Starting server at %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
