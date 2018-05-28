package kv

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test_set_a_key(t *testing.T) {
	kv := New()
	kv.Set("name", "helbing")
	val, ok := kv.Get("name")

	if ok == false {
		t.Fatal("Get key was failure，return is false")
	}

	if val != "helbing" {
		t.Fatal("Get key was failure，return is", val)
	}
}

func Test_key_exipre(t *testing.T) {
	kv := New()
	kv.Set("name", "helbing", 1000)
	val, ok := kv.Get("name")

	if ok == false {
		t.Fatal("Get key was failure，return is false")
	}

	if val != "helbing" {
		t.Fatal("Get key was failure，return is", val)
	}

	time.Sleep(2 * time.Second)

	val, ok = kv.Get("name")

	if ok != false {
		t.Fatal("The key not expire")
	}
}

func Test_is_lru(t *testing.T) {
	kv := New()
	for i := 0; i < 1200; i++ {
		kv.Set(fmt.Sprintf("test_%d", i), strings.Repeat("test", 256))
	}

	if kv.currentMemory >= kv.maxMemory {
		t.Fatalf("Lru not doing, current memory is %d, max memory is %d", kv.currentMemory, kv.maxMemory)
	}
}

func Test_key_invalid(t *testing.T) {
	kv := New()
	_, err := kv.Set(strings.Repeat("helbing", 10), "helbing")

	if err != ErrKeyInValid {
		t.Fatal("The key is not invalid")
	}
}

func Test_value_invalid(t *testing.T) {
	kv := New()
	_, err := kv.Set("helbing", strings.Repeat("helbing", 1000))

	if err != ErrValueInvalid {
		t.Fatal("The value is not invalid")
	}
}

func Test_memory_size(t *testing.T) {
	kv := New()
	kv.Set("name", "helbing")

	if kv.currentMemory != int64(len("name")+len("helbing")) {
		t.Errorf("Memory size is error, the error size is %d, the true size is %d", kv.currentMemory, int64(len("helbing")))
	}
}

func Test_parse_size(t *testing.T) {

	size, err := parseSizeStr("2GB")

	if err != nil {
		t.Fatal("Has Fatal", err)
	}

	if size != 2*1024*1024*1024 {
		t.Fatal("Return size is Fatal，the Fatal size is", size)
	}
}
