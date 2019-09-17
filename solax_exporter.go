package main

import (
  "net/http"
  "io/ioutil"
  // "fmt"
  "encoding/json"
  "time"
  "regexp"

  "go.uber.org/zap"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"

  "gitlab.monarch-ares.com/go-includes/golog"

  "github.com/alexflint/go-arg"

)

//  metricNames[0]  = "pv_pv1_current"
//  metricNames[1]  = "pv_pv2_current"
//  metricNames[2]  = "pv_pv1_voltage"
//  metricNames[3]  = "pv_pv2_voltage"
//  metricNames[11] = "pv_pv1_input_power"
//  metricNames[12] = "pv_pv2_input_power"
//  metricNames[4]  = "grid_output_current"
//  metricNames[5]  = "grid_network_voltage"
//  metricNames[6]  = "grid_power"
//  metricNames[10] = "grid_feed_power"
//  metricNames[41] = "grid_frequency"
//  metricNames[42] = "grid_exported"
//  metricNames[50] = "grid_imported"
//  metricNames[13] = "battery_voltage"
//  metricNames[14] = "dis_charge_current"
//  metricNames[15] = "battery_power"
//  metricNames[16] = "battery_temperature"
//  metricNames[17] = "remaining_capacity"
//  metricNames[8]  = "interter_yeild_today"
//  metricNames[9]  = "interter_yeild_month"
//  metricNames[19] = "battery_yeild_total"

type Domain struct {
  Version         string
  ScrapeErrCount  int
  Registry        *prometheus.Registry
  Metrics         Metrics
}

type Metrics struct {
  IsRegistered          bool
  PVCurrent             prometheus.GaugeVec 
  PVVoltage             prometheus.GaugeVec
  PVPower               prometheus.GaugeVec
  GridOutputCurrent     prometheus.Gauge
  GridNetworkVoltage    prometheus.Gauge
  GridPower             prometheus.Gauge
  GridFeedPower         prometheus.Gauge
  GridFrequency         prometheus.Gauge
  GridExported          prometheus.Gauge
  GridImported          prometheus.Gauge
  BatteryVoltage        prometheus.Gauge
  BatteryCurrent        prometheus.Gauge
  BatteryPower          prometheus.Gauge
  BatteryTemp           prometheus.Gauge
  BatteryCap            prometheus.Gauge
  InvertDaily           prometheus.Gauge 
  InvertMonthly         prometheus.Gauge
  InvertYearly          prometheus.Gauge
  CoreScrapeFailCount   prometheus.Counter
}

var args struct {
  Port        string  `arg:"-p" help:"Port to listen on localhost (default :4444)"`
  IP          string  `help:"Interface to listen on (default: 0.0.0.0)"`
  Verbose     bool    `arg:"-v" help:"Provide information on the actions running"`
  Silent      bool    `arg:"-q" help:"Dont log anything - just use return code (for piping)"`
  Endpoint    string  `arg:"-t" help:"Solax Endpoint target" (default: http://11.11.11.1/api/realTimeData.htm)"`
}

// var core Metrics
var root Domain

var target = "http://11.11.11.1/api/realTimeData.htm"
// var target = "http://api:2015/realTimeData.htm"

// logger
var logger *zap.Logger

