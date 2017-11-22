package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crawlerclub/dl"
	"github.com/golang/glog"
	"github.com/liuzl/filestore"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	searchUrlOld1 = `http://search.ccgp.gov.cn/oldsearch?searchtype=1&page_index=`
	searchUrlOld2 = `&bidSort=0&buyerName=&projectId=&pinMu=0&bidType=0&dbselect=bidx&kw=&start_time=%d%%3A01%%3A01&end_time=%d%%3A12%%3A31&timeType=6&displayZone=&zoneId=&agentName=`
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
	urlCh := make(chan string)
	recordCh := make(chan string)
	exitCh := make(chan int)
	go Dispatch(urlCh, *t)
	go SaveRecord(recordCh, exitCh)
	var wg sync.WaitGroup
	for i := 0; i < *j; i++ {
		wg.Add(1)
		go List(urlCh, recordCh, &wg, i)
	}
	wg.Wait()
	close(recordCh)
	<-exitCh
	glog.Info("Done!")
}

func getTotalPages(url string) int {
	req := &dl.HttpRequest{Url: url, Method: "GET", UseProxy: true, Platform: "pc", Retry: 20,
		ValidFuncs: []func(resp *dl.HttpResponse) bool{func(resp *dl.HttpResponse) bool {
			if strings.Contains(resp.Text, "国家级政府采购专业网站") {
				return true
			}
			return false
		}}}
	res := dl.Download(req)
	if res.Error != nil {
		glog.Error(res.Error)
		glog.Error("failed to get page nums")
		glog.Error(url)
		return -1
	} else {
		//re := regexp.MustCompile(`size:\s*(\d+),`)
		re := regexp.MustCompile(`Pager\({\r\n\s*size:\s*(\d+),`)
		ret := re.FindAllStringSubmatch(res.Text, -1)
		if len(ret) <= 0 || len(ret[0]) <= 1 {
			glog.Error("failed to parse page nums")
			glog.Error(url)
			return -1
		}
		pageNum, err := strconv.Atoi(ret[0][1])
		if err != nil {
			glog.Error("failed to parse page nums")
			glog.Error(url)
			return -1
		}
		return pageNum
	}

}

func getUrl(year, page, t int) string {
	var url string
	if t == 1 {
		url = searchUrlBx1 + fmt.Sprintf("%d", page) + searchUrlBx2
	} else {
		url = searchUrlOld1 + fmt.Sprintf("%d", page) + fmt.Sprintf(searchUrlOld2, year, year)
	}
	return url
}

func Dispatch(urlCh chan string, t int) {
	var years []int
	if t == 1 {
		// don't need year when t=1
		years = append(years, 2017)
	} else {
		for year := 2001; year <= 2012; year++ {
			years = append(years, year)
		}
	}
	for _, year := range years {
		initUrl := getUrl(year, 1, t)
		pageNum := getTotalPages(initUrl)
		glog.Info(fmt.Sprintf("get total page of url: %s, %d pages", initUrl, pageNum))
		if pageNum <= 0 {
			glog.Error(fmt.Sprintf("failed to get total page of year: %d", year))
			continue
		}
		for i := 0; i <= pageNum; i++ {
			url := getUrl(year, i, t)
			urlCh <- url
		}
	}
	close(urlCh)
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

func List(urlCh, recordCh chan string, wg *sync.WaitGroup, id int) {
	glog.Info("start worker ", id)
	defer glog.Info("finish worker ", id)
	defer wg.Done()
	for url := range urlCh {
		//url := oldSearch1 + fmt.Sprintf("%d", i) + oldSearch2
		glog.Info(url)
		req := &dl.HttpRequest{Url: url, Method: "GET", UseProxy: true, Platform: "pc", Retry: 10,
			ValidFuncs: []func(resp *dl.HttpResponse) bool{func(resp *dl.HttpResponse) bool {
				if strings.Contains(resp.Text, "国家级政府采购专业网站") {
					return true
				}
				return false
			}}}
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
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func Detail(url string, recordCh chan string) {
	req := &dl.HttpRequest{Url: url, Method: "GET", UseProxy: true, Platform: "pc", Retry: 10,
		ValidFuncs: []func(resp *dl.HttpResponse) bool{func(resp *dl.HttpResponse) bool {
			if strings.Contains(resp.Text, "国家级政府采购专业网站") {
				return true
			}
			return false
		}}}
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
