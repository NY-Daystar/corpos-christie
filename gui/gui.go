// Copyright 2016 The corpos-christie author
// Licensed under GPLv3.

// Package gui defines component and script to launch gui application
package gui

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/LucasNoga/corpos-christie/config"
	"github.com/LucasNoga/corpos-christie/gui/settings"
	"github.com/LucasNoga/corpos-christie/gui/themes"
	"github.com/LucasNoga/corpos-christie/tax"
	"github.com/LucasNoga/corpos-christie/user"
	"github.com/LucasNoga/corpos-christie/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

// GUI represents the program parameters to launch in gui the application
type GUI struct {
	Config   *config.Config    // Config to use correctly the program
	User     *user.User        // User param to use program
	App      fyne.App          // Fyne application
	Window   fyne.Window       // Fyne window
	Settings settings.Settings // Settings of the app
	Logger   *zap.Logger       // Logger of GUI

	// Settings
	Theme    themes.Theme   // Fyne theme for the application
	Language settings.Yaml  // Yaml struct with all language data
	Currency binding.String // Currency to display

	// Widgets
	entryIncome    *widget.Entry       // Input Entry to set income
	radioStatus    *widget.RadioGroup  // Input Radio buttons to get status
	selectChildren *widget.SelectEntry // Input Select to know how children

	// buttonSave *widget.Button // Label for save button

	// Bindings
	Tax                binding.String     // Bind for tax value
	Remainder          binding.String     // Bind for remainder value
	Shares             binding.String     // Bind for shares value
	labelShares        binding.String     // Bind for shares label
	labelIncome        binding.String     // Bind for income label
	labelStatus        binding.String     // Bind for status label
	labelChildren      binding.String     // Bind for children label
	labelTax           binding.String     // Bind for tax label
	labelRemainder     binding.String     // Bind for remainder label
	labelsAbout        binding.StringList // List of label in about modal
	labelsTaxHeaders   binding.StringList // List of label for tax details headers
	labelsMinTranche   binding.StringList // List of labels for min tranche in grid
	labelsMaxTranche   binding.StringList // List of labels for max tranche in grid
	labelsTrancheTaxes binding.StringList // List of tranches tax label
}

// Start Launch GUI application
func (gui GUI) Start() {
	gui.App = app.New()
	gui.Window = gui.App.NewWindow(config.APP_NAME)

	// Set Logger
	gui.setLogger()
	gui.Logger.Info("Launch application")

	// Load settings
	gui.setAppSettings()

	// Size and Position
	const WIDTH = 1100
	const HEIGHT = 540
	gui.Window.Resize(fyne.NewSize(WIDTH, HEIGHT))
	gui.Window.CenterOnScreen()
	gui.Logger.Info("Load window", zap.Int("height", HEIGHT), zap.Int("width", WIDTH))

	// Set menu
	var menu *fyne.MainMenu = gui.setMenu()
	gui.Window.SetMainMenu(menu)

	// Create layouts and widgets
	gui.setLayouts()
	gui.Logger.Info("Layout and widgets loaded")

	// Create layouts and widgets
	gui.setEvents()
	gui.Logger.Info("Event loaded")

	// Set Icon
	var iconName string = "logo.ico"
	var iconPath string = fmt.Sprintf("%s/%s", config.ASSETS_PATH, iconName)
	var icon fyne.Resource = settings.GetIcon(iconPath)
	gui.Logger.Info("Load icon", zap.String("name", iconName), zap.String("path", iconPath))
	gui.Window.SetIcon(icon)

	gui.Window.ShowAndRun()
}

// setSettings get and configure app settings
func (gui *GUI) setAppSettings() {
	gui.Settings, _ = settings.Load(gui.Logger)

	gui.Logger.Info("Settings loaded",
		zap.Int("theme", gui.Settings.Theme),
		zap.String("language", gui.Settings.Language),
		zap.String("theme", gui.Settings.Currency),
	)

	gui.setTheme(gui.Settings.Theme)
	gui.setLanguage(gui.Settings.Language)
	gui.Currency = binding.BindString(&gui.Settings.Currency)

}

// SetTheme change theme of the application
// (if param = 0 then dark if 1 then light)
func (gui *GUI) setTheme(theme int) {
	var t themes.Theme
	if theme == settings.DARK {
		t = themes.DarkTheme{}
	} else {
		t = themes.LightTheme{}
	}
	gui.Logger.Info("Set theme", zap.Int("theme", theme))
	gui.App.Settings().SetTheme(t)
}

// SetLanguage change language of the application
func (gui *GUI) setLanguage(code string) {
	gui.Logger.Info("Set language", zap.String("code", code))

	var languageFile string = fmt.Sprintf("%s/%s.yaml", config.LANGUAGES_PATH, code)
	gui.Logger.Debug("Load file for language", zap.String("file", languageFile))

	yamlFile, _ := os.ReadFile(languageFile)
	err := yaml.Unmarshal(yamlFile, &gui.Language)

	gui.Language.Code = code

	if err != nil {
		gui.Logger.Sugar().Fatalf("Unmarshal language file %s: %v", languageFile, err)
	}

	gui.Logger.Sugar().Debugf("Language Yaml %v", gui.Language)
}

