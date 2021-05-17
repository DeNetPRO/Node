package server

import (
	"bytes"
	"dfile-secondary-node/shared"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var AccountAddress string

func Start(address, port string) {

	AccountAddress = address

	r := mux.NewRouter()

	r.HandleFunc("/upload", SaveFiles).Methods("POST")
	r.HandleFunc("/download/{name}", ServeFiles).Methods("GET")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
		},

		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
	})

	fmt.Println("Dfile node is ready and started listening to port: " + port)

	err := http.ListenAndServe(":"+port, corsOpts.Handler(checkSignature(r)))
	if err != nil {
		panic(err)
	}

}

// ====================================================================================

func checkSignature(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// splittedPath := strings.Split(r.URL.Path, "/")
		// signature := splittedPath[len(splittedPath)-1]
		// splittedPath = splittedPath[:len(splittedPath)-1]
		// reqURL := strings.Join(splittedPath, "/")

		// verified, err := verifySignature(sessionKeyBytes, reqURL, signature)
		// if err != nil {
		// 	http.Error(w, "session key verification error", 500)
		// 	return
		// }

		// if !verified {
		// 	http.Error(w, "wrong session key", http.StatusForbidden)
		// }

		h.ServeHTTP(w, r)
	})
}

// ========================================================================================================

func SaveFiles(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(1 << 20) // maxMemory 32MB
	if err != nil {
		http.Error(w, "Parse multiform problem", 400)
		return
	}

	fhs := r.MultipartForm.File["name"]

	for _, fh := range fhs {
		var buf bytes.Buffer

		fhFile, err := fh.Open()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "File saving problem", 400)

		}
		defer fhFile.Close()

		_, err = io.Copy(&buf, fhFile)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "File saving problem", 400)
		}

		path := filepath.Join(shared.AccDir, AccountAddress, "storage", fh.Filename)

		newFile, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "File saving problem", 400)

		}
		defer newFile.Close()

		_, err = newFile.Write(buf.Bytes())
		if err != nil {
			http.Error(w, "File saving problem", 400)
		}
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// ====================================================================================

func ServeFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if strings.Trim(name, " ") == "" {
		http.Error(w, "address or name is not provided", 400)
		return
	}

	pathToFile := filepath.Join(shared.AccDir, AccountAddress, "storage", name)
	http.ServeFile(w, r, pathToFile)
}
