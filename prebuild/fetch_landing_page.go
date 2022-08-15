package main

import (
	"fmt"
	"path/filepath"

	ButterCMS "github.com/ButterCMS/buttercms-go"
)

const defaultLandingPageSlug = "landing-page-with-components"

type LandingPageSectionBase struct {
	ScrollAnchorId string
	Headline       string
	SubHeadline    string
	IsSet          bool
}

type HeroSection struct {
	LandingPageSectionBase

	ButtonUrl   string
	ButtonLabel string
	Image       string
}

type ImageWithTextSection struct {
	HeroSection

	ImagePosition string
}

type SEOMetadata struct {
	Title       string
	Description string
}

type LandingPageFile struct {
	SEOMetadata

	HeroSection  HeroSection
	AboutSection ImageWithTextSection
	TryIt        ImageWithTextSection
}

func FetchLandingPages(pathToFiles string) {
	pageType := "landing-page"

	response, pageFetchingErr := ButterCMS.GetPages(pageType, map[string]string{})
	HandleErr(pageFetchingErr)

	for _, page := range response.PageList {
		processLandingPage(pathToFiles, page)
	}
}

func processLandingPage(pathToFile string, page ButterCMS.Page) {
	data := LandingPageFile{
		SEOMetadata: processSEOMetadata(page),
	}

	if body, err := GetValue[[]interface{}](page.Fields, "body"); err == nil {
		for _, untypedSection := range body {
			section := untypedSection.(map[string]interface{})

			scrollAnchorId, _ := getSectionFieldsValue[string](section, "scroll_anchor_id")
			headline, _ := getSectionFieldsValue[string](section, "headline")

			fmt.Printf("\n%s", scrollAnchorId)
			switch scrollAnchorId {
			case "home":
				data.HeroSection = processHeroSection(section, headline, scrollAnchorId)
			case "about":
				data.AboutSection = processImageWithTextSection(section, headline, scrollAnchorId)
			case "tryit":
				data.TryIt = processImageWithTextSection(section, headline, scrollAnchorId)
			}

		}
	}

	if defaultLandingPageSlug == page.Slug {
		CreateFile(data, filepath.Join(pathToFile, "_index.md"))
	}

	CreateFile(data, filepath.Join(pathToFile, fmt.Sprintf("%s.md", page.Slug)))
}

func processImageWithTextSection(section map[string]interface{}, headline string, scrollAnchorId string) ImageWithTextSection {
	heroSection := processHeroSection(section, headline, scrollAnchorId)

	imagePosition, _ := getSectionFieldsValue[string](section, "image_position")

	return ImageWithTextSection{
		HeroSection: heroSection,

		ImagePosition: imagePosition,
	}
}

func processHeroSection(section map[string]interface{}, headline string, scrollAnchorId string) HeroSection {
	buttonLabel, _ := getSectionFieldsValue[string](section, "button_label")
	buttonUrl, _ := getSectionFieldsValue[string](section, "button_url")
	image, _ := getSectionFieldsValue[string](section, "image")
	subHeadline, _ := getSectionFieldsValue[string](section, "subheadline")

	return HeroSection{
		ButtonLabel: buttonLabel,
		ButtonUrl:   buttonUrl,
		Image:       image,

		LandingPageSectionBase: LandingPageSectionBase{
			SubHeadline:    subHeadline,
			IsSet:          true,
			Headline:       headline,
			ScrollAnchorId: scrollAnchorId,
		},
	}
}

func processSEOMetadata(page ButterCMS.Page) SEOMetadata {
	result := SEOMetadata{}

	if seo, err := GetValue[map[string]interface{}](page.Fields, "seo"); err == nil {
		title, _ := GetValue[string](seo, "title")
		description, _ := GetValue[string](seo, "description")

		result.Title = title
		result.Description = description
	}

	return result
}

func getSectionFieldsValue[T any](input map[string]interface{}, name string) (T, error) {
	fields, err := GetValue[map[string]interface{}](input, "fields")
	if err != nil {
		var null T
		return null, err
	}

	return GetValue[T](fields, name)
}
