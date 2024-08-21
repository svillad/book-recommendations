package models

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	MinPages = 1
	MaxPages = 10000
	MinYear  = 1800
	MaxYear  = 2100
	MinBooks = 1
	MaxBooks = 1000
)

type Book struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	YearPublished int64   `json:"yearPublished"`
	Rating        float64 `json:"rating"`
	Pages         int64   `json:"pages"`
	Genre         Genre   `json:"genre"`
	Author        Author  `json:"author"`
}

type BookRequest struct {
	Authors  string `json:"authors"`
	Genres   string `json:"genres"`
	MinPages string `json:"min-pages"`
	MaxPages string `json:"max-pages"`
	MinYear  string `json:"min-year"`
	MaxYear  string `json:"max-year"`
	Limit    string `json:"limit"`
}

var (
	idsRules = []validation.Rule{
		validation.Match(regexp.MustCompile("^([0-9]+,)*[0-9]+$")).Error("should be for exaple: 123,456,789"),
	}
	pagesRules = []validation.Rule{
		is.Int.Error("should be numeric"),
		validation.By(validateMinMax(MinPages, MaxPages)),
	}
	yearRules = []validation.Rule{
		is.Int.Error("should be numeric"),
		validation.By(validateMinMax(MinYear, MaxYear)),
	}
	limitRules = []validation.Rule{
		is.Int.Error("should be numeric"),
		validation.By(validateMinMax(MinBooks, MaxBooks)),
	}
)

func (br BookRequest) Validate() error {
	reqCopy := br

	return validation.ValidateStruct(&reqCopy,
		validation.Field(&reqCopy.Authors, idsRules...),
		validation.Field(&reqCopy.Genres, idsRules...),
		validation.Field(&reqCopy.MinPages, pagesRules...),
		validation.Field(&reqCopy.MaxPages, pagesRules...),
		validation.Field(&reqCopy.MinYear, yearRules...),
		validation.Field(&reqCopy.MaxYear, yearRules...),
		validation.Field(&reqCopy.Limit, limitRules...),
	)
}

func validateMinMax(minValue, maxValue int) validation.RuleFunc {
	return func(value interface{}) error {
		if value.(string) != "" {
			s, _ := strconv.Atoi(value.(string))
			if s < minValue || s > maxValue {
				errorLimit := fmt.Sprintf("should be between %v and %v", minValue, maxValue)
				return errors.New(errorLimit)
			}
		}
		return nil
	}
}
