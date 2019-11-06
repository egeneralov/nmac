package main

import (
  "io"
  "fmt"
  "flag"
  "strings"
  "net/http"
  "encoding/json"
  
  parser "github.com/egeneralov/nmac/parser"
)


func handleRoot (w http.ResponseWriter, r *http.Request) {
  pageIdString := strings.Replace(r.URL.Path, "/", "", 2)
  
  if pageIdString == "" { pageIdString = "1" }
  w.Header().Set("Content-Type", "application/json")
  
  var NmacURL string
  var result []*parser.NmacItem
  
  if pageIdString == "1" {
    NmacURL = "https://nmac.to"
  } else {
    NmacURL = fmt.Sprintf("https://nmac.to/page/%s/", pageIdString)
  }
  
  indexPageUrl := fmt.Sprintf(NmacURL)
  
  for _, itemPageUrl := range parser.ExtractIndexPage(indexPageUrl) {
    result = append(result, parser.ExtractItemPage(itemPageUrl))
  }
  
  NIjson, err := json.Marshal(result)
  if err != nil { panic(err) }
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, string(NIjson))
  io.WriteString(w, "\n")

}


func main() {
  bindTo := flag.String("bind", "0.0.0.0:8018", "golang net/http bind")
  flag.Parse()
  
  http.HandleFunc("/", handleRoot)
  
  fmt.Println("serving on " + *bindTo)
  err := http.ListenAndServe(*bindTo, nil)
  if err != nil { panic(err) }
}
