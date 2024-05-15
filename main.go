package main

import (
	"encoding/xml"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/thedatashed/xlsxreader"
	"go.reizu.org/greek"
	html "golang.org/x/net/html"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	saveUs()
}

func saveUs() {
	englishArticles := parseProducts("en")
	greekArticles := parseProducts("el")
	aggregatedArticles := make(map[string][]Article)
	aggregatedSaveArticles := make(map[string][]HelpArticle)

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

	counter := 0
	for k, v := range aggregatedArticles {
		if counter == 100 {
			break
		}
		counter++
		var helpArticle HelpArticle
		if strings.Contains(v[0].ContentHtml, "toggle-content") {
			wraper := handleToggle(v[0].ContentHtml)
			mItems := make([]ToggleItem, 0)
			for _, v := range wraper.items {
				mItems = append(mItems, ToggleItem{
					Title: v.title,
					Body:  v.items,
				})
			}
			toggle := Toggle{
				Description: wraper.description,
				Component:   "reusables.toggles",
				Items:       mItems,
			}
			toggles := []Toggle{toggle}

			helpArticle = HelpArticle{
				Id:          nil,
				Title:       v[0].Title,
				Keywords:    v[0].Keywords,
				Locale:      v[0].Locale,
				Reusables:   toggles,
				ArticleId:   v[0].ArticleId,
				ReferenceId: k,
			}

		} else {
			itm := Reusable{
				Title:     v[0].Title,
				Body:      v[0].ContentHtml,
				Component: "reusables.html-reusable",
			}
			inter := []Reusable{itm}
			helpArticle = HelpArticle{
				Id:          nil,
				Title:       v[0].Title,
				Keywords:    v[0].Keywords,
				Locale:      v[0].Locale,
				Reusables:   inter,
				ArticleId:   v[0].ArticleId,
				ReferenceId: k,
			}
		}
		aggregatedSaveArticles[k] = make([]HelpArticle, 0)
		aggregatedSaveArticles[k] = append(aggregatedSaveArticles[k], helpArticle)
		if len(v) == 1 {
			continue
		}
		if strings.Contains(v[1].ContentHtml, "toggle-content") {

			wraper := handleToggle(v[1].ContentHtml)
			mItems := make([]ToggleItem, 0)
			for _, v := range wraper.items {
				mItems = append(mItems, ToggleItem{
					Title: v.title,
					Body:  v.items,
				})
			}
			toggle := Toggle{
				Description: wraper.description,
				Component:   "reusables.toggles",
				Items:       mItems,
			}
			toggles := []Toggle{toggle}
			helpArticle = HelpArticle{
				Id:          nil,
				Title:       v[1].Title,
				Keywords:    v[1].Keywords,
				Locale:      v[1].Locale,
				Reusables:   toggles,
				ArticleId:   v[1].ArticleId,
				ReferenceId: k,
			}

		} else {
			itm := Reusable{
				Title:     v[1].Title,
				Body:      v[1].ContentHtml,
				Component: "reusables.html-reusable",
			}
			inter := []Reusable{itm}
			helpArticle = HelpArticle{
				Id:          nil,
				Title:       v[1].Title,
				Keywords:    v[1].Keywords,
				Locale:      v[1].Locale,
				Reusables:   inter,
				ArticleId:   v[1].ArticleId,
				ReferenceId: k,
			}
		}

		artV := make([]HelpArticle, 0)
		artV = append(artV, aggregatedSaveArticles[k][0])
		artV = append(artV, helpArticle)
		aggregatedSaveArticles[k] = artV
	}
	adapter := NewAdapter()
	savedArticles := make(map[string][]HelpArticle)
	for key, rts := range aggregatedSaveArticles {
		newArticles := make([]HelpArticle, 0)
		res, _ := adapter.Create(rts[0])
		newArticles = append(newArticles, HelpArticle{
			Id:          &res.Data.Id,
			ReferenceId: rts[0].ReferenceId,
			Locale:      rts[0].Locale,
		})
		if len(rts) > 1 {
			res2, _ := adapter.Localization(rts[1], res.Data.Id)
			newArticles = append(newArticles, HelpArticle{
				Id:          &res2.Id,
				ReferenceId: rts[1].ReferenceId,
				Locale:      rts[1].Locale,
			})
		}
		savedArticles[key] = newArticles
	}
	categories, _ /*categoryArticles*/ := parseArticleCategories()
	categoriesSavedEl := make(map[string]int)
	categoriesSavedEn := make(map[string]int)
	for cat, catv := range categories {
		nc := HelpArticleCategory{
			Name:        catv.name,
			CategoryId:  catv.code,
			Description: catv.name,
			Locale:      "el",
		}
		res, _ := adapter.CreateHelpArticleCategory(nc)
		categoriesSavedEl[cat] = res.Data.Id
		ncEn := HelpArticleCategory{
			Name:        catv.name,
			CategoryId:  catv.code,
			Description: catv.name,
			Locale:      "en",
		}
		res2, _ := adapter.LocalizationsHelpArticleCategory(ncEn, res.Data.Id)
		categoriesSavedEn[cat] = res2.Id
	}

	savedArticlesEl := make(map[string]HelpArticle)
	savedArticlesEn := make(map[string]HelpArticle)
	for k, v := range savedArticles {
		for _, ar := range v {
			if ar.Locale == "el" {
				savedArticlesEl[k] = ar
			}
			if ar.Locale == "en" {
				savedArticlesEn[k] = ar
			}
		}
	}
	for cat, catv := range categories {

		if catv.childCategories != nil || catv.parentCategories != nil {
			elId, _ := categoriesSavedEl[cat]
			enId, _ := categoriesSavedEn[cat]
			elIdsP := make([]string, 0)
			enIdsP := make([]string, 0)

			elIdsC := make([]string, 0)
			enIdsC := make([]string, 0)
			for _, v := range catv.parentCategories.Values() {
				if v1, ok := categoriesSavedEl[fmt.Sprintf("%s", v)]; ok {
					elIdsP = append(elIdsP, fmt.Sprintf("%d", v1))
				}
				if v1, ok := categoriesSavedEn[fmt.Sprintf("%s", v)]; ok {
					enIdsP = append(enIdsP, fmt.Sprintf("%d", v1))
				}
			}

			for _, v := range catv.childCategories.Values() {
				if v1, ok := categoriesSavedEl[fmt.Sprintf("%s", v)]; ok {
					elIdsC = append(elIdsC, fmt.Sprintf("%d", v1))
				}
				if v1, ok := categoriesSavedEn[fmt.Sprintf("%s", v)]; ok {
					enIdsC = append(enIdsC, fmt.Sprintf("%d", v1))
				}
			}

			adapter.UpdateHelpArticleCategory(HelpArticleCategory{
				ChildCategories:  elIdsC,
				ParentCategories: elIdsP,
			}, elId)
			adapter.UpdateHelpArticleCategory(HelpArticleCategory{
				ChildCategories:  enIdsC,
				ParentCategories: enIdsP,
			}, enId)
			/*for _, v := range catv.parentCategories.Values() {
				fmt.Println(">", v)
			}*/
			//HelpArticleCategory{ParentCategories: }
		}
	}

	/*for cat, catv := range categories {
		catv.childCategories
	}*/
	/*for rtId, cats := range categoryArticles {
		if svd, ok := savedArticles[rtId]; ok {
			for _, cat := range cats {
				//fmt.Println(vat, svd, categories[vat].name)
				if itm, ok1 := categories[cat]; ok1 {
					if elItm, ok2 := categoriesSavedEl[cat]; ok2 {

					}
					if enItm, ok3 := categoriesSavedEn[cat]; ok3 {

					}

				}

			}
		}
	}*/
	/*categories, categoryArticles := parseArticleCategories()
	videos, videoArticles := parseVideoCategories()
	englishArticles := parseProducts("en")
	greekArticles := parseProducts("el")*/

	// 1 create englishArticles
	// 2 from categoryArticles enhance with id for each
}

