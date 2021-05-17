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
	"testing"
	"time"

	"github.com/gorilla/mux"
)

var (
	accAddress     string
	fileSize       int64
	firstFileName  = "file.txt"
	secondFileName = "text.doc"
)

func TestMain(m *testing.M) {

	shared.CreateIfNotExistAccDirs()

	address, err := account.CreateAccount("12345")
	if err != nil {
		log.Fatal(err)
	}

	accAddress = address
	server.AccountAddress = address

	exitVal := m.Run()

	err = os.RemoveAll(shared.WorkDir)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitVal)
}

// tests =======================================

func TestFilesUpload(t *testing.T) {

	rand.Seed(time.Now().Unix())

	fileSize = int64(rand.Intn(10000000))

	b := make([]byte, fileSize)

	firstFilePath := filepath.Join(shared.WorkDir, firstFileName)

	f, err := os.Create(firstFilePath)
	if err != nil {
		fmt.Println(err)
	}

	rand.Read(b)

	f.Write(b)
	f.Close()

	b1 := make([]byte, fileSize)

	secondFilePath := filepath.Join(shared.WorkDir, secondFileName)

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

	stat, err := os.Stat(filepath.Join(shared.AccDir, accAddress, "storage", firstFileName))
	if err != nil {
		t.Error(err)
	}

	stat1, err := os.Stat(filepath.Join(shared.AccDir, accAddress, "storage", secondFileName))
	if err != nil {
		t.Error(err)
	}

	if stat == nil || stat1 == nil {
		t.Error("files aren't saved")
	}

	if stat.Size() != fileSize || stat1.Size() != fileSize {
		t.Error("files size is incorrect")
	}
}

// ====================================================================================

func TestFileDownload(t *testing.T) {

	endpoint := "/download/" + firstFileName

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/download/{fileKey}", server.ServeFiles)
	router.ServeHTTP(rr, req)
	if int64(rr.Body.Len()) != fileSize {
		t.Error("downloaded file size is wrong")
	}
}
