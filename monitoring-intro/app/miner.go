package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
)

func handleMineKey(rw http.ResponseWriter, r *http.Request) {
	// parse request
	params := r.URL.Query()
	keyBytes := 4096 // default key keyBytes
	if rawSize := params.Get("keyBytes"); rawSize != "" {
		parsed, err := strconv.Atoi(rawSize)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = rw.Write([]byte(fmt.Sprintf("Invalid keyBytes parameter: %q", rawSize)))
			return
		}
		keyBytes = parsed
	}
	zeroBits := 0 // default required zero bits
	if rawZeroBits := params.Get("zeroBits"); rawZeroBits != "" {
		parsed, err := strconv.Atoi(rawZeroBits)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = rw.Write([]byte(fmt.Sprintf("Invalid zeroBits parameter: %q", rawZeroBits)))
			return
		}
		zeroBits = parsed
	}

	// handle request
	key, nonce := mineKey(keyBytes, zeroBits)

	// send response
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(map[string]string{
		"key":   key,
		"nonce": nonce,
	})

}

// mineKey generates a random key of keyBytes and a 64-bin nonce such that SHA-256(concat(nonce, key)) ends
// in zeroBits bits. The generated key is returned in Base64/URL format and nonce is returned in base16 format.
func mineKey(keyBytes int, zeroBits int) (keyB64 string, nonceB16 string) {
	key, _ := io.ReadAll(io.LimitReader(crand.Reader, int64(keyBytes)))

	buf := make([]byte, 8+keyBytes)
	copy(buf[8:], key)

	var nonce uint64
	for nonce = 0; nonce <= math.MaxUint64; nonce++ {
		binary.BigEndian.PutUint64(buf[:8], nonce)
		hash := sha256.Sum256(buf)
		if hasLeadingZeroBits(hash, zeroBits) {
			break
		}
	}

	return base64.URLEncoding.EncodeToString(key), fmt.Sprintf("%x", nonce)
}

func hasLeadingZeroBits(hash [32]byte, bits int) bool {
	leadingBits := 0
outer:
	for _, curByte := range hash {
		for i := 7; i >= 0; i-- {
			if (1<<i)&curByte != 0 {
				break outer
			}
			leadingBits++
		}
	}
	return leadingBits >= bits
}
