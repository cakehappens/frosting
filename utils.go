package frosting

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func newULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
