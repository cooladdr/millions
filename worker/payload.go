package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	_ "os"
	"time"
)

type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

type Payload struct {
	StorageFolder string
}

type Job struct {
	Payload Payload
}

var (
	MaxQueue  = 5 //os.Getenv("MAX_QUEUE")
	JobQueue  chan Job
	MaxLength int64 = 1024 * 1024 * 4
)

func init() {
	JobQueue = make(chan Job, MaxQueue)
}

func (p *Payload) UploadToS3() error {
	storage_path := fmt.Sprintf("%v/%v", p.StorageFolder, time.Now().UnixNano())

	b := new(bytes.Buffer)
	encodeErr := json.NewEncoder(b).Encode(p)
	if encodeErr != nil {
		return encodeErr
	}

	fmt.Println(storage_path)

	return nil
}

func PayloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var content = &PayloadCollection{}
	err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, payload := range content.Payloads {
		job := Job{Payload: payload}
		JobQueue <- job
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", "your request is processed")
}
