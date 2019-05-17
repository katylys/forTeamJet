package main

import (
	"fmt"
	"sort"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
)

type Result struct {
	Array []int `json:"array"`
	Uniq bool `json:"uniq"`
}

type MainResult struct {
	Array []int `json:"array"`
}

func Time(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		time := time.Now().UTC()
		w.Write([]byte(time.String()))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Method not implemented"))
	}
}

func unique(intSlice []int) []int {
    keys := make(map[int]bool)
    list := []int{}
    for _, entry := range intSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

func Sort(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error", 400)
			return
		}
		var result Result
		err = json.Unmarshal(body, &result)
		if err != nil || len(result.Array) > 100 || len(result.Array) == 0 {
			http.Error(w, "error", 400)
			return
		}
		sort.Ints(result.Array)
		if result.Uniq == true {
			t := result.Array
			result.Array = unique(t)
		}
		var main MainResult
		main.Array = result.Array
		b, _ := json.Marshal(main)
		w.Write([]byte(b))
	}
}

func Weather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
        case "GET":
		keys, ok := r.URL.Query()["city"]
		if !ok || len(keys[0]) < 1 {
			http.Error(w, "Url Param 'key' is missing", 400)
			return
		}
		key := keys[0]
		url := "https://api.openweathermap.org/data/2.5/weather?q=" + key + "&appid=e43304ba235e11ae949e4556aa58ae12"
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "error", 400)
                        return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
                        http.Error(w, "error", 400)
                }
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		if result["cod"] == "404" {
			http.Error(w, "error", 404)
                        return
		}
		main := result["main"].(map[string]interface{})
		str := fmt.Sprint(main["temp"])
		res := "Temperature in " + key + " : " + str
                w.Write([]byte(res))
	}
}

func main() {
	http.HandleFunc("/api/now", Time)
	http.HandleFunc("/api/sort", Sort)
	http.HandleFunc("/api/weather", Weather)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}














