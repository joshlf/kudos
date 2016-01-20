package kudos

import (
	"crypto/rand"
	"encoding/hex"
)

type Anonymizer map[string]string

func NewAnonymizer() Anonymizer { return make(Anonymizer) }

func (a Anonymizer) NewToken(uid string) (string, error) {
	var nonce [8]byte
	var nstr string
	for {
		_, err := rand.Read(nonce[:])
		if err != nil {
			return "", err
		}
		nstr = hex.EncodeToString(nonce[:])
		if _, ok := a[nstr]; !ok {
			break
		}
	}
	a[nstr] = uid
	return nstr, nil
}

func (a Anonymizer) LookupToken(token string) (uid string, ok bool) {
	uid, ok = a[token]
	return
}
