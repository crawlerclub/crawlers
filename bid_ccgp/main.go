package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crawlerclub/dl"
	"github.com/golang/glog"
	"github.com/liuzl/filestore"
	"regexp"
	"sync"
	"time"
)

const (
	searchUrlOld1 = `http://search.ccgp.gov.cn/oldsearch?searchtype=1&page_index=`
	searchUrlOld2 = `&bidSort=0&buyerName=&projectId=&pinMu=0&bidType=0&dbselect=bidx&kw=&start_time=2001%3A10%3A10&end_time=2012%3A12%3A31&timeType=6&displayZone=&zoneId=&agentName=`
	searchUrlBx1  = `http://search.ccgp.gov.cn/bxsearch?searchtype=1&page_index=`
	searchUrlBx2  = `&bidSort=0&buyerName=&projectId=&pinMu=0&bidType=0&dbselect=bidx&kw=&start_time=2013%3A01%3A01&end_time=2017%3A08%3A17&timeType=6&displayZone=&zoneId=&pppStatus=0&agentName=`
)

var (
	start = flag.Int("start", 1, "start page number")
	end   = flag.Int("end", 3034, "end page number")
	j     = flag.Int("j", 100, "thread count")
	out   = flag.String("out", "./data", "output dir")
	t     = flag.Int("t", 0, "search type: 0 oldsearch, 1 bxsearch")
)

var (
	ReLink = regexp.MustCompile(`(?ims)href="(/oldsearch/detail\?docId=\d+?)"`)
)

type Record struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	Html string `json:"html"`
}

func main() {
	flag.Parse()
	if *j <= 0 || *j > 1000 {
		glog.Error("thread count must be between 1 and 1000")
		return
	}
	pageIdCh := make(chan int)
	recordCh := make(chan string)
	exitCh := make(chan int)
	go Dispatch(pageIdCh, *start, *end)
	go SaveRecord(recordCh, exitCh)
	var wg sync.WaitGroup
	for i := 0; i < *j; i++ {
		wg.Add(1)
		go List(pageIdCh, recordCh, &wg, i)
	}
	wg.Wait()
	close(recordCh)
	<-exitCh
	glog.Info("Done!")
}

func Dispatch(pageIdCh chan int, start, end int) {
	for i := start; i <= end; i++ {
		pageIdCh <- i
	}
	close(pageIdCh)
}

func SaveRecord(recordCh chan string, exitCh chan int) {
	fs, err := filestore.NewFileStore(*out)
	if err != nil {
		glog.Error(err)
	} else {
		defer fs.Close()
		for record := range recordCh {
			fs.WriteLine([]byte(record))
		}
	}
	exitCh <- 0
}

func List(pageIdCh chan int, recordCh chan string, wg *sync.WaitGroup, id int) {
	glog.Info("start worker ", id)
	defer glog.Info("finish worker ", id)
	defer wg.Done()
	for i := range pageIdCh {
		url := searchUrlOld1 + fmt.Sprintf("%d", i) + searchUrlOld2
		if *t == 1 {
			url = searchUrlBx1 + fmt.Sprintf("%d", i) + searchUrlBx2
		}
		glog.Info(url)
		req := &dl.HttpRequest{Url: url, Method: "GET", UseProxy: false, Platform: "pc"}
		res := dl.Download(req)
		if res.Error != nil {
			glog.Error(res.Error)
			glog.Error(url)
		} else {
			rec := &Record{Url: url, Type: "list", Html: res.Text}
			line, err := json.Marshal(rec)
			if err != nil {
				glog.Error(err)
			} else {
				recordCh <- string(line)
			}
			ret := ReLink.FindAllStringSubmatch(res.Text, -1)
			for _, link := range ret {
				u, err := MakeAbsoluteUrl(link[1], url)
				if err != nil {
					glog.Error(err)
					continue
				}
				Detail(u, recordCh)
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func Detail(url string, recordCh chan string) {
	req := &dl.HttpRequest{Url: url, Method: "GET", UseProxy: false, Platform: "pc"}
	res := dl.Download(req)
	if res.Error != nil {
		glog.Error(res.Error)
		glog.Error(url)
	} else {
		rec := &Record{Url: url, Type: "detail", Html: res.Text}
		line, err := json.Marshal(rec)
		if err != nil {
			glog.Error(err)
		} else {
			recordCh <- string(line)
		}
	}
}
