package solax_prometheus

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "regexp"
  "encoding/json"
  "sync"
)

type Scraper struct {
  Target      string
  Client      *http.Client
  Com         *ScraperController
}

type ScraperController struct {
  Running int
  mtx     *sync.WaitGroup
  Command chan int
  Result  chan *APIStructuredData
}

type APIStructuredData struct {
  PVCurrentOne          float64
  PVCurrentTwo          float64
  PVVoltageOne          float64
  PVVoltageTwo          float64
  PVPowerOne            float64
  PVPowerTwo            float64
  GridOutputCurrent     float64
  GridNetworkVoltage    float64
  GridPower             float64
  GridFeedPower         float64
  GridFrequency         float64
  GridExported          float64
  GridImported          float64
  BatteryVoltage        float64
  BatteryCurrent        float64
  BatteryPower          float64
  BatteryTemp           float64
  BatteryCap            float64
  InvertDaily           float64
  InvertMonthly         float64
  InvertYearly          float64
  CoreScrapeVersion     float64
  CoreScrapeFailCount   float64
}

func NewScraper(t string) (*Scraper){
  return &Scraper{
    Target: t,
    Client: &http.Client{},
    Com:    &ScraperController{
      Running: 0,
      mtx: &sync.WaitGroup{},
      Command: make(chan int, 1),
      Result:  make(chan *APIStructuredData, 1),
    },
  }
}

func (s *Scraper) Get( t string ) ( *http.Response, error ) {
  return s.Client.Get(t)
}

func (s *Scraper) Run() ( *APIStructuredData, error ) {
  resp, err := s.Get(s.Target)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll( resp.Body )
  if err != nil {
    return nil, err
  }
  if resp.StatusCode < 200 || resp.StatusCode > 299 {
    return nil, fmt.Errorf("bad response code: %d", resp.StatusCode )
  }
  data := APIStructuredData{}

  if err := marshalAPI(body,&data); err != nil {
    return nil, err
  }
  return &data, nil
}

func marshalAPI( b []byte, d *APIStructuredData ) error {
  re := regexp.MustCompile(`,,`)
  intermediate := re.ReplaceAllLiteralString( string(b), ",0,")
  cData := re.ReplaceAllLiteralString( intermediate, ",0," )

  var xtract struct {
    Data  []float64 `json:"data"`
  }

  if err := json.Unmarshal( []byte(cData), &xtract ); err != nil {
    return err
  }
  if err := structureAPIData(&xtract.Data,d); err != nil {
    return err
  }
  return nil
}

func structureAPIData(d *[]float64, r *APIStructuredData ) ( error ) {
  if len(*d) != 68 {
    return fmt.Errorf("Bad length of API raw data")
  }
  r.PVCurrentOne = (*d)[0]
  r.PVCurrentTwo = (*d)[1]
  r.PVVoltageOne = (*d)[2]
  r.PVVoltageTwo = (*d)[3]
  r.PVPowerOne = (*d)[11]
  r.PVPowerTwo = (*d)[12]
  r.GridOutputCurrent = (*d)[4]
  r.GridNetworkVoltage = (*d)[5]
  r.GridPower = (*d)[6]
  r.GridFeedPower = (*d)[10]
  r.GridFrequency = (*d)[50]
  r.GridExported = (*d)[41]
  r.GridImported = (*d)[42]
  r.BatteryVoltage = (*d)[13]
  r.BatteryCurrent = (*d)[14]
  r.BatteryPower = (*d)[15]
  r.BatteryTemp = (*d)[16]
  r.BatteryCap = (*d)[17]
  r.InvertDaily = (*d)[8]
  r.InvertMonthly = (*d)[9]
  r.InvertYearly = (*d)[19]
  // fmt.Println(r)
  return nil
}
