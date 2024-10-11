package storage

import (
	"math/rand"
	"strconv"
	"testing"
)

type testCase struct {
	name  string
	key   string
	value string
}

func TestSetGet(t *testing.T) {
	cases := []testCase{
		{"hello world", "hello", "world"},
		{"hello world1", "hello", "world1"},
		{"H12317 w09482", "H12317", "w09482"},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)

			sValue := s.Get(c.key)

			if *sValue != c.value {
				t.Errorf("values not equal")
			}
		})
	}
}

type testType struct {
	name  string
	key   string
	value string
	t     Type
}

func TestSetGetType(t *testing.T) {
	cases := []testType{
		{"hello world", "hello", "world", TypeString},
		{"int value", "key", "666667778", TypeInt},
		{"int value", "we123112", "890273451", TypeInt},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)

			sValue := s.Get(c.key)

			if *sValue != c.value {
				t.Errorf("values not equal")
			}

			if getType(*sValue) != getType(c.value) {
				t.Errorf("value kinds not equal")
			}
		})
	}
}

type bench struct {
	name    string
	valleng int
	keylen  int
	testlen int
}

var cases = []bench{
	{"1 test", 10, 5, 10},
	{"2 test", 100, 10, 20},
	{"3 test", 1000, 15, 30},
	{"4 test", 10000, 20, 50},
	{"5 test", 100000, 25, 80},
}

func BenchmarkGet(b *testing.B) {
	for _, tCase := range cases {
		b.Run(tCase.name, func(b *testing.B) {

			x := strconv.Itoa(rand.Intn(tCase.valleng))
			y := strconv.Itoa(rand.Intn(tCase.keylen))

			n, _ := NewStorage()

			n.Set(x, y)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				n.Get(x)
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	for _, tCase := range cases {
		b.Run(tCase.name, func(b *testing.B) {
			listx := make([]string, 0)
			listy := make([]string, 0)
			for i := 0; i < tCase.testlen; i++ {
				x := strconv.Itoa(rand.Intn(tCase.valleng))
				y := strconv.Itoa(rand.Intn(tCase.keylen))

				listx = append(listx, x)
				listy = append(listy, y)

			}
			n, _ := NewStorage()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				n.Set(listx[i%tCase.testlen], listy[i%tCase.testlen])
			}
		})
	}
}

func BenchmarkSetGet(b *testing.B) {
	for _, tCase := range cases {
		b.Run(tCase.name, func(b *testing.B) {
			listx := make([]string, 0)
			listy := make([]string, 0)
			for i := 0; i < tCase.testlen; i++ {
				x := strconv.Itoa(rand.Intn(tCase.valleng))
				y := strconv.Itoa(rand.Intn(tCase.keylen))

				listx = append(listx, x)
				listy = append(listy, y)

			}
			n, _ := NewStorage()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				n.Set(listx[i%tCase.testlen], listy[i%tCase.testlen])
				n.Get(listx[i%tCase.testlen])
			}
		})
	}
}
