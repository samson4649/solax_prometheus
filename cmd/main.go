package main

import (
  "context"
  "sync"
  "os"
  "os/signal"
  "syscall"
  "time"
  "net/http"
  "github.com/alexflint/go-arg"
  solax "github.com/samson4649/solax_prometheus"
  "go.uber.org/zap"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var args struct {
  Port        string  `arg:"-p,env:SERVER_PORT" help:"Port to listen on localhost"`
  IP          string  `arg:"env:SERVER_HOST" help:"Interface to listen on"`
  Verbose     bool    `arg:"-v,env:SERVER_VERBOSE" help:"Provide information on the actions running"`
  Endpoint    string  `arg:"-t,required,env:SERVER_ENDPOINT" help:"Solax Endpoint target"`
  LogPath     string  `arg:"--log-path,env:SERVER_LOG_PATH" help:"Path to log file"`
}

func main(){
  // defer logger.Sync()

  sigc := make(chan os.Signal, 1)
  signal.Notify(sigc,
      syscall.SIGHUP,
      syscall.SIGINT,
      syscall.SIGTERM,
      syscall.SIGQUIT)

  args.Port     = "4444"
  args.IP       = "0.0.0.0"
  args.LogPath  = "stdout"
  // args.Endpoint = "http://api:2015/realTimeData.htm"
  arg.MustParse(&args)

  logLevel := "info"
  if args.Verbose {
    logLevel = "verbose"
  }
  log,_ := solax.InitLogger(args.LogPath, logLevel)

  log.Info("Starting application")

  var tracker = struct {
    Scraper int
    Server  int
  }{
    Scraper: 1,
    Server:  1,
  }

  metrics := solax.Initialise( )
  scraper := solax.NewScraper(args.Endpoint)
  var wg sync.WaitGroup

  wg.Add(1)
  go func (){
    defer wg.Done()
    for {
      if tracker.Scraper != 1 {
        break
      }
      apiData, err := scraper.Run()
      if err != nil {
        log.Error("Bad scrape run", zap.String( "err", err.Error() ))
        metrics.Clear()
      } else {
        metrics.Update(apiData)
      }
      for x := 0; x < 10; x++ {
        if tracker.Scraper != 1 {
          break
        }
        time.Sleep( time.Duration(1) * time.Second )
      }
    }
    log.Info("scraper: closed")
  }()

  handler := promhttp.HandlerFor( metrics.Registry, promhttp.HandlerOpts{} )
  http.Handle("/metrics", middleware(log,handler) )

  serveString := args.IP + ":" + args.Port

  srv := &http.Server{ Addr: serveString }

  go func() {
      <-sigc
      log.Info("SIGINT Caught - shutting down...")
      ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
      srv.Shutdown(ctx)
      tracker.Scraper = 0
      tracker.Server  = 0
  }()

  wg.Add(1)
  go func(){
    err := srv.ListenAndServe( )
    if err != nil {
      log.Error( err.Error() )
    }
    wg.Done()
  }()
  wg.Wait()
}

func middleware(log *zap.Logger, next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    ts := time.Now()
    next.ServeHTTP(w,r)
    log.With(
      zap.Namespace("Request Data"),
      zap.String("path", r.URL.RequestURI() ),
      zap.String("time", time.Since(ts).String() ),
    ).Info("New Request")
  })
}

























// end
