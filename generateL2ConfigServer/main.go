package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
)

// The `json:"whatever"` bit is a way to tell the JSON
// encoder and decoder to use those names instead of the
// capitalised names
type person struct {
        Name string `json:"name"`
        Age int `json:"age"`
}

var tom *person = &person{
        Name: "Tom",
        Age:  28,
}

func tomHandler(w http.ResponseWriter, r *http.Request) {

        switch r.Method {
        case "POST":
                // Just send out the JSON version of 'tom'
                j, _ := json.Marshal(tom)
		j = []byte(`
		{
    "id": 67,
    "jsonrpc": "2.0",
    "result": {
        "baseFeePerGas":"0x0",
        "difficulty": "0x1",
        "extraData": "0x",
        "gasLimit": "0x7ffffffffff",
        "gasUsed": "0x0",
        "hash": "0x47db260f84102f0c1f5bf7d929cb8910081323324cc80a4847ab33cc01c5d7b8",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "miner": "0xbbbbbbbbbbbbbbbbbbbbbbbb5530ea015b900000",
        "mixHash": "0x0651d23338760ccfd9f992b431ad6ce17003ffd1bb816c160ea7efaa7ef2bec6",
        "nonce": "0x0000000000000000",
        "number": "0x125798b",
        "parentHash": "0x969c6a8f7fad53acd5cf4ff77bc4b62cceb3ae94254166340b43c14193314522",
        "receiptsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0x203",
        "stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "timestamp": "0x65412c3e",
        "totalDifficulty": "0x125798c",
        "transactions": [],
        "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "uncles": []
    }
}
`)

                w.Write(j)
        default:
                w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprintf(w, "I can't do that.")
        }
}
func tomHandler2(w http.ResponseWriter, r *http.Request) {

        switch r.Method {
        case "GET":
                // Just send out the JSON version of 'tom'
                j, _ := json.Marshal(tom)
                w.Write(j)
        case "POST":
                // Decode the JSON in the body and overwrite 'tom' with it
                d := json.NewDecoder(r.Body)
                p := &person{}
                err := d.Decode(p)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                }
                tom = p
        default:
                w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprintf(w, "I can't do that.")
        }
}

func main() {
        http.HandleFunc("/tom", tomHandler)
        http.HandleFunc("/", tomHandler)

        log.Println("Go!")
        http.ListenAndServe(":10002", nil)
}
