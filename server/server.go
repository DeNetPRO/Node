package server

import (
	"bytes"
	"dfile-secondary-node/common"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func StartServer(port string) {

	r := mux.NewRouter()

	r.HandleFunc("/upload/{address}", saveFiles).Methods("POST")
	r.HandleFunc("/download/{address}/{name}", serveFiles).Methods("GET")

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

	fmt.Println("Starting server at port: " + port)

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

func saveFiles(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	address := vars["address"]

	if strings.Trim(address, " ") == "" {
		http.Error(w, "address is not specified", 400)
		return
	}

	accountDir, err := common.GetAccountDirectory()
	if err != nil {
		http.Error(w, "Account directiry is not found", 400)
		return
	}

	err = r.ParseMultipartForm(1 << 20) // maxMemory 32MB
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

		path := filepath.Join(accountDir, address, fh.Filename)

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

func serveFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	name := vars["name"]

	if strings.Trim(address, " ") == "" || strings.Trim(name, " ") == "" {
		http.Error(w, "address or name is not provided", 400)
		return
	}

	accountDir, err := common.GetAccountDirectory()
	if err != nil {
		http.Error(w, "account directiry is not found", 400)
		return
	}

	pathToFile := filepath.Join(accountDir, address, name)
	http.ServeFile(w, r, pathToFile)
}
