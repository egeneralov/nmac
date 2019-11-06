package main


import (
  "os"
  "fmt"
  "time"
  "syscall"
  "os/signal"
  "encoding/json"
  
  parser "github.com/egeneralov/nmac/parser"
)

const (
  maxQueueLength = 100
  maxFill = 635
  maxTasks = 5
)
var (
  doneTasks = 0
)


type application struct {
  Queue chan int
}

func (app *application) Init() {
  app.Queue = make(chan int, maxQueueLength)
}

func (app *application) Get() int {
  return <- app.Queue
}

func (app *application) Put(i int) {
  app.Queue <- i
}

func (app *application) FillIt() {
  for i := 1; i < maxFill; i++ {
    app.Put(i)
  }
}

func (app *application) Worker(id int) {
  for {
    if len(app.Queue) > 0 {
      
      pageNumber := app.Get()
      var NmacURL string
      
      if pageNumber == 1 {
        NmacURL = "https://nmac.to"
      } else {
        NmacURL = fmt.Sprintf("https://nmac.to/page/%d/", pageNumber)
      }
      indexPageUrl := fmt.Sprintf(NmacURL)
      for _, itemPageUrl := range parser.ExtractIndexPage(indexPageUrl) {
        NI := parser.ExtractItemPage(itemPageUrl)
        NIjson, err := json.Marshal(NI)
        if err != nil { panic(err) }
        fmt.Println(string(NIjson))
      }

    } else {
      doneTasks = doneTasks + 1
      if doneTasks == maxTasks {
        os.Exit(0)
      }
      break
    }
  }
}


func main() {
  myApp := new(application)
  
  myApp.Init()
  go myApp.FillIt()
  time.Sleep(time.Second * 1)
  for i := 0; i < maxTasks; i++ {
    go myApp.Worker(i)
  }

  sigchan := make(chan os.Signal, 1)
  signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
  func() {
    select {
      case _ = <-sigchan:
       os.Exit(0)
    }
  }()

}
