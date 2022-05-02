package main

import (
	"context"

	"github.com/liumingmin/goutils/log"
	"github.com/xuri/excelize/v2"

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

	f, err := excelize.OpenFile("address.xlsx")
	if err != nil {
		log.Error(ctx, "OpenFile excel err: %v", err)
		return
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Error(ctx, "Close file excel err: %v", err)
		}
	}()

	for _, sheetName := range f.GetSheetList() {
		//if sheetName != "Sheet1" {
		//	return
		//}
		dataTable := ReadExcel(ctx, f, sheetName)
		if dataTable == nil {
			return
		}

		importData(ctx, dataTable, "城市", "地址", []string{""})
	}
}