func (d *Domain) initMetrics() (bool) {

  core := &d.Metrics

  //  metricNames[0]  = "pv_pv1_current"
  //  metricNames[1]  = "pv_pv2_current"
  core.PVCurrent = *prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_current",
    Help: "PV current now",
    },
    []string{"pv"},
  )
  
  //  metricNames[2]  = "pv_pv1_voltage"
  //  metricNames[3]  = "pv_pv2_voltage"
  core.PVVoltage = *prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_voltage",
    Help: "PV voltage now",
    },
    []string{"pv"},
  )

  //  metricNames[11] = "pv_pv1_input_power"
  //  metricNames[12] = "pv_pv2_input_power"
  core.PVPower = *prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_input_power",
    Help: "PV input power now",
    },
    []string{"pv"},
  )

  //  metricNames[4]  = "grid_output_current"
  core.GridOutputCurrent = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_output_current",
    Help: "Grid Output Current",
  })

  //  metricNames[5]  = "grid_network_voltage"
  core.GridNetworkVoltage = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_network_voltage",
    Help: "Grid Network Voltage",
  })

  //  metricNames[6]  = "grid_power"
  core.GridPower = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_power",
    Help: "Grid Power",
  })

  //  metricNames[10] = "grid_feed_power"
  core.GridFeedPower = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_feed_power",
    Help: "Grid Feed Power",
  })

  //  metricNames[41] = "grid_frequency"
  core.GridFrequency = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_frequency",
    Help: "Grid Frequency",
  })

  //  metricNames[42] = "grid_exported"
  core.GridExported = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_exported",
    Help: "Grid Exported",
  })

  //  metricNames[50] = "grid_imported"
  core.GridImported = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_imported",
    Help: "Grid Imported",
  })

  //  metricNames[13] = "battery_voltage"
  core.BatteryVoltage = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_voltage",
    Help: "Battery Voltage",
  })

  //  metricNames[14] = "dis_charge_current"
  core.BatteryCurrent = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_charge_discharge_current",
    Help: "Current Charge / Discharge",
  })

  //  metricNames[15] = "battery_power"
  core.BatteryPower = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_power",
    Help: "Battery Power",
  })

  //  metricNames[16] = "battery_temperature"
  core.BatteryTemp = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_temp",
    Help: "Battert Temp",
  })

  //  metricNames[17] = "remaining_capacity"
  core.BatteryCap = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_remaining_capacity",
    Help: "Remaining capacity of system",
  })
  //  metricNames[8]  = "inverter_yeild_today"
  core.InvertDaily = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_inverter_yeild_today",
    Help: "Inverter Yeild Today",
  })

  //  metricNames[9]  = "inverter_yeild_month"
  core.InvertMonthly = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_inverter_yeild_month",
    Help: "Inverter Yeild for the Month",
  })

  //  metricNames[19] = "battery_yeild_total"
  core.InvertYearly = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_yeild_total",
    Help: "Inverter Battery Yeild Total",
  })

  core.CoreScrapeFailCount = prometheus.NewCounter( prometheus.CounterOpts{
    Name: "solax_scrape_error",
    Help: "Count of errors from backend scraping",
  })
  return true
}

type SolarData struct {
  Method string     `json:"method"`
  Version string    `json:"version"`
  Type string       `json:"type"`
  Serial string     `json:"SN"`
  Status string     `json:"Status"`
  Data []float64    `json:"Data"`
}

func (d *Domain) registerCoreMetrics(){
  d.Registry.MustRegister( d.Metrics.CoreScrapeFailCount )
}

func (d *Domain) registerSolaxMetrics() {
  // metrics have to be registered to be exposed
  d.Registry.Register( d.Metrics.PVCurrent )
  d.Registry.Register( d.Metrics.PVVoltage )
  d.Registry.Register( d.Metrics.PVPower )
  d.Registry.Register( d.Metrics.GridOutputCurrent )
  d.Registry.Register( d.Metrics.GridNetworkVoltage )
  d.Registry.Register( d.Metrics.GridPower )
  d.Registry.Register( d.Metrics.GridFeedPower )
  d.Registry.Register( d.Metrics.GridFrequency )
  d.Registry.Register( d.Metrics.GridExported )
  d.Registry.Register( d.Metrics.GridImported )
  d.Registry.Register( d.Metrics.BatteryVoltage )
  d.Registry.Register( d.Metrics.BatteryCurrent )
  d.Registry.Register( d.Metrics.BatteryPower )
  d.Registry.Register( d.Metrics.BatteryTemp )
  d.Registry.Register( d.Metrics.BatteryCap )
  d.Registry.Register( d.Metrics.InvertDaily )
  d.Registry.Register( d.Metrics.InvertMonthly )
  d.Registry.Register( d.Metrics.InvertYearly )
}

