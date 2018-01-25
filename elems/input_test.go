package elems_test

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

func TestInputInsertDelete2(t *testing.T) {
	err := quick.Check(func(x Int100, y Int100) bool {
		input := elems.NewInput()
		for i := 0; i < int(x); i++ {
			input.Insert('k')
		}
		for i := 0; i < int(y); i++ {
			input.Delete()
		}
		return len(input.String()) == max(int(x-y), 0)
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestInputInsertMove(t *testing.T) {
	err := quick.Check(func(x Int100, y Int100, z Int100) bool {
		input := elems.NewInput()
		for i := 0; i < int(x); i++ {
			input.Insert('k')
		}
		input.Move(-int(y))
		for i := 0; i < int(z); i++ {
			input.Insert('k')
		}
		return len(input.String()) == int(x+z) &&
			input.Cursor() == max(int(x-y), 0)+int(z)
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
