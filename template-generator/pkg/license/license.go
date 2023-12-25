package license

import (
	"fmt"
	"time"
)

type Header struct {
  SerialNo string
  Creator string
  CreatedTime string
  Country string
  Custom string
  Office string
  Sign string
}

type Opt struct {
  SerialNo string
  Creator string
  FeatureHeader *FeatureHeader
}

func NewHeader(opts *Opt) *Header {
  time := time.Now()
  createdTime := time.Format("2006-01-02 15:04:05")
  serialNo := fmt.Sprintf("LIC%d%02d%s", time.Year(), time.Month(), NewSign(6))
  creator := "Example Technologies Co., Ltd."
  sign := NewSign(512)

  if opts.SerialNo != "" {
    serialNo = opts.SerialNo
  }

  if opts.Creator != "" {
    creator = opts.Creator
  }

  return &Header{
    SerialNo: serialNo,
    Creator: creator,
    CreatedTime: createdTime,
    Country: "Unknow",
    Custom: "No relevant customer information",
    Office: "Unknow",
    Sign: sign,
  }
}

type FeatureHeader struct {
  Product string
  Feature string
  Esn string
  Attrib string
  Version string
  Libver string
  Sign string
  comment string
}

func DefaultFeatureHeader() *FeatureHeader {
  return &FeatureHeader{
    Sign: NewSign(128),
  }
}

func NewFeatureHeader(product string, feature string, esn string, attrib string, version string, libver string, comment string) *FeatureHeader {
  return &FeatureHeader{
      Product: product,
      Feature: feature,
      Esn: esn,
      Attrib: attrib,
      Libver: libver,
      Version: version,
      Sign: NewSign(128),
      comment: comment,
  }
}

type Feature struct {
  Product string
  Feature string
  Esn string
  Attrib string
  Function string
  Comment string
  Sign string
}

type License struct {
  Header Header
  FeatureHeader FeatureHeader
  Features map[string]*Feature
}

func (l *License) AddFeature(feature string, attrib string, function string) {
  feat := &Feature{
    Product: l.FeatureHeader.Product,
    Feature: feature,
    Esn: l.FeatureHeader.Esn,
    Attrib: attrib,
    Function: function,
    Comment: l.FeatureHeader.comment,
    Sign: NewSign(128),
  }
  l.Features[feature] = feat
}


func New(opts *Opt) *License {
  var fh *FeatureHeader
  if opts.FeatureHeader != nil {
    fh = opts.FeatureHeader
    fh.Sign = NewSign(128)
  } else {
    fh = DefaultFeatureHeader()
  }
  return &License{
    Header: *NewHeader(opts),
    FeatureHeader: *fh,
    Features: map[string]*Feature{},
  }
}


