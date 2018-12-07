package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "path/filepath"
  "regexp"
  "strconv"
  "strings"
)

func Scrape() ([]string) {
  var urls []string
  base := "https://standardebooks.org/ebooks/?page="
  i := 0
  flag := false
  for flag == false {
    i++
    flag = true
	url := base + strconv.Itoa(i)
	
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	html, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile("(/ebooks/.*)\">")
	found := re.FindAllStringSubmatch(string(html), -1)
	for f := range found {
	  u := found[f][1]
	  if (!strings.Contains(u, "?") && !strings.Contains(u, "method=") && u != "/ebooks/") {
	    flag = false
	    u := "https://standardebooks.org" + u
		flag2 := false
		for _, v := range urls {
		  if (v == u) {
		    flag2 = true
		  }
		}
		if (flag2 == false) {
		  urls = append(urls, u)
		}
	  }
	}
  }
  
  return urls
}

func Download(urls []string) {
  os.Mkdir("download", os.ModePerm)
  for _, v := range urls {
    resp, _ := http.Get(v)
	defer resp.Body.Close()
	html, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile("(https://standardebooks.org/ebooks/.*)\"")
	found := re.FindAllStringSubmatch(string(html), -1)
	for _, f := range found {
	  if (strings.HasSuffix(f[1], ".epub") || strings.HasSuffix(f[1], ".epub3") || strings.HasSuffix(f[1], ".azw3")) && !strings.Contains(f[1], "kepub") {
	    filename := "download/" + filepath.Base(f[1])
		out, _ := os.Create(filename)
		resp, _ := http.Get(f[1])
		defer resp.Body.Close()
		io.Copy(out, resp.Body)
		out.Close()
		fmt.Println(f[1])
	  }
	}
  }
}

func main() {
  urls := Scrape()
  Download(urls)
}