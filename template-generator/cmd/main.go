package main

import (
	"log"
	"os"
	"template-generator/pkg/license"
	"template-generator/pkg/storage"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var DbPath = "product.dat"

func main() {

  h := storage.NewH()
  if err := h.InitDb(DbPath); err != nil {
    log.Fatalf("Failed to initialize database: %v", err)
  }

  opts := &license.Opt{
    FeatureHeader: license.NewFeatureHeader(
      "5000 V3",
      "Service",
      "2102351CMA10J8000008",
      "COMM, NULL, NULL, NULL, NULL, NULL",
      "V300R002",
      "1.2",
      "$***#SWID=1810-23B7-7ND7$***#,,,CX4LAJAFLA6X",
    ),
  }

  lic := license.New(opts)

  tmpl, err := template.ParseGlob("templates/*.tmpl")
  if err != nil {
    log.Fatalf("Failed parsing template: %v", err)
  }
  tmpl.Name()

  // products, err := h.GetProducts()
  // if err != nil {
  //   log.Fatalf("Failed to get products: %v", err)
  // }

  product := "5000 V3"
  features, err := h.GetProductFeatures(product)
  if err != nil {
    log.Fatalf("Failed to get features of %s product: %v", product, err)
  }

  for _, f := range features {
    lic.AddFeature(f.Name, "?", f.FuncName)
  }


  err = tmpl.Lookup("main").Execute(os.Stdout, lic)
  if err != nil {
    panic(err)
  }


}
