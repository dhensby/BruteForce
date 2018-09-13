package hashs

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"strconv"
)

type hasherJd struct {
	cache hash.Hash
	saltHex []byte
	saltByt []byte
}

func NewHasherJd() Hasher {
	saltHex, _ := hex.DecodeString(getSaltString())
	saltBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(saltBytes, hexToInt(getSaltString()))
	return &hasherJd{md5.New(), saltHex, saltBytes }
}

func getSaltString() string {
	return "ce4f5046"
}

func hexToInt(data string) uint32 {
	n, err := strconv.ParseUint(data, 16, 32);
	if err == nil {
		return uint32(n)
	}
	return 0
}

func (h *hasherJd) Hash(data string) []byte {
	return append(
		h.saltHex[:],
		h.binaryHash(
		append(h.saltByt[:],
			data[:]...))[:]...)
}

func (h *hasherJd) IsValid(data string) bool {
	return len(data) == 8 + 32 &&
		genericBase64Validator(h, data)
}

func (h *hasherJd) binaryHash(data []byte) []byte {
	h.cache.Reset()
	h.cache.Write([]byte(data))
	return h.cache.Sum(nil)
}

func (h *hasherJd) convert(s string) []byte {
	return []byte(s)
}
