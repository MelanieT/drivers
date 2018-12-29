package ph_board

import (
	"encoding/json"
	"fmt"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

const driverName = "ph-board"

type Config struct {
	Address byte `json:"address"`
}

type driver struct {
	channels []hal.ADCChannel
	meta     hal.Metadata
}

func HalAdapter(c []byte, bus i2c.Bus) (hal.ADCDriver, error) {
	var config Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}
	if err := bus.WriteBytes(config.Address, []byte{0x06}); err != nil {
		return nil, err
	}
	if err := bus.WriteBytes(config.Address, []byte{0x08}); err != nil {
		return nil, err
	}

	ch := &channel{
		bus:  bus,
		addr: config.Address,
	}
	return &driver{
		channels: []hal.ADCChannel{ch},
		meta: hal.Metadata{
			Name:         "ph-board",
			Description:  "An ADS115 based analog to digital converted with onboard female BNC connector",
			Capabilities: []hal.Capability{hal.PH},
		},
	}, nil
}
func (d *driver) Metadata() hal.Metadata {
	return d.meta
}

func (d *driver) ADCChannels() []hal.ADCChannel {
	return d.channels
}

func (d *driver) ADCChannel(n string) (hal.ADCChannel, error) {
	if n != chName {
		return nil, fmt.Errorf("ph board has only a signle channel named %s", chName)
	}
	return d.channels[0], nil
}

func (d *driver) Close() error {
	return nil
}
