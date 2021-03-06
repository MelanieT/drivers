package ph_board

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

const chName = "0"

type channel struct {
	bus        i2c.Bus
	addr       byte
	calibrator hal.Calibrator
}

func NewChannel(b i2c.Bus, addr byte) (*channel, error) {
	c, err := hal.CalibratorFactory([]hal.Measurement{})
	if err != nil {
		return nil, err
	}
	return &channel{
		bus:        b,
		addr:       addr,
		calibrator: c,
	}, nil
}

func (c *channel) Name() string {
	return chName
}

func (c *channel) Calibrate(points []hal.Measurement) error {
	cal, err := hal.CalibratorFactory(points)
	if err != nil {
		return err
	}
	c.calibrator = cal
	return nil
}

func (c *channel) Read() (float64, error) {
	if err := c.bus.WriteBytes(c.addr, []byte{0x10}); err != nil {
		return -1, err
	}
	buf := make([]byte, 2)
	if err := c.bus.ReadFromReg(c.addr, 0x0, buf); err != nil {
		return -1, err
	}
	var v int16
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &v); err != nil {
		return -1, err
	}
	return float64(v), nil
}

func (c *channel) Measure() (float64, error) {
	v, err := c.Read()
	if err != nil {
		return 0, err
	}
	if c.calibrator == nil {
		return 0, fmt.Errorf("Not calibrated")
	}
	return c.calibrator.Calibrate(v), nil
}
