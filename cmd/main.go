package main

import (
	"flag"
	"github.com/matzew/khue/pkg/connector"
)

var (
	hueBridge string
	username  string
	lights    string
)

func init() {
	flag.StringVar(&hueBridge, "hueBridge", "", "URL of the Hue Bridge (CloudEvents)")
	flag.StringVar(&username, "username", "", "account on the Brdige")
	flag.StringVar(&lights, "lights", "", "account on the Brdige")
}

func main() {
	flag.Parse()

	ha := connector.NewHueAdapter(hueBridge, username)
	ha.ObserveLightState(lights)
}
