package main

func main(){

}

var geoQuery =`
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
var baiduPick =`http://api.map.baidu.com/lbsapi/getpoint/index.html`
