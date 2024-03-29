package settings

// Handle the languages in GUI settings

import (
	"reflect"
)

// Enum for languages
const (
	FRENCH  string = "fr"
	ENGLISH string = "en"
)

// Languages yaml struct for theme's app
type LanguageYaml struct {
	English string `yaml:"english"`
	French  string `yaml:"french"`
}

// About text yaml struct for theme's app
type AboutYaml struct {
	Text1 string `yaml:"text_1"`
	Text2 string `yaml:"text_2"`
	Text3 string `yaml:"text_3"`
	Text4 string `yaml:"text_4"`
	Text5 string `yaml:"text_5"`
}

// Headers yaml for tax detail
type TaxHeadersYaml struct {
	Header1 string `yaml:"header_1"`
	Header2 string `yaml:"header_2"`
	Header3 string `yaml:"header_3"`
	Header4 string `yaml:"header_4"`
	Header5 string `yaml:"header_5"`
}

// Handle all data about language data
type Yaml struct {
	Code         string         // code of the language (fr, en, etc...)
	Theme        ThemeYaml      `yaml:"themes"`
	Languages    LanguageYaml   `yaml:"languages"`
	Abouts       AboutYaml      `yaml:"abouts"`
	TaxHeaders   TaxHeadersYaml `yaml:"tax_headers"`
	File         string         `yaml:"file"`
	Settings     string         `yaml:"settings"`
	Income       string         `yaml:"income"`
	Status       string         `yaml:"status"`
	Children     string         `yaml:"children"`
	Tax          string         `yaml:"tax"`
	Remainder    string         `yaml:"remainder"`
	Share        string         `yaml:"share"`
	Save         string         `yaml:"save"`
	ThemeCode    string         `yaml:"theme"`
	LanguageCode string         `yaml:"language"`
	Currency     string         `yaml:"currency"`
	Logs         string         `yaml:"logs"`
	Help         string         `yaml:"help"`
	About        string         `yaml:"about"`
	Author       string         `yaml:"author"`
	Close        string         `yaml:"close"`
	Quit         string         `yaml:"quit"`
}

// GetLanguage get value of last language selected (fr, en)
func GetDefaultLanguage() string {
	return ENGLISH
}

// GetThemes parse ThemeYaml struct to get value of each field
func (yaml *Yaml) GetThemes() []string {
	v := reflect.ValueOf(yaml.Theme)
	themes := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		themes[i] = v.Field(i).String()
	}
	return themes
}

// GetLanguages parse LanguageYaml struct to get value of each field
func (yaml *Yaml) GetLanguages() []string {
	v := reflect.ValueOf(yaml.Languages)
	languages := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		languages[i] = v.Field(i).String()
	}
	return languages
}

// GetAbouts parse AboutYaml struct to get value of each field
func (yaml *Yaml) GetAbouts() []string {
	v := reflect.ValueOf(yaml.Abouts)
	abouts := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		abouts[i] = v.Field(i).String()
	}
	return abouts
}

// GetTaxHeaders parse TaxHeadersYaml struct to get value of each field
func (yaml *Yaml) GetTaxHeaders() []string {
	v := reflect.ValueOf(yaml.TaxHeaders)
	taxHeaders := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		taxHeaders[i] = v.Field(i).String()
	}
	return taxHeaders
}
