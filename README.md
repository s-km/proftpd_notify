# proftpd_notify

ftpmail wasn't working for me so I built a scuffed version in go because fuck perl :^)

I  realize you're supposed to toss pid files into `/var/run`, but because this is intended for use on a shared server without root access I had to be a little creative with where I put it.

## Installation

Clone the repo, then run `go mod download` followed by `go build -o proftpd_notify`.
Move the resulting binary into a folder included in your `$PATH` and make it executable (eg. `chmod +x proftpd_notify && mv proftpd_notify ~/bin`).

## Configuration

```
Required
--------
smtp_user - username used to authenticate with your SMTP server
smtp_pass - password used to authenticate with your SMTP server
smtp_host - SMTP server hostname
mail_to - Recipient email
mail_from - Sender email


Optional
--------
smtp_port - SMTP server port (default: 587)
log_dir - Directory containing your proftpd transfer log (default: $HOME/.config/proftpd)
log_name - Name of the log file (default: transfer.log)
```

Configuration is done via the `notify_config.json` file, either in your `$HOME` directory or under `$HOME/.config/proftpd`. `smtp_user` and `smtp_pass` can alternatively be set using the `PROFTPD_NOTIFY_SMTP_USER` and `PROFTPD_NOTIFY_SMTP_PASS` environment variables.


## Usage

After creating the config file, just run `proftpd_notify` to start the daemon.  The pid file is stored in the same directory as the transfer log, so you can stop the daemon using `kill $(cat /your/transferlog/dir/proftpd_notify.pid)`.
The order in this case doesn't matter - you can run `proftpd_notify` before or after you've started `proftpd`.
