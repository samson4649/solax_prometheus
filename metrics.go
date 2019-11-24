package solax_prometheus

import (
  "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
  IsRegistered          bool
  PVCurrent             *prometheus.GaugeVec
  PVVoltage             *prometheus.GaugeVec
  PVPower               *prometheus.GaugeVec
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
  CoreScrapeVersion     prometheus.Gauge
  CoreScrapeFailCount   prometheus.Counter
}

func (s *Server) initMetrics() (bool) {

  //  metricNames[0]  = "pv_pv1_current"
  //  metricNames[1]  = "pv_pv2_current"
  s.Metrics.PVCurrent = prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_current",
    Help: "PV current now",
    },
    []string{"pv"},
  )

  //  metricNames[2]  = "pv_pv1_voltage"
  //  metricNames[3]  = "pv_pv2_voltage"
  s.Metrics.PVVoltage = prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_voltage",
    Help: "PV voltage now",
    },
    []string{"pv"},
  )

  //  metricNames[11] = "pv_pv1_input_power"
  //  metricNames[12] = "pv_pv2_input_power"
  s.Metrics.PVPower = prometheus.NewGaugeVec( prometheus.GaugeOpts{
    Name: "solax_pv_input_power",
    Help: "PV input power now",
    },
    []string{"pv"},
  )

  //  metricNames[4]  = "grid_output_current"
  s.Metrics.GridOutputCurrent = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_output_current",
    Help: "Grid Output Current",
  })

  //  metricNames[5]  = "grid_network_voltage"
  s.Metrics.GridNetworkVoltage = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_network_voltage",
    Help: "Grid Network Voltage",
  })

  //  metricNames[6]  = "grid_power"
  s.Metrics.GridPower = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_power",
    Help: "Grid Power",
  })

  //  metricNames[10] = "grid_feed_power"
  s.Metrics.GridFeedPower = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_feed_power",
    Help: "Grid Feed Power",
  })

  //  metricNames[41] = "grid_frequency"
  s.Metrics.GridFrequency = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_frequency",
    Help: "Grid Frequency",
  })

  //  metricNames[42] = "grid_exported"
  s.Metrics.GridExported = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_exported",
    Help: "Grid Exported",
  })

  //  metricNames[50] = "grid_imported"
  s.Metrics.GridImported = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_grid_imported",
    Help: "Grid Imported",
  })

  //  metricNames[13] = "battery_voltage"
  s.Metrics.BatteryVoltage = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_voltage",
    Help: "Battery Voltage",
  })

  //  metricNames[14] = "dis_charge_current"
  s.Metrics.BatteryCurrent = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_charge_discharge_current",
    Help: "Current Charge / Discharge",
  })

  //  metricNames[15] = "battery_power"
  s.Metrics.BatteryPower = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_power",
    Help: "Battery Power",
  })

  //  metricNames[16] = "battery_temperature"
  s.Metrics.BatteryTemp = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_temp",
    Help: "Battert Temp",
  })

  //  metricNames[17] = "remaining_capacity"
  s.Metrics.BatteryCap = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_remaining_capacity",
    Help: "Remaining capacity of system",
  })
  //  metricNames[8]  = "inverter_yeild_today"
  s.Metrics.InvertDaily = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_inverter_yeild_today",
    Help: "Inverter Yeild Today",
  })

  //  metricNames[9]  = "inverter_yeild_month"
  s.Metrics.InvertMonthly = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_inverter_yeild_month",
    Help: "Inverter Yeild for the Month",
  })

  //  metricNames[19] = "battery_yeild_total"
  s.Metrics.InvertYearly = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_battery_yeild_total",
    Help: "Inverter Battery Yeild Total",
  })

  s.Metrics.CoreScrapeFailCount = prometheus.NewCounter( prometheus.CounterOpts{
    Name: "solax_scrape_error",
    Help: "Count of errors from backend scraping",
  })

  s.Metrics.CoreScrapeVersion = prometheus.NewGauge( prometheus.GaugeOpts{
    Name: "solax_scrape_version",
    Help: "Solax Scraper Version running",
  })

  return true
}


func (s *Server) registerCoreMetrics(){
  s.Registry.MustRegister( s.Metrics.CoreScrapeFailCount )
  s.Registry.MustRegister( s.Metrics.CoreScrapeVersion )
  s.Metrics.CoreScrapeVersion.Set( s.Version )
}

