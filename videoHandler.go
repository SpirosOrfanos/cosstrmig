package main

import "fmt"

func videos() {
	englishArticles := parseProducts("en")
	greekArticles := parseProducts("el")
	aggregatedArticles := make(map[string][]Article)
	aggregatedSaveArticles := make(map[string][]HelpVideo)

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

		itm := VideoComponent{
			VideoID:     "",
			Description: v[0].ContentHtml,
			Component:   "reusables.help-video-item",
		}
		inter := []VideoComponent{itm}
		helpArticle := HelpVideo{
			Id:          nil,
			Title:       v[0].Title,
			Locale:      v[0].Locale,
			ArticleId:   v[0].ArticleId,
			ReferenceId: k,
			Videos:      inter,
		}
		aggregatedSaveArticles[k] = make([]HelpVideo, 0)
		aggregatedSaveArticles[k] = append(aggregatedSaveArticles[k], helpArticle)
		if len(v) == 1 {
			continue
		}

		itm2 := VideoComponent{
			VideoID:     "",
			Description: v[1].ContentHtml,
			Component:   "reusables.help-video-item",
		}
		inter2 := []VideoComponent{itm2}
		helpArticle2 := HelpVideo{
			Id:          nil,
			Title:       v[1].Title,
			Locale:      v[1].Locale,
			ArticleId:   v[1].ArticleId,
			ReferenceId: k,
			Videos:      inter2,
		}
		artV := make([]HelpVideo, 0)
		artV = append(artV, aggregatedSaveArticles[k][0])
		artV = append(artV, helpArticle2)
		aggregatedSaveArticles[k] = artV
	}

	categories, categoryArticles := parseVideoCategories()

	adapter := NewAdapter()
	savedArticles := make(map[string][]HelpVideo)
	for key, rts := range aggregatedSaveArticles {
		if _, ok := categoryArticles[key]; !ok {
			continue
		}
		newArticles := make([]HelpVideo, 0)
		res, _ := adapter.CreateVideo(rts[0])
		newArticles = append(newArticles, HelpVideo{
			Id:          &res.Data.Id,
			ReferenceId: rts[0].ReferenceId,
			Locale:      rts[0].Locale,
		})
		if len(rts) > 1 {
			res2, _ := adapter.LocalizationVideo(rts[1], res.Data.Id)
			newArticles = append(newArticles, HelpVideo{
				Id:          &res2.Id,
				ReferenceId: rts[1].ReferenceId,
				Locale:      rts[1].Locale,
			})
		}
		savedArticles[key] = newArticles
	}

	categoriesSavedEl := make(map[string]int)
	categoriesSavedEn := make(map[string]int)

	for cat, catv := range categories {
		nc := HelpVideoCategory{
			Name:       catv.name,
			CategoryId: catv.code,
			Locale:     "el",
		}
		res, _ := adapter.CreateHelpVideoCategory(nc)
		categoriesSavedEl[cat] = res.Data.Id
		ncEn := HelpVideoCategory{
			Name:       catv.name,
			CategoryId: catv.code,
			Locale:     "en",
		}
		res2, _ := adapter.LocalizationsHelpVideoCategory(ncEn, res.Data.Id)
		categoriesSavedEn[cat] = res2.Id
	}

	savedArticlesEl := make(map[string]HelpVideo)
	savedArticlesEn := make(map[string]HelpVideo)
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
			adapter.UpdateHelpVideoCategory(HelpVideoCategory{
				HelpVideos: elIdsC,
			}, elId)
			adapter.UpdateHelpVideoCategory(HelpVideoCategory{
				HelpVideos: enIdsC,
			}, enId)

		}
	}
	for articleId, categoryCodes := range categoryArticles {
		if elId, ok := savedArticlesEl[articleId]; ok {
			categoryIds := make([]string, 0)
			for _, categoryCode := range categoryCodes {
				if catgr, ok := categoriesSavedEn[categoryCode]; ok {
					categoryIds = append(categoryIds, fmt.Sprintf("%d", catgr))
				}
			}
			adapter.UpdateHelpVideo(HelpVideo{
				HelpVideoCategories: categoryIds,
			}, *elId.Id)
		}
		if enId, ok := savedArticlesEn[articleId]; ok {
			categoryIds := make([]string, 0)
			for _, categoryCode := range categoryCodes {
				if catgr, ok := categoriesSavedEn[categoryCode]; ok {
					categoryIds = append(categoryIds, fmt.Sprintf("%d", catgr))
				}
			}
			adapter.UpdateHelpVideo(HelpVideo{
				HelpVideoCategories: categoryIds,
			}, *enId.Id)
		}
	}

}
