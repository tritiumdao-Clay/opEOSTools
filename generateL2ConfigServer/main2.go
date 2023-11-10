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
		j = []byte(`{
    "id": 67,
    "jsonrpc": "2.0",
    "result": {
        "baseFeePerGas":"0x0",
        "difficulty": "0x1",
        "extraData": "0x",
        "gasLimit": "0x7ffffffffff",
        "gasUsed": "0x0",
        "hash": "0xb2d922e0572d7654851919a0ca59b7390a174820005c6a4faaadf3da26c1015a",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "miner": "0xbbbbbbbbbbbbbbbbbbbbbbbb5530ea015b900000",
        "mixHash": "0x06443f4e5c53ee406bc6da3d138059059d15e5841cd7941aa89eaaede34c916e",
        "nonce": "0x0000000000000000",
        "number": "0x11eaff6",
        "parentHash": "0xe0bffa69b6c9554bfe61a1030915d57fd44d467bf9e8c1f94e52e4204a157e34",
        "receiptsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0x203",
        "stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "timestamp": "0x653a62a9",
        "totalDifficulty": "0x11eaff7",
        "transactions": [],
        "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "uncles": []
    }
}`)
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
