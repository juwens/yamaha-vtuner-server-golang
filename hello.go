package main
 
import (
    "fmt"
    "io"
    "log"
    "net/http"
    "encoding/xml"
    "strings"
//    "html/template"
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

func loginxml(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "todo");
}

func favxml(w http.ResponseWriter, r *http.Request) {
    i1 := &Item {StationName: "Flux FM", StationType: "Station", StationUrl: "www.irgendwas.de" };
    i2 := &Item {StationName: "radio eins", StationType: "Station", StationUrl: "www.irgendwas.de" };
    los := &ListOfItems {};
    los.Item = []Item {*i1, *i2};
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

func main() {
    mux := http.NewServeMux()
    http.Handle("/", CaselessMatcher(mux))

    mux.HandleFunc("/setupapp/yamaha/asp/browsexml/favxml.asp", favxml);
    mux.HandleFunc("/setupapp/yamaha/asp/browsexml/loginxml.asp", loginxml);

    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
