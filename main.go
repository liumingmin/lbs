package main

import (
	"context"

	"github.com/liumingmin/goutils/db/mongo"
)

var amap *AMap
var collection *mongo.CompCollectionOp

func main() {
	ctx := context.Background()
	amap = &AMap{Key: ""}

	mongo.InitClients()
	c, _ := mongo.MgoClient("lbs")

	collection = mongo.NewCompCollectionOp(c, "lbs", "location")

	dataTable := ReadExcel(ctx, "address.xlsx", "Sheet1")
	if dataTable == nil {
		return
	}

	importData(ctx, dataTable, "城市", "地址", []string{""})
}
