package data

import (
	"crypto/x509"
)

var (
	Roots *x509.CertPool = x509.NewCertPool()
	Proxy Proxys
	Con   Config
	Csrf  string
)

type Config struct {
	ROBLOSECURITY string  `json:"roblosecurity."`
	Percentage    float64 `json:"percentage_range"`
	BuyersProfit  float64 `json:"profitpercentage_determinedtobuy"`
	PreferredIds  []int64 `json:"preferredassets"`
}

type SellerInfo struct {
	PreviousPageCursor interface{} `json:"previousPageCursor"`
	NextPageCursor     string      `json:"nextPageCursor"`
	Data               []BuyerInfo `json:"data"`
}

type BuyerInfo struct {
	UserAssetID  int64       `json:"userAssetId"`
	Seller       Seller      `json:"seller"`
	Price        int         `json:"price"`
	SerialNumber interface{} `json:"serialNumber"`
}

type Seller struct {
	HasVerifiedBadge bool   `json:"hasVerifiedBadge"`
	ID               int    `json:"id"`
	Type             string `json:"type"`
	Name             string `json:"name"`
}

type Purchased struct {
	Purchased        bool   `json:"purchased"`
	Reason           string `json:"reason"`
	ProductID        int    `json:"productId"`
	StatusCode       int    `json:"statusCode"`
	Title            string `json:"title"`
	ErrorMsg         string `json:"errorMsg"`
	ShowDivID        string `json:"showDivId"`
	ShortfallPrice   int    `json:"shortfallPrice"`
	BalanceAfterSale int    `json:"balanceAfterSale"`
	ExpectedPrice    int    `json:"expectedPrice"`
	Currency         int    `json:"currency"`
	Price            int    `json:"price"`
	AssetID          int64  `json:"assetId"`
}
