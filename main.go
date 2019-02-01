package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"bufio"
)

var cwd, _ =  os.Getwd();

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", uri, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	return r, err
}

func expand_path(path string) string{
	if path[0:1] == "/" {
		return path;
	} else if path[0:1] == "~" {
		return UserHomeDir() + path[1:];
	}
	return cwd + "/" + path
}

func main() {
	fmt.Println(cwd);
	file, err := os.Open("/home/jan/pp")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1{
		if os.Args[1] == "add"{

		}
	}

	extraParams := map[string]string{}
	request, err := newfileUploadRequest("http://127.0.0.1:5000/uploader", extraParams, "file", "/home/jan/pp")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
}
