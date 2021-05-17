package server_test

import (
	"bytes"
	"dfile-secondary-node/account"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	accAddress     string
	workingDir     string
	firstFileName  = "file.txt"
	secondFileName = "text.doc"
)

func TestMain(m *testing.M) {

	workDir, err := shared.GetOrCreateWorkDir()
	if err != nil {
		log.Fatal(err)
	}

	workingDir = workDir

	address, err := account.CreateAccount("12345")
	if err != nil {
		log.Fatal(err)
	}

	accAddress = address
	server.AccountAddress = address

	exitVal := m.Run()

	err = os.RemoveAll(workDir)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitVal)
}

// tests =======================================

func TestFilesUpload(t *testing.T) {

	rand.Seed(time.Now().Unix())

	fileLen := int64(rand.Intn(10000000))

	b := make([]byte, fileLen)

	firstFilePath := filepath.Join(workingDir, firstFileName)

	f, err := os.Create(firstFilePath)
	if err != nil {
		fmt.Println(err)
	}

	rand.Read(b)

	f.Write(b)
	f.Close()

	fileLen = int64(rand.Intn(10000000))

	b1 := make([]byte, fileLen)

	secondFilePath := filepath.Join(workingDir, secondFileName)

	f1, err := os.Create(secondFilePath)
	if err != nil {
		fmt.Println(err)
	}

	rand.Read(b1)

	f1.Write(b)
	f1.Close()

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	fw, err := writer.CreateFormFile("name", firstFileName)
	if err != nil {
		t.Error(err)
	}

	_, err = io.Copy(fw, bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}

	fw1, err := writer.CreateFormFile("name", secondFileName)
	if err != nil {
		t.Error(err)
	}

	_, err = io.Copy(fw1, bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}

	writer.Close()

	endpoint := "/upload"
	handler := http.HandlerFunc(server.SaveFiles)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqBody.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	fmt.Println(rr)
}

// ====================================================================================

// func TestFileDownload(t *testing.T) {

// 	endpoint := "/file/download/" + secondTestQKey

// 	req, err := http.NewRequest("GET", endpoint, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()

// 	router := mux.NewRouter()
// 	router.HandleFunc("/file/download/{fileKey}", server.FileDownload)
// 	router.ServeHTTP(rr, req)
// }

// ========================================================================================================

func testHandler(t *testing.T, r *strings.Reader, e string, h http.HandlerFunc) *httptest.ResponseRecorder {
	req, err := http.NewRequest("POST", e, r)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	return rr
}