func (d *Domain) unRegisterSolaxMetrics(){
  core := &d.Metrics
  // metrics unregistered when missing data
  d.Registry.Unregister( core.PVCurrent )
  d.Registry.Unregister( core.PVVoltage )
  d.Registry.Unregister( core.PVPower )
  d.Registry.Unregister( core.GridOutputCurrent )
  d.Registry.Unregister( core.GridNetworkVoltage )
  d.Registry.Unregister( core.GridPower )
  d.Registry.Unregister( core.GridFeedPower )
  d.Registry.Unregister( core.GridFrequency )
  d.Registry.Unregister( core.GridExported )
  d.Registry.Unregister( core.GridImported )
  d.Registry.Unregister( core.BatteryVoltage )
  d.Registry.Unregister( core.BatteryCurrent )
  d.Registry.Unregister( core.BatteryPower )
  d.Registry.Unregister( core.BatteryTemp )
  d.Registry.Unregister( core.BatteryCap )
  d.Registry.Unregister( core.InvertDaily )
  d.Registry.Unregister( core.InvertMonthly )
  d.Registry.Unregister( core.InvertYearly )
  // d.Registry.Unregister( d.Metrics.CoreScrapeFails )
}

func (d *Domain) registerIfNotAlready() {
  if ! d.Metrics.IsRegistered {
    d.registerSolaxMetrics()
    d.Metrics.IsRegistered = false
  }
}

func (d *Domain) unregisterIfAlready() {
  if d.Metrics.IsRegistered {
    d.unRegisterSolaxMetrics()
    d.Metrics.IsRegistered = true
  }
}