// setCurrency change language of the application
func (gui *GUI) setCurrency(currency string) {
	gui.Logger.Info("Set currency", zap.String("currency", currency))
	gui.Currency.Set(currency)
}

// setEvents Set the events/trigger of gui widgets
func (gui *GUI) setEvents() {
	gui.entryIncome.OnChanged = func(input string) {
		gui.calculate()
	}
	gui.radioStatus.OnChanged = func(input string) {
		gui.calculate()
	}
	gui.selectChildren.OnChanged = func(input string) {
		gui.calculate()
	}

}

// getIncome Get value of widget entry
func (gui *GUI) getIncome() int {
	intVal, err := strconv.Atoi(gui.entryIncome.Text)
	if err != nil {
		return 0
	}
	return intVal
}

// getStatus Get value of widget radioGroup
func (gui *GUI) getStatus() bool {
	return gui.radioStatus.Selected == "Couple"
}

// getChildren get value of widget select
func (gui *GUI) getChildren() int {
	children, err := strconv.Atoi(gui.selectChildren.Entry.Text)
	if err != nil {
		return 0
	}
	return children
}

// reload Refresh widget who needed specially when language changed
func (gui *GUI) Reload() {
	// Simple data bind
	gui.labelIncome.Set(gui.Language.Income)
	gui.labelStatus.Set(gui.Language.Status)
	gui.labelChildren.Set(gui.Language.Children)
	gui.labelTax.Set(gui.Language.Tax)
	gui.labelRemainder.Set(gui.Language.Remainder)
	gui.labelShares.Set(gui.Language.Share)

	// Handle widget
	// gui.buttonSave.SetText(gui.Language.Save) // TODO

	// Reload about content
	gui.labelsAbout.Set(gui.Language.GetAbouts())

	// Reload header tax details
	gui.labelsTaxHeaders.Set(gui.Language.GetTaxHeaders())

	// Reload grid header
	currency, _ := gui.Currency.Get()
	gui.labelsTrancheTaxes.Set(*createTrancheTaxesLabels(gui.labelsTrancheTaxes.Length(), currency))

	// Reload grid min tranches
	var minList []string
	for index := 0; index < gui.labelsMinTranche.Length(); index++ {
		var min string = utils.ConvertIntToString(gui.Config.Tax.Tranches[index].Min) + " " + currency
		minList = append(minList, min)
	}
	gui.labelsMinTranche.Set(minList)

	// Reload grid max tranches
	var maxList []string
	for index := 0; index < gui.labelsMaxTranche.Length(); index++ {
		var max string = utils.ConvertIntToString(gui.Config.Tax.Tranches[index].Max) + " " + currency
		if gui.Config.Tax.Tranches[index].Max == math.MaxInt64 {
			max = "-"
		}
		maxList = append(maxList, max)
	}
	gui.labelsMaxTranche.Set(maxList)
}

// calculate Get values of gui to calculate tax
func (gui *GUI) calculate() {
	gui.User.Income = gui.getIncome()
	gui.User.IsInCouple = gui.getStatus()
	gui.User.Children = gui.getChildren()

	result := tax.CalculateTax(gui.User, gui.Config)
	gui.Logger.Sugar().Debugf("Result taxes %#v", result)

	var tax string = utils.ConvertInt64ToString(int64(result.Tax))
	var remainder string = utils.ConvertInt64ToString(int64(result.Remainder))
	var shares string = utils.ConvertInt64ToString(int64(result.Shares))

	// Set data in tax layout
	gui.Tax.Set(tax)
	gui.Remainder.Set(remainder)
	gui.Shares.Set(shares)

	// Set Tax details
	currency, _ := gui.Currency.Get()
	for index := 0; index < gui.labelsTrancheTaxes.Length(); index++ {
		var taxTranche string = utils.ConvertIntToString(int(result.TaxTranches[index].Tax))
		gui.labelsTrancheTaxes.SetValue(index, taxTranche+" "+currency)
	}
}

// createMenu create mainMenu for window
func (gui *GUI) setMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		gui.createFileMenu(),
		gui.createHelpMenu(),
	)
}

// createFileMenu create file item in toolbar to handle app settings
func (gui *GUI) createFileMenu() *fyne.Menu {
	fileMenu := fyne.NewMenu(gui.Language.File,
		fyne.NewMenuItem(gui.Language.Settings, func() {
			dialog.ShowCustom(gui.Language.Settings, gui.Language.Close,
				container.NewVBox(
					gui.createSelectTheme(),
					widget.NewSeparator(),
					gui.createSelectLanguage(),
					widget.NewSeparator(),
					gui.createSelectCurrency(),
					widget.NewSeparator(),
					gui.createLabelLogs(),
				), gui.Window)
		}),
		fyne.NewMenuItem(gui.Language.Quit, func() { gui.App.Quit() }),
	)
	return fileMenu
}

