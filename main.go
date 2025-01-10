package main

import (
	"fmt"
	"io"
	"net/http"
)

// ダウンロード用データのサイズ（MB）
const dataSizeMB = 10

// ダウンロード速度測定用エンドポイント
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, dataSizeMB*1024*1024) // 10MBのデータを作成
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data) // データをクライアントに送信
}

// アップロード速度測定用エンドポイント
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	bytesReceived, err := io.Copy(io.Discard, r.Body) // クライアントから送信されたデータを無視して読み取る
	if err != nil {
		http.Error(w, "Failed to read upload data", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Received %d bytes\n", bytesReceived) // コンソールに受信したデータ量を出力
	w.WriteHeader(http.StatusOK)
}

func main() {
	// エンドポイントを登録
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/upload", uploadHandler)

	// サーバーを起動
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