func (d *Domain) Run(){
  go func(){
    logger.Info("Metrics Spider Started.")
    for {
      
      //example data
      // r_data := string( `{"method":"uploadsn","version":"Solax_SI_CH_2nd_20160912_DE02","type":"AL_SE","SN":"829E4BBD","Data":[0.0,0.0,0.0,0.0,4.0,240.0,937,38,21.6,9162.7,-748,0,0,53.18,-19.71,-1050,19,67,0.0,3482.1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0.00,0.00,0,0,0,0,0,0,0,49.96,0,0,0.0,0.0,0,0.00,0,0,0,0.00,0,8,0,0,0.00,0,8],"Status":"2"}` )

      client := &http.Client{}
      // url := target
      // url := "http://api-fault"
      url := args.Endpoint
      resp,err := client.Get(url)
      
      if err != nil {
        logger.Warn("Not able to connect to host",
          zap.String( "url", url ),
        )
        d.unRegisterSolaxMetrics()
        d.Metrics.CoreScrapeFailCount.Inc()
        // fmt.Println( "connecting to host " + url + err.Error() )
      } else {
      
        defer resp.Body.Close()

        if resp.StatusCode < 200 || resp.StatusCode > 299 {
        // if true {
          d.unRegisterSolaxMetrics()
          d.Metrics.CoreScrapeFailCount.Inc()
          logger.Warn("Not a valid response from server",
            zap.String( "url", url ),
            zap.Int( "statusCode", resp.StatusCode ),
            zap.String( "statusCodeString", http.StatusText(resp.StatusCode) ),
          )
        } else {
        
          html, err := ioutil.ReadAll(resp.Body)
          if err == nil {

            re := regexp.MustCompile(`,,`)
            tmp := re.ReplaceAllLiteralString( string(html), ",0,")
            r_data := re.ReplaceAllLiteralString( tmp, ",0," )
            
            solarData := SolarData{}
            if err := json.Unmarshal([]byte(r_data), &solarData); err != nil {
              d.unRegisterSolaxMetrics()
              d.Metrics.CoreScrapeFailCount.Inc()
              logger.Warn("Not able to Unmarshal JSON data",
                zap.String("unmarshal_err", err.Error() ),
              )
            } else {

              if len(solarData.Data) == 68 {
                logger.Info("step")
                d.registerIfNotAlready()
                logger.Info("step_2")
                d.Metrics.PVCurrent.With( prometheus.Labels{"pv":"1"}).Set( solarData.Data[0] )
                d.Metrics.PVCurrent.With( prometheus.Labels{"pv":"2"}).Set( solarData.Data[1] )
                d.Metrics.PVVoltage.With( prometheus.Labels{"pv":"1"}).Set( solarData.Data[2] )
                d.Metrics.PVVoltage.With( prometheus.Labels{"pv":"2"}).Set( solarData.Data[3] )
                d.Metrics.PVPower.With( prometheus.Labels{"pv":"1"}).Set( solarData.Data[11] )
                d.Metrics.PVPower.With( prometheus.Labels{"pv":"2"}).Set( solarData.Data[12] )
                d.Metrics.GridOutputCurrent.Set( solarData.Data[4] )
                d.Metrics.GridNetworkVoltage.Set( solarData.Data[5] )
                d.Metrics.GridPower.Set( solarData.Data[6] )
                d.Metrics.GridFeedPower.Set( solarData.Data[10] )
                d.Metrics.GridFrequency.Set( solarData.Data[50] )
                d.Metrics.GridExported.Set( solarData.Data[41] )
                d.Metrics.GridImported.Set( solarData.Data[42] )
                d.Metrics.BatteryVoltage.Set( solarData.Data[13] )
                d.Metrics.BatteryCurrent.Set( solarData.Data[14] )
                d.Metrics.BatteryPower.Set( solarData.Data[15] )
                d.Metrics.BatteryTemp.Set( solarData.Data[16] )
                d.Metrics.BatteryCap.Set( solarData.Data[17] )
                d.Metrics.InvertDaily.Set( solarData.Data[8] )
                d.Metrics.InvertMonthly.Set( solarData.Data[9] )
                d.Metrics.InvertYearly.Set( solarData.Data[19] )
              } else {
                logger.Warn("Unmarshalled data at bad length. Aborting.",
                  zap.Int("unmarshal_len", len(solarData.Data) ),
                )
                d.unRegisterSolaxMetrics()
                d.Metrics.CoreScrapeFailCount.Inc()
              }
            }

          } else {
            d.unRegisterSolaxMetrics()
            d.Metrics.CoreScrapeFailCount.Inc()
            logger.Warn("HTML response parse error",
              zap.String("err", err.Error() ),
            )
          }
        }
      }

      // sleep to prevent over polling
      time.Sleep( 10 * time.Second )
    }
  }()
}

func Initialise() ( *Domain ){
  d := Domain{
    Version: "0.0.1",
    Registry: prometheus.NewRegistry(),
  }
  d.initMetrics()
  d.registerCoreMetrics()

  return &d
}

func main() {

  // application argument defaults and parse
  args.Port     = "4444"
  args.IP       = "0.0.0.0"
  args.Endpoint = "http://11.11.11.1/api/realTimeData.htm"
  arg.MustParse(&args)

  // logger construction
  logger = golog.New("xen-safe", golog.LogLevelDetermine( args.Silent, args.Verbose ) )

  logger.Warn("Listening on " + args.IP +":"+ string(args.Port) )

  solax := Initialise()
  solax.registerIfNotAlready()
  solax.Run()

  logger.Info("Serving metrics on 0.0.0.0:4444/metrics")

  handler := promhttp.HandlerFor( solax.Registry, promhttp.HandlerOpts{} )
  http.Handle("/metrics", handler )

  serveString := args.IP + ":" + args.Port
  err := http.ListenAndServe( serveString , nil )
  if err != nil {
    logger.Fatal(err.Error() )
  }
}

























