# Auto Printer

Auto Printer is a process which helps to automatically print any file received in your Gmail that matches certain criteria.

For example, you can configure this program so that prints every attachment received from myworkaddress@gmail.com with the subject "Print".

## Installation

Clone this repository and build:

```bash
go build
```

## Configuration

Authenticate using (following the instructions):

```bash
auto-printer auth
```

Then configure the program with the following command (writing in your own config):

```bash
sudo nano /etc/auto-printer.toml
```

Then create a crontab using:

```bash
crontab -e
```

And finally configure the process however you want, for example:

```bash
0,30 * * * *  cd /path/to/auto-printer && ./auto-printer service
```

## Configuration file

The configuration file is located at /etc/auto-printer.toml

The structure of this file is the following:

```
allowed_emails []string
allowed_email_subjects []string
printed_label_name string
```

**Note:** Subjects cannot have spaces.
