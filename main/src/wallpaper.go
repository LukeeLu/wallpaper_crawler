package main

import (
	fmt1 "fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func handlee(err error, reason string) {
	if err != nil {
		println(err, reason)
	}

}
func download(url string, filename string) (ok bool) {
	resp, err := http.Get(url)
	handlee(err, "http error")
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	handlee(err, "resp.body")
	filename = "***your path****" + filename
	err = ioutil.WriteFile(filename, bytes, 0666)
	if err != nil {
		return false
	} else {
		return true
	}
}

var (
	chanImageurls chan string
	waitgroup     sync.WaitGroup
	chantask      chan string
	reImag        = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

func downloadimage() {
	for url := range chanImageurls {
		filename := getfilename(url)
		ok := download(url, filename)
		if ok {
			fmt1.Printf("%s download success\n", filename)
		} else {
			fmt1.Printf("%s download fail\n", filename)
		}

	}
	waitgroup.Done()
}
func getfilename(url string) (filename string) {
	lastindex := strings.LastIndex(url, "/")
	filename = url[lastindex+1:]
	timeprefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timeprefix + "_" + filename
	return

}

func checkok() {
	var count int
	for {
		url := <-chantask
		fmt1.Printf("%s finished\n", url)
		count++
		if count == 26 {
			close(chanImageurls)
			break
		}
	}
	waitgroup.Done()

}
func getimageurls(url string) {
	urls := getimages(url)
	for _, url := range urls {
		chanImageurls <- url
	}
	chantask <- url
	waitgroup.Done()

}
func getimages(url string) (urls []string) {
	pagestr := getpagestr(url)
	re := regexp.MustCompile(reImag)
	resultimage := re.FindAllStringSubmatch(pagestr, -1)
	fmt1.Printf("find%dimages\n", len(resultimage))
	for _, data := range resultimage {
		url := data[0]
		urls = append(urls, url)
	}
	return
}
func getpagestr(url string) (pagestr string) {
	resp, err := http.Get(url)
	handlee(err, "http.get error")
	defer resp.Body.Close()
	pagebytes, err := ioutil.ReadAll(resp.Body)
	handlee(err, "ioutil read error")
	pagestr = string(pagebytes)
	return pagestr
}
func main() {
	chanImageurls = make(chan string, 1000000)
	chantask = make(chan string, 26)
	for i := 1; i < 27; i++ {
		waitgroup.Add(1)
		go getimageurls("https://www.bizhizu.cn/shouji/tag-简约/" + strconv.Itoa(i) + ".html")
	}

	waitgroup.Add(1)
	go checkok()

	for i := 0; i < 5; i++ {
		waitgroup.Add(1)
		go downloadimage()
	}
	waitgroup.Wait()
}