func handleToggle(data string) TGItemWrapper {

	toggleDescription := ""
	toggleItems := make([]TGItem, 0)
	wrapper := TGItemWrapper{}
	myReader := strings.NewReader(data)
	z := html.NewTokenizer(myReader)
	toggleStart := false
	tTitle := ""
	titleFound := false
	subFound := false
	subDesc := ""
	counter := 0
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if len(subDesc) > 0 {
				//fmt.Println(">>" + subDesc)
				toggleItems[counter-1].items = subDesc
			}
			wrapper.description = toggleDescription
			wrapper.items = toggleItems
			//fmt.Println("----> ", len(toggleItems), toggleDescription)
			return wrapper

		}

		token := z.Token()
		if strings.Contains(token.String(), "toggle-content") {
			toggleStart = true
		}
		if !toggleStart {
			if token.Type == 2 {
				toggleDescription = toggleDescription + fmt.Sprintf("%s", token)
			} else if token.Type == 3 {
				toggleDescription = toggleDescription + fmt.Sprintf("%s", token)
			} else if token.Type == 1 {
				toggleDescription = toggleDescription + token.Data
			}
		} else {
			if strings.Contains(token.String(), "toggle-content") {
				titleFound = true
				subFound = false
				toggleItems = append(toggleItems, TGItem{
					title: "",
					items: "",
				})
			}

			if strings.Contains(token.String(), "hidden-content") {
				subFound = true
				titleFound = false
				subDesc = ""
			}

			if subFound {
				subDesc = subDesc + fmt.Sprintf("%s", token)
			}
			if titleFound && token.Type == 1 {
				tTitle = token.Data
				if len(subDesc) > 0 {
					//fmt.Println(">>" + subDesc)
					toggleItems[counter-1].items = subDesc
				}
				toggleItems[counter].title = tTitle
				counter++
				titleFound = false
				subFound = false
			}
		}
	}
	//fmt.Println(wrapper)
	return wrapper
}

