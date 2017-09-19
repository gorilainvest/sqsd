package sqsd

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"net/http"
	"net/http/httptest"
	"time"
)

type MockClient struct {
	sqsiface.SQSAPI
	Resp             *sqs.ReceiveMessageOutput
	RecvRequestCount int
	DelRequestCount  int
	Err              error
}

func NewMockClient() *MockClient {
	return &MockClient{
		Resp: &sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{},
		},
	}
}

func (c *MockClient) ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	c.RecvRequestCount++
	return c.Resp, c.Err
}

func (c *MockClient) DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	c.DelRequestCount++
	return &sqs.DeleteMessageOutput{}, nil
}

func MockJobServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "text")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "no goood")
	})
	mux.HandleFunc("/long", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "text")
		fmt.Fprintf(w, "goood")
		time.Sleep(1 * time.Second)
	})
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "text")
		fmt.Fprintf(w, "goood")
	})
	return httptest.NewServer(mux)
}

type MockResponseWriter struct {
	http.ResponseWriter
	header     http.Header
	ResBytes   []byte
	StatusCode int
	Err        error
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		header:     http.Header{},
		ResBytes:   []byte{},
		StatusCode: http.StatusOK,
	}
}

func (w *MockResponseWriter) Header() http.Header {
	return w.header
}

func (w *MockResponseWriter) Write(b []byte) (int, error) {
	w.ResBytes = b
	return len(b), w.Err
}

func (w *MockResponseWriter) WriteHeader(s int) {
	w.StatusCode = s
}

func (w *MockResponseWriter) ResponseString() string {
	return bytes.NewBuffer(w.ResBytes).String()
}
