package nipplebot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"io"
	"io/ioutil"
	"log"
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
		var err error
		loopCount := 0
		workingDirectory, err := ioutil.TempDir("", "ygf")
		if err != nil {
			workingDirectory = DefaultWorkingDirectory
		}
		// create instagram object
		insta := goinsta.New(userName, password)
		//CheckErr(err)
		db := GetMongoDb("ygf")

		workingCollection := db.Collection(DbCollectionWorking)
		_, err = workingCollection.DeleteMany(nil, nil)
		CheckErr(err)

		if useProxy {
			shuffledProxies := funk.ShuffleString(ProxyList)
			err = insta.SetProxy(shuffledProxies[0], false)
			CheckErr(err)
			log.WithFields(logrus.Fields{
				"proxy": shuffledProxies[0],
			}).Infof("*** USING PROXIED CONNECTION")
		}

		logFile, err := os.OpenFile(LogDirectory, os.O_WRONLY|os.O_CREATE, 0644)
		CheckErr(err)
		if silent {
			log.SetOutput(logFile)
		} else {
			mw := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(mw)
		}

		// --------------------------------------

		log.WithFields(logrus.Fields{
			"count": loopCount + 1,
			"total": runs,
		}).Infof("*** RUN START")

		err = Login(workingDirectory, insta, db)
		CheckErr(err)
		log.WithFields(logrus.Fields{
			"step": "login",
		}).Info("***")

		// print out important information about us before going to town
		meBefore, err := insta.Profiles.ByID(int64(insta.Account.ID))
		CheckErr(err)
		err = meBefore.Sync(true)
		CheckErr(err)
		log.WithFields(logrus.Fields{
			"id":        meBefore.ID,
			"username":  meBefore.Username,
			"follower":  meBefore.FollowerCount,
			"following": meBefore.FollowingCount,
		}).Info("our stats before")

		if refreshDb {
			// fetch & store all possible targets (to-be-followed users)
			err = FetchTargets(starsUserNames, insta, db)
			CheckErr(err)
			log.WithFields(logrus.Fields{
				"step": "fetchTargets",
			}).Info("***")
		} else {
			log.WithFields(logrus.Fields{
				"step": "fetchTargets",
			}).Infof("[NOT UPDATING]: not fetching targets!")
		}

		for loopCount < runs {
			// prepare fetched targets (to-be-followed users)
			err := PrepareTargets(insta, db)
			CheckErr(err)
			log.WithFields(logrus.Fields{
				"step": "prepareTargets",
			}).Info("***")

			// try to follow all cleaned up targets and store target if successful
			err = GoFollow(insta, db)
			CheckErr(err)
			log.WithFields(logrus.Fields{
				"step": "goFollow",
			}).Info("***")

			// --------------------------------------

			// wait until hopefully all newFollowings are following us
			if dryRun != true {
				WaitMinutes(60, 180)
				log.WithFields(logrus.Fields{
					"step": "waitMinutes",
				}).Info("***")
			} else {
				log.WithFields(logrus.Fields{
					"step": "waitMinutes",
				}).Infof("[DRY RUN MODE]: not waiting!")
			}

			// --------------------------------------

			// try to unfollow all previously followed users and store target if successful
			err = GoUnfollow(insta, db)
			CheckErr(err)
			log.WithFields(logrus.Fields{
				"step": "goUnfollow",
			}).Info("***")

			// print out important information about us after all is done
			meAfter, err := insta.Profiles.ByID(int64(insta.Account.ID))
			CheckErr(err)
			err = meAfter.Sync(true)
			CheckErr(err)
			log.WithFields(logrus.Fields{
				"id":                meAfter.ID,
				"username":          meAfter.Username,
				"follower":          meAfter.FollowerCount,
				"changed follower":  meAfter.FollowerCount - meBefore.FollowerCount,
				"following":         meAfter.FollowingCount,
				"changed following": meAfter.FollowingCount - meBefore.FollowingCount,
			}).Info("our stats after")

			// --------------------------------------

			log.SetOutput(os.Stderr)
			log.WithFields(logrus.Fields{
				"count": loopCount + 1,
				"total": runs,
			}).Infof("*** RUN DONE")

			if dryRun != true {
				WaitMinutes(60, 180)
				log.WithFields(logrus.Fields{
					"step": "waitMinutes",
				}).Info("***")
			} else {
				log.WithFields(logrus.Fields{
					"step": "waitMinutes",
				}).Infof("[DRY RUN MODE]: not waiting!")
			}
			loopCount++
		}

		err = Cleanup(insta)
		CheckErr(err)
		log.WithFields(logrus.Fields{
			"step": "cleanup",
		}).Infof("***")
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
	rootCmd.PersistentFlags().StringArrayVar(&starsUserNames, "stars", []string{}, "stars user names")

	var err error
	err = viper.BindPFlag("use-proxy", rootCmd.PersistentFlags().Lookup("use-proxy"))
	err = viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))
	err = viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))
	err = viper.BindPFlag("runs", rootCmd.PersistentFlags().Lookup("runs"))
	err = viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	err = viper.BindPFlag("pass", rootCmd.PersistentFlags().Lookup("pass"))
	err = viper.BindPFlag("stars", rootCmd.PersistentFlags().Lookup("stars"))
	CheckErr(err)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
