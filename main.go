package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	root string
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
	var address, root string
	flag.StringVar(&address, "a", "10.10.10.3:80", "address for this host")
	flag.StringVar(&root, "d", "./", "directory for this hfs")
	flag.Parse()
	log.Printf(" server info{address:%s directory:%s}", address, root)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/delete/", delete)
	http.Handle("/", http.FileServer(http.Dir(root)))

	log.Println("server is running")

	err := http.ListenAndServe(address, nil)

	if err != nil {
		log.Fatal("server listen error", err)
	}

}
