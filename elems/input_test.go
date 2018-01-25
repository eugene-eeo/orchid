package elems_test

import "strings"
import "math/rand"
import "reflect"
import "testing"
import "testing/quick"
import "github.com/eugene-eeo/orchid/elems"
import "github.com/stretchr/testify/assert"

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type Int100 int

func (i Int100) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(Int100(rand.Intn(100)))
}

// Insert(a), Delete(b) => max(a-b, 0)
// Insert(a), Move(b), Insert(c) => a+c
// Insert(a), Move(b), Delete(c) => max(a-b-c,0) + min(a,b)

func TestInputInsertDelete2(t *testing.T) {
	err := quick.Check(func(x Int100, y Int100) bool {
		input := elems.NewInput()
		for i := 0; i < int(x); i++ {
			input.Insert('k')
		}
		for i := 0; i < int(y); i++ {
			input.Delete()
		}
		return input.String() == strings.Repeat("k", max(int(x-y), 0))
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestInputInsertMove(t *testing.T) {
	err := quick.Check(func(x Int100, y Int100, z Int100) bool {
		a := int(x)
		b := int(y)
		c := int(z)
		input := elems.NewInput()
		for i := 0; i < a; i++ {
			input.Insert('a')
		}
		input.Move(-int(b))
		for i := 0; i < c; i++ {
			input.Insert('c')
		}
		target := strings.Repeat("a", max(a-b, 0)) +
			strings.Repeat("c", c) +
			strings.Repeat("a", min(b, a))
		return input.String() == target && input.Cursor() == max(a-b, 0)+int(c)
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestInputInsertMoveDelete(t *testing.T) {
	err := quick.Check(func(a Int100, b Int100, c Int100) bool {
		input := elems.NewInput()
		for i := 0; i < int(a); i++ {
			input.Insert('k')
		}
		input.Move(-int(b))
		for i := 0; i < int(c); i++ {
			input.Delete()
		}
		return len(input.String()) == max(int(a-b-c), 0)+min(int(a), int(b))
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertNatural(t *testing.T) {
	input := elems.NewInput()
	input.Insert('a')
	input.Insert('b')
	input.Insert('c')
	assert.Equal(t, "abc", input.String())
	input.Move(-1)
	input.Delete()
	input.Insert('k')
	assert.Equal(t, "akc", input.String())
}
