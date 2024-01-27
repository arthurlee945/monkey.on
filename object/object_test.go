package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "Monkey is Name"}
	diff2 := &String{Value: "Monkey is Name"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content has different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content has different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("string with different content has same hash keys")
	}
}
