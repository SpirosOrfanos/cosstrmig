package main

import (
	"fmt"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	excelize "github.com/xuri/excelize/v2"
	"strings"
)

func createMigration(pi string) {
	englishArticles := parseProducts("en")
	greekArticles := parseProducts("el")
	aggregatedArticles := make(map[string][]Article)
	aggregatedSaveArticles := make(map[string][]HelpArticleRaw)

	for k, v := range greekArticles {
		_, ok := aggregatedArticles[k]
		if !ok {
			aggregatedArticles[k] = make([]Article, 0)
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		} else {
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		}
	}

	for k, v := range englishArticles {
		_, ok := aggregatedArticles[k]
		if !ok {
			aggregatedArticles[k] = make([]Article, 0)
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		} else {
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		}
	}

	for k, v := range aggregatedArticles {
		var helpArticle HelpArticleRaw

		helpArticle = HelpArticleRaw{
			Id:          nil,
			Title:       v[0].Title,
			Keywords:    v[0].Keywords,
			Locale:      v[0].Locale,
			Reusables:   contentHtml(v[0].Title, v[0].ContentHtml),
			ArticleId:   v[0].ArticleId,
			ReferenceId: k,
		}
		aggregatedSaveArticles[k] = make([]HelpArticleRaw, 0)
		aggregatedSaveArticles[k] = append(aggregatedSaveArticles[k], helpArticle)
		if len(v) == 1 {
			continue
		}

		helpArticle2 := HelpArticleRaw{
			Id:          nil,
			Title:       v[1].Title,
			Keywords:    v[1].Keywords,
			Locale:      v[1].Locale,
			Reusables:   contentHtml(v[1].Title, v[1].ContentHtml),
			ArticleId:   v[1].ArticleId,
			ReferenceId: k,
		}

		artV := make([]HelpArticleRaw, 0)
		artV = append(artV, aggregatedSaveArticles[k][0])
		artV = append(artV, helpArticle2)
		aggregatedSaveArticles[k] = artV
	}

	_, categoryArticles := parseArticleCategories()
	_, videoArticles := parseVideoCategories()
	stackEl := lls.New()
	stackEn := lls.New()
	for k, vals := range aggregatedSaveArticles {
		if _, ok2 := videoArticles[k]; ok2 {
			//fmt.Println("video article")
			//continue
		} else {
			if catd, ok := categoryArticles[k]; ok {
				for _, artic := range vals {
					if artic.Locale == "en" {
						addItemToList(k, artic, stackEn, catd)
					} else {
						addItemToList(k, artic, stackEl, catd)
					}

				}
			}
		}

	}

	writeToFile(stackEl, "el", pi)
	writeToFile(stackEn, "en", pi)

}

func createMigrationVideo(pi string) {
	englishArticles := parseProducts("en")
	greekArticles := parseProducts("el")
	aggregatedArticles := make(map[string][]Article)
	aggregatedSaveArticles := make(map[string][]HelpArticleRaw)

	for k, v := range greekArticles {
		_, ok := aggregatedArticles[k]
		if !ok {
			aggregatedArticles[k] = make([]Article, 0)
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		} else {
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		}
	}

	for k, v := range englishArticles {
		_, ok := aggregatedArticles[k]
		if !ok {
			aggregatedArticles[k] = make([]Article, 0)
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		} else {
			aggregatedArticles[k] = append(aggregatedArticles[k], v)
		}
	}

	for k, v := range aggregatedArticles {
		var helpArticle HelpArticleRaw

		helpArticle = HelpArticleRaw{
			Id:          nil,
			Title:       v[0].Title,
			Keywords:    v[0].Keywords,
			Locale:      v[0].Locale,
			Reusables:   contentHtml(v[0].Title, v[0].ContentHtml),
			ArticleId:   v[0].ArticleId,
			ReferenceId: k,
		}
		aggregatedSaveArticles[k] = make([]HelpArticleRaw, 0)
		aggregatedSaveArticles[k] = append(aggregatedSaveArticles[k], helpArticle)
		if len(v) == 1 {
			continue
		}

		helpArticle2 := HelpArticleRaw{
			Id:          nil,
			Title:       v[1].Title,
			Keywords:    v[1].Keywords,
			Locale:      v[1].Locale,
			Reusables:   contentHtml(v[1].Title, v[1].ContentHtml),
			ArticleId:   v[1].ArticleId,
			ReferenceId: k,
		}

		artV := make([]HelpArticleRaw, 0)
		artV = append(artV, aggregatedSaveArticles[k][0])
		artV = append(artV, helpArticle2)
		aggregatedSaveArticles[k] = artV
	}

	//_, categoryArticles := parseArticleCategories()
	_, videoArticles := parseVideoCategories()
	stackEl := lls.New()
	stackEn := lls.New()
	for k, vals := range aggregatedSaveArticles {
		if vss, ok2 := videoArticles[k]; ok2 {
			for _, artic := range vals {
				if artic.Locale == "en" {
					addItemToList(k, artic, stackEn, vss)
				} else {
					addItemToList(k, artic, stackEl, vss)
				}

			}
		} else {

		}

	}

	writeToFile(stackEl, "el", pi)
	writeToFile(stackEn, "en", pi)

}

func writeToFile(stack *lls.Stack, lang string, pi string) {
	filename := lang + "" + pi + ".xlsx"
	f := excelize.NewFile()
	_, err := f.NewSheet("sheet1")
	if err != nil {

	}
	f.SetCellValue("Sheet1", "A1", "doc_key")
	f.SetCellValue("Sheet1", "B1", "name")
	f.SetCellValue("Sheet1", "C1", "title")
	f.SetCellValue("Sheet1", "D1", "contentHtml")
	f.SetCellValue("Sheet1", "E1", "content_type")
	f.SetCellValue("Sheet1", "F1", "keywords")
	f.SetCellValue("Sheet1", "G1", "is_top")
	f.SetCellValue("Sheet1", "H1", "top_position")
	f.SetCellValue("Sheet1", "I1", "categories")
	f.SetCellValue("Sheet1", "J1", "status")
	f.SetCellValue("Sheet1", "K1", "segments")
	f.SetCellValue("Sheet1", "L1", "segment_boost")

	for ind, v := range stack.Values() {
		itm := v.(XlsItem)

		row := ind + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "A", row), itm.doc_key)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "B", row), itm.name)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "C", row), itm.title)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "D", row), itm.content_html)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "E", row), itm.content_type)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "F", row), itm.keywords)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "G", row), itm.is_top)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "H", row), itm.top_position)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "I", row), itm.categories)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "J", row), itm.status)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "K", row), itm.segments)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "L", row), itm.segment_boost)
	}
	f.SaveAs(filename)
}

func contentHtml(title string, body string) string {
	return strings.Trim(title, " ") + " " + strings.Trim(body, " ")
}

func addItemToList(key string, article HelpArticleRaw, stack *lls.Stack, categories []string) {

	stack.Push(XlsItem{
		doc_key:       key,
		name:          article.ArticleId,
		title:         article.Title,
		content_type:  "Raw",
		keywords:      article.Keywords,
		is_top:        "N",
		top_position:  "",
		categories:    strings.Join(categories, ","),
		status:        "published",
		segments:      "",
		segment_boost: "1",
		content_html:  fmt.Sprintf("%s", article.Reusables),
	})
}

type XlsItem struct {
	doc_key       string
	name          string
	title         string
	content_type  string
	keywords      string
	is_top        string
	top_position  string
	categories    string
	status        string
	segments      string
	segment_boost string
	content_html  string
}
