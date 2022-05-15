package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Required struct {
	host string
	port []string
}

type Optional struct {
	UDP        bool
	knockDelay bool
	knockMs    int
	connDelay  bool
	connMs     int
	verbose    bool
}

func showHelp() {
	fmt.Printf(`
Usage:
knock [-u] [-d=ms] [-c=ms] [-v] [-h] host port1 port2...

Required:
host   host name or address
port   port numbers to knock

Optional:
-u     knock to UDP ports (Default TCP)
-d=ms  delay between knocks (Default 250 ms)
-c=ms  connect timeout delay (Default 250 ms)
-v     verbose output
-h     show this help
`)
}

func timeout(v string) (delay bool, msec int, err error) {
	split := strings.Split(v, "=")
	if len(split) == 2 {
		msec, err = strconv.Atoi(split[1])
		if err != nil {
			return false, 0, errors.New("no int")
		}
	} else {
		return false, 0, errors.New("wrong key")
	}
	return true, msec, nil
}

func checkout(args []string) (req Required, opt Optional, err error) {
	var tail []string
	if len(args) == 0 {
		showHelp()
		os.Exit(0)
	}
	for _, v := range args {
		switch v[0] {
		case '-':
			switch v[1] {
			case 'h':
				showHelp()
				os.Exit(0)
			case 'u':
				opt.UDP = true
			case 'v':
				opt.verbose = true
			case 'd':
				// just dropping wrong keys
				opt.knockDelay, opt.knockMs, err = timeout(v)
				if err != nil {
					log.Printf("invalid -d key dropped\n")
					continue
				}
			case 'c':
				opt.connDelay, opt.connMs, err = timeout(v)
				if err != nil {
					log.Printf("invalid -c key dropped\n")
					continue
				}
			default:
				showHelp()
				os.Exit(0)
			}
		default:
			tail = append(tail, v)
		}
	}
	if len(tail) < 2 {
		showHelp()
		os.Exit(0)
	}
	req.host = tail[0]
	ip := net.ParseIP(req.host)
	if ip == nil {
		return Required{}, Optional{}, errors.New("invalid IP")
	}
	for _, port := range tail[1:] {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return Required{}, Optional{}, errors.New("invalid ports")
		}
		req.port = append(req.port, port)
	}
	return req, opt, nil
}

func main() {
	req := Required{}
	opt := Optional{}

	req, opt, err := checkout(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	var proto string
	if opt.UDP {
		proto = "udp"
	} else {
		proto = "tcp"
	}

	if !opt.knockDelay {
		opt.knockMs = 250
	}
	if !opt.connDelay {
		opt.connMs = 250
	}

	for _, port := range req.port {
		address := fmt.Sprintf("%s:%s", req.host, port)
		if opt.verbose {
			fmt.Printf("knocking to %s/%s\n", address, proto)
		}
		conn, _ := net.DialTimeout(proto, address,
			time.Duration(opt.connMs)*time.Millisecond)
		if conn != nil {
			if opt.UDP {
				_, _ = conn.Write([]byte("\xFF"))
			}
			_ = conn.Close()
		}
		time.Sleep(time.Duration(opt.knockMs) * time.Millisecond)
	}
}
