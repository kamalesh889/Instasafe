package Rout

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

func Router() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/transactions", Transaction).Methods("POST")
	r.HandleFunc("/statistics", Statistics).Methods("GET")
	r.HandleFunc("/delete", Delete).Methods("DELETE")
	return r

}

var Requestbody Request

var Reqq []Request

func Transaction(res http.ResponseWriter, req *http.Request) {

	err := json.NewDecoder(req.Body).Decode(&Requestbody)
	if err != nil {
		fmt.Println("error in parsing", err)
		res.WriteHeader(http.StatusUnprocessableEntity) //Fields are not parsable
		return
	}

	v := validator.New()

	a := Request{
		Amount:    Requestbody.Amount,
		Timestamp: Requestbody.Timestamp,
	}

	err = v.Struct(a)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("validation error", e) //Validation error
		}
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	fmt.Println("Location : ", loc, " Time : ", now)

	if now.Sub(Requestbody.Timestamp).Seconds() > 60 {
		res.WriteHeader(http.StatusNoContent) //If transaction is older than 60 seconds
		return
	}

	if now.Sub(Requestbody.Timestamp).Seconds() < 0 { //If transaction time is in future
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	fmt.Println(Requestbody.Timestamp)
	Reqq = append(Reqq, Requestbody)
	fmt.Println("listt", Reqq)
	res.WriteHeader(http.StatusCreated)
}

func Statistics(res http.ResponseWriter, req *http.Request) {

	var Stats Stats
	if len(Reqq) > 0 {
		var Arr []float64
		Counter := 0
		loc, _ := time.LoadLocation("Asia/Kolkata")
		now := time.Now().In(loc)
		fmt.Println("Location : ", loc, " Time : ", now)

		fmt.Println("Gett", Reqq)
		for _, value := range Reqq {
			if now.Sub(value.Timestamp).Seconds() < 60 {
				Arr = append(Arr, value.Amount)
				Counter = Counter + 1
				Stats.Sum = Stats.Sum + value.Amount //Total amount
			}
		}
		fmt.Println("Total Amount is ", Stats.Sum)
		Stats.Avg = Stats.Sum / float64(Counter) //Avearge amount
		fmt.Println("Avearge Amount is", Stats.Avg)

		Stats.Min, Stats.Max = Find(Arr) //Max & min amount
		fmt.Println(Stats.Max, Stats.Min)

		Stats.Count = len(Arr) //No of transactions
		fmt.Println("Transaction count", Stats.Count)

		Final_resp, err := json.Marshal(Stats)
		if err != nil {
			fmt.Println(err)
		}

		_, err = res.Write(Final_resp)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		Stats.Sum = 0
		Stats.Avg = 0
		Stats.Max = 0
		Stats.Min = 0
		Stats.Count = 0
		Final_resp, err := json.Marshal(Stats)
		if err != nil {
			fmt.Println(err)
		}

		_, err = res.Write(Final_resp)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Delete(res http.ResponseWriter, req *http.Request) {
	Reqq = nil
	fmt.Println("All txn clear", Reqq)
	res.WriteHeader(http.StatusNoContent)
}

func Find(arr []float64) (float64, float64) {
	min := arr[0]
	max := arr[0]

	for _, value := range arr {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}

	}
	return min, max

}

// con, err := json.Marshal(Reqq)
// if err != nil {
// 	fmt.Println(err)
// }
// err = ioutil.WriteFile("instasfe.json", con, 0644)
// if err != nil {
// 	fmt.Println(err)
// }
