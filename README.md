# proftpd_notify

ftpmail wasn't working for me so i built a scuffed version in go because fuck perl

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

Similar to `ftpmail`, after configuring `proftpd_notify` just run it as a background process (eg. `proftpd_notify &` assuming it is in your `$PATH`).
The order in this case doesn't matter - you can run `proftpd_notify` before or after you've started `proftpd`.


