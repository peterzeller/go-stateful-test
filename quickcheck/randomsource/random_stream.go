package randomsource

import (
	"encoding/binary"
	"github.com/peterzeller/go-fun/iterable"
	"math/rand"
)

// RandomStream is the stream of randomness.
type RandomStream = iterable.Iterator[byte]

// RandomSource is the source of randomness.
type RandomSource = iterable.Iterable[byte]

// RandomSourceFromBytes creates a RandomSource that cycles through the byte buffer
func RandomSourceFromBytes(buf []byte) RandomSource {
	return iterable.IterableFun[byte](func() iterable.Iterator[byte] {
		pos := 0
		return iterable.Fun[byte](func() (byte, bool) {
			//if len(buf) == 0 {
			//	return 0, true
			//}
			//v := buf[pos]
			//pos = (pos + 1) % len(buf)
			//return v, true
			if pos >= len(buf) {
				// when at end of buffer, return only zeros
				return 0, true
			}
			v := buf[pos]
			pos++
			return v, true
		})
	})
}

// FromSeed returns a RandomSource created from a seed
func FromSeed(seed int64) RandomSource {
	return iterable.IterableFun[byte](func() iterable.Iterator[byte] {
		r := rand.New(rand.NewSource(seed))
		return iterable.Fun[byte](func() (byte, bool) {
			v := r.Intn(256)
			return byte(v), true
		})
	})
}

// ZeroRandomSource returns only zeros
func ZeroRandomSource() RandomSource {
	return iterable.Generate(0, func(prev byte) byte {
		return 0
	})
}

func Uint64B(r RandomStream, numBytes int) uint64 {
	buf := make([]byte, 8)
	for i := 0; i < numBytes; i++ {
		b, ok := r.Next()
		if !ok {
			break
		}
		buf[i+8-numBytes] = b
	}
	return binary.BigEndian.Uint64(buf)
}

func Uint64(r RandomStream) uint64 {
	buf := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b, ok := r.Next()
		if !ok {
			break
		}
		buf[i] = b
	}
	return binary.BigEndian.Uint64(buf)
}

func Int64B(r RandomStream, byteSize int) int64 {
	return int64(Uint64B(r, byteSize))
}

func Int64(r RandomStream) int64 {
	return int64(Uint64(r))
}

// Float64 returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0).
func Float64(r RandomStream) float64 {
	// There are exactly 1<<53 float64s in [0,1). Use Intn(1<<53) / (1<<53).
	return float64(Uint64(r)<<11>>11) / (1 << 53)
}

func Uint64N(r RandomStream, n uint64) uint64 {
	if n <= 0 {
		panic("invalid argument to Int64N")
	}
	i := Uint64(r)
	// modulo is not really random, but for testing purposes we do not care
	return i % n
}

func Int64N(r RandomStream, n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int64N")
	}
	return int64(Uint64N(r, uint64(n)))
}

func IntN(r RandomStream, n int) int {
	if n <= 0 {
		panic("invalid argument to Int64N")
	}
	return int(Uint64N(r, uint64(n)))
}
