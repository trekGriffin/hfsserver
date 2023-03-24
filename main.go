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

const (
	version = 20230324
)

type progressReader struct {
	reader io.Reader
	total  int64
	count  int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if err != nil {
		return n, err
	}
	pr.count += int64(n)
	fmt.Fprintf(*writer, "Uploaded  %d %%  of %d M\r", (pr.count*100/pr.total + 1), pr.total/1024/1024)
	fmt.Printf("Uploaded  %d %%  of %d M\r", (pr.count*100/pr.total + 1), pr.total/1024/1024)
	return n, nil
}

var writer *http.ResponseWriter

func upload(w http.ResponseWriter, r *http.Request) {

	log.Println("new upload:", r.Header)
	filename := strings.Replace(r.URL.Path, "/upload/", "", 1)

	//save file

	file, err := os.Create(root + filename)
	if err != nil {
		log.Println("create error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	reader := &progressReader{
		reader: r.Body,
		total:  r.ContentLength,
	}
	writer = &w
	// Copy the uploaded data into the new file
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Println("copy error", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("")
	fmt.Println("upload over")
	fmt.Fprintln(w, "")
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
	if os.Args[1] == "-v" {
		log.Println("version:", version)
		return
	}
	if len(os.Args) != 3 {
		log.Fatalf("usage:%s listen-port serving-direcotry, example: %s 7878 d:/", os.Args[0], os.Args[0])
	}
	//check port
	port, err := strconv.ParseInt(os.Args[1], 0, 16)
	if err != nil {
		log.Fatalf("port is not correct %s", os.Args[1])
	}
	//check path
	root = os.Args[2]
	if root == "" || strings.LastIndex(root, "/") != len(root)-1 {
		log.Fatalf("directory is not correct %s or not end with / ", os.Args[1])
	}
	if _, err = os.Stat(root); err != nil {
		log.Fatalf(" check direct stat error:%s", err.Error())
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
