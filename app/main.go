package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"bufio"
)

var cwd, _ =  os.Getwd();
var config []string
var user string = ""
var password string = ""

func assert(condition bool, msg string){
	if !condition{
		fmt.Println(msg)
		os.Exit(1)
	}
}

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
		return path
	} else if path[0:1] == "~" {
		return UserHomeDir() + path[1:]
	}
	return cwd + "/" + path
}

func read_config(){
	file, err := os.OpenFile(UserHomeDir() + "/.syncfgrc", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if len(s) > 1{
			if s[:1] == "$" {
				if user == "" {
					user = s[1:]
				} else {
					password = s[1:]
				}
			} else {
				config = append(config, scanner.Text())
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
func write_config(){
	file, err := os.OpenFile((UserHomeDir() + "/.syncfgrc"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < len(config); i++{
		fmt.Println(config[i])
		fmt.Fprintf(file, config[i]+"\n")
	}
}

func send_file(cur_conf string, cc string) {

	names := strings.Split(cc,":")

	extraParams := map[string]string{
		"conf_name": cur_conf,
		"user_name": user,
		"pretty_name": names[0],
		"file_name": names[1],
	}
	request, err := newfileUploadRequest("http://127.0.0.1:5000/uploader", extraParams, "file", names[1])
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp.StatusCode)
	}
}

func main() {

	if len(os.Args) == 1 || os.Args[1] == "help" {
		fmt.Println("help")
		os.Exit(0)
	}

	if len(os.Args) > 1{
		read_config()
		if os.Args[1] == "add"{
			assert(len(os.Args) == 5, "wrong number of arguments \n maybe call syncfg help")

			found := false
			for i := 0; i < len(config); i++ {
				if config[i][:1] == "#" && config[i][1:] == os.Args[2]{
					for j := i+1; j < len(config); j++{
						if strings.Split(config[j], ":")[0] == os.Args[3] {
							config[j] = os.Args[3] + ":" + expand_path(os.Args[4])
							found = true
							break
						}
					}
					if !found{
						config = append(config, "")
						copy(config[i+2:], config[i+1:])
						config[i+1] = os.Args[3] + ":" + expand_path(os.Args[4])
						found = true
					}
					break
				}
			}
			if !found {
				config = append(config, "#" + os.Args[2])
				config = append(config, os.Args[3] + ":" + expand_path(os.Args[4]))
			}
			write_config()
		}

		if os.Args[1] == "commit" {
			cur_conf := ""
			for i := 0; i < len(config); i++ {
				if config[i][:1] == "#"{
					cur_conf = config[i][1:]
				} else if cur_conf != "" && config[i] != ""{
					send_file(cur_conf, config[i])
				}
			}
		}
	}
}