func parseArticleCategories() (map[string]Category, map[string][]string) {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	f, err := xlsxreader.OpenFile("articles.xlsx")
	if err != nil {
		panic(err)
	}
	categories := make(map[string]Category)
	counter := 0

	articleKWs := make(map[string][]string)
	for row := range f.ReadRows("articles") {
		if counter < 1 {
			counter++
			continue
		}

		cells := row.Cells
		length := len(cells)
		currentArticle := ""
		for index, cell := range cells {
			if index == 0 {
				articleKWs[cell.Value] = make([]string, 0)
				currentArticle = cell.Value
				continue
			}
			product := cell.Value
			if product == "NONE" {
				counter++
				continue
			}
			name := product
			name = strings.TrimSpace(name)
			code := greek.Greeklish(strings.ToLower(name))
			code = nonAlphanumericRegex.ReplaceAllString(code, "")
			code = strings.Replace(code, " ", "-", -1)
			code = strings.Replace(code, "&", "", -1)
			articleKWs[currentArticle] = append(articleKWs[currentArticle], code)
			if _, ok := categories[code]; !ok {
				categories[code] = Category{
					code:             code,
					name:             name,
					childCategories:  hashset.New(),
					parentCategories: hashset.New(),
				}

			}
			if index >= 1 {
				if index > 1 {
					parent := cells[index-1].Value
					parent = greek.Greeklish(strings.ToLower(parent))
					parent = nonAlphanumericRegex.ReplaceAllString(parent, "")
					parent = strings.Replace(parent, " ", "-", -1)
					parent = strings.Replace(parent, "&", "", -1)
					if parent != "none" {
						categories[code].parentCategories.Add(parent)
					}
				}

				if index < length-1 {
					child := cells[index+1].Value
					child = greek.Greeklish(strings.ToLower(child))
					child = nonAlphanumericRegex.ReplaceAllString(child, "")
					child = strings.Replace(child, " ", "-", -1)
					child = strings.Replace(child, "&", "", -1)
					if child != "none" {
						categories[code].childCategories.Add(child)
					}
				}
			}
		}

	}

	marshaled := make([]CategoryDto, 0)
	for _, v := range categories {
		nm := v.name
		cd := v.code
		cat := CategoryDto{
			Name:             nm,
			Code:             cd,
			ChildCategories:  v.childCategories.Values(),
			ParentCategories: v.parentCategories.Values(),
		}
		marshaled = append(marshaled, cat)
	}

	articleMarshal := make([]ArticleDto, 0)

	for k, v := range articleKWs {
		articleMarshal = append(articleMarshal, ArticleDto{
			Code:       k,
			Categories: v,
		})
	}
	return categories, articleKWs
}

