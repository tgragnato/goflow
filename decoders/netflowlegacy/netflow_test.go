package netflowlegacy

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeNetFlowV5(t *testing.T) {
	t.Parallel()

	data := []byte{
		0x00, 0x05, 0x00, 0x06, 0x00, 0x82, 0xc3, 0x48, 0x5b, 0xcd, 0xba, 0x1b, 0x05, 0x97, 0x6d, 0xc7,
		0x00, 0x00, 0x64, 0x3d, 0x08, 0x08, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x79, 0x0a, 0x80, 0x02, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x02, 0x4e,
		0x00, 0x82, 0x9b, 0x8c, 0x00, 0x82, 0x9b, 0x90, 0x1f, 0x90, 0xb9, 0x18, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x77, 0x0a, 0x81, 0x02, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x94,
		0x00, 0x82, 0x95, 0xa9, 0x00, 0x82, 0x9a, 0xfb, 0x1f, 0x90, 0xc1, 0x2c, 0x00, 0x12, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x81, 0x02, 0x01, 0x0a, 0x80, 0x02, 0x77,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0xc2,
		0x00, 0x82, 0x95, 0xa9, 0x00, 0x82, 0x9a, 0xfc, 0xc1, 0x2c, 0x1f, 0x90, 0x00, 0x16, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x01, 0x0a, 0x80, 0x02, 0x79,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x09, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x01, 0xf1,
		0x00, 0x82, 0x9b, 0x8c, 0x00, 0x82, 0x9b, 0x8f, 0xb9, 0x18, 0x1f, 0x90, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x01, 0x0a, 0x80, 0x02, 0x79,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x09, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x02, 0x2e,
		0x00, 0x82, 0x9b, 0x90, 0x00, 0x82, 0x9b, 0x9d, 0xb9, 0x1a, 0x1f, 0x90, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x79, 0x0a, 0x80, 0x02, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x0b, 0xac,
		0x00, 0x82, 0x9b, 0x90, 0x00, 0x82, 0x9b, 0x9d, 0x1f, 0x90, 0xb9, 0x1a, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	buf := bytes.NewBuffer(data)

	var decNfv5 PacketNetFlowV5
	assert.Nil(t, DecodeMessageVersion(buf, &decNfv5))
	assert.Equal(t, uint16(5), decNfv5.Version)
	assert.Equal(t, uint16(9), decNfv5.Records[0].Input)
}
