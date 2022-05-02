package main

import (
	"context"

	"github.com/liumingmin/goutils/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var locationOp *mongo.CompCollectionOp

func InitOps() {
	mongo.InitClients()
	c, _ := mongo.MgoClient("lbs")
	locationOp = mongo.NewCompCollectionOp(c, "lbs", "location")
}

type LocationReq struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
	Size int     `json:"size"`
	Min  int     `json:"min"`
	Max  int     `json:"max"`
}

type GeoLocation struct {
	Address string         `json:"address" bson:"address"`
	City    string         `json:"city" bson:"city"`
	DataId  string         `json:"dataId" bson:"dataId"`
	Source  interface{}    `json:"source" bson:"source"`
	Geo     *LocationPoint `json:"geo" bson:"geo"`
}

type LocationPoint struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

func QueryByLoc(ctx context.Context, reqVal interface{}) (interface{}, error) {
	req := reqVal.(*LocationReq)
	if req.Size <= 0 || req.Size > 200 {
		req.Size = 20
	}

	if req.Min < 0 {
		req.Min = 0
	}

	if req.Max < 0 || req.Max > 100000 {
		req.Max = 1000
	}

	var locations []*GeoLocation
	err := locationOp.Find(ctx, mongo.FindModel{
		Query: bson.M{"geo": bson.M{"$near": bson.M{
			"$geometry": bson.M{
				"type":        "Point",
				"coordinates": []float64{req.Long, req.Lat},
			},
			"$maxDistance": req.Max,
			"$minDistance": req.Min,
		}}},
		Cursor:  0,
		Size:    req.Size,
		Results: &locations,
	})

	return locations, err
}

var geoQuery = `
db.getCollection('location').find({"geo":{
    $near: {
       $geometry: {
          type: "Point" ,
          coordinates: [ 117.139436,31.646295 ]
       },
       $maxDistance: 40000,
       $minDistance: 0
     }
    }})
`

var mgoRun = `H:\mongodb4.2>bin\mongo.exe --port 27017`
var baiduPick = `http://api.map.baidu.com/lbsapi/getpoint/index.html`
