package main

import (
	"io/ioutil"
    	"fmt"
    	"html/template"
    	"log"
	"net/http"
	"encoding/xml"
	"encoding/json"
	"math/rand"
	"time"
)

type Hotels struct {
	Hotel    []hotel  `xml:"hotel" json:"hotels"`
}
type hotel struct {
	Name    string   `xml:"name,attr" json:"hotel_name"`
    Rooms    []room  `xml:"rooms>room" json:"rooms"`
}

type room struct {
	Type    string  `xml:"type,attr" json:"room_type"`
	View    string  `xml:"view,attr" json:"room_view"`
	//XMLName  xml.Name `xml:"price"`
	Price   float64  `xml: "amount,attr" json:"price"`
	Currency   string `xml: "currency,attr" json:"currency"`
}



func roomsearch(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("./main.html")
    t.Execute(w, nil)
}
func RandomURL() string {
	URL := []string{
				   "http://mocky.io/v2/5e581e9b3000003129fd4064",
				   "http://mocky.io/v2/5e581ed83000003129fd4066",
				   "http://mocky.io/v2/5e581f14300000ec2bfd4069",
				   "http://mocky.io/v2/5e581f4c300000440cfd406a",
				   "http://mocky.io/v2/5e581f7a300000ec2bfd406c",
				   }
 	return URL[rand.Intn(len(URL))]
}

func numberOfDay(checkin, checkout time.Time) int {
	return int(checkout.Sub(checkin).Hours() / 24)
  }

func searchresult(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        t, _ := template.ParseFiles("./searchresult.html")
        t.Execute(w, nil)
    } else {
		r.ParseForm()
		client := &http.Client{}
		req, _:= http.NewRequest("POST", RandomURL(), nil)
		resp, err2 := client.Do(req)
		if err2 != nil {
			t, _ := template.ParseFiles("./error.html")
        	t.Execute(w, nil)
		}else{
			defer resp.Body.Close()
			byteValue, _ := ioutil.ReadAll(resp.Body)
			var hotels Hotels
			xml.Unmarshal(byteValue, &hotels)
			jsonByte, _:= json.Marshal(&hotels)
			json.Unmarshal(jsonByte, &hotels)
			checkin, _ := time.Parse("2006-01-02T15:04:05Z07:00", r.Form["checkin"][0]+"T00:00:00Z")
			checkout, _ := time.Parse("2006-01-02T15:04:05Z07:00", r.Form["checkout"][0]+"T00:00:00Z")

			numberOfDay := numberOfDay(checkin, checkout)
			fmt.Println("numberOfDay",numberOfDay)

			//price'ı xml'den çekemediğim için bir price belirleyip devam etmeye çalıştım.
			price := 225
			nightly_price := price / numberOfDay
			fmt.Println("nightly Price:",nightly_price)


			myMap := make(map[time.Time]int)
			date := checkin
			for i := 0;i<numberOfDay;i++ {
				myMap[date] = nightly_price
				date = date.AddDate(0, 0, 1)
			}
			myMap[checkout] = nightly_price + (price % numberOfDay)
			//myMap Objesi her otel fiyatı için hesaplanıp hotels objesine eklenmeli. 
		}
	}
}


func main() {
    http.HandleFunc("/", roomsearch)
	http.HandleFunc("/main", roomsearch)
	http.HandleFunc("/searchresult", searchresult)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
