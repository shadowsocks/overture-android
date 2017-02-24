// Copyright (c) 2016 holyshawn. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"runtime"
    "syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/holyshawn/overture/core"
)

func main() {

	var (
		configPath      string
		logVerbose      bool
		processorNumber int
        vpn             bool
	)

	flag.StringVar(&configPath, "c", "./config.json", "config file path")
	flag.BoolVar(&logVerbose, "v", false, "verbose mode")
	flag.IntVar(&processorNumber, "p", runtime.NumCPU(), "number of processor to use")
    flag.BoolVar(&vpn, 'V', true, 'VPN mode')
	flag.Parse()

	if logVerbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if vpn {
		path := "protect_path"

		callback := func(fd int, sotype int) {

            if sotype == syscall.SOCK_STREAM {
                return
            }

			socket, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
			if err != nil {
				log.Println(err)
				return
			}
			defer syscall.Close(socket)

			C.set_timeout(C.int(socket))

			err = syscall.Connect(socket, &syscall.SockaddrUnix{Name: path})
			if err != nil {
				log.Println(err)
				return
			}

			C.ancil_send_fd(C.int(socket), C.int(fd))

			dummy := []byte{1}
			n, err := syscall.Read(socket, dummy)
			if err != nil {
				log.Println(err)
				return
			}
			if n != 1 {
				log.Println("Failed to protect fd: ", fd)
				return
			}
		}

		SetNetCallback(callback)
	}

	log.Info("If you need any help or want to check update, please visit the project repository: https://github.com/holyshawn/overture")

	runtime.GOMAXPROCS(processorNumber)

	core.Init(configPath)
}
