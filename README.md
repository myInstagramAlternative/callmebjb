## API for using gsm module over serial

Initially, I wrote this because I don't want to pay for a crappy Twilio service. I just needed a way to make calls to my phone to wake me up in case of an alarm on Home Assistant or some other service, especially if I'm not connected to the internet.

This API should work best with the SIM868 module and probably will work with most Waveshare GSM modules, as well as other AT command-based modems. If not, you will need to change some code.

The code is not great but it gets the job done. There isn't a queue or any special processor that watches serial output, so it's not meant to be used in production.

Used with [this](https://www.waveshare.com/gsm-gprs-gnss-hat.htm) GSM module and prepaid SIM card.

### TODO

- [ ] Simple auth
- [ ] Call endpoint with audio file

  I may add new endpoint for making calls with uploaded audio since most of the dev GSM modules support audio in/out so you can connect it your RPi and emulate micrphone.

- [ ] Read SMS inbox (if possible or create one in memory)

### Config example

```yaml
modem:
  port: "/dev/tty.usbserial-0001" # Find your port with `ls /dev/tty.*`
  baudrate: 115200 # RTFM
server:
  listen: "localhost"
  port: "8081"
```

btw you can set CONFIG_PATH env var

### Building

For RPi 0w

```sh
GOOS=linux GOARCH=arm GOARM=6 go build
```

See other GOARM values for RPi [here](https://zchee.github.io/golang-wiki/GoArm/) and [here](https://en.wikipedia.org/wiki/Raspberry_Pi#Specifications)

### Development

Generates OPENAPI spec with [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen)

```sh
go generate
```
