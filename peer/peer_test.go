package peer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialize1(t *testing.T) {
	p := Peer{
		IP:   "10.10.10.10",
		Port: 55555,
	}
	r, err := p.BTSerialize()
	assert.Nil(t, err)
	assert.Equal(t, r, "d2:ip11:10.10.10.104:porti55555ee")
}

func TestSerialize2(t *testing.T) {
	p := Peer{
		IP:   "10.10.10.10",
		IA:   "1-ff00:0:1",
		Port: 55555,
	}
	r, err := p.BTSerialize()
	assert.Nil(t, err)
	assert.Equal(t, "d2:ia10:1-ff00:0:12:ip11:10.10.10.104:porti55555ee", r)
}

func TestSerialize3(t *testing.T) {
	p := Peer{
		ID:        "1000",
		IP:        "10.10.10.10",
		Port:      55555,
		InfoHash:  "deadbeef",
		Key:       "secret_key",
		BytesLeft: 10000,
	}
	r, err := p.BTSerialize()
	assert.Nil(t, err)
	assert.Equal(t, "d10:bytes_lefti10000e2:id4:10009:info_hash8:deadbeef2:ip11:10.10.10.103:key10:secret_key4:porti55555ee", r)
}

func TestSerialize4(t *testing.T) {
	p := Peer{
		ID:        "1000",
		IP:        "10.10.10.10",
		Port:      55555,
		InfoHash:  "deadbeef",
		Key:       "secret_key",
		BytesLeft: 10000,
	}
	r, err := p.BTSerialize()
	assert.Nil(t, err)
	assert.Equal(t, r, "d10:bytes_lefti10000e2:id4:10009:info_hash8:deadbeef2:ip11:10.10.10.103:key10:secret_key4:porti55555ee")
}

func TestSerialize5(t *testing.T) {
	p := Peer{
		ID:        "1000",
		IP:        "10.10.10.10",
		IA:        "1-ff00:0:1",
		Port:      55555,
		InfoHash:  "deadbeef",
		Key:       "secret_key",
		BytesLeft: 10000,
	}
	r, err := p.BTSerialize()
	assert.Nil(t, err)
	assert.Equal(t, "d10:bytes_lefti10000e2:ia10:1-ff00:0:12:id4:10009:info_hash8:deadbeef2:ip11:10.10.10.103:key10:secret_key4:porti55555ee", r)
}

func TestSerializeAndDeserialize1(t *testing.T) {
	peer := Peer{
		IP:   "10.10.10.10",
		Port: 55555,
	}
	serialized, err := peer.BTSerialize()
	assert.Nil(t, err)
	deserialized, err := BTDeserialize([]byte(serialized))
	assert.Nil(t, err)
	assert.Equal(t, peer, *deserialized)
}

func TestSerializeAndDeserialize2(t *testing.T) {
	peer := Peer{
		ID:        "1000",
		IP:        "10.10.10.10",
		IA:        "1-ff00:0:1",
		Port:      55555,
		InfoHash:  "deadbeef",
		Key:       "secret_key",
		BytesLeft: 10000,
	}
	serialized, err := peer.BTSerialize()
	fmt.Println(serialized)
	assert.Nil(t, err)
	deserialized, err := BTDeserialize([]byte(serialized))
	assert.Nil(t, err)
	assert.Equal(t, peer, *deserialized)
}
