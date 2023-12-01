package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type item struct {
	Photo    string `json:"photo"`
	Title    string `json:"title"`
	Gender   string `json:"gender"`
	Price    string `json:"price"`
	Discount string `json:"discount"`
}

func main() {
	// Permitir domínios para realizar o Scrap
	call := colly.NewCollector(
		colly.AllowedDomains("https://www.amazon.com.br/", "amazon.com.br", "https://www.amazon.com.br/gp/bestsellers/?ref_=nav_cs_bestsellers"),
	)

	// Array de Itens
	var items []item

	// Busca dos dados pelos elementos
	call.OnHTML("div.a-carousel-card", func(h *colly.HTMLElement) {
		item := item{
			Photo:    h.ChildAttr("img.a-dynamic-image", "`src`"),
			Title:    h.ChildText("p.Typography-styled__StyledParagraph-sc-8f9244e7-2 "),
			Gender:   h.ChildText("p.Typography-styled__StyledParagraph-sc-8f9244e7-2"),
			Price:    h.ChildText("p.Typography-styled__StyledParagraph-sc-8f9244e7-2"),
			Discount: h.ChildText("p.Typography-styled__StyledParagraph-sc-8f9244e7-2"),
		}

		items = append(items, item)
	})

	// // Buscar a próxima página
	// call.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
	// 	next_page := h.Request.AbsoluteURL(h.Attr("href"))
	// 	call.Visit(next_page)
	// })

	call.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	// Fetch na URL
	err := call.Visit("https://www.amazon.com.br/gp/bestsellers/?ref_=nav_cs_bestsellers")
	if err != nil {
		log.Fatal(err)
	}

	// Transformar os dados do Web Scraping em JSON
	content, err := json.Marshal(items)

	if err != nil {
		log.Fatal(err)
	}

	// Vai gerar os produtos encontrados dentro um arquivo products.json internamente
	os.WriteFile("products.json", content, 0644)
}
