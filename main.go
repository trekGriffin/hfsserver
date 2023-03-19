package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	root string
)

func upload(w http.ResponseWriter, r *http.Request) {

	log.Println("all headers", r.Header)
	filename := strings.Replace(r.URL.Path, "/upload/", "", 1)

	//save file

	file, err := os.Create(root + filename)
	if err != nil {
		log.Println("create error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Copy the uploaded data into the new file
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Println("copy error", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	fmt.Fprintln(w, "File uploaded successfully!")

}

func delete(w http.ResponseWriter, r *http.Request) {
	filename := strings.Replace(r.URL.Path, "/delete/", "", 1)
	// if _, err := os.Stat(root + path); err != nil {
	// 	fmt.Fprintf(w, "delete %s error: file not exist", root+path)
	// 	return
	// }

	err := os.Remove(root + filename)
	if err != nil {
		log.Println("delete  error: ", filename, err.Error())
		fmt.Fprintf(w, "delete %s error: %s", filename, err.Error())
		return
	}
	//todo delete file operation
	fmt.Fprintf(w, " %s delted success", filename)
}
func main() {
	if len(os.Args) != 3 {

		log.Fatalf("usage:%s listen-port serving-direcotry, example: %s 7878 d:/", os.Args[0], os.Args[0])
	}
	port, err := strconv.ParseInt(os.Args[1], 0, 16)
	if err != nil {
		log.Fatalf("port is not correct %s", os.Args[1])
	}
	root = os.Args[2]
	if root == "" {
		log.Fatalf("directory is not correct %s", os.Args[1])
	}

	log.Printf(" server info{port:%d directory:%s}", port, root)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/delete/", delete)
	http.Handle("/", http.FileServer(http.Dir(root)))

	log.Println("server is running")

	err = http.ListenAndServe(":"+strconv.FormatInt(port, 10), nil)

	if err != nil {
		log.Fatal("server listen error", err)
	}

}
