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

type ProfileTemplateNoStream struct {
	User bool
}

func (p *ProfileTemplateNoStream) Render(w http.ResponseWriter) error {
	err := userTmpl.ExecuteTemplate(w, "profile", struct {
		User bool
	}{
		User: p.User,
	})
	if err != nil {
		return err
	}

	return nil
}

type Product struct {
	Name        string
	Description string
	Image       string
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

type ProductsTemplateNoStream struct {
	Products []Product
}

func (p *ProductsTemplateNoStream) Render(w http.ResponseWriter) error {
	err := productsTmpl.ExecuteTemplate(w, "products", struct {
		Products []Product
	}{
		Products: p.Products,
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

type CategoriesTemplateNoStream struct {
	Categories []string
}

func (c *CategoriesTemplateNoStream) Render(w http.ResponseWriter) error {
	err := categoriesTmpl.ExecuteTemplate(w, "categories", struct {
		Categories []string
	}{
		Categories: c.Categories,
	})
	if err != nil {
		return err
	}
	return nil
}

func slowBoolFetch(t time.Duration) chan bool {
	c := make(chan bool)
	go func() {
		time.Sleep(t * time.Millisecond)
		c <- true
	}()

	return c
}

func slowFetchProducts(t time.Duration) chan []Product {
	c := make(chan []Product)
	go func() {
		time.Sleep(t * time.Millisecond)
		c <- []Product{
			{
				Name:        "Product 1",
				Description: "Product 1 Description",
				Image:       "/static/camera.jpeg",
			},
			{
				Name:        "Product 2",
				Description: "Product 2 Description",
				Image:       "/static/cbd.jpeg",
			},
			{
				Name:        "Product 3",
				Description: "Product 3 Description",
				Image:       "/static/coke.jpeg",
			},
			{
				Name:        "Product 4",
				Description: "Product 4 Description",
				Image:       "/static/cream.webp",
			},
			{
				Name:        "Product 5",
				Description: "Product 5 Description",
				Image:       "/static/laptop.jpg",
			},
			{
				Name:        "Product 6",
				Description: "Product 6 Description",
				Image:       "/static/phone.jpeg",
			},
		}
	}()

	return c
}

func slowFetchCategories(t time.Duration) chan []string {
	c := make(chan []string)
	go func() {
		time.Sleep(t * time.Millisecond)
		c <- []string{
			"Category 1",
			"Category 2",
		}
	}()

	return c
}

func index(w http.ResponseWriter, req *http.Request) {
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
				User: slowBoolFetch(time.Duration(rand.Intn(1000))),
			},
			W:       w,
			Flusher: flusher,
		},
		Products: &TemplateLoader{
			ID: "products",
			Template: &ProductsTemplate{
				Products: slowFetchProducts(time.Duration(rand.Intn(4000))),
			},
			W:       w,
			Flusher: flusher,
		},
		Categories: &TemplateLoader{
			ID: "categories",
			Template: &CategoriesTemplate{
				Categories: slowFetchCategories(time.Duration(rand.Intn(3000))),
			},
			W:       w,
			Flusher: flusher,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func indexNoStream(w http.ResponseWriter, req *http.Request) {
	c1 := <-slowBoolFetch(time.Duration(rand.Intn(1000)))
	c2 := <-slowFetchProducts(time.Duration(rand.Intn(4000)))
	c3 := <-slowFetchCategories(time.Duration(rand.Intn(3000)))
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
			Template: &ProfileTemplateNoStream{
				User: c1,
			},
			W:       w,
			Flusher: flusher,
		},
		Products: &TemplateLoader{
			ID: "products",
			Template: &ProductsTemplateNoStream{
				Products: c2,
			},
			W:       w,
			Flusher: flusher,
		},
		Categories: &TemplateLoader{
			ID: "categories",
			Template: &CategoriesTemplateNoStream{
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
	http.HandleFunc("/stream", index)
	http.HandleFunc("/nostream", indexNoStream)

	fmt.Println("Stream version available at http://localhost:1337/stream")
	fmt.Println("No stream version available at http://localhost:1337/nostream")
	if err := http.ListenAndServe("localhost:1337", nil); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
