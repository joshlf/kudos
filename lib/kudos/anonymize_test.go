// +build ignore

package kudos

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestAnonymize(t *testing.T) {
	old := newAnonymizerUseGlobalRand
	newAnonymizerUseGlobalRand = true
	defer func() { newAnonymizerUseGlobalRand = old }()

	rand.Seed(2220431966)

	a := NewAnonymizer()

	for i := 0; i < 1000; i++ {
		uid := fmt.Sprint(uint16(rand.Uint32()))
		ciphertext := a.Anonymize(uid)
		got := a.Deanonymize(ciphertext)
		if got != uid {
			t.Errorf("unexpected uid: got %v; want %v (ciphertext: %v)", got, uid, ciphertext)
		}
	}
}
