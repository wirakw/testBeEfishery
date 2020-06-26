package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	hp "testEFishery/helper"
	md "testEFishery/model"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type m map[string]interface{}

const refresh = 10 * time.Second

func registrasiUser(ir iris.Context) {
	request := md.User{}
	err := ir.ReadJSON(&request)
	status := md.ResponseStatus{}
	if err != nil {
		ir.StatusCode(iris.StatusBadRequest)
		ir.WriteString(err.Error())
		return
	}
	if (request.Name == "") || (request.Role == "") || (request.Phone == "") {
		ir.StatusCode(iris.StatusBadRequest)
		status.Status = "BadRequest"
		status.Message = "value from name or role should't empty"
		return
	}
	users := md.User{}
	isUserNotExist := db.Where("users.name = ?", request.Name).First(&users).RecordNotFound()
	if isUserNotExist == true {
		request.Password = hp.GenerateID(4)
		db.Create(&request)
		a := db.NewRecord(request)
		if a == true {
			ir.StatusCode(iris.StatusBadRequest)
			status.Status = "BadRequest"
			status.Message = "Fail Create User"
			json.NewEncoder(ir).Encode(status)
			return
		}
		token := hp.GenerateToken(request.Name, request.Role, request.Password)
		request.Token = token
		response := m{
			"status": "success",
			"data":   request,
		}
		json.NewEncoder(ir).Encode(response)
	} else {
		ir.StatusCode(iris.StatusConflict)
		status.Status = "Conflict"
		status.Message = "User with this name already exist"
		json.NewEncoder(ir).Encode(status)
	}
}

func getTokenHandler(ir iris.Context) {
	request := md.Login{}
	err := ir.ReadJSON(&request)
	status := md.ResponseStatus{}
	if err != nil {
		ir.StatusCode(iris.StatusBadRequest)
		ir.WriteString(err.Error())
		return
	}
	if (request.Phone == "") || (request.Password == "") {
		ir.StatusCode(iris.StatusBadRequest)
		status.Status = "BadRequest"
		status.Message = "value from name or role should't empty"
		return
	}
	user := md.User{}
	a := db.Where("users.phone = ?", request.Phone).Where("users.password = ?", request.Password).First(&user)
	userNotFound := a.RecordNotFound()
	if userNotFound == true {
		ir.StatusCode(iris.StatusNotFound)
		status.Status = "Not Found"
		status.Message = "wrong user phone or password"
		json.NewEncoder(ir).Encode(status)
		return
	}
	token := hp.GenerateToken(user.Name, user.Role, user.Password)
	user.Token = token
	response := m{
		"status": "success",
		"data":   user,
	}
	json.NewEncoder(ir).Encode(response)
}

