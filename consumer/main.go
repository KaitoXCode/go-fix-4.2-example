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

// Consumer implements quickfix.Application interface
type Consumer struct {
  Env string
  SessionID quickfix.SessionID
  TargetCompID string
  SenderCompID string
}

// OnCreate
func (p Consumer) OnCreate(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnCreate]: Creating session: %s\n", sessionID))
}

// OnLogon
func (p Consumer) OnLogon(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnLogon]: Starting session: %s\n", sessionID))
}

// OnLogout
func (p Consumer) OnLogout(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnLogout]: Terminating session: %s\n", sessionID))
}

// FromAdimn
func (p Consumer) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [FromAdmin]: Received: %s\n", msg.String()))
  return nil
}

// ToAdimn
func (p Consumer) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [ToAdmin]: Sending: %s\n", msg.String()))
}

// ToApp
func (p Consumer) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [ToApp]: Sending: %s\n", msg.String()))
  return
}

// FromApp
func (p Consumer) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [FromApp]: Received: %s\n", msg.String()))
  return
}

func (p Consumer) NewApp(env string) {
  p.Env = env
}

func main() {
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
  shutdownChannel := make(chan os.Signal, 1)
  signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
  <-shutdownChannel
  acceptor.Stop()
}
