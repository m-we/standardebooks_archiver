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

func isin(list []string, element string) bool {
  for _, v := range list {
    if element == v {
      return true;
    }
  }
  return false;
}

func find() []string {
  baseUrl := "https://standardebooks.org/ebooks/?page="
  var urls []string
  page := 1

  for true {
    fmt.Print("Checking page " + strconv.Itoa(page) + "\r")
    urls_len := len(urls)
    resp, _ := http.Get(baseUrl + strconv.Itoa(page))
    defer resp.Body.Close()
    html, _ := ioutil.ReadAll(resp.Body)

    re := regexp.MustCompile("(/ebooks/.*/.+?)\"")
    for _, url := range re.FindAllStringSubmatch(string(html), -1) {
      urlSplit := "https://standardebooks.org" + strings.Split(url[0], "\"")[0]
      if !isin(urls, urlSplit) {
        urls = append(urls, urlSplit)
      }
    }

    if urls_len == len(urls) {
      break
    }
    page += 1
  }

  fmt.Println("\n\tdone\n")
  return urls
}

func download(urls []string, frmt string, folder string) {
  _, err := os.Stat(folder)
  if os.IsNotExist(err) {
    os.MkdirAll(folder, os.ModePerm)
  }

  for _, url := range urls {
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    html, _ := ioutil.ReadAll(resp.Body)
    text := strings.Replace(string(html), "\\", "", -1)

    re := regexp.MustCompile("(https://standardebooks.org/ebooks/.*/.+?)\"")
    for _, dLink := range re.FindAllStringSubmatch(text, -1) {
      dUrl := strings.Replace(dLink[0], "\"", "", -1)
      if strings.HasSuffix(dUrl, frmt) {
        if (frmt == "epub" || frmt == ".epub") && (strings.Contains(dUrl, "kepub") || strings.Contains(dUrl, "_advanced")) {
            continue
        }
        fn := folder + filepath.Base(dUrl)

        _, err := os.Stat(fn)
        if os.IsNotExist(err) {
          fmt.Println(dUrl)
          out, _ := os.Create(fn)
          resp, _ := http.Get(dUrl)
          defer resp.Body.Close()
          io.Copy(out, resp.Body)
          out.Close()
        }
      }
    }
  }
}

func main() {
  download(find(), os.Args[1], os.Args[2] + "/")
}
