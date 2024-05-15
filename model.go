package main

import (
	"encoding/xml"
	"github.com/emirpasic/gods/sets/hashset"
)

type Category struct {
	id               *int         `json:"id,omitempty"`
	name             string       `json:"name,omitempty"`
	code             string       `json:"code,omitempty"`
	childCategories  *hashset.Set `json:"childCategories,omitempty"`
	parentCategories *hashset.Set `json:"parentCategories,omitempty"`
	Videos           *hashset.Set `json:"videos,omitempty"`
}
type CategoryDto struct {
	Name             string        `json:"name,omitempty"`
	Code             string        `json:"code,omitempty"`
	ChildCategories  []interface{} `json:"childCategories,omitempty"`
	ParentCategories []interface{} `json:"parentCategories,omitempty"`
	Videos           []interface{} `json:"videos,omitempty"`
}

type ArticleDto struct {
	Code       string   `json:"code,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Id         *int     `json:"id,omitempty"`
}
type EGainExportKBTranslation struct {
	XMLName       xml.Name        `xml:"eGainExportKBTranslation"`
	Text          string          `xml:",chardata"`
	KBtranslation []KBtranslation `xml:"KBtranslation"`
}

type KBtranslation struct {
	Text          string        `xml:",chardata"`
	EgplKbArticle EgplKbArticle `xml:"egpl_kb_article"`
}
type Article struct {
	ArticleId   string `json:"articleId,omitempty"`
	Title       string `json:"title,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	ContentHtml string `json:"contentHtml,omitempty"`
	ContentText string `json:"contentText,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Kw          string `json:"kw,omitempty"`
}

type HtmlReusable struct {
	Component string `json:"__component,omitempty"`
	Title     string `json:"name,omitempty"`
	Body      string `json:"body,omitempty"`
}

type EgplKbArticle struct {
	Text                       string `xml:",chardata"`
	IUUID                      string `xml:"iUUID"`
	IArticleId                 string `xml:"iArticleId"`
	IArticleReferenceId        string `xml:"iArticleReferenceId"`
	IArticleVersion            string `xml:"iArticleVersion"`
	IArticleAvailabilityDate   string `xml:"iArticleAvailabilityDate"`
	IArticleAvailabilityDateTZ string `xml:"iArticleAvailabilityDateTZ"`
	IArticleExpiryDate         string `xml:"iArticleExpiryDate"`
	IArticleExpiryDateTZ       string `xml:"iArticleExpiryDateTZ"`
	IArticleLabel              string `xml:"iArticleLabel"`
	IName                      string `xml:"iName"`
	IDescription               string `xml:"iDescription"`
	IArticleType               string `xml:"iArticleType"`
	IKeywords                  string `xml:"iKeywords"`
	ISummary                   string `xml:"iSummary"`
	IContent                   string `xml:"iContent"`
	IContentText               string `xml:"iContentText"`
	IAddtionalInfo             string `xml:"iAddtionalInfo"`
	IComment                   string `xml:"iComment"`
	INewFlag                   string `xml:"iNewFlag"`
}

type HelpArticleCategoryWrapper struct {
	Data HelpArticleCategory `json:"data,omitempty"`
}
type HelpArticleCategory struct {
	id               *int     `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	CategoryId       string   `json:"categoryId,omitempty"`
	Description      string   `json:"description,omitempty"`
	Locale           string   `json:"locale,omitempty"`
	ChildCategories  []string `json:"child_categories,omitempty"`
	ParentCategories []string `json:"parent_categories,omitempty"`
	Articles         []int    `json:"articles,omitempty"`
}

type HelpArticleWrapper struct {
	Data HelpArticle `json:"data,omitempty"`
}
type HelpArticle struct {
	Id          *int        `json:"id,omitempty"`
	Title       string      `json:"title,omitempty"`
	Keywords    string      `json:"keywords,omitempty"`
	Locale      string      `json:"locale,omitempty"`
	Reusables   interface{} `json:"reusables,omitempty"`
	ArticleId   string      `json:"articleId,omitempty"`
	ReferenceId string      `json:"referenceId,omitempty"`
}

type Reusable struct {
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	Component string `json:"__component,omitempty"` //html-reusable
}

type Toggle struct {
	Description string       `json:"description,omitempty"`
	Component   string       `json:"__component,omitempty"` //toggles
	Items       []ToggleItem `json:"Item,omitempty"`
}

type ToggleItem struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type HelpVideoCategoryWrapper struct {
	Data HelpVideoCategory `json:"data,omitempty"`
}
type HelpVideoCategory struct {
	Name       string `json:"name,omitempty"`
	CategoryId string `json:"categoryId,omitempty"`
	Locale     string `json:"locale,omitempty"`
	HelpVideos []int  `json:"help_videos,omitempty"`
}

type HelpVideoWrapper struct {
	Data HelpVideoCategory `json:"data,omitempty"`
}
type HelpVideo struct {
	Id          *int             `json:"id,omitempty"`
	Title       string           `json:"name,omitempty"`
	ArticleId   string           `json:"articleId,omitempty"`
	ReferenceId string           `json:"referenceId,omitempty"`
	Locale      string           `json:"locale,omitempty"`
	Videos      []VideoComponent `json:"videos,omitempty"`
}

type VideoComponent struct {
	VideoID     string `json:"videoID,omitempty"`
	Description string `json:"description,omitempty"`
	Component   string `json:"__component,omitempty"` //help-video-item
}
type ResponseWrapper struct {
	Data Response `json:"data,omitempty"`
}
type Response struct {
	Id int `json:"id,omitempty"`
}

type TGItem struct {
	title string
	items string
}
type TGItemWrapper struct {
	description string
	items       []TGItem
}
