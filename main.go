// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/zhgo/config"
	"github.com/zhgo/console"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		c string

		replace = make(map[string]string)
	)

	flag.StringVar(&c, "c", console.WorkingDir+"/dogo.json", "Usage: dogo -c=/path/to/dogo.json")
	flag.Parse()

	buf, err := ioutil.ReadFile(c)
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	finder, err := regexp.Compile("{(.*?)}")
	if err != nil {
		log.Fatal(err)
	}
	replacer := strings.NewReplacer("{", "", "}", "")
	for _, key := range finder.FindAllString(string(buf), -1) {
		envVar := replacer.Replace(key)
		if value := os.Getenv(envVar); len(value) > 0 {
			if _, found := replace[key]; !found {
				fmt.Printf("[dogo] Expanding %s => %s\n", key, value)
				replace[key] = value
			}
		}
	}

	var dogo Dogo

	gopath := console.Getenv("GOPATH")
	c = strings.Replace(c, "{GOPATH}", gopath, -1)

	err = config.NewConfig(c).Replace(replace).Parse(&dogo)
	if err != nil {
		fmt.Printf("[dogo] Warning: no configuration file loaded.\n")
	} else {
		fmt.Printf("[dogo] Loaded configuration file:\n")
		fmt.Printf("       %s\n", c)
	}

	dogo.NewMonitor()

	l := len(dogo.Files)
	if l > 0 {
		fmt.Printf("[dogo] Ready. %d files to be monitored.\n\n", l)
		dogo.BuildAndRun()
		dogo.Monitor()
	} else {
		fmt.Printf("[dogo] Error: Did not find any files. Press any key to exit.\n\n")
		var a string
		fmt.Scanf("%s", &a)
	}
}
