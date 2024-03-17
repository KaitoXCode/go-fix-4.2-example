package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/config"
)

// Initiator implements quickfix.Application interface
type Initiator struct {
  Env string
  SessionID quickfix.SessionID
  TargetCompID string
  SenderCompID string
}

// OnCreate
func (i Initiator) OnCreate(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [INITIATOR] [INFO] [OnCreate]: Creating session: %s", sessionID))
  i.SessionID = sessionID
}

// OnLogon
func (i Initiator) OnLogon(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [INITIATOR] [INFO] [OnLogon]: Starting session: %s\n", sessionID))
}

// OnLogout
func (i Initiator) OnLogout(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [INITIATOR] [INFO] [OnLogout]: Terminating session: %s\n", sessionID))
}

// FromAdimn
func (i Initiator) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [INITIATOR] [INFO] [FromAdmin]: Received: %s\n", msg.String()))
  return nil
}

// ToAdimn
func (i Initiator) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [INITIATOR] [INFO] [ToAdmin]: Sending: %s\n", msg.String()))
}

// ToApp
func (i Initiator) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
  println(fmt.Sprintf(" [INITIATOR] [INFO] [ToApp]: Sending: %s\n", msg.String()))
  return
}

// FromApp
func (i Initiator) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [INITIATOR] [INFO] [FromApp]: Received: %s\n", msg.String()))
  return
}

func (i Initiator) NewApp(env string) {
  i.Env = env
}

func main() {
  log.Println(" [INITIATOR] [INFO] [*]: Starting FIX initiator...")
  app := Initiator{}
  app.NewApp("DEVELOPMENT")
  cfg, err := os.Open("initiator/quickfix.cfg")
  if err != nil {
    log.Fatalf(" [INITIATOR] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  appSettings, err := quickfix.ParseSettings(cfg)
  if err != nil {
    log.Fatalf(" [INITIATOR] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  for _, sessionSettings := range appSettings.SessionSettings() {
    targetCompID, err := sessionSettings.Setting(config.TargetCompID)
    if err != nil {
      log.Fatalf(" [INITIATOR] [ERROR] [Startup]: Failed: %s\n", err)
      return
    }
    app.TargetCompID = targetCompID
    senderCompID, err := sessionSettings.Setting(config.SenderCompID)
    if err != nil {
      log.Fatalf(" [INITIATOR] [ERROR] [Startup]: Failed: %s\n", err)
      return
    }
    app.SenderCompID = senderCompID
  }
  storeFactory := quickfix.NewMemoryStoreFactory()
  logFactory, err := quickfix.NewFileLogFactory(appSettings)
  if err != nil {
    log.Fatalf(" [INITIATOR] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  initiator, err := quickfix.NewInitiator(app, storeFactory, appSettings, logFactory)
  if err != nil {
    log.Fatalf(" [INITIATOR] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  if err := initiator.Start(); err != nil {
    log.Fatalf(" [INITIATOR] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  shutdownChannel := make(chan os.Signal, 1)
  signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
  log.Println(" [INITIATOR] [INFO] [*]: Stopping FIX initiator...")
  <-shutdownChannel
  initiator.Stop()
  log.Println(" [INITIATOR] [INFO] [*]: Stopped FIX initiator")
}
