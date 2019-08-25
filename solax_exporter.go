package main

import (
  "net/http"
  "io/ioutil"
  "fmt"
  "encoding/json"
  "time"
  "regexp"
  "github.com/prometheus/client_golang/prometheus"
  //"github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
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

var (
  //  metricNames[0]  = "pv_pv1_current"
  //  metricNames[1]  = "pv_pv2_current"
  m_pv_current = prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_current",
    Help: "PV current now",
    },
    []string{"pv"},
  )
  
  //  metricNames[2]  = "pv_pv1_voltage"
  //  metricNames[3]  = "pv_pv2_voltage"
  m_pv_voltage = prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_voltage",
    Help: "PV voltage now",
    },
    []string{"pv"},
  )

  //  metricNames[11] = "pv_pv1_input_power"
  //  metricNames[12] = "pv_pv2_input_power"
  m_pv_input_power = prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_input_power",
    Help: "PV input power now",
    },
    []string{"pv"},
  )

  //  metricNames[4]  = "grid_output_current"
  m_grid_output_current = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_output_current",
    Help: "Grid Output Current",
  })

  //  metricNames[5]  = "grid_network_voltage"
  m_grid_network_voltage = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_network_voltage",
    Help: "Grid Network Voltage",
  })

  //  metricNames[6]  = "grid_power"
  m_grid_power = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_power",
    Help: "Grid Power",
  })

  //  metricNames[10] = "grid_feed_power"
  m_grid_feed_power = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_feed_power",
    Help: "Grid Feed Power",
  })

  //  metricNames[41] = "grid_frequency"
  m_grid_frequency = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_frequency",
    Help: "Grid Frequency",
  })

  //  metricNames[42] = "grid_exported"
  m_grid_exported = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_exported",
    Help: "Grid Exported",
  })

  //  metricNames[50] = "grid_imported"
  m_grid_imported = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_imported",
    Help: "Grid Imported",
  })

  //  metricNames[13] = "battery_voltage"
  m_battery_voltage = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_voltage",
    Help: "Battery Voltage",
  })

  //  metricNames[14] = "dis_charge_current"
  m_charge_discharge_current = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_charge_discharge_current",
    Help: "Current Charge / Discharge",
  })

  //  metricNames[15] = "battery_power"
  m_battery_power = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_power",
    Help: "Battery Power",
  })

  //  metricNames[16] = "battery_temperature"
  m_battery_temp = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_temp",
    Help: "Battert Temp",
  })

  //  metricNames[17] = "remaining_capacity"
  m_remaining_cap = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_remaining_capacity",
    Help: "Remaining capacity of system",
  })
  //  metricNames[8]  = "inverter_yeild_today"
  m_inverter_yeild_today = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_inverter_yeild_today",
    Help: "Inverter Yeild Today",
  })

  //  metricNames[9]  = "inverter_yeild_month"
  m_inverter_yeild_month = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_inverter_yeild_month",
    Help: "Inverter Yeild for the Month",
  })

  //  metricNames[19] = "battery_yeild_total"
  m_battery_yeild_total = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_yeild_total",
    Help: "Inverter Battery Yeild Total",
  })

)

type SolarData struct {
  Method string     `json:"method"`
  Version string    `json:"version"`
  Type string       `json:"type"`
  Serial string     `json:"SN"`
  Status string     `json:"Status"`
  Data []float64    `json:"Data"`
}

func init(){
  // metrics have to be registered to be exposed
  prometheus.MustRegister( m_pv_current )
  prometheus.MustRegister( m_pv_voltage )
  prometheus.MustRegister( m_pv_input_power )
  prometheus.MustRegister( m_grid_output_current )
  prometheus.MustRegister( m_grid_network_voltage )
  prometheus.MustRegister( m_grid_power )
  prometheus.MustRegister( m_grid_feed_power )
  prometheus.MustRegister( m_grid_frequency )
  prometheus.MustRegister( m_grid_exported )
  prometheus.MustRegister( m_grid_imported )
  prometheus.MustRegister( m_battery_voltage )
  prometheus.MustRegister( m_charge_discharge_current )
  prometheus.MustRegister( m_battery_power )
  prometheus.MustRegister( m_battery_temp )
  prometheus.MustRegister( m_remaining_cap )
  prometheus.MustRegister( m_inverter_yeild_today )
  prometheus.MustRegister( m_inverter_yeild_month )
  prometheus.MustRegister( m_battery_yeild_total )
}

