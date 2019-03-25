package nipplebot

import (
  "fmt"
  "github.com/bwmarrin/discordgo"
  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "os"
  "os/signal"
  "syscall"
)

var email string
var password string
var token string
var googleKey string

var rootCmd = &cobra.Command{
  Use:     "nipplebot",
  Aliases: []string{"nb"},
  Short:   "nipplebot is a discord bot playing sounds in voice channels",
  Long: `A real
            discord
            voice
            bot`,
  Run: func(cmd *cobra.Command, args []string) {
    if token == "" {
      fmt.Println("No token provided")
      return
    }

    dg, err := discordgo.New("Bot " + token)
    CheckErr(err)

    // Register ready as a callback for the ready events.
    dg.AddHandler(ready)

    // Register messageCreate as a callback for the messageCreate events.
    dg.AddHandler(messageCreate)

    // Open the websocket and begin listening.
    err = dg.Open()
    CheckErr(err)

    fmt.Println("Listening...")
    lock := make(chan os.Signal, 1)
    signal.Notify(lock, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-lock

    // Cleanly close down the Discord session.
    err = dg.Close()
    CheckErr(err)
  },
}

func initConfig() {
  var log = logrus.New()
  // Don't forget to read config either from cfgFile or from home directory!
  viper.SetConfigFile(CfgFile)

  if err := viper.ReadInConfig(); err != nil {
    log.WithError(err).Warn("Can't read config")
    //os.Exit(1)
  }
}

func Execute() {
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringVar(&email, "user", "", "userName")
  rootCmd.PersistentFlags().StringVar(&password, "pass", "", "password")
  rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "token")
  rootCmd.PersistentFlags().StringVar(&googleKey, "googleKey", "", "googleKey")

  var err error
  err = viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
  err = viper.BindPFlag("pass", rootCmd.PersistentFlags().Lookup("pass"))
  err = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
  err = viper.BindPFlag("googleKey", rootCmd.PersistentFlags().Lookup("googleKey"))
  CheckErr(err)

  //err = rootCmd.MarkFlagRequired("token")
  //CheckErr(err)

  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
