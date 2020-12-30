package peer

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/jackpal/bencode-go"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
)

// Peer represents a scion bittorrent peer
type Peer struct {
	ID        string `json:"id,omitempty" bencode:"id,omitempty"`
	IP        string `json:"ip,omitempty" bencode:"ip,omitempty"`
	IA        string `json:"ia,omitempty" bencode:"ia,omitempty"`
	Port      int    `json:"port,omitempty" bencode:"port,omitempty"`
	InfoHash  string `json:"info_hash,omitempty" bencode:"info_hash,omitempty"`
	Key       string `json:"key,omitempty" bencode:"key,omitempty"`
	BytesLeft uint64 `json:"bytes_left,omitempty" bencode:"bytes_left,omitempty"`

	computedHash string `bencode:"-"`
}

// IsSeed returns true if the peer has no more bytes left to receive
func (p *Peer) IsSeed() bool {
	return p.BytesLeft == 0
}

// BTSerialize returns the peer's information serialized in the the bencoding format
func (p *Peer) BTSerialize() (string, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *p)
	return buf.String(), err
}

func BTDeserialize(d []byte) (*Peer, error) {
	p := new(Peer)
	buf := bytes.NewBuffer(d)
	err := bencode.Unmarshal(buf, &p)
	return p, err
}

// PeerFromRequest returns a peer from an http GET request
func PeerFromRequest(r *http.Request) (*Peer, error) {
	v := r.URL.Query()
	fmt.Printf("FINDME: hello\n")
	port, err := strconv.Atoi(v.Get("port"))
	if err != nil {
		return nil, err
	}

	left, err := strconv.ParseUint(v.Get("left"), 10, 64)
	if err != nil {
		return nil, err
	}

	ia, ip, err := getIaAndIpFromRequest(r)
	if err != nil {
		logrus.Debugln("could not parse scion ia / ip from request")
		return nil, err
	}

	p := &Peer{
		IP:        ip.String(),
		IA:        ia.String(),
		Port:      port,
		ID:        v.Get("peer_id"),
		InfoHash:  v.Get("info_hash"),
		Key:       v.Get("key"),
		BytesLeft: left,
	}

	return p, nil
}

func getIaAndIpFromRequest(r *http.Request) (*addr.IA, *net.IP, error) {
	udpAddr, err := snet.ParseUDPAddr(r.RemoteAddr)
	if err != nil {
		return nil, nil, err
	}
	return &udpAddr.IA, &udpAddr.Host.IP, nil
}

// Hash returns a sha1 of the peer ID and InfoHash
func (p *Peer) Hash() string {
	if p.computedHash == "" {
		hash := sha1.New()
		fmt.Fprint(hash, p.ID, p.InfoHash)

		p.computedHash = hex.EncodeToString(hash.Sum(nil))
	}

	return p.computedHash
}
