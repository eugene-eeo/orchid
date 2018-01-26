package liborchid_test

import "math/rand"
import "testing"
import "testing/quick"
import "github.com/eugene-eeo/orchid/elems"
import "github.com/stretchr/testify/assert"

// update([a,b), i) => [a',b')
// where b' >= i >= a' >= 0,
//       b' - a' <= Height, and
//            b' <= Max

func TestViewboxUpdate(t *testing.T) {
	err := quick.Check(func(max, height Int100) bool {
		m := int(max) + 1
		h := int(height) + 1
		viewbox := elems.NewViewbox(m, h)
		for j := 0; j < 100; j++ {
			i := rand.Intn(m)
			a, b := viewbox.Update(i)
			if !(b >= i && i >= a && a >= 0 &&
				b-a <= h &&
				b <= m) {
				return false
			}
		}
		return true
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestViewBoxLoHi(t *testing.T) {
	v := elems.NewViewbox(10, 10)
	a, b := v.Update(1)
	assert.Equal(t, a, v.Lo())
	assert.Equal(t, b, v.Hi())
}
