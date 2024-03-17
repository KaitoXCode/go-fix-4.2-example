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

// Producer implements quickfix.Application interface
type Producer struct {
  Env string
  SessionID quickfix.SessionID
  TargetCompID string
  SenderCompID string
}

// OnCreate
func (p Producer) OnCreate(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [OnCreate]: Creating session: %s", sessionID))
  p.SessionID = sessionID
}

// OnLogon
func (p Producer) OnLogon(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [OnLogon]: Starting session: %s\n", sessionID))
}

// OnLogout
func (p Producer) OnLogout(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [OnLogout]: Terminating session: %s\n", sessionID))
}

// FromAdimn
func (p Producer) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [FromAdmin]: Received: %s\n", msg.String()))
  return nil
}

// ToAdimn
func (p Producer) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [ToAdmin]: Sending: %s\n", msg.String()))
}

// ToApp
func (p Producer) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [ToApp]: Sending: %s\n", msg.String()))
  return
}

// FromApp
func (p Producer) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [FromApp]: Received: %s\n", msg.String()))
  return
}

func (p Producer) NewApp(env string) {
  p.Env = env
}

func main() {
  app := Producer{}
  app.NewApp("DEVELOPMENT")
  cfg, err := os.Open("producer/quickfix.cfg")
  if err != nil {
    log.Fatalf(" [PRODUCER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  appSettings, err := quickfix.ParseSettings(cfg)
  if err != nil {
    log.Fatalf(" [PRODUCER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  for _, sessionSettings := range appSettings.SessionSettings() {
    targetCompID, err := sessionSettings.Setting(config.TargetCompID)
    if err != nil {
      log.Fatalf(" [PRODUCER] [ERROR] [Startup]: Failed: %s\n", err)
      return
    }
    app.TargetCompID = targetCompID
    senderCompID, err := sessionSettings.Setting(config.SenderCompID)
    if err != nil {
      log.Fatalf(" [PRODUCER] [ERROR] [Startup]: Failed: %s\n", err)
      return
    }
    app.SenderCompID = senderCompID
  }
  storeFactory := quickfix.NewMemoryStoreFactory()
  logFactory, err := quickfix.NewFileLogFactory(appSettings)
  if err != nil {
    log.Fatalf(" [PRODUCER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  initiator, err := quickfix.NewInitiator(app, storeFactory, appSettings, logFactory)
  if err != nil {
    log.Fatalf(" [PRODUCER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  if err := initiator.Start(); err != nil {
    log.Fatalf(" [PRODUCER] [ERROR] [Startup]: Failed: %s\n", err)
    return
  }
  shutdownChannel := make(chan os.Signal, 1)
  signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
  <-shutdownChannel
  initiator.Stop()
}
