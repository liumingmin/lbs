package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/liumingmin/goutils/log"

	"github.com/liumingmin/goutils/container"
	"github.com/mozillazg/go-pinyin"
	"go.mongodb.org/mongo-driver/bson"
)

//db.getCollection('location').createIndex({"geo":"2dsphere"})

func importData(ctx context.Context, dataTable *container.DataTable, cityCol, addressCol string, expandCols []string) {
	lastCity := ""
	cityRows := make([]*container.DataRow, 0)

	for i := 0; i < len(dataTable.Rows()); i++ {
		row := dataTable.Rows()[i]

		city := row.String(cityCol)
		city = strings.Split(city, "-")[0]

		address := row.String(addressCol)
		if strings.TrimSpace(address) == "" {
			continue
		}

		if city == lastCity && len(cityRows) < 10 {
			cityRows = append(cityRows, row)
		} else {
			if len(cityRows) > 0 && lastCity != "" {
				locData(ctx, dataTable, cityCol, addressCol, expandCols, cityRows, lastCity)
			}

			lastCity = city
			cityRows = make([]*container.DataRow, 0)
			cityRows = append(cityRows, row)
		}
	}
	if len(cityRows) > 0 {
		locData(ctx, dataTable, cityCol, addressCol, expandCols, cityRows, lastCity)
	}

}

func locData(ctx context.Context, dataTable *container.DataTable, cityCol, addressCol string, expandCols []string,
	cityRows []*container.DataRow, city string) {
	addresses := make([]string, 0, len(cityRows))
	for _, row := range cityRows {
		address := strings.TrimSpace(row.String(addressCol))

		if address == "" {
			continue
		}
		addresses = append(addresses, address)
	}

	locations := amap.GeoByAddresses(ctx, addresses, city)
	if len(locations) == 0 {
		return
	}

	addressIndex := 0

	for _, row := range cityRows {
		address := strings.TrimSpace(row.String(addressCol))
		if address == "" {
			continue
		}

		if addressIndex >= len(locations) {
			continue
		}

		location := strings.TrimSpace(locations[addressIndex])
		addressIndex++
		if location == "" {
			continue
		}

		item := bson.M{}
		item["dataId"] = row.String(dataTable.PkCol())
		item["city"] = row.String(cityCol)
		item["address"] = address

		for _, expandCol := range expandCols {
			item[pinyinToString(expandCol)] = row.String(expandCol)
		}

		item["source"] = sliceToMap(row.Data(), dataTable.Cols())

		item["geo"] = &LocationPoint{
			Type:       "Point",
			Cordinates: sliceAtof(strings.Split(location, ",")),
		}

		err := collection.Insert(ctx, item)
		if err != nil {
			log.Error(ctx, "insert to mongo err: %v", err)
		}
	}
}

func pinyinToString(expandCol string) string {
	s := pinyin.Pinyin(expandCol, pinyinArgs)

	result := ""
	for _, c := range s {
		for _, l := range c {
			result += l
		}
	}
	return result
}

func sliceToMap(row []string, cols []string) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(row); i++ {
		if i>= len(cols){
			break
		}
		result[cols[i]] = row[i]
	}
	return result
}

func sliceAtof(a []string) []float64 {
	fs := make([]float64, 0)
	for _, s := range a {
		f, _ := strconv.ParseFloat(s, 64)
		fs = append(fs, f)
	}
	return fs
}

type LocationPoint struct {
	Type       string    `bson:"type"`
	Cordinates []float64 `bson:"coordinates"`
}

var pinyinArgs pinyin.Args

func init() {
	pinyinArgs = pinyin.NewArgs()
	pinyinArgs.Style = pinyin.FIRST_LETTER
}
