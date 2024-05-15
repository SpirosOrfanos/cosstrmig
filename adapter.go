package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Adapter struct {
	BaseUrl    string
	httpClient *http.Client
}

func (srv *Adapter) CreateHelpArticleCategory(category HelpArticleCategory) (ResponseWrapper, error) {
	var resp ResponseWrapper
	data := HelpArticleCategoryWrapper{Data: category}
	res2B, _ := json.Marshal(data)
	uri, _ := url.JoinPath(srv.BaseUrl, "api", "help-article-categories")
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
		return resp, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	response, err := srv.httpClient.Do(req)
	if err != nil {
		fmt.Println("CreateHelpArticleCategory", string(res2B))
		panic(err)
		return resp, err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		fmt.Println(response.StatusCode)
		fmt.Println(response)
		fmt.Println("CreateHelpArticleCategory", string(res2B))
		panic(err)
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
		return resp, err
	}
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func (srv *Adapter) UpdateHelpArticleCategory(category HelpArticleCategory, id int) (ResponseWrapper, error) {
	var resp ResponseWrapper
	data := HelpArticleCategoryWrapper{Data: category}
	res2B, _ := json.Marshal(data)
	uri, _ := url.JoinPath(srv.BaseUrl, "api", "help-article-categories", fmt.Sprintf("%d", id))
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
		return resp, err
	}
	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	response, err := srv.httpClient.Do(req)
	if err != nil {
		fmt.Println("UpdateHelpArticleCategory", string(res2B))
		panic(err)
		return resp, err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		fmt.Println(response.StatusCode)
		fmt.Println(response)
		fmt.Println("UpdateHelpArticleCategory", id, string(res2B))
		panic(err)
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
		return resp, err
	}
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func (srv *Adapter) LocalizationsHelpArticleCategory(category HelpArticleCategory, id int) (Response, error) {
	var resp Response

	res2B, _ := json.Marshal(category)
	uri, _ := url.JoinPath(srv.BaseUrl, "api", "help-article-categories", fmt.Sprintf("%d", id), "localizations")
	jsonData, err := json.Marshal(category)
	if err != nil {
		panic(err)
		return resp, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	response, err := srv.httpClient.Do(req)
	if err != nil {
		fmt.Println("CreateHelpArticleCategory", string(res2B))
		panic(err)
		return resp, err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		fmt.Println(response.StatusCode)
		fmt.Println(response)
		fmt.Println("LocalizationsHelpArticleCategory", string(res2B))
		panic(err)
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
		return resp, err
	}
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func (srv *Adapter) Create(article HelpArticle) (ResponseWrapper, error) {
	var resp ResponseWrapper
	data := HelpArticleWrapper{Data: article}
	res2B, _ := json.Marshal(data)
	uri, _ := url.JoinPath(srv.BaseUrl, "api", "help-articles")
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
		return resp, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	response, err := srv.httpClient.Do(req)
	if err != nil {
		fmt.Println("Create", string(res2B))
		panic(err)
		return resp, err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		fmt.Println(response.StatusCode)
		fmt.Println(response)
		fmt.Println("Create", string(res2B))
		panic(err)
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
		return resp, err
	}
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func (srv *Adapter) Localization(article HelpArticle, id int) (Response, error) {
	var resp Response

	res2B, _ := json.Marshal(article)

	uri, _ := url.JoinPath(srv.BaseUrl, "api", "help-articles", fmt.Sprintf("%d", id), "localizations")
	jsonData, err := json.Marshal(article)
	if err != nil {
		panic(err)
		return resp, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	response, err := srv.httpClient.Do(req)
	if err != nil {
		fmt.Println("Localization", string(res2B))
		panic(err)
		return resp, err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		fmt.Println("StatusCode", response.StatusCode)
		fmt.Println("response", response)
		fmt.Println("Localization", string(res2B))
		panic(err)
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
		return resp, err
	}
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func NewAdapter() *Adapter {
	return &Adapter{
		BaseUrl:    os.Getenv("CMS_URL"),
		httpClient: makeClient(),
	}
}

func makeClient() *http.Client {
	client := &http.Client{Timeout: 60 * time.Second}
	return client
}
