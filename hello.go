package main
 
import (
    "fmt"
    "io"
    "log"
    "net/http"
    "encoding/xml"
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
 
func main() {
    http.HandleFunc("/setupapp/yamaha/asp/browsexml/FavXML.asp", favxml);
 
    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
