package main

import (
	"corpos-christie/colors"
	"corpos-christie/config"
	"corpos-christie/utils"
	"fmt"
	"log"
	"os"
)

var cfg *config.Config

// Start tax calculator from input user
func start(cfg *config.Config) bool {
	fmt.Print("Enter your income (Revenu net imposable): ")
	var input string = utils.ReadValue()
	r, err := utils.ConvertStringToInt(input)
	if err != nil {
		log.Printf("Error: Tax income is not convertible in int, details: %v", err)
		return false
	}

	log.Printf("income to calculate tax: %v", colors.Red(r))

	//TODO faire un fichier process pour faire les calculs d'impot
	//TODO Creer un fichier de config.json avec une struct tranches et une struct tranche qui gere ca

	//TODO Calculer l'imposition
	//TODO Faire ensuite le differentiel pour savoir ce qui nous reste

	return true
}

// Ask user if he wants to restart program
func askRestart() bool {
	for {
		fmt.Print("Would you want to enter a new income (Y/n): ")
		var input string = utils.ReadValue()
		if input == "Y" || input == "y" || input == "Yes" || input == "yes" {
			log.Printf("Restarting program...")
			return true
		} else {
			return false
		}
	}
}

// Init configuration file
func init() {
	cfg = new(config.Config)
	config.LoadConfiguration(cfg)
}

func main() {
	log.Printf("Project: %v", colors.Yellow(cfg.Name))
	log.Printf("Version %v", colors.Yellow(cfg.Version))

	var keep bool
	for ok := true; ok; ok = keep {
		status := start(cfg)
		log.Printf("Status of operation: %v", status)
		keep = askRestart()
	}

	log.Printf("Program exited...")
	os.Exit(0)
}