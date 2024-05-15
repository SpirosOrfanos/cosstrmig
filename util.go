package main

import (
	"encoding/xml"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/thedatashed/xlsxreader"
	"go.reizu.org/greek"
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

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
