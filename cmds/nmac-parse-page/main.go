package main

import (
  "fmt"
  "flag"
  "encoding/json"
  
  parser "github.com/egeneralov/nmac/parser"
)

func main() {
  pageNumber := flag.Int("page", 1, "page number to parse")
  flag.Parse()
  
  var NmacURL string
  var result []*parser.NmacItem
  
  if *pageNumber == 1 {
    NmacURL = "https://nmac.to"
  } else {
    NmacURL = fmt.Sprintf("https://nmac.to/page/%d/", pageNumber)
  }
  indexPageUrl := fmt.Sprintf(NmacURL)
  for _, itemPageUrl := range parser.ExtractIndexPage(indexPageUrl) {
    NI := parser.ExtractItemPage(itemPageUrl)
    result = append(result, NI)
  }
  NIjson, err := json.Marshal(result)
  if err != nil { panic(err) }
  fmt.Println(string(NIjson))
}
