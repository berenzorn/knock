# knock

## How to use it.

```
knock [-u] [-d=ms] [-c=ms] [-v] [-h] host port1 port2...
```

```
Required:
host   host name or address
port   port numbers to knock

Optional:
-u     knock to UDP ports (Default TCP)
-d=ms  delay between knocks (Default 250 ms)
-c=ms  connect timeout delay (Default 250 ms)
-v     verbose output
```

## Example
```
knock -u -d=500 -c=500 -v 12.3.4.56 1234 1235 1236 1237
```
Knocking to *UDP* ports *1234-1237* of IP *12.3.4.56* with *500 ms delay* between knocks and maximum *500 ms timeout* on every try to connect and some *verbose messages*.
