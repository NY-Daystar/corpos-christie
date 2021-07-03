package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/LucasNoga/corpos-christie/lib/utils"
)

// Define the program config
type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Tax     Tax
	TaxList []Tax `json:"tax"`
}

// Define the tranche on a specific year
type Tax struct {
	Year     int       `json:"year"`     // Year of the tax specifications
	Tranches []Tranche `json:"tranches"` // Tranches list of tax
}

// Define one of the tranch of tax
type Tranche struct {
	Min  int    `json:"min"`  // Minimun in euros to get in the tranche
	Max  int    `json:"max"`  // Maximum in euros to get in the tranche
	Rate string `json:"rate"` // Rate taxable in euros in this tranche
}

// Load configuration from config.json file
func (cfg *Config) LoadConfiguration(file string) (bool, error) {
	// default path
	if file == "" {
		file = "./config.json"
	}

	jsonFile, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer jsonFile.Close()

	// Load json file
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(cfg)
	if err != nil {
		return false, err
	}

	// Define tax of the current year as reference
	cfg.loadTaxYear()

	return true, nil
}

// Define the tax of the current year
func (cfg *Config) loadTaxYear() {
	for _, tax := range cfg.TaxList {
		if tax.Year == utils.GetCurrentYear() { // get tax of current year
			cfg.Tax = tax
			break
		}
	}

	// If no tax tranches are defined from current year load default tax 2021
	if len(cfg.Tax.Tranches) == 0 {
		cfg.Tax = cfg.TaxList[0]
	}
}

// Load default configuration file if we don't have a 'config.json file'
func (cfg *Config) LoadDefaultConfiguration() {
	log.Printf("Loading Default configuration...")
	cfg.Name = "Corpos-Christie"
	cfg.Version = "1.0.0"
	cfg.Tax = Tax{
		Year: 2021,
		Tranches: []Tranche{
			{Min: 0, Max: 10084, Rate: "0%"},
			{Min: 10085, Max: 25710, Rate: "11%"},
			{Min: 25711, Max: 73516, Rate: "30%"},
			{Min: 73517, Max: 158122, Rate: "41%"},
			{Min: 158123, Max: 1000000, Rate: "45%"},
		},
	}
	cfg.TaxList = append(cfg.TaxList, cfg.Tax)
}
