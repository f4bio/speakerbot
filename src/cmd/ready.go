package nipplebot

import (
  "fmt"
  "github.com/bwmarrin/discordgo"
  "time"
)

func ready(s *discordgo.Session, m *discordgo.Ready) {
  fmt.Printf("%20s %20s %d\n", m.SessionID, time.Now().Format(time.Stamp), m.Version)

  err := s.UpdateStatus(0, "")

  CheckErr(err)
  return
}
