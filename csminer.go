// Copyright 2020 cryptonote.social. All rights reserved. Use of this source code is governed by
// the license found in the LICENSE file.
package csminer

import (
	"flag"
	"fmt"
	"github.com/cryptonote-social/csminer/crylog"
	"strconv"
	"strings"
)

const (
	APPLICATION_NAME = "bgminer"
	VERSION_STRING   = "0.3.3"
	STATS_WEBPAGE    = "https://google.com"
	DONATE_USERNAME  = "donate-getmonero-org"

	INVALID_EXCLUDE_FORMAT_MESSAGE = "invalid format for exclude specified. Specify XX-YY, e.g. 11-16 for 11:00am to 4:00pm."
)

var (
	saver   = flag.Bool("saver", false, "run only when screen is locked")
	t       = flag.Int("threads", 3, "number of threads")
	uname   = flag.String("user", "sm4rt", "your pool username")
	rigid   = flag.String("rigid", "csminer", "your rig id")
	tls     = flag.Bool("tls", false, "whether to use TLS when connecting to the pool")
	exclude = flag.String("exclude", "", "pause mining during these hours, e.g. -exclude=11-16 will pause mining between 11am and 4pm")
	config  = flag.String("config", "", "advanced pool configuration options, e.g. start_diff=1000;donate=1.0")
	wallet  = flag.String("wallet", "", "your wallet id. only specify this when establishing a new username, or specifying a 'secure' config change such as a change in donation amount")
	dev     = flag.Bool("dev", false, "whether to connect to dev server")
)

func MultiMain(s MachineStater, agent string) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "==== %s %s ====\n", APPLICATION_NAME, VERSION_STRING)
		fmt.Fprint(flag.CommandLine.Output(),
			`Usage of ./csminer
  -user <string>
    	your pool username (default "donate-getmonero-org")
  -saver=<bool>
    	mine only when screen is locked (default true)
  -exclude <string>
        pause mining during the specified hours. Format is XX-YY where XX and YY are hours of
        the day designated in 24 hour time. For example, -exclude=11-16 will pause mining betwen
        11:00am and 4:00pm. This can be used, for example, to pause mining during times of high
        machine usage or high electricity rates.
  -threads <int>
    	number of threads (default 1)
  -rigid <string>
    	your rig id (default "csminer")
  -tls <bool>
       whether to use TLS when connecting to the pool (default false)
  -config <string>
        advanced pool config option string, for specifying starting diff, donation percentage,
        email address for notifications, and more. See "advanced configuration options" under Get
        Started on the pool site for details. Some options will require you to also specify your
        wallet id (see below) in order to be changed.
  -wallet <string>
        your wallet id. You only need to specify this when establishing a new username, or if
        specifying a 'secure' config parameter change such as a new pool donation amount or email
        address. New usernames will be established upon submitting at least one valid share.
`)
		fmt.Fprintf(flag.CommandLine.Output(), "\nMonitor your miner progress at: %s\n", STATS_WEBPAGE)
		fmt.Fprint(flag.CommandLine.Output(), "Send feedback to: cryptonote.social@gmail.com\n")
	}
	flag.Parse()

	var hr1, hr2 int
	hr1 = -1
	var err error
	if len(*exclude) > 0 {
		hrs := strings.Split(*exclude, "-")
		if len(hrs) != 2 {
			crylog.Fatal(INVALID_EXCLUDE_FORMAT_MESSAGE)
			return
		}
		hr1, err = strconv.Atoi(hrs[0])
		if err != nil {
			crylog.Fatal(INVALID_EXCLUDE_FORMAT_MESSAGE, err)
			return
		}
		hr2, err = strconv.Atoi(hrs[1])
		if err != nil {
			crylog.Fatal(INVALID_EXCLUDE_FORMAT_MESSAGE, err)
			return
		}
		if hr1 > 24 || hr1 < 0 || hr2 > 24 || hr2 < 0 {
			crylog.Fatal("INVALID_EXCLUDE_FORMAT_MESSAGE", ": XX and YY must each be between 0 and 24")
			return
		}
	}
	fmt.Printf("==== %s v%s ====\n", APPLICATION_NAME, VERSION_STRING)
	if *uname == DONATE_USERNAME {
		fmt.Printf("\nNo username specified, mining on behalf of donate.getmonero.org.\n")
	}
	if *saver {
		fmt.Printf("\nNOTE: Mining only when screen is locked. Specify -saver=false to mine always.\n")
	}
	if *t == 1 {
		fmt.Printf("\nMining with only one thread. Specify -threads=X to use more.\n")
		fmt.Printf("Or use the [i] keyboard command to add threads dynamically.\n")
	}
	if hr1 != -1 {
		fmt.Printf("\nMining will be paused between the hours of %v:00 and %v:00.\n", hr1, hr2)
	}
	fmt.Printf("\nMonitor your mining progress at: %s\n", STATS_WEBPAGE)
	fmt.Printf("\nSend feedback to: cryptonote.social@gmail.com\n")

	fmt.Println("\n==== Status/Debug output follows ====")
	crylog.Info("Miner username:", *uname)
	crylog.Info("Threads:", *t)

	if hr1 == -1 || hr2 == -1 {
		hr1 = 0
		hr2 = 0
	}

	config := MinerConfig{
		MachineStater:  s,
		Threads:        *t,
		Username:       *uname,
		RigID:          *rigid,
		Wallet:         *wallet,
		Agent:          agent,
		Saver:          *saver,
		ExcludeHrStart: hr1,
		ExcludeHrEnd:   hr2,
		UseTLS:         *tls,
		AdvancedConfig: *config,
		Dev:            *dev,
	}
	if err = Mine(&config); err != nil {
		crylog.Fatal("Miner failed:", err)
	}
}
