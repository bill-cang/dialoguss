package configInit

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

type UssdOpt struct {
	Opt config `toml:"config"`
}

type config struct {
	Network       string
	Address       string
	Ports         []int
	ConnectionMax int `toml:"connection_max"`
}

func GetUSSDConfig() (uconfig *UssdOpt, err error) {
	uconfig = new(UssdOpt)
	file := *flag.String("config", "../config/ussdConfig.toml", "Path to toml config file.")
	flag.Parse()
	if _, err = toml.DecodeFile(file, uconfig); err != nil {
		fmt.Printf("DecodeFile err =%+v\n", err)
		return nil, err
	}
	//fmt.Printf("[GetUSSDConfig] uconfig =%+v", uconfig)
	return
}
