package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	//getAutoHomeCarSeries()

	hans := "Acc想想"

	// 默认
	a := pinyin.NewArgs()
	a.Separator = ""
	letters := pinyin.Slug(hans, a)
	fmt.Println(letters)

	// 默认
	b := pinyin.NewArgs()
	a.Separator = ""
	letters2 := pinyin.LazyPinyin(hans, b)

	//fmt.Println(letters2[0])
	fmt.Println(letters2)
	// [[zhong] [guo] [ren]]
}

func getAutoHomeCarSeries() {
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	for _, letter := range letters {
		// Request the HTML page.
		res, err := http.Get("https://www.autohome.com.cn/grade/carhtml/" + letter + ".html")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Find the review items
		doc.Find("li").Each(func(i int, li *goquery.Selection) {
			brandCode, _ := li.Parent().Parent().Parent().Attr("id")
			//brandLogo, _ := li.Parent().Parent().Prev().Find("a").Find("img").Attr("src")
			brandName := li.Parent().Parent().Prev().Find("div").Find("a").Text()
			tmpBrandName, _ := GbkToUtf8([]byte(brandName))
			brandName = string(tmpBrandName)
			brandSeries := li.Parent().Prev().Find("a").Text()
			tmpBrandSeries, _ := GbkToUtf8([]byte(brandSeries))
			brandSeries = string(tmpBrandSeries)

			seriesCode, exists := li.Attr("id")
			if seriesCode != "" && exists {
				seriesCode = seriesCode[1:]
				seriesName := li.Find("h4").Find("a").Text()
				tmpSeries, _ := GbkToUtf8([]byte(seriesName))
				seriesName = string(tmpSeries)
				fmt.Printf("firstLetter: %s, brandCode: %s, brandName: %s, brandSeries: %s, seriesCode: %s, seriesName: %s\n", letter, brandCode, brandName, brandSeries, seriesCode, seriesName)
			}
		})
	}
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
