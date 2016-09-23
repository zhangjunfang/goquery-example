package main

import (
	"bufio"
	//"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	subUrls chan string
	sign    sync.WaitGroup
)

func ExampleScrape() {
	subUrls = make(chan string, 20)
	path := "http://www.oschina.net/"
	doc, err := goquery.NewDocument(path)
	if err != nil {
		log.Fatal(err)
	}
	go doc.Find("#IndustryNews ul li ").Each(func(i int, s *goquery.Selection) {
		sign.Add(1)
		href, _ := s.Find("a").Attr("href")
		subUrls <- (path + href)
	})
	file, _ := os.OpenFile("./aa.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	go func() {
		for {
			select {
			case mm, ok := <-subUrls:
				{
					doc, err := goquery.NewDocument(mm)
					if err != nil {
						log.Fatal(err)
					}
					doc.Find("img").Each(func(i int, g *goquery.Selection) {
						img, _ := g.Attr("src")
						if strings.HasPrefix(img, "/") {
							xx, _ := url.Parse(mm)
							img = xx.Scheme + "://" + xx.Host + img
						} else if strings.HasPrefix(img, "http://") {

						} else {
							img = mm + img
						}
						c, _ := g.Html()
						fmt.Println(c)
						file.WriteString(c)
					})
					sign.Done()
					if !ok && mm == "" {
						fmt.Println("############", "chan is closed ")
					}
					break
				}
			}
		}
	}()

	sign.Wait()
}

func CssCapture() []string {
	path := "http://www.oschina.net/"
	//path := "https://www.jd.com/"
	doc, err := goquery.NewDocument(path)
	if err != nil {
		log.Fatal(err)
	}

	//	doc.Find("style[type='text/css'][rel='stylesheet']").Each(func(i int, g *goquery.Selection) {
	//		css, _ := g.Html()
	//		buff := bufio.NewReader(strings.NewReader(css))
	//		for {
	//			line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
	//			if err != nil || io.EOF == err {
	//				break
	//			}
	//			if !strings.Contains(line, "*") && strings.Contains(line, "background-image") {
	//				//fmt.Print(strings.Split(line, ":")[1]) //可以对一行进行处理
	//				url := strings.Split(line, ":")[1]
	//				if strings.Contains(url, "url('") {
	//					url = string([]rune(url)[(strings.Index(url, "url('") + 5):strings.Index(url, "');")])
	//					//fmt.Println(path + url) //可以对一行进行处理
	//				}
	//				if strings.Contains(url, "url(\"") {
	//					url = string([]rune(url)[(strings.Index(url, "url(\"") + 5):(len(url) - 4)])
	//					//fmt.Println("-----", path+url) //可以对一行进行处理
	//				}

	//			}
	//			//fmt.Print(line) //可以对一行进行处理
	//		}
	//	})
	last := make([]string, 32)
	doc.Find("link[type='text/css']").Each(func(i int, g *goquery.Selection) {
		href, b := g.Attr("href")
		if b {
			temp := strings.Split(href, "?")[0]
			if strings.HasSuffix(temp, "css") {
				csspath := path + temp
				resp, err := http.Get(csspath)
				if err != nil {
					resp.Body.Close()
					return
				}
				buff := bufio.NewReader(resp.Body)
				for {
					line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
					if err != nil || io.EOF == err {
						break
					}
					if !strings.Contains(line, "*") && strings.Contains(line, "background-image") {
						//fmt.Print(strings.Split(line, ":")[1]) //可以对一行进行处理
						url := strings.Split(line, ":")[1]
						if strings.Contains(url, "url('") {
							url = string([]rune(url)[(strings.Index(url, "url('") + 5):strings.Index(url, "');")])
							//fmt.Println(path + url) //可以对一行进行处理
							last = append(last, path+url)
						}
						if strings.Contains(url, "url(\"") {
							url = string([]rune(url)[(strings.Index(url, "url(\"") + 5):(len(url) - 4)])
							fmt.Println("-----", path+url) //可以对一行进行处理
							last = append(last, path+url)
						}

					}
					//fmt.Print(line) //可以对一行进行处理
				}
				resp.Body.Close()
			}
		}
	})
	return last
}

type T struct {
}

func (*T) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*T) Done() <-chan struct{} {
	return nil
}

func (*T) Err() error {
	return nil
}

func (*T) Value(key interface{}) interface{} {
	return nil
}

func StudyContext() {
	var ctx = new(T)
	go func(ct *T) {
		ct.Deadline()
	}(ctx)

	fmt.Println(ctx)
}

func main() {

	if 1 != 2 {
		fmt.Println("-----------------------")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	//CssCapture()
	//ExampleScrape()

	//StudyContext()

}
