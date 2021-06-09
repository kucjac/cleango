package uniqueid

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	prefix     string
	rs         string
	generators []uint64
	contexts   []string
	l          sync.Mutex
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	rs = b64[0:10]
	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// Generator gets the next unique identifier at given generator context.
type Generator interface {
	NextId() string
}

type prefixGenerator int

// NextId implements Generator interface.
func (g prefixGenerator) NextId() string {
	thisID := atomic.AddUint64(&generators[g], 1)
	return fmt.Sprintf("%s-%s-%06d", prefix, contexts[g], thisID)
}

func NextNoPrefixGenerator(name string) Generator {
	l.Lock()
	defer l.Unlock()
	gen := len(generators)
	generators = append(generators, 0)
	contexts = append(contexts, name)
	return simpleGenerator(gen)
}

// NextGenerator gets the next generator with given package name context.
func NextGenerator(name string) Generator {
	l.Lock()
	defer l.Unlock()
	gen := len(generators)
	generators = append(generators, 0)
	contexts = append(contexts, name)
	return prefixGenerator(gen)
}

// NextBaseGenerator creates a new generator that gets the random entries in a format:
// RANDOM_BYTES-CONTEXT_NAME-CURRENT_ID
func NextBaseGenerator(name string) Generator {
	l.Lock()
	defer l.Unlock()
	gen := len(generators)
	generators = append(generators, 0)
	contexts = append(contexts, name)
	return baseGenerator(gen)
}

type simpleGenerator int

// NextId implements Generator interface.
func (g simpleGenerator) NextId() string {
	thisID := atomic.AddUint64(&generators[g], 1)
	return fmt.Sprintf("%s-%06d", contexts[g], thisID)
}

type baseGenerator int

// NextId implements Generator interface.
func (r baseGenerator) NextId() string {
	thisID := atomic.AddUint64(&generators[r], 1)
	return fmt.Sprintf("%s-%s-%06d", rs, contexts[r], thisID)
}
