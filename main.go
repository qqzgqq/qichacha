package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// ExampleScrape 企查查爬虫股东和上市公司
func ExampleScrape() {
	// go run main.go 传参数
	flag.Parse()
	var test string = flag.Arg(0)
	//connection mysql_database data_acquisition_development
	DB, _ := sqlx.Open("mysql", "root:1@tcp(127.0.0.1:3306)/data_acquisition_development?charset=utf8")

	//goquery解析html
	doc, err := goquery.NewDocument(test)
	if err != nil {
		log.Fatal(err)
	}
	//清空表（guang）info
	DB.MustExec("delete from guang")
	// Find the review items(Find方法中查找类前面用.；查找id等用#) each遍历查找信息
	doc.Find(".tab-content").Each(func(i int, s *goquery.Selection) {
		//查找类中包含active的类中信息
		h := s.Find(".active")
		//each 类active中tr信息
		h.Each(func(i int, s *goquery.Selection) {
			tr := s.Find("tr")
			//遍历tr中信息以a标签中信息
			tr.Each(func(i int, s *goquery.Selection) {
				aNode := s.Find("a").First()
				href, _ := aNode.Attr("href")
				if href == "" {
				} else {
					href2 := "http://www.qichacha.com" + href
					fmt.Println(aNode.Text(), href2)

					DB.MustExec("insert into guang(zhu,www) values(?,?)", aNode.Text(), href2)
				}
			})
		})
	})

	var zhu, ww, zhu2, quan string
	DB.Get(&zhu, "select count(*) from guang")
	zhushu, _ := strconv.Atoi(zhu)
	for ii := 0; ii < zhushu; ii++ {

		DB.Get(&ww, "select www from guang limit ?,1", ii)
		DB.Get(&zhu2, "select zhu from guang limit ?,1", ii)
		doc, err := goquery.NewDocument(ww)

		if err != nil {
			log.Fatal(err)
		}
		//遍历id为Sockinfo中信息
		doc.Find("#Sockinfo").Each(func(i int, s *goquery.Selection) {
			title2 := s.Find("td")
			var aaa string
			var k = 0
			title2.Each(func(i int, s *goquery.Selection) {
				b, _ := s.Attr("width")
				if b == "43%" {
					if k == 0 {
						aaa += s.Find("a").First().Text()
					} else {
						aaa += "、" + s.Find("a").First().Text()
					}
					k++
				}
				DB.MustExec("update guang set chong = ? where www = ?", aaa, ww)

			})

			fmt.Println(zhu2, "--", aaa)
			//strings.Trim方法是去掉zhu2中的字符串的空格
			wancheng := strings.Trim(zhu2, " ") + "--" + aaa
			// 追加wancheng后回车
			quan += wancheng + "\n"

		})
		//停顿6秒执行
		time.Sleep(6 * time.Second)
	}
	//将quan中字符串中追加到output2.txt中
	var d1 = []byte(quan)
	ioutil.WriteFile("./output2.txt", d1, 0666) //写入文件(字节数组)

}

func main() {
	ExampleScrape()

}
