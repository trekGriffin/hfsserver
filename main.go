package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

var (
	root      string
	startTime int64
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
	pr.count += int64(n)
	duration := time.Now().Unix() - startTime
	if duration != 0 {
		speed := bytefmt.ByteSize(uint64(pr.count / duration))
		output := fmt.Sprintf("Uploaded  %d %%  %s of %s, speed:%s duration:%s  \r",
			(pr.count * 100 / pr.total), bytefmt.ByteSize(uint64(pr.count)), bytefmt.ByteSize(uint64(pr.total)), speed, time.Duration(duration*1000*1000*1000))
		fmt.Println(output)
		fmt.Fprint(*writer, output)
	}

	return n, err
}

var writer *http.ResponseWriter

func upload(w http.ResponseWriter, r *http.Request) {

	filename := strings.Replace(r.URL.Path, "/upload/", "", 1)
	log.Println("new upload file name", filename)
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
	startTime = time.Now().Unix()
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Println("copy error", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("")
	fmt.Println("upload over")
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "File uploaded successfully!")

}

func delete(w http.ResponseWriter, r *http.Request) {
	filename := strings.Replace(r.URL.Path, "/delete/", "", 1)

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
	showVersion := false
	var address string
	flag.BoolVar(&showVersion, "v", false, "show verison")
	flag.StringVar(&address, "a", "10.10.10.3:80", "bind address and port")
	flag.StringVar(&root, "p", "./", " serving directory")

	flag.Parse()
	if showVersion {
		fmt.Println("version:", version)
		return
	}

	//check path

	if _, err := os.Stat(root); err != nil {
		log.Fatalf(" check direct stat error:%s", err.Error())
	}

	log.Printf(" server info{port:%s directory:%s}", address, root)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/delete/", delete)
	http.Handle("/", http.FileServer(http.Dir(root)))

	log.Println("server is running")

	log.Fatal(http.ListenAndServe(address, nil))

}
