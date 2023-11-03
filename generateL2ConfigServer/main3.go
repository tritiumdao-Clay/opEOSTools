package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// The `json:"whatever"` bit is a way to tell the JSON
// encoder and decoder to use those names instead of the
// capitalised names

func tomHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		dataA := make([]byte, 4000)
		r.Body.Read(dataA)
		fmt.Println("debug0", string(dataA))

		{
			client := &http.Client{}
			var data = strings.NewReader(string(dataA))
			req, err := http.NewRequest("POST", "https://api.testnet.evm.eosnetwork.com", data)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			bodyText, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			//bodyText = bodyText[:len(bodyText)-1]
			w.Header().Set("Content-Type", "application/json")
			fmt.Println("----------")
			fmt.Printf("%s", bodyText[:len(bodyText)-1])
			//fmt.Printf(string(bodyText))
			fmt.Println("----------")
			type JsonResult struct {
				Id      uint64 `json:"id"`
				Jsonrpc string `json:"jsonrpc"`
				Result  string `json:"result"`
			}
			var jsonResult JsonResult
			err = json.Unmarshal(bodyText, &jsonResult)
			if err != nil {
				fmt.Println("debug2", err.Error())
			}
			bod, err := json.Marshal(jsonResult)
			if err != nil {
				fmt.Println("debug3", err.Error())
			}
			fmt.Println("------2----")
			fmt.Printf("%s", bod)
			//fmt.Printf(string(bodyText))
			fmt.Println("------2----")
			//io.WriteString(w, string(bod))
			io.WriteString(w, string(bodyText)+"\n")
			fmt.Println("------3----")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	http.HandleFunc("/tom", tomHandler)
	http.HandleFunc("/", tomHandler)

	log.Println("Go!")
	http.ListenAndServe(":8080", nil)
}