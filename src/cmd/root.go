package nipplebot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var email string
var password string
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

		discord, err := discordgo.New(email, password)
		CheckErr(err)

		discord.AddHandler(messageCreate)

		// Open the websocket and begin listening.
		err = discord.Open()
		CheckErr(err)

		fmt.Println("Listening...")
		lock := make(chan int)
		<-lock

		// Cleanly close down the Discord session.
		err = discord.Close()
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
	rootCmd.PersistentFlags().StringVar(&googleKey, "googleKey", "", "google key")

	var err error
	err = viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	err = viper.BindPFlag("pass", rootCmd.PersistentFlags().Lookup("pass"))
	err = viper.BindPFlag("googleKey", rootCmd.PersistentFlags().Lookup("googleKey"))
	CheckErr(err)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
