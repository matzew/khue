package connector

import (
	"context"
	"fmt"
	"github.com/amimof/huego"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"time"
)

type HueAdapter struct {
	bridge *huego.Bridge
}

func NewHueAdapter(host string, username string) *HueAdapter {
	return &HueAdapter{
		bridge: huego.New(host, username),
	}
}

func (ha *HueAdapter) GetLight(name string) (*huego.Light, error) {
	lights, err := ha.bridge.GetLights()

	if err != nil {
		return nil, err
	}
	for _, light := range lights {
		if light.Name == name {
			return &light, nil
		}
	}
	return nil, err
}

func (ha *HueAdapter) ObserveLightState(name string) {
	light, err := ha.GetLight(name)

	if err != nil {
		panic(err)
	}

	if light != nil {

		k_sink := os.Getenv("K_SINK")

		p, err := cloudevents.NewHTTP(cloudevents.WithTarget(k_sink))
		if err != nil {
			log.Fatalf("failed to create http protocol: %s", err.Error())
		}

		c, err := cloudevents.NewClient(p, cloudevents.WithUUIDs(), cloudevents.WithTimeNow())
		if err != nil {
			log.Fatalf("failed to create client: %s", err.Error())
		}

		// stash it
		initialState := light.State.On

		ticker := time.NewTicker(500 * time.Millisecond)
		for {

			light, _ := ha.GetLight(name)

			if light.State.On != initialState {
				initialState = light.State.On

				event := cloudevents.NewEvent()
				event.SetType("net.wessendorf.hue.light")
				event.SetSource(name)

				if err := event.SetData(cloudevents.ApplicationJSON, light); err != nil {
					log.Printf("failed to set cloudevents data: %s", err.Error())
				}

				if res := c.Send(context.Background(), event); !cloudevents.IsACK(res) {
					log.Printf("failed to send cloudevent: %v", res)
				}
			}

			// Wait for next tick
			<-ticker.C
		}
	} else {
		fmt.Println("None found...")
	}
}