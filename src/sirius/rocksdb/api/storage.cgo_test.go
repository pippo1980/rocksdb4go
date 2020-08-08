package api

import (
	"fmt"
	"github.com/facebookgo/ensure"
	"strconv"
	"testing"
)

func TestLifecycle(t *testing.T) {

	storage, err := NewCGOStorage("test", "/tmp/test")
	if err != nil {
		t.Error(err)
	} else {
		defer storage.Close(nil)
	}

	for i := 0; i < 1000; i++ {
		if err := storage.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i))); err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 1000; i++ {
		if val, err := storage.Get([]byte(strconv.Itoa(i))); err != nil {
			t.Error(err)
		} else {
			ensure.DeepEqual(t, []byte(strconv.Itoa(i)), val)
		}
	}

	for i := 0; i < 1000; i++ {
		if err := storage.Delete(nil, []byte(strconv.Itoa(i))); err != nil {
			t.Error(err)
		}
	}
}

func TestCGOStorage_GetBatch(t *testing.T) {

	storage, err := NewCGOStorage("test", "/tmp/test")
	if err != nil {
		t.Error(err)
	} else {
		defer storage.Close(nil)
	}

	for i := 0; i < 1000; i++ {
		if err := storage.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i))); err != nil {
			t.Error(err)
		}
	}

	if values, err := storage.GetBatch([]byte("11"), []byte("22"), []byte("33")); err != nil {
		t.Error(err)
	} else {
		ensure.DeepEqual(t, values[1], []byte("22"))
	}

	if values, err := storage.GetBatch([]byte("ee"), []byte("ff"), []byte("dd")); err != nil {
		t.Error(err)
	} else {
		ensure.DeepEqual(t, values[1], []byte(nil))
	}

}

func TestCGOStorage_Prefix(t *testing.T) {
	storage, err := NewCGOStorage("test", "/tmp/test")
	if err != nil {
		t.Error(err)
	} else {
		defer storage.Close(nil)
	}

	for i := 0; i < 1024; i++ {
		if err := storage.Put([]byte("pippo#"+strconv.Itoa(i)), []byte(strconv.Itoa(i))); err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 2048; i++ {
		if err := storage.Put([]byte("hippo#"+strconv.Itoa(i)), []byte(strconv.Itoa(i))); err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 1024; i++ {
		if err := storage.Put([]byte("pi#"+strconv.Itoa(i)), []byte(strconv.Itoa(i))); err != nil {
			t.Error(err)
		}
	}

	if keys, values, err := storage.PrefixGet([]byte("pippo"), 10000); err != nil {
		t.Error(err)
	} else {
		ensure.True(t, len(keys) == 1024)
		ensure.True(t, len(values) == 1024)
		ensure.DeepEqual(t, keys[1], []byte("pippo#1"))
	}

	if keys, values, err := storage.PrefixGet([]byte("hippo"), 10000); err != nil {
		t.Error(err)
	} else {
		ensure.True(t, len(keys) == 2048)
		ensure.True(t, len(values) == 2048)
		ensure.DeepEqual(t, keys[1], []byte("hippo#1"))
	}

	if keys, values, err := storage.PrefixGet([]byte("p"), 10000); err != nil {
		t.Error(err)
	} else {
		ensure.True(t, len(keys) == 2048)
		ensure.True(t, len(values) == 2048)
	}

	if keys, err := storage.PrefixDelete([]byte("hippo"), 10000); err != nil {
		t.Error(err)
	} else {
		fmt.Println(len(keys))
		ensure.True(t, len(keys) == 2048)
	}
}
