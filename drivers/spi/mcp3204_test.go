package spi

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

// this ensures that the implementation is based on spi.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MCP3204Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3204Driver)(nil)

func initTestMCP3204DriverWithStubbedAdaptor() (*MCP3204Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor()
	d := NewMCP3204Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMCP3204Driver(t *testing.T) {
	var di interface{} = NewMCP3204Driver(newSpiTestAdaptor())
	d, ok := di.(*MCP3204Driver)
	if !ok {
		require.Fail(t, "NewMCP3204Driver() should have returned a *MCP3204Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "MCP3204"))
}

func TestMCP3204Read(t *testing.T) {
	tests := map[string]struct {
		chanNum     int
		simRead     []byte
		want        int
		wantWritten []byte
		wantErr     error
	}{
		"number_negative_error": {
			chanNum: -1,
			wantErr: fmt.Errorf("Invalid channel '-1' for read"),
		},
		"number_0_ok": {
			chanNum:     0,
			simRead:     []byte{0xFF, 0xFF, 0xFF},
			wantWritten: []byte{0x06, 0x00, 0x00},
			want:        0x0FFF,
		},
		"number_1_ok": {
			chanNum:     1,
			simRead:     []byte{0xFF, 0xFE, 0xFF},
			wantWritten: []byte{0x06, 0x40, 0x00},
			want:        0x0EFF,
		},
		"number_3_ok": {
			chanNum:     3,
			simRead:     []byte{0xFF, 0xF3, 0x21},
			wantWritten: []byte{0x06, 0xC0, 0x00},
			want:        0x0321,
		},
		"number_4_error": {
			chanNum: 4,
			wantErr: fmt.Errorf("Invalid channel '4' for read"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestMCP3204DriverWithStubbedAdaptor()
			a.spi.SetSimRead(tc.simRead)
			// act
			got, err := d.Read(tc.chanNum)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantWritten, a.spi.Written())
		})
	}
}

func TestMCP3204ReadWithError(t *testing.T) {
	// arrange
	d, a := initTestMCP3204DriverWithStubbedAdaptor()
	a.spi.SetReadError(true)
	// act
	got, err := d.Read(0)
	// assert
	require.ErrorContains(t, err, "error while SPI read in mock")
	assert.Equal(t, 0, got)
}
