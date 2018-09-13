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
}

func NewHasherJd() Hasher {
	return &hasherJd{md5.New()}
}

func createSalt() string {
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
	saltHash := hexToInt(createSalt());
	saltBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(saltBytes, saltHash)
	s, _ := hex.DecodeString(createSalt())
	return append(
		s[:],
		h.binaryHash(
		append(saltBytes[:],
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
