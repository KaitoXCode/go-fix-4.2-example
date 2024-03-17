package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix42/executionreport"
	"github.com/shopspring/decimal"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/config"
)

// Consumer implements quickfix.Application interface
type Consumer struct {
  Env string
  SessionID quickfix.SessionID
  TargetCompID string
  SenderCompID string
}

// OnCreate
func (c Consumer) OnCreate(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnCreate]: Creating session: %s\n", sessionID))
}

// OnLogon
func (c Consumer) OnLogon(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnLogon]: Starting session: %s\n", sessionID))
}

// OnLogout
func (c Consumer) OnLogout(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnLogout]: Terminating session: %s\n", sessionID))
}

// FromAdimn
func (c Consumer) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [FromAdmin]: Received: %s\n", msg.String()))
  return nil
}

// ToAdimn
func (c Consumer) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [ToAdmin]: Sending: %s\n", msg.String()))
}

// ToApp
func (c Consumer) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [ToApp]: Sending: %s\n", msg.String()))
  return
}

// FromApp
func (c Consumer) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [FromApp]: Received: %s\n", msg.String()))
  return
}

func (c Consumer) NewApp(env string) {
  c.Env = env
}

func (c Consumer) sendExampleReport() {
  execReport := executionreport.New(
    field.NewOrderID("orderid#1"), 
    field.NewExecID("execid#1"), 
    field.NewExecTransType(enum.ExecTransType_NEW), 
    field.NewExecType(enum.ExecType("0")), 
    field.NewOrdStatus("0"), 
    field.NewSymbol("EUR/USD"), 
    field.NewSide(enum.Side_BUY), 
    field.NewLeavesQty(decimal.NewFromFloat(float64(100)), 2), 
    field.NewCumQty(decimal.NewFromFloat(float64(100)), 2), 
    field.NewAvgPx(decimal.NewFromFloat(float64(1.1)), 2),
  )
  execReport.Header.SetTargetCompID(c.TargetCompID)
  execReport.Header.SetSenderCompID(c.SenderCompID)
  // send example report every 10s
  ticker := time.NewTicker(10 * time.Second)
  defer ticker.Stop()
  for {
    select {
    case <-ticker.C:
      // send exec report
      if err := quickfix.Send(execReport); err != nil {
        log.Fatalf(" [CONSUMER] [ERROR] [Send]: Failed: %s\n", err)
      }
    }
  }
}

func main() {
  log.Println(" [CONSUMER] [INFO] [*]: Starting FIX acceptor...")
  app := Consumer{}
  app.NewApp("DEVELOPMENT")
  cfg, err := os.Open("consumer/quickfix.cfg")
  if err != nil {
    log.Fatalf(" [CONSUMER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  appSettings, err := quickfix.ParseSettings(cfg)
  if err != nil {
    log.Fatalf(" [CONSUMER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  for _, sessionSettings := range appSettings.SessionSettings() {
    targetCompID, err := sessionSettings.Setting(config.TargetCompID)
    if err != nil {
      log.Fatalf(" [CONSUMER] [ERROR] [Startup]: Failed: %s\n", err)
      return
    }
    app.TargetCompID = targetCompID
    senderCompID, err := sessionSettings.Setting(config.SenderCompID)
    if err != nil {
      log.Fatalf(" [CONSUMER] [ERROR] [Startup]: Failed: %s\n", err)
      return
    }
    app.SenderCompID = senderCompID
  }
  storeFactory := quickfix.NewMemoryStoreFactory()
  logFactory, err := quickfix.NewFileLogFactory(appSettings)
  if err != nil {
    log.Fatalf(" [CONSUMER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  acceptor, err := quickfix.NewAcceptor(app, storeFactory, appSettings, logFactory)
  if err != nil {
    log.Fatalf(" [CONSUMER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  if err := acceptor.Start(); err != nil {
    log.Fatalf(" [CONSUMER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  go app.sendExampleReport()
  shutdownChannel := make(chan os.Signal, 1)
  signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
  <-shutdownChannel
  log.Println(" [CONSUMER] [INFO] [*]: Stopping FIX acceptor...")
  acceptor.Stop()
  log.Println(" [CONSUMER] [INFO] [*]: Stopped FIX acceptor")
}
