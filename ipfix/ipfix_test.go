package ipfix

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"ipfix-gen/util"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestCheckE(t *testing.T) {
	fmt.Println(0x80)
	fmt.Println(1 << 7)
	fmt.Println(0x00 == uint8(0))
	fmt.Println(reflect.TypeOf(0x00).Size(), reflect.TypeOf(uint8(0)).Size())
	fmt.Println(reflect.TypeOf(0x00).Name(), reflect.TypeOf(uint8(0)).Name())
	fmt.Println([]byte{0x00}, []byte{uint8(0)})
	fmt.Println(len([]byte{0x00}), len([]byte{uint8(0)}))
}

func TestParseData(t *testing.T) {
	ip := net.ParseIP("127.0.3.8").To4()
	fmt.Println(len(ip), ip, reflect.TypeOf(ip).Name(), reflect.TypeOf(ip).Size())
	ip2 := net.IP{127, 0, 3, 8}
	fmt.Println(len(ip2), ip2, reflect.TypeOf(ip2).Name(), reflect.TypeOf(ip2).Size())

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, ip2)
	fmt.Println(buf.Bytes())

	//uint16
	var a uint16 = 23
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, a)
	fmt.Println("uint16 a:", buf.Bytes())

	//uint32
	var b uint32 = 23
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, b)
	fmt.Println("b uint32:", buf.Bytes())

	ip4 := net.ParseIP("2400:dd01:12:1028:0:3316:12b2:979d").To16()
	fmt.Println(len(ip4), ip4)
	//ip4 := net.IP{2400,0xdd01,0x1001,0x1028,0x9999,0x3316,0x12b2,0x979d}
	//fmt.Println(len(ip4),ip4)
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, ip4)
	fmt.Println(buf.Bytes())

	fmt.Println(len(mac))
}

func TestGetLength(t *testing.T) {
	fmt.Println(InfoModel[ElementKey{
		EnterpriseNo: 0,
		ElementID:    5,
	}].Type.minLen())
}

var mac, _ = net.ParseMAC("00-FF-1D-3C-84-D4")
var t, _ = time.Parse("yyyy-MM-dd HH:mm:ss", "2018-5-6 13:30:00")

//Error,can not parse every type.
var testMap = map[int]interface{}{
	4:  uint16(34),
	7:  uint32(32768),
	8:  net.ParseIP("10.10.29.8").To4(),
	62: net.ParseIP("2400:dd01:1001:1028:9999:3316:12b2:979d").To16(),
	21: t.UnixNano(),
	56: mac,
}

//need bytes arr before build
var testBytesMap = map[int]interface{}{
	4:  []byte{34},
	7:  util.HostTo2Net(32768),
	8:  net.ParseIP("10.10.29.8").To4(),
	62: net.ParseIP("2400:dd01:1001:1028:9999:3316:12b2:979d").To16(),
	21: util.HostTo4Net(uint32(t.UnixNano())),
	56: mac,
}

var IDs = []uint16{4, 7, 8, 62, 21, 56}
var Vals = []interface{}{[]byte{34}, util.HostTo2Net(32768), net.ParseIP("10.10.29.8").To4(), net.ParseIP("2400:dd01:1001:1028:9999:3316:12b2:979d").To16(),
	util.HostTo4Net(uint32(t.UnixNano())), mac}

func TestBuildIPFIX(t *testing.T) {
	tID := getTemplateID()
	//msg := buildMap(testBytesMap, tID)
	msg := BuildArr(IDs, Vals, tID)

	Filling(msg)
	js, _ := json.Marshal(msg)
	fmt.Println(bytes.NewBuffer(js).String())
	bs := Encode(*msg, 234234234)
	fmt.Println(bs)
	//send(bs,"127.0.0.1:2055","127.0.0.1:2055")
	send(bs, "10.10.28.139:8088", "159.226.26.107:4739")
}

func send(message []byte, srcAddr, dstAddr string) {
	laddr, err := net.ResolveUDPAddr("udp", srcAddr)
	if err != nil {
		panic(err)
	}
	udpDialer := net.Dialer{
		Timeout:   time.Second * time.Duration(10),
		LocalAddr: laddr,
	}
	conn, err := udpDialer.Dial("udp", dstAddr)
	if err != nil {
		panic(err)
	}

	count := 100
	for count > 0 {
		count--
		conn.Write(message)
	}
	conn.Close()
}

var templateID uint16 = 257

func getTemplateID() uint16 {
	templateID++
	if templateID >= 1<<15-1 {
		templateID = 257
	}
	return templateID
}
