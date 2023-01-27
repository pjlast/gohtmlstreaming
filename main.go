package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

//go:embed templates/base.tmpl
var baseTemplate string
var baseTmpl = template.Must(template.New("base").Parse(baseTemplate))

//go:embed templates/load.tmpl
var loaderTemplate string
var loadTmpl = template.Must(template.New("loader").Parse(loaderTemplate))

//go:embed templates/user.tmpl
var userTemplate string
var userTmpl = template.Must(template.New("user").Parse(userTemplate))

//go:embed templates/products.tmpl
var productsTemplate string
var productsTmpl = template.Must(template.New("products").Parse(productsTemplate))

//go:embed templates/categories.tmpl
var categoriesTemplate string
var categoriesTmpl = template.Must(template.New("categories").Parse(categoriesTemplate))

type LoadableTemplate interface {
	Render(http.ResponseWriter) error
}

type TemplateLoader struct {
	ID       string
	W        http.ResponseWriter
	Flusher  http.Flusher
	Template LoadableTemplate
}

func (t *TemplateLoader) Load() error {
	err := loadTmpl.ExecuteTemplate(t.W, "load", t)
	if err != nil {
		return err
	}
	return nil
}

func (t *TemplateLoader) Flush() error {
	t.Flusher.Flush()
	return nil
}

type ProfileTemplate struct {
	User chan bool
}

func (p *ProfileTemplate) Render(w http.ResponseWriter) error {
	err := userTmpl.ExecuteTemplate(w, "profile", struct {
		User bool
	}{
		User: <-p.User,
	})
	if err != nil {
		return err
	}

	return nil
}

type Product struct {
	Name        string
	Description string
}

type ProductsTemplate struct {
	Products chan []Product
}

func (p *ProductsTemplate) Render(w http.ResponseWriter) error {
	err := productsTmpl.ExecuteTemplate(w, "products", struct {
		Products []Product
	}{
		Products: <-p.Products,
	})
	if err != nil {
		return err
	}

	return nil
}

type CategoriesTemplate struct {
	Categories chan []string
}

func (c *CategoriesTemplate) Render(w http.ResponseWriter) error {
	err := categoriesTmpl.ExecuteTemplate(w, "categories", struct {
		Categories []string
	}{
		Categories: <-c.Categories,
	})
	if err != nil {
		return err
	}
	return nil
}

func slowBoolFetch(t time.Duration, c chan bool) {
	time.Sleep(t * time.Millisecond)
	c <- true
}

func slowFetchProducts(t time.Duration, c chan []Product) {
	time.Sleep(t * time.Millisecond)
	c <- []Product{
		{
			Name:        "Product 1",
			Description: "Product 1 Description",
		},
		{
			Name:        "Product 2",
			Description: "Product 2 Description",
		},
		{
			Name:        "Product 3",
			Description: "Product 3 Description",
		},
		{
			Name:        "Product 4",
			Description: "Product 4 Description",
		},
		{
			Name:        "Product 5",
			Description: "Product 5 Description",
		},
		{
			Name:        "Product 6",
			Description: "Product 6 Description",
		},
	}
}

func slowFetchCategories(t time.Duration, c chan []string) {
	time.Sleep(t * time.Millisecond)
	c <- []string{
		"Category 1",
		"Category 2",
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	c1 := make(chan bool)
	c2 := make(chan []Product)
	c3 := make(chan []string)
	go slowBoolFetch(time.Duration(rand.Intn(1000)), c1)
	go slowFetchProducts(time.Duration(rand.Intn(4000)), c2)
	go slowFetchCategories(time.Duration(rand.Intn(3000)), c3)
	flusher, _ := w.(http.Flusher)
	err := baseTmpl.ExecuteTemplate(w, "base", struct {
		Title      string
		Profile    *TemplateLoader
		Products   *TemplateLoader
		Categories *TemplateLoader
	}{
		Title: "Hello World",
		Profile: &TemplateLoader{
			ID: "profile",
			Template: &ProfileTemplate{
				User: c1,
			},
			W:       w,
			Flusher: flusher,
		},
		Products: &TemplateLoader{
			ID: "products",
			Template: &ProductsTemplate{
				Products: c2,
			},
			W:       w,
			Flusher: flusher,
		},
		Categories: &TemplateLoader{
			ID: "categories",
			Template: &CategoriesTemplate{
				Categories: c3,
			},
			W:       w,
			Flusher: flusher,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", index)

	if err := http.ListenAndServe("localhost:1337", nil); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
