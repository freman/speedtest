package speedtest

import (
	"math/rand"
	"net/url"
)

func randomBytes(length int) []byte {
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
		letters       = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	)

	b := make([]byte, length)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func randomURLComponent() string {
	return string(randomBytes(16))
}

func makeURLRandom(u *url.URL) *url.URL {
	q := u.Query()
	q.Set("x", randomURLComponent())
	return u.ResolveReference(&url.URL{RawQuery: q.Encode()})
}
