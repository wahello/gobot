//nolint:forcetypeassert // ok here
package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*DirectPinDriver)(nil)

func initTestDirectPinDriver() *DirectPinDriver {
	a := newGpioTestAdaptor()
	a.digitalReadFunc = func(string) (int, error) {
		return 1, nil
	}
	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	a.pwmWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	a.servoWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	return NewDirectPinDriver(a, "1")
}

func TestDirectPinDriver(t *testing.T) {
	var ret map[string]interface{}
	var err interface{}

	d := initTestDirectPinDriver()
	assert.Equal(t, "1", d.Pin())
	assert.NotNil(t, d.Connection())

	ret = d.Command("DigitalRead")(nil).(map[string]interface{})

	assert.Equal(t, 1, ret["val"].(int))
	assert.Nil(t, ret["err"])

	err = d.Command("DigitalWrite")(map[string]interface{}{"level": "1"})
	require.EqualError(t, err.(error), "write error")

	err = d.Command("PwmWrite")(map[string]interface{}{"level": "1"})
	require.EqualError(t, err.(error), "write error")

	err = d.Command("ServoWrite")(map[string]interface{}{"level": "1"})
	require.EqualError(t, err.(error), "write error")
}

func TestDirectPinDriverStart(t *testing.T) {
	d := initTestDirectPinDriver()
	require.NoError(t, d.Start())
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver()
	require.NoError(t, d.Halt())
}

func TestDirectPinDriverOff(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.Off())

	a := newGpioTestAdaptor()
	d = NewDirectPinDriver(a, "1")
	require.NoError(t, d.Off())
}

func TestDirectPinDriverOffNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.Off(), "DigitalWrite is not supported by this platform")
}

func TestDirectPinDriverOn(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	require.NoError(t, d.On())
}

func TestDirectPinDriverOnError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.On())
}

func TestDirectPinDriverOnNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.On(), "DigitalWrite is not supported by this platform")
}

func TestDirectPinDriverDigitalWrite(t *testing.T) {
	adaptor := newGpioTestAdaptor()
	d := NewDirectPinDriver(adaptor, "1")
	require.NoError(t, d.DigitalWrite(1))
}

func TestDirectPinDriverDigitalWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.DigitalWrite(1), "DigitalWrite is not supported by this platform")
}

func TestDirectPinDriverDigitalWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.DigitalWrite(1))
}

func TestDirectPinDriverDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver()
	ret, err := d.DigitalRead()
	assert.Equal(t, 1, ret)
	require.NoError(t, err)
}

func TestDirectPinDriverDigitalReadNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	_, e := d.DigitalRead()
	require.EqualError(t, e, "DigitalRead is not supported by this platform")
}

func TestDirectPinDriverPwmWrite(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	require.NoError(t, d.PwmWrite(1))
}

func TestDirectPinDriverPwmWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.PwmWrite(1), "PwmWrite is not supported by this platform")
}

func TestDirectPinDriverPwmWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.PwmWrite(1))
}

func TestDirectPinDriverServoWrite(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	require.NoError(t, d.ServoWrite(1))
}

func TestDirectPinDriverServoWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.ServoWrite(1), "ServoWrite is not supported by this platform")
}

func TestDirectPinDriverServoWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.ServoWrite(1))
}

func TestDirectPinDriverDefaultName(t *testing.T) {
	d := initTestDirectPinDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Direct"))
}

func TestDirectPinDriverSetName(t *testing.T) {
	d := initTestDirectPinDriver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
