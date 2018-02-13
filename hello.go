package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "encoding/xml"
    "encoding/json"
    "strings"
    "time"
    "io/ioutil"
)

type Item struct {
    StationName string
    StationType string
    StationUrl string
}

type ListOfItems struct {
    ItemCount uint8
    Item []Item
}

type VtunerConfig struct {
    EncryptedToken string
    HttpPort uint16
    DnsServer string
    VtunerServerOne string
    VtunerServerTwo string
    DnsPort uint16
}

type FirebaseConfig struct {
    databaseURL string
    baseRef string
    dbSecret string
}

type FirebaseItem struct {
    Item Item
    Key string
}

var (
    Items []Item
)

func loginxml(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "todo");
}

func favxml(w http.ResponseWriter, r *http.Request) {
    los := &ListOfItems {};
    los.Item = Items;
    los.ItemCount = 9;
    output, err := xml.MarshalIndent(los, "", "  ");
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }

    io.WriteString(w, xml.Header);
    w.Write(output);
}

func CaselessMatcher(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.URL.Path = strings.ToLower(r.URL.Path)
        h.ServeHTTP(w, r)
    })
}

func loadItems() {
    dbSecret := "<...>"
    urlBase := "https://<...>.firebaseio.com/production/.json"
    url := urlBase + "?auth=" + dbSecret
    myClient := &http.Client{Timeout: 10 * time.Second}
    resp, err := myClient.Get(url)
    if err != nil {
        fmt.Println(err)
        fmt.Println("request failed")
	return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    var data []FirebaseItem
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("decoding failed")
        fmt.Println(err)
	return
    }

    Items = make([]Item, len(data))
    for i, item := range data {
        Items[i] = item.Item
    }
}

func main() {
    loadItems()

    mux := http.NewServeMux()
    http.Handle("/", CaselessMatcher(mux))

    mux.HandleFunc("/setupapp/yamaha/asp/browsexml/favxml.asp", favxml);
    mux.HandleFunc("/setupapp/yamaha/asp/browsexml/loginxml.asp", loginxml);

    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
