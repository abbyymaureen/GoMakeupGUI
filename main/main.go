package main

/**
@author abbybrown
@date 03.20.24
@filename main.go

	This Go program loads in a json file to store as a Makeup Product struct. The user can then
	select various makeup products and learn more about what they are, how much they cost, and how to use
	them, making the makeup user's life easier!

Sources:
	https://www.sohamkamani.com/golang/json/#structured-data-decoding-json-into-structs
	https://stackoverflow.com/questions/75206234/for-go-ioutil-readall-ioutil-readfile-ioutil-readdir-deprecated
	https://docs.fyne.io/started/
	https://github.com/fyne-io/fyne/issues/3337
	https://pkg.go.dev/golang.org/x/text/currency
	ChatGPT
*/

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"os"
	"strconv"
)

/*
Makeup Struct
-------------

	Structure that defines the elements within a makeup product. Important features include
	Brand, Name, Price, and Description. This struct corresponds with the products.json file.
*/
type Makeup struct {
	ID            int
	Brand         string
	Name          string
	Price         string
	PriceSign     string
	Currency      string
	ImageLink     string
	ProductLink   string
	WebsiteLink   string
	Description   string
	Rating        float32
	Category      string
	ProductType   string
	TagList       []string
	CreatedAt     string
	UpdatedAt     string
	ProductAPIURL string
	APIImage      string
	ProductColors []struct {
		HexValue  string
		ColorName string
	}
}

/*
ReadFile
--------

@param filename : string

	The filename for reading in the json file

@returns products : []Makeup

	A list of makeup products
*/
func readFile(filename string) []Makeup {
	// Open the JSON file
	file, err := os.ReadFile(filename)

	// Catch any potential errors and print these on the console (sudo-logging)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return nil
	}

	// Define a variable to hold the data
	var products []Makeup

	// 'Unmarshal' JSON data into the struct
	err = json.Unmarshal(file, &products)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	// Return a list of all the Makeup products
	return products
}

/*
The Main Function for running our program

Generates a GUI displaying numerous makeup products for user selection
*/
func main() {
	products := readFile("main/products.json")

	// Create a GUI window and set the title
	a := app.New()
	w := a.NewWindow("Makeup Database")

	w.SetContent(widget.NewLabel("Products"))

	// Create a list object that contains all the makeup products
	list := widget.NewList(
		func() int {
			return len(products)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Products")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(products[i].Name)
		})

	// Create details widgets
	nameLabel := widget.NewLabel("")
	brandLabel := widget.NewLabel("")
	priceLabel := widget.NewLabel("")
	descriptionLabel := widget.NewLabel("")
	descriptionLabel.Wrapping = fyne.TextWrapWord // Enable text wrapping for description
	colorLabel := widget.NewLabel("")

	// Create a container for the details widgets
	detailsContainer := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			widget.NewLabel("Name: "), nameLabel),
		container.New(layout.NewHBoxLayout(),
			widget.NewLabel("Brand: "), brandLabel),
		container.New(layout.NewHBoxLayout(),
			widget.NewLabel("Price: "), priceLabel),
		container.New(layout.NewHBoxLayout(),
			widget.NewLabel("Description: "), container.NewVScroll(descriptionLabel)),
		container.New(layout.NewHBoxLayout(),
			widget.NewLabel("Color: "), colorLabel),
	)

	// Handle selection event
	list.OnSelected = func(id widget.ListItemID) {
		selectedProduct := products[id]

		nameLabel.SetText(selectedProduct.Name)
		brandLabel.SetText(selectedProduct.Brand)

		// Parse the price string into a float64
		price, err := strconv.ParseFloat(selectedProduct.Price, 64)
		if err != nil {
			// Handle error if price cannot be parsed
			fmt.Println("Error parsing price:", err)
			return
		}

		// Format the price label with a dollar sign
		formattedPrice := fmt.Sprintf("$%.2f", price)

		// Set the text of the price label
		priceLabel.SetText(formattedPrice)
		descriptionLabel.SetText(selectedProduct.Description)
		if len(selectedProduct.ProductColors) > 0 {
			colorLabel.SetText(selectedProduct.ProductColors[0].ColorName)
		} else {
			colorLabel.SetText("N/A")
		}
	}

	// Create a grid layout with the list on the left and details on the right
	grid := container.New(layout.NewGridLayout(2), list, detailsContainer)

	w.SetContent(grid)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