func parseVideoCategories() (map[string]*Category, map[string][]string) {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	f, err := xlsxreader.OpenFile("videos.xlsx")
	if err != nil {
		panic(err)
	}
	categories := make(map[string]*Category)
	counter := 0
	aricleKWs := make(map[string][]string)
	for row := range f.ReadRows("videos") {
		if counter < 1 {
			counter++
			continue
		}

		cells := row.Cells
		length := len(cells)
		currentArticle := ""
		for index, cell := range cells {
			if index == 0 {
				aricleKWs[cell.Value] = make([]string, 0)
				currentArticle = cell.Value
				continue
			}
			product := cell.Value
			if product == "NONE" {
				counter++
				continue
			}
			name := product
			name = strings.TrimSpace(name)
			code := greek.Greeklish(strings.ToLower(name))
			code = nonAlphanumericRegex.ReplaceAllString(code, "")
			code = strings.Replace(code, " ", "-", -1)
			code = strings.Replace(code, "&", "", -1)
			aricleKWs[currentArticle] = append(aricleKWs[currentArticle], code)
			if _, ok := categories[code]; !ok {
				categories[code] = &Category{
					code:             code,
					name:             name,
					childCategories:  hashset.New(),
					parentCategories: hashset.New(),
				}

			}
			if index >= 1 {
				if index > 1 {
					parent := cells[index-1].Value
					parent = greek.Greeklish(strings.ToLower(parent))
					parent = nonAlphanumericRegex.ReplaceAllString(parent, "")
					parent = strings.Replace(parent, " ", "-", -1)
					parent = strings.Replace(parent, "&", "", -1)
					if parent != "none" {
						categories[code].parentCategories.Add(parent)
					}
				}

				if index < length-1 {
					child := cells[index+1].Value
					child = greek.Greeklish(strings.ToLower(child))
					child = nonAlphanumericRegex.ReplaceAllString(child, "")
					child = strings.Replace(child, " ", "-", -1)
					child = strings.Replace(child, "&", "", -1)
					if child != "none" {
						categories[code].childCategories.Add(child)
					}
				}
			}
		}
	}
	marshaled := make([]CategoryDto, 0)
	for _, v := range categories {
		nm := v.name
		cd := v.code
		cat := CategoryDto{
			Name:             nm,
			Code:             cd,
			ChildCategories:  v.childCategories.Values(),
			ParentCategories: v.parentCategories.Values(),
		}
		marshaled = append(marshaled, cat)
	}
	articleMarshal := make([]ArticleDto, 0)
	for k, v := range aricleKWs {
		articleMarshal = append(articleMarshal, ArticleDto{
			Code:       k,
			Categories: v,
		})
	}
	return categories, aricleKWs
}

func parseProducts(locale string) map[string]Article {
	var translations EGainExportKBTranslation
	xmlf := "Service_euk.xml"
	if locale == "en" {
		xmlf = "Service_grk.xml"
	}
	xmlFile, err := os.Open(xmlf)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(data, &translations)
	if err != nil {
		panic(err)
	}
	res := make(map[string]Article)
	for i, _ := range translations.KBtranslation {
		article := translations.KBtranslation[i].EgplKbArticle
		kw := extractKW(article.IName)
		articleId := extractArticleId(article)
		if len(kw) == 0 || len(articleId) == 0 {
			continue
		}
		title := extractTitle(article.IName)
		keywords := ""
		if len(translations.KBtranslation[i].EgplKbArticle.IKeywords) > 0 {
			keywords = translations.KBtranslation[i].EgplKbArticle.IKeywords
		}
		body := article.IContent
		bodyText := article.IContentText
		articleDto := Article{
			ArticleId:   "PROD-" + articleId,
			Title:       title,
			Keywords:    keywords,
			ContentHtml: body,
			ContentText: bodyText,
			Locale:      locale,
			Kw:          kw,
		}
		res[kw] = articleDto
	}
	defer xmlFile.Close()
	return res
}

func extractArticleId(article EgplKbArticle) string {
	if len(article.IArticleId) == 0 {
		return ""
	}
	arId := article.IArticleId
	arId = strings.Replace(arId, "201600", "", -1)
	for i := 0; i < len(arId); i++ {
		if strings.HasPrefix(arId, "0") {
			arId = arId[1:]
		}
	}
	return arId
}
func extractKW(kw string) string {
	res := ""
	if strings.HasPrefix(kw, "[KW") {
		ind := strings.Index(kw, "]")
		if ind > 0 {
			res = kw[1:ind]
		}

	}
	return res
}

func extractTitle(kw string) string {
	res := ""
	if strings.HasPrefix(kw, "[KW") {
		ind := strings.Index(kw, "]")
		if ind > 0 {
			res = kw[ind+1:]
		}

	}
	return res
}
