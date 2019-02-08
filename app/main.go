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
	"strconv"
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

	fmt.Fprintf(file, "$"+user+"\n")
	fmt.Fprintf(file, "$"+password+"\n")
	for i := 0; i < len(config); i++{
		fmt.Fprintf(file, config[i]+"\n")
	}
}

func send_file(cur_conf string, cc string) {

	names := strings.Split(cc,":")

	extraParams := map[string]string{
		"conf_name": cur_conf,
		"user_name": user,
		"pass": password,
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

func get_request(url string)  string {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		return string(contents)
	}
	return ""
}

func insert_into_config(conf string, name string, path string){
	found := false
	for i := 0; i < len(config); i++ {
		if config[i][:1] == "#" && config[i][1:] == conf{
			for j := i+1; j < len(config); j++{
				if strings.Split(config[j], ":")[0] == name{
					config[j] = name + ":" + expand_path(path)
					found = true
					break
				}
			}
			if !found{
				config = append(config, "")
				copy(config[i+2:], config[i+1:])
				config[i+1] = name + ":" + expand_path(path)
				found = true
			}
			break
		}
	}
	if !found {
		config = append(config, "#" + conf)
		config = append(config, name + ":" + expand_path(path))
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
			insert_into_config(os.Args[2], os.Args[3], os.Args[4])
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

		if os.Args[1] == "get" {
			assert(len(os.Args) == 3, "wrong number of arguments \n maybe call syncfg help")

			ur := "http://localhost:5000/getfiles?user_name="+user+"&conf_name="+os.Args[2]+"&pass="+password
			files_str := get_request(ur)
			files := strings.Split(files_str, " ")
			//all_except := []int{};
			x := 1;
			for i := 0; i < len(files); i++{
				if files[i] != ""{
					fmt.Println("["+strconv.Itoa(x)+"] "+strings.Split(files[i], "/")[2] + " ")
					x += 1
				}
			}
			for i := 0; i < len(files); i++{
				if files[i] != ""{
					path := get_request("http://localhost:5000/file/"+files[i]+"/config")
					file_str := get_request("http://localhost:5000/file/"+files[i]+"/file")

					file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()

					fmt.Fprintf(file, file_str)
					path_details := strings.Split(files[i], "/")
					insert_into_config(os.Args[2], path_details[2], path)
				}
			}
			write_config()
		}

		if os.Args[1] == "register"{
			ur := "http://localhost:5000/register?user_name="+os.Args[2]+"&pass="+os.Args[3]
			pass := get_request(ur)
			fmt.Println(pass);
			if pass == "done" {
				user = os.Args[2]
				password = os.Args[3]
				write_config()
			}
		}
	}
}
