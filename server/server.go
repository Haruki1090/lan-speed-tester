package main

import (
	"fmt"
	"net/http"
)

const dataSizeMB = 10 // ダウンロード用のデータサイズ（MB）

func main() {
	// ダウンロード用エンドポイント
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		data := make([]byte, dataSizeMB*1024*1024) // ダミーデータ作成
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		fmt.Printf("Sent %d MB of data to client\n", dataSizeMB)
	})

	// アップロード用エンドポイント
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		data := make([]byte, r.ContentLength)
		_, err := r.Body.Read(data)
		if err != nil {
			http.Error(w, "Failed to read data", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Received %d bytes from client\n", r.ContentLength)
		w.WriteHeader(http.StatusOK)
	})

	// サーバーの起動
	fmt.Println("Starting server on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
