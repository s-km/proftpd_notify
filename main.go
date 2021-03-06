package main

import (
	"fmt"
	"log"
	"net/smtp"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/sevlyar/go-daemon"
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

	log.Println("Successfully sent email")
}

func handleLogEntry(xferCh chan fsnotify.Event, logFile string) {
	for {
		select {
		case <-xferCh:
			xfer := ParseXferLogEntry(GetLatestXfer(logFile))
			log.Println("Parsed log entry; attempting to notify...")
			notify(&xfer)
		}
	}
}

func watchTransferLog(xferCh chan fsnotify.Event, logFile string, watcher *fsnotify.Watcher) {
	md5Cache := map[string]string{}

	for {
		select {
		case event := <-watcher.Events:
			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				if event.Name == logFile {
					checksum := Hash(event.Name)
					if md5Cache[event.Name] != checksum {
						xferCh <- event
						if checksum != "" {
							md5Cache[event.Name] = checksum
						}
					}
				}
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
	dir := Expand(viper.GetString("log_dir"))
	name := viper.GetString("log_name")
	logFile := filepath.Join(dir, name)

	ctx := &daemon.Context{
		PidFileName: dir + "/proftpd_notify.pid",
		PidFilePerm: 0644,
		LogFileName: dir + "/proftpd_notify.log",
		LogFilePerm: 0640,
	}

	d, err := ctx.Reborn()
	HandleErr("Failed to start: ", err)
	if d != nil {
		return
	}
	defer ctx.Release()

	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	HandleErr("Failed to create new watcher: ", err)

	err = watcher.Add(dir)
	HandleErr("Failed to watch directory: ", err)

	log.Println("Listening for changes to", logFile)
	go watchTransferLog(xferCh, logFile, watcher)
	go handleLogEntry(xferCh, logFile)
	<-done
}