func updateMetrics(){
  go func(){
    fmt.Println("Metric Parser Started...")
    for {
      
      //example data
      // r_data := string( `{"method":"uploadsn","version":"Solax_SI_CH_2nd_20160912_DE02","type":"AL_SE","SN":"829E4BBD","Data":[0.0,0.0,0.0,0.0,4.0,240.0,937,38,21.6,9162.7,-748,0,0,53.18,-19.71,-1050,19,67,0.0,3482.1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0.00,0.00,0,0,0,0,0,0,0,49.96,0,0,0.0,0.0,0,0.00,0,0,0,0.00,0,8,0,0,0.00,0,8],"Status":"2"}` )

      client := &http.Client{}
      url := "http://11.11.11.1/api/realTimeData.htm"
      resp,err := client.Get(url)
      
      if err != nil {
       fmt.Println( "connecting to host " + url + err.Error() )
      }
      
      defer resp.Body.Close()
      
      html, err := ioutil.ReadAll(resp.Body)
      if err == nil {

        re := regexp.MustCompile(`,,`)
        tmp := re.ReplaceAllLiteralString( string(html), ",0,")
        r_data := re.ReplaceAllLiteralString( tmp, ",0," )
        
        solarData := SolarData{}
        if err := json.Unmarshal([]byte(r_data), &solarData); err != nil {
          fmt.Println(err.Error())
        }

        if len(solarData.Data) == 68 {
          m_pv_current.With( prometheus.Labels{"pv":"1"}).Set( solarData.Data[0] )
          m_pv_current.With( prometheus.Labels{"pv":"2"}).Set( solarData.Data[1] )
          m_pv_voltage.With( prometheus.Labels{"pv":"1"}).Set( solarData.Data[2] )
          m_pv_voltage.With( prometheus.Labels{"pv":"2"}).Set( solarData.Data[3] )
          m_pv_input_power.With( prometheus.Labels{"pv":"1"}).Set( solarData.Data[11] )
          m_pv_input_power.With( prometheus.Labels{"pv":"2"}).Set( solarData.Data[12] )
          m_grid_output_current.Set( solarData.Data[4] )
          m_grid_network_voltage.Set( solarData.Data[5] )
          m_grid_power.Set( solarData.Data[6] )
          m_grid_feed_power.Set( solarData.Data[10] )
          m_grid_frequency.Set( solarData.Data[50] )
          m_grid_exported.Set( solarData.Data[41] )
          m_grid_imported.Set( solarData.Data[42] )
          m_battery_voltage.Set( solarData.Data[13] )
          m_charge_discharge_current.Set( solarData.Data[14] )
          m_battery_power.Set( solarData.Data[15] )
          m_battery_temp.Set( solarData.Data[16] )
          m_remaining_cap.Set( solarData.Data[17] )
          m_inverter_yeild_today.Set( solarData.Data[8] )
          m_inverter_yeild_month.Set( solarData.Data[9] )
          m_battery_yeild_total.Set( solarData.Data[19] )
        } else {
          fmt.Println("data in json at bad length")
        }

      } else {
       fmt.Println( "HTML response parse error: " + err.Error() )
      }

      // sleep to prevent over polling
      time.Sleep( 10 * time.Second )
    }
  }()
}

func main() {
  updateMetrics()
  fmt.Println("Serving metrics on 0.0.0.0:4444/metrics")
  http.Handle("/metrics", promhttp.Handler() )
  http.ListenAndServe(":4444",nil)
}

