func (s *Server) registerMetrics() {
  // metrics have to be registered to be exposed
  s.Registry.Register( s.Metrics.PVCurrent )
  s.Registry.Register( s.Metrics.PVVoltage )
  s.Registry.Register( s.Metrics.PVPower )
  s.Registry.Register( s.Metrics.GridOutputCurrent )
  s.Registry.Register( s.Metrics.GridNetworkVoltage )
  s.Registry.Register( s.Metrics.GridPower )
  s.Registry.Register( s.Metrics.GridFeedPower )
  s.Registry.Register( s.Metrics.GridFrequency )
  s.Registry.Register( s.Metrics.GridExported )
  s.Registry.Register( s.Metrics.GridImported )
  s.Registry.Register( s.Metrics.BatteryVoltage )
  s.Registry.Register( s.Metrics.BatteryCurrent )
  s.Registry.Register( s.Metrics.BatteryPower )
  s.Registry.Register( s.Metrics.BatteryTemp )
  s.Registry.Register( s.Metrics.BatteryCap )
  s.Registry.Register( s.Metrics.InvertDaily )
  s.Registry.Register( s.Metrics.InvertMonthly )
  s.Registry.Register( s.Metrics.InvertYearly )
}

func (s *Server) unregisterMetrics(){
  core := s.Metrics
  // metrics unregistered when missing data
  s.Registry.Unregister( core.PVCurrent )
  s.Registry.Unregister( core.PVVoltage )
  s.Registry.Unregister( core.PVPower )
  s.Registry.Unregister( core.GridOutputCurrent )
  s.Registry.Unregister( core.GridNetworkVoltage )
  s.Registry.Unregister( core.GridPower )
  s.Registry.Unregister( core.GridFeedPower )
  s.Registry.Unregister( core.GridFrequency )
  s.Registry.Unregister( core.GridExported )
  s.Registry.Unregister( core.GridImported )
  s.Registry.Unregister( core.BatteryVoltage )
  s.Registry.Unregister( core.BatteryCurrent )
  s.Registry.Unregister( core.BatteryPower )
  s.Registry.Unregister( core.BatteryTemp )
  s.Registry.Unregister( core.BatteryCap )
  s.Registry.Unregister( core.InvertDaily )
  s.Registry.Unregister( core.InvertMonthly )
  s.Registry.Unregister( core.InvertYearly )
}

func (s *Server) updateMetrics( d *APIStructuredData ) error {
  s.registerIfNotAlready()
  s.Metrics.PVCurrent.With( prometheus.Labels{"pv":"1"}).Set( d.PVCurrentOne )
  s.Metrics.PVCurrent.With( prometheus.Labels{"pv":"2"}).Set( d.PVCurrentTwo )
  s.Metrics.PVVoltage.With( prometheus.Labels{"pv":"1"}).Set( d.PVVoltageOne )
  s.Metrics.PVVoltage.With( prometheus.Labels{"pv":"2"}).Set( d.PVVoltageTwo )
  s.Metrics.PVPower.With( prometheus.Labels{"pv":"1"}).Set( d.PVPowerOne )
  s.Metrics.PVPower.With( prometheus.Labels{"pv":"2"}).Set( d.PVPowerTwo )
  s.Metrics.GridOutputCurrent.Set( d.GridOutputCurrent )
  s.Metrics.GridNetworkVoltage.Set( d.GridNetworkVoltage )
  s.Metrics.GridPower.Set( d.GridPower )
  s.Metrics.GridFeedPower.Set( d.GridFeedPower )
  s.Metrics.GridFrequency.Set( d.GridFrequency )
  s.Metrics.GridExported.Set( d.GridExported )
  s.Metrics.GridImported.Set( d.GridImported )
  s.Metrics.BatteryVoltage.Set( d.BatteryVoltage )
  s.Metrics.BatteryCurrent.Set( d.BatteryCurrent )
  s.Metrics.BatteryPower.Set( d.BatteryPower )
  s.Metrics.BatteryTemp.Set( d.BatteryTemp )
  s.Metrics.BatteryCap.Set( d.BatteryCap )
  s.Metrics.InvertDaily.Set( d.InvertDaily )
  s.Metrics.InvertMonthly.Set( d.InvertMonthly )
  s.Metrics.InvertYearly.Set( d.InvertYearly )
  return nil
}

func (s *Server) registerIfNotAlready() {
  if ! s.Metrics.IsRegistered {
    s.registerMetrics()
    s.Metrics.IsRegistered = false
  }
}

func (s *Server) unregisterIfAlready() {
  if s.Metrics.IsRegistered {
    s.unregisterMetrics()
    s.Metrics.IsRegistered = true
  }
}
