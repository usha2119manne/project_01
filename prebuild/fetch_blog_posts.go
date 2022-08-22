package main

import (
	"fmt"
	"path/filepath"

	ButterCMS "github.com/ButterCMS/buttercms-go"
)

type BlogPostFile struct {
	ButterCMS.Post

	CategoriesSlugs []string
	ImageMeta       string
}

func FetchBlogPosts(pathToFiles string) {
	response, err := ButterCMS.GetPosts(map[string]string{})
	HandleErr(err)

	for _, post := range response.PostList {
		post.URL = "" // Conflict with Hugo URL definition in content file. + not used at all

		categoriesSlugs := []string{}
		for _, category := range post.CategoryList {
			categoriesSlugs = append(categoriesSlugs, category.Slug)
		}

		data := BlogPostFile{
			Post:            post,
			CategoriesSlugs: categoriesSlugs,
		}

		if post.FeaturedImageAlt != "" {
			data.ImageMeta = post.FeaturedImageAlt
		}

		CreateFile(data, filepath.Join(pathToFiles, fmt.Sprintf("%s.md", post.Slug)))
	}
}