func fetchAllData(ir iris.Context) {
	status := md.ResponseStatus{}
	var storages []md.Storage
	response, err := http.Get("https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		ir.StatusCode(iris.StatusInternalServerError)
		status.Status = "Internal Server Error"
		status.Message = "Fail fetch data from https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list"
		json.NewEncoder(ir).Encode(status)
		return
	}

	usdCurrency, err := http.Get("https://free.currconv.com/api/v7/convert?q=IDR_USD&compact=ultra&apiKey=b866a289f3c0936aba88")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		ir.StatusCode(iris.StatusInternalServerError)
		status.Status = "Internal Server Error"
		status.Message = "Fail load api free.currconv.com"
	}
	usd, _ := ioutil.ReadAll(usdCurrency.Body)
	var Kurs struct {
		IdrUsd float64 `json:"IDR_USD"`
	}
	if err := json.Unmarshal(usd, &Kurs); err != nil {
		fmt.Println(err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(data, &storages); err != nil {
		fmt.Println(err)
	}
	for i := range storages {
		price, _ := strconv.ParseFloat(storages[i].Price.String, 64)
		usdPrice := price * Kurs.IdrUsd
		storages[i].USDPrice = fmt.Sprintf("%f", usdPrice)
	}
	json.NewEncoder(ir).Encode(storages)

}

func aggregatDataByArea(ir iris.Context) {
	isAdmin := ir.Values().GetString("role")

	status := md.ResponseStatus{}
	if isAdmin != "admin" {
		ir.StatusCode(iris.StatusUnauthorized)
		status.Status = "Unathorized"
		status.Message = "Role is not Admin"
		json.NewEncoder(ir).Encode(status)
		return
	}
	var storages []md.Storage
	response, err := http.Get("https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		ir.StatusCode(iris.StatusInternalServerError)
		status.Status = "Internal Server Error"
		status.Message = "Fail fetch data from https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list"
		json.NewEncoder(ir).Encode(status)
		return
	}
	data, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(data, &storages); err != nil {
		fmt.Println(err)
	}
	var newStorages []md.AggregateStorage
	for i := range storages {
		singleStorage := md.AggregateStorage{}
		dataProv := md.DataProv{}
		singleStorage.AreaProvinsi = storages[i].AreaProvinsi
		singleStorage.TglParsed = storages[i].TglParsed
		size, _ := strconv.ParseFloat(storages[i].Size.String, 64)
		var isExist bool = false
		for j := 0; j < len(newStorages); j++ {
			var same bool = false
			same = hp.DateEqual(storages[i].TglParsed.String, newStorages[j].TglParsed.String)
			if (newStorages[j].AreaProvinsi == storages[i].AreaProvinsi) && (same == true) {
				isExist = true
				if newStorages[j].Min > size {
					newStorages[j].Min = size
				}

				if newStorages[j].Max < size {
					newStorages[j].Max = size
				}
				dataProv.AreaKota = storages[i].AreaKota
				dataProv.Komoditas = storages[i].Komoditas
				dataProv.Price = storages[i].Price
				dataProv.Size = storages[i].Size
				dataProv.Timestamp = storages[i].Timestamp
				dataProv.UUID = storages[i].UUID
				newStorages[j].Data = append(newStorages[j].Data, dataProv)
				break
			}
		}

		if isExist == false {
			dataProv.AreaKota = storages[i].AreaKota
			dataProv.Komoditas = storages[i].Komoditas
			dataProv.Price = storages[i].Price
			dataProv.Size = storages[i].Size
			dataProv.Timestamp = storages[i].Timestamp
			dataProv.UUID = storages[i].UUID
			singleStorage.Data = append(singleStorage.Data, dataProv)
			singleStorage.Min = size
			singleStorage.Max = size
			newStorages = append(newStorages, singleStorage)
		}

		for i := range newStorages {
			var count float64
			var total float64
			sort.Slice(newStorages[i].Data, func(x, y int) bool {
				a, _ := strconv.ParseFloat(newStorages[i].Data[x].Size.String, 64)
				b, _ := strconv.ParseFloat(newStorages[i].Data[y].Size.String, 64)
				return a < b
			})
			middle := len(newStorages[i].Data) / 2

			if len(newStorages[i].Data)%2 == 1 {
				newStorages[i].Median, _ = strconv.ParseFloat(newStorages[i].Data[middle].Size.String, 64)
			} else {
				a, _ := strconv.ParseFloat(newStorages[i].Data[middle-1].Size.String, 64)
				b, _ := strconv.ParseFloat(newStorages[i].Data[middle].Size.String, 64)
				newStorages[i].Median = (a + b) / 2
			}
			for j := range newStorages[i].Data {
				count++
				size, _ := strconv.ParseFloat(newStorages[i].Data[j].Size.String, 64)
				total = total + size
			}
			newStorages[i].Avg = total / count
		}
	}
	json.NewEncoder(ir).Encode(newStorages)
}

func getClaimPriavte(ir iris.Context) {
	jwtToken := ir.Values().Get("jwt").(*jwt.Token)

	foobar := jwtToken.Claims.(jwt.MapClaims)
	response := m{
		"status": "success",
		"data": m{
			"name":      foobar["nama"],
			"password":  foobar["password"],
			"role":      foobar["role"],
			"timestamp": foobar["timestamp"],
		},
	}
	json.NewEncoder(ir).Encode(response)
}

func main() {
	app := iris.Default()
	app.Use(recover.New())
	app.Use(iris.Cache304(refresh)) // uncomment this for cache all route in api

	initDb()
	migrate()
	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type")
		ctx.Header("Content-Type", "application/json")
		ctx.Next()
	}
	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return hp.MySecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	// ctx := context.Background()
	v1 := app.Party("/api/v1/", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Post("registrasi", registrasiUser)
		v1.Post("login", getTokenHandler)
		v1.Get("secured", j.Serve, getClaimPriavte)
		// v1.Get("data/withcache", cache.Handler(10*time.Second), fetchAllData) //fetch with cache
		v1.Get("data", j.Serve, fetchAllData)                                    //fetch without cache
		v1.Get("aggregate", j.Serve, myAuthenticatedHandler, aggregatDataByArea) //fetch without cache
	}

	app.Run(iris.Addr(":8080"),
		iris.WithoutPathCorrection,
		iris.WithoutPathCorrectionRedirection,
		iris.WithoutServerError(iris.ErrServerClosed))
}