// createSelectTheme create select to change theme
func (gui *GUI) createSelectTheme() *fyne.Container {
	selectTheme := widget.NewSelect(gui.Language.GetThemes(), nil)

	selectTheme.OnChanged = func(s string) {
		index := selectTheme.SelectedIndex()
		gui.setTheme(index)
		gui.Settings.Set("theme", index)
	}
	selectTheme.SetSelectedIndex(gui.Settings.Theme)
	return container.NewHBox(
		widget.NewLabel(gui.Language.ThemeCode),
		selectTheme,
	)
}

// createSelectLanguage create select to change language
func (gui *GUI) createSelectLanguage() *fyne.Container {
	selectLanguage := widget.NewSelect(gui.Language.GetLanguages(), nil)
	selectLanguage.SetSelectedIndex(getLanguageIndex(gui.Language.Code))
	selectLanguage.OnChanged = func(s string) {
		index := selectLanguage.SelectedIndex()
		var getLanguage = func() string {
			switch index {
			case 0:
				return settings.ENGLISH
			case 1:
				return settings.FRENCH
			default:
				return settings.ENGLISH
			}
		}

		language := getLanguage()
		gui.setLanguage(language)
		gui.Settings.Set("language", language)
		gui.Reload()
	}

	return container.NewHBox(
		widget.NewLabel(gui.Language.LanguageCode),
		selectLanguage,
	)
}

// createSelectCurrency create select to change currency
func (gui *GUI) createSelectCurrency() *fyne.Container {
	selectCurrency := widget.NewSelect(settings.GetCurrencies(), func(currency string) {
		gui.setCurrency(currency)
		gui.Settings.Set("currency", currency)
		gui.Reload()
	})
	currency, _ := gui.Currency.Get()
	selectCurrency.SetSelected(currency)
	return container.NewHBox(
		widget.NewLabel(gui.Language.Currency),
		selectCurrency,
	)
}

// createLabelLogs create label to show logs
func (gui *GUI) createLabelLogs() *fyne.Container {
	return container.NewHBox(
		widget.NewLabel(gui.Language.Logs),
		widget.NewLabel(config.LOGS_PATH),
	)
}

// createHelpMenu create help item in toolbar to show about app
func (gui *GUI) createHelpMenu() *fyne.Menu {
	url, _ := url.Parse(config.APP_LINK)

	gui.labelsAbout = binding.NewStringList()
	gui.labelsAbout.Set(gui.Language.GetAbouts())
	var labels []binding.DataItem
	for index := range gui.Language.GetAbouts() {
		about, _ := gui.labelsAbout.GetItem(index)
		labels = append(labels, about)
	}

	// Setup layouts with data
	firstLine := container.NewHBox(
		widget.NewLabelWithData(labels[0].(binding.String)),
		widget.NewLabel(config.APP_NAME),
		widget.NewLabelWithData(labels[1].(binding.String)),
	)
	secondLine := container.NewHBox(
		widget.NewLabelWithData(labels[2].(binding.String)),
		widget.NewHyperlink("GitHub Project", url),
		widget.NewLabelWithData(labels[3].(binding.String)),
	)
	thirdLine := widget.NewLabelWithData(labels[4].(binding.String))
	fourthLine := container.NewHBox(
		widget.NewLabel("Version:"),
		canvas.NewText(fmt.Sprintf("v%s", config.APP_VERSION), color.NRGBA{R: 218, G: 20, B: 51, A: 255}),
	)
	fifthLine := widget.NewLabel(fmt.Sprintf("%s: %s", gui.Language.Author, config.APP_AUTHOR))

	helpMenu := fyne.NewMenu(gui.Language.Help,
		fyne.NewMenuItem(gui.Language.About, func() {
			dialog.ShowCustom(gui.Language.About, gui.Language.Close,
				container.NewVBox(
					firstLine,
					secondLine,
					thirdLine,
					fourthLine,
					fifthLine,
				), gui.Window)
		}))
	return helpMenu
}

// getLanguageIndex get index to selectLanguage in settings from language of the app
func getLanguageIndex(langue string) int {
	switch langue {
	case settings.ENGLISH:
		return 0
	case settings.FRENCH:
		return 1
	default:
		return 0
	}
}

// setLogger create logger with zap librairy
func (gui *GUI) setLogger() {
	configZap := zap.NewProductionEncoderConfig()
	configZap.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(configZap)

	// Create logs folder if not exists
	path := "logs"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	logger := lumberjack.Logger{
		Filename:   config.LOGS_PATH, // File path
		MaxSize:    500,              // 500 megabytes per files
		MaxBackups: 3,                // 3 files before rotate
		MaxAge:     15,               // 15 days
	}

	writer := zapcore.AddSync(&logger)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, zapcore.DebugLevel),
	)

	gui.Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	gui.Logger.Info("Zap logger set",
		zap.String("path", logger.Filename),
		zap.Int("filesize", logger.MaxSize),
		zap.Int("backupfile", logger.MaxBackups),
		zap.Int("fileage", logger.MaxAge),
	)
}
