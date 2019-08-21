# AOS Listener

## Disclaimer
WARNING: This plugin is experimental and should not be used in production.

## Support 
This plugin supports AOS up to version 3.1.0.

## Overview
Input Plugins for Apstra AOS Telemetry Streaming
 - Configure Streaming Session on AOS Server
 - Listen on TCP Port and decode AOS GPB
 - Construct Time Series Information
 - Collect Additional info from AOS over REST API

## Configuration
```
[[inputs.aos]]
  # TCP Port to listen for incoming sessions from the AOS Server
  port = 7777                   # mandatory

  # Address of the server running Telegraf, it needs to be reacheable from AOS
  address = 192.168.59.1        # mandatory

  # Interval to refresh content from the AOS server (in sec)
  refresh_interval = 30         # Default 30

  # Streaming Type Can be "perfmon", "alerts" or "events"
  streaming_type = [ "events" ]

  # Define parameter to join the AOS Server using the REST API
  aos_server = 192.168.59.250   # mandatory
  aos_port = 8888               # Default 8888
  aos_login = admin             # Default admin
  aos_password = admin          # Default admin
```

