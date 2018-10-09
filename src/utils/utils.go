package utils

import (
	"encoding/base64"
	"encoding/binary"
	"hash/fnv"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

/******************************************************************************
 *  use example:
 *	type A struct {}
 *	Register(A{})
 *	a, _ = New("A")
 *	o, ok := a.(A)
*******************************************************************************/
var _registerTypes = make(map[string]reflect.Type, 0)

func Register(i interface{}) {
	t := reflect.TypeOf(i)
	_registerTypes[t.Name()] = t
}

func New(name string) (interface{}, bool) {
	t, ok := _registerTypes[name]
	if !ok {
		return nil, false
	}
	v := reflect.New(t)
	return v.Interface(), true
}

/*
			+----------------------------------------------------------------+
	    		| seq(4) | msgId(4) | msgData(SIZE-6)            			|
	    		+----------------------------------------------------------------+
*/
func UnpackMsg(buf []byte) (int32, uint32, []byte) {
	if len(buf) <= 6 {
		return 0, 0, nil
	}
	seq := binary.BigEndian.Uint32(buf[0:4])
	msgId := binary.BigEndian.Uint32(buf[4:8])
	msgData := buf[8:]
	return int32(seq), msgId, msgData
}

func PackMsg(seq int32, msgId uint32, msgData []byte) []byte {
	if msgId == 0 || len(msgData) == 0 {
		return nil
	}
	msg := make([]byte, len(msgData)+8)
	binary.BigEndian.PutUint32(msg[0:4], uint32(seq))
	binary.BigEndian.PutUint32(msg[4:8], msgId)
	copy(msg[8:], msgData)
	return msg
}

const (
	base64Table = "lmABdefghCDJNOP012QRSTUVWXYZabcijknopqrstuvGHwxyz3456KLM789+/EFI"
)

var _coder = base64.NewEncoding(base64Table)

func base64Encode(src []byte) ([]byte, error) {
	return []byte(_coder.EncodeToString(src)), nil
}

func base64Decode(src []byte) ([]byte, error) {
	return _coder.DecodeString(string(src))
}

func Base64EncodeV1(s string) string {
	b := []byte(s)
	return base64.StdEncoding.EncodeToString(b)
}

func Base64DecodeV1(s string) string {
	ret, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatal("Base64EncodeV1 faile")
	}
	return string(ret)
}

func Base64EncodeV2(s string) string {
	b := []byte(s)
	ret, err := base64Encode(b)
	if err != nil {
		log.Fatal("base64 encode fail")
	}
	return string(ret)
}

func Base64DecodeV2(s string) string {
	b := []byte(s)
	ret, err := base64Decode(b)
	if err != nil {
		log.Fatal("base64 decode fail")
	}
	return string(ret)
}

func CreateUUID() string {
	u1, err := uuid.NewV1()
	if err != nil {
		log.Fatal("Create UUID fail")
		return ""
	}
	return Base64EncodeV2(u1.String())
}

func CreateServiceId(serviceType string, addr string) string {
	ret := serviceType + "__" + Base64EncodeV1(addr)
	return string(ret[0 : len(ret)-1])
}

func GetServiceType(serviceId string) string {
	strs := strings.Split(serviceId, "__")
	if len(strs) != 2 {
		return ""
	}
	return strs[0]
}

var _isSeed = false

func Random(max int) int {
	if _isSeed == false {
		rand.Seed(time.Now().UnixNano())
		_isSeed = true
	}
	return rand.Intn(max)
}

func HashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
