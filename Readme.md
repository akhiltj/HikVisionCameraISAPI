# HikVision IP Camera Alarm Server

Alarm Server for your Hikvision IP cameras!

Supported Cameras 📸:
  - Hikvision (Annke/Alarm.com/etc.)

Supported Delivery 📬:
  - MQTT
  - Webhooks

## Configuration

Create file `config.yaml` in the folder where the application binary lives.

Also see [sample config](docs/config.yaml).

When alarm server is coming online, it will also send a status message to `/camera/alarmserver` topic with its status.

#### Hikvision

Alarm Server uses HTTP streaming to connect to each camera individually and subscribe to all the events.

Some lower-end cameras, especially doorbells and intercoms, have broken HTTP streaming implementation that can't open more that 1 connection and "close" the http response, but leave TCP connection open (without sending keep-alive header!). For those, Alarm Server has an alternative streaming implementation. To use it, set `rawTcp: true` in camera's config file.

```yaml
hikvision:
  enabled: true              # if not enabled, it won't connect to any hikvision cams
  cams:
    myCam:                   # name of your camera
      address: 192.168.1.69  # ip address or domain name
      https: false           # if your camera supports ONLY https - set to true
      username: admin        # username that you use to log in to camera's web panel 
      password: admin1234    # password that you use to log in to camera's web panel
      rawTcp: false          # some cams have broken streaming. Set to true if normal HTTP streaming doesn't work 
```


## Tested cameras:

- HikVision DS-2CD2047G2-LU
- HikVision DS-2CD2043G2-IU

If your camera works with this - create an issue with some details about it and a picture, and we'll post it here. 

## Feedback

This project was created for author's personal use and while it will probably work for you too, author has no idea if you use it the same way as author. To fit everyone's usage scenario better, it would be beneficial if you could describe how YOU use the Alarm Server.

If you just started to use Alarm Server in your network (or use it for long time!), if you like it (or hate it!) - feel free to create an issue and just share how it works for you and what cameras you have. I would be curious to know if it fits your use case well (or not at all!).

## License

MIT License


## Special Thanks to toxuin:
this code is ported from https://github.com/toxuin/alarmserver. you can refer full code hear.