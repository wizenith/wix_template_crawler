package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lithammer/shortuuid/v3"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	for page := 1; page <= 81; page++ {
		url := fmt.Sprintf("https://www.wix.com/website/templates/html/all/%v", page)

		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		var list []string
		// Find the review items
		doc.Find("source[srcset]:first-child").Each(func(i int, s *goquery.Selection) {

			item_url, ok := s.Attr("srcset")
			if ok {

				list = append(list, regexp.MustCompile(`^(.*?)\/v1`).FindStringSubmatch(item_url)[1])

			}
		})

		base_pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		img_dir := filepath.Join(base_pwd, "img")
		fmt.Println("img_dir", img_dir)

		if _, err := os.Stat(img_dir); os.IsNotExist(err) {
			os.Mkdir(img_dir, 0755)
		}

		for _, item_url := range list {
			// fmt.Println("itme:", item_url)

			res, err := http.Get("https:" + item_url)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			_ = item_url
			u := shortuuid.New()

			// img_dir := fmt.Sprintf("/root/golab/img/%v.jpg", u)
			img_path := fmt.Sprintf("%v/%v.jpg", img_dir, u)
			fmt.Println("img_path", img_path)

			file, err := os.Create(img_path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			// Save file as a image
			// Method I:
			// file.ReadFrom(res.Body)

			// Method II:
			_, err = io.Copy(file, res.Body)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Success!")
		}

	}

}
