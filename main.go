package main

import (
	"fmt"
	"log"
	"net/smtp"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func notify(t *Transfer) {
	host := viper.GetString("smtp_host")
	port := viper.GetString("smtp_port")
	to := viper.GetString("mail_to")
	dest := fmt.Sprintf("%s:%s", host, port)

	msg := fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: New file %s\r\n", t.Direction) +
		"\r\n" +
		fmt.Sprintf("%s was %s by %s (%s) at %s.\nTransferred %s in %s.\r\n",
			t.Filename,
			t.Direction,
			t.User,
			t.RemoteHost,
			t.Date,
			t.FileSize,
			t.Duration)

	auth := smtp.PlainAuth("", viper.GetString("smtp_user"), viper.GetString("smtp_pass"), host)
	err := smtp.SendMail(dest, auth, viper.GetString("mail_from"), []string{to}, []byte(msg))
	HandleErr("Failed to send email:", err)
}

func handleLogEntry(xferCh chan fsnotify.Event) {
	dir := Expand(viper.GetString("log_dir"))
	name := viper.GetString("log_name")
	logFile := filepath.Join(dir, name)

	for {
		select {
		case ev := <-xferCh:
			if ev.Name != logFile {
				continue
			}

			xfer := ParseXferLogEntry(GetLatestXfer(logFile))
			notify(&xfer)
		}
	}
}

func watchTransferLog(xferCh chan fsnotify.Event, watcher *fsnotify.Watcher) {
	for {
		select {
		case event := <-watcher.Events:
			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				xferCh <- event
			}

		case err := <-watcher.Errors:
			log.Println(err)
		}
	}
}

func init() {
	viper.SetDefault("log_dir", "$HOME/.config/proftpd")
	viper.SetDefault("log_name", "transfer.log")
	viper.SetDefault("smtp_port", "587")

	viper.SetConfigName("notify_config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.config/proftpd")
	err := viper.ReadInConfig()
	HandleErr("Failed to read config file", err)

	viper.SetEnvPrefix("proftpd_notify")
	viper.AutomaticEnv()
}

func main() {
	done := make(chan bool)
	xferCh := make(chan fsnotify.Event)

	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	HandleErr("Failed to create new watcher:", err)

	dir := Expand(viper.GetString("log_dir"))
	err = watcher.Add(dir)
	HandleErr("Failed to watch directory:", err)

	go watchTransferLog(xferCh, watcher)
	go handleLogEntry(xferCh)
	<-done
}
