<div id="top"></div>

<div align="center">
  <a href="https://github.com/MarcoTomasRodriguez/eprinter">
    <img src="assets/logo.svg" alt="Logo" width="80" height="80">
  </a>
  <h2 align="center">eprinter</h2>
  <p align="center">
    Automatically print email attachments.
    <br />
    <a href="https://github.com/MarcoTomasRodriguez/eprinter"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/MarcoTomasRodriguez/eprinter">View Demo</a>
    ·
    <a href="https://github.com/MarcoTomasRodriguez/eprinter/issues">Report Bug</a>
    ·
    <a href="https://github.com/MarcoTomasRodriguez/eprinter/issues">Request Feature</a>
  </p>
</div>

## Description

Auto Printer is a process which helps to automatically print any file received in your Gmail that matches certain criteria.

For example, you can configure this program so that prints every attachment received from myworkaddress@gmail.com with the subject "Print".

<p align="right">(<a href="#top">back to top</a>)</p>

## Installation

Clone this repository and install:

```bash
go install github.com/MarcoTomasRodriguez/eprinter@latest
```

Download the credentials.json from your Google dev account and paste it in the program folder.

Authenticate using (following the instructions):

```bash
eprinter setup
```

Configure the program with the following command (writing in your own config):

```bash
sudo nano ~/.eprinter/config.toml
```

Create a crontab using:

```bash
crontab -e
```

Configure the process however you want, for example:

```bash
0,30 * * * *  eprinter
```

<p align="right">(<a href="#top">back to top</a>)</p>

## Configuration file

The configuration file is located at $HOME/.eprinter/config.toml

### Structure

```
allowed_emails []string
allowed_email_subjects []string
printed_label_name string
```

### Example

```toml
# Accept emails comming from "myemail@provider.com" and "someemail@provider.com".
allowed_emails = ["myemail@provider.com", "someemail@provider.com"]

# Accept emails with the title "Print", "Drucken" or "Imprimir".
allowed_email_subjects = ["Print", "Drucken", "Imprimir"]

# After the email was printed, add the label "Printed" so we do no print it twice.
printed_label = "Printed"
```

**Note:** Subjects cannot have spaces.

<p align="right">(<a href="#top">back to top</a>)</p>
