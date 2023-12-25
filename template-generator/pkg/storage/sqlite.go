package storage

import (
	"database/sql"
)

type H struct {
  Db *sql.DB
}

func NewH() *H {
  return &H{}
}

func (h *H) InitDb(filename string) error {
  db, err := sql.Open("sqlite3", filename)
  if err != nil {
    return err
  }

  h.Db = db
  return nil
}

func (h *H) GetProducts() ([]string, error) {
  rows, err := h.Db.Query("SELECT DISTINCT(Product) FROM LicFeatureCfg ORDER BY Product")
  if err != nil {
    return nil, err
  }

  var products []string
  for rows.Next() {
    var product string
    rows.Scan(&product)
    products = append(products, product)
  }

  return products, nil
}

type Feature struct {
  ID int64
  Name string
  FuncName string
}

func (h *H) GetProductFeatures(product string) ([]Feature, error) {
  rows, err := h.Db.Query("SELECT FeatureID, FeatureName, FeatureFuncName FROM LicFeatureCfg WHERE Product = ? ORDER BY FeatureID", product)
  if err != nil {
    return nil, err
  }

  var features []Feature
  for rows.Next() {
    f := Feature{}
    rows.Scan(&f.ID, &f.Name, &f.FuncName)
    features = append(features, f)
  }

  return features, nil

}
