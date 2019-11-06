package parser

import (
  "fmt"
  "strings"
  "regexp"
  "encoding/base64"
  "github.com/antchfx/htmlquery"
)

type NmacLinks struct {
  Uploaded     string `json:"Uploaded"`
  Turbobit     string `json:"Turbobit"`
  Nitroflare   string `json:"Nitroflare"`
  OneFichier   string `json:"1fichier"`
  Letsupload   string `json:"Letsupload"`
  Torrent      string `json:"Torrent"`
  Speed4up     string `json:"Speed4up"`
  Uptobox      string `json:"Uptobox"`
  SenditCloud  string `json:"Sendit.cloud"`
  Userscloud   string `json:"Sendit.cloud"`
  Filescdn     string `json:"Filescdn"`
  Depositfiles string `json:"Depositfiles"`
  Openload     string `json:"Openload"`
  Dailyuploads string `json:"Dailyuploads"`
  Uploadocean  string `json:"Uploadocean"`
  Kingfiles    string `json:"Kingfiles"`
}

type NmacItem struct {
  Title       string    `json:"title"`
  Author      string    `json:"author"`
  Version     string    `json:"version"`
  Description string    `json:"description"`
  NmacLinks   NmacLinks `json:"NmacLinks"`
}



func ExtractTitleFromString(text string) (string, string, string) {
  re := regexp.MustCompile(`^(.+) .*?((?:\d+\.?)+\d+) ?(.+)? ?[–|-] (.+)$`)
  matches := re.FindStringSubmatch(text)
  description := ""
  description = fmt.Sprintf(`%s %s`, matches[3], matches[4])
  description = strings.Trim(description, " ")
  return matches[1], matches[2], description
}




func ExtractIndexPage(url string) ([]string) {
  
  var result []string
  
  doc, err := htmlquery.LoadURL(url)
  if err != nil { panic(err) }

  list, err := htmlquery.QueryAll(doc, `//div[@class="panel-wrapper"]/div[@class="panel"]/div[@class="article-excerpt-wrapper"]/div[@class="article-excerpt"]`)
  if err != nil { panic(err) }

  for _, n := range list {
    a := htmlquery.FindOne(n, "//a")
    
    result = append(result, htmlquery.SelectAttr(a, "href"))
  }
  
  return result

}


func ExtractItemPage(url string) (*NmacItem) {
  doc, err := htmlquery.LoadURL(url)
  if err != nil { panic(err) }
  list, err := htmlquery.QueryAll(doc, `//div[@class="main-loop-content"]`)
  if err != nil { panic(err) }
  
  var title string
  var author string
  var downloadHrefNmac string
  NI := &NmacItem{}
  
  for _, n := range list {
    titleNode := htmlquery.FindOne(n, "//h1")
    title = htmlquery.InnerText(titleNode)
    title = strings.TrimSpace(title)

    authorNode := htmlquery.FindOne(n, `//span[@class="author"]/span[@itemprop="author"]`)
    author = htmlquery.InnerText(authorNode)
    author = strings.TrimSpace(author)
    NI.Author = author
    
    for _, downloadNode := range htmlquery.Find(n, `//div[@class="the-content"]/div/p/a[@target="_blank"]`) {
      downloadNameNmac := strings.TrimSpace(htmlquery.InnerText(downloadNode))
      downloadHrefNmac = strings.TrimSpace(htmlquery.SelectAttr(downloadNode, "href"))
      downloadHrefNmac = strings.Replace(downloadHrefNmac, "https://nmac.to/dl/", "", 1)
      downloadHrefNmacDecoded, _ := base64.URLEncoding.DecodeString(downloadHrefNmac)
      downloadHrefNmac = strings.TrimSpace(string(downloadHrefNmacDecoded))
      switch downloadNameNmac {
        case "Kingfiles":
          NI.NmacLinks.Kingfiles = downloadHrefNmac
        case "Uploadocean":
          NI.NmacLinks.Uploadocean = downloadHrefNmac
        case "Dailyuploads":
          NI.NmacLinks.Dailyuploads = downloadHrefNmac
        case "Openload":
          NI.NmacLinks.Openload = downloadHrefNmac
        case "Depositfiles":
          NI.NmacLinks.Depositfiles = downloadHrefNmac
        case "Filescdn":
          NI.NmacLinks.Filescdn = downloadHrefNmac
        case "Userscloud":
          NI.NmacLinks.Userscloud = downloadHrefNmac
        case "Sendit.cloud":
          NI.NmacLinks.SenditCloud = downloadHrefNmac
        case "Uptobox":
          NI.NmacLinks.Uptobox = downloadHrefNmac
        case "Speed4up":
          NI.NmacLinks.Speed4up = downloadHrefNmac
        case "Torrent":
          NI.NmacLinks.Torrent = downloadHrefNmac
        case "1fichier":
          NI.NmacLinks.OneFichier = downloadHrefNmac
        case "Nitroflare":
          NI.NmacLinks.Nitroflare = downloadHrefNmac
        case "Turbobit":
          NI.NmacLinks.Turbobit = downloadHrefNmac
        case "Uploaded":
          NI.NmacLinks.Uploaded = downloadHrefNmac
        case "Letsupload":
          NI.NmacLinks.Letsupload = downloadHrefNmac
        default:
          panic(downloadNameNmac)
      }
    }

  }
  
  title, version, description := ExtractTitleFromString(title)
  NI.Title = title
  NI.Version = version
  NI.Description = description
  return NI
}

