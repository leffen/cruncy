package cruncy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	os.Remove("dbtest.db")
	db, err := Open("dbtest.db")
	if err != nil {
		t.Fatal(err)
	}
	bucket := "test_bucket"

	db.CreateBucket(bucket)

	// put a key
	if err := db.Put(bucket, "key1", "value1"); err != nil {
		t.Fatal(err)
	}
	// get it back
	var val string
	if err := db.Get(bucket, "key1", &val); err != nil {
		t.Fatal(err)
	} else if val != "value1" {
		t.Fatalf("got \"%s\", expected \"value1\"", val)
	}
	// put it again with same value
	if err := db.Put(bucket, "key1", "value1"); err != nil {
		t.Fatal(err)
	}
	// get it back again
	if err := db.Get(bucket, "key1", &val); err != nil {
		t.Fatal(err)
	} else if val != "value1" {
		t.Fatalf("got \"%s\", expected \"value1\"", val)
	}
	// get something we know is not there
	if err := db.Get(bucket, "no.such.key", &val); err != ErrNotFound {
		t.Fatalf("got \"%s\", expected absence", val)
	}
	// delete our key
	if err := db.Delete(bucket, "key1"); err != nil {
		t.Fatal(err)
	}
	// delete it again
	if err := db.Delete(bucket, "key1"); err != ErrNotFound {
		t.Fatalf("delete returned %v, expected ErrNotFound", err)
	}
	// done
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestMoreNotFoundCases(t *testing.T) {
	os.Remove("dbtest.db")
	db, err := Open("dbtest.db")
	if err != nil {
		t.Fatal(err)
	}

	bucket := "test_bucket"

	db.CreateBucket(bucket)

	var val string
	if err := db.Get(bucket, "key1", &val); err != ErrNotFound {
		t.Fatal(err)
	}
	if err := db.Put(bucket, "key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(bucket, "key1"); err != nil {
		t.Fatal(err)
	}
	if err := db.Get(bucket, "key1", &val); err != ErrNotFound {
		t.Fatal(err)
	}
	if err := db.Get(bucket, "", &val); err != ErrNotFound {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

type aStruct struct {
	Numbers *[]int
}

func testGetPut(t *testing.T, inval string, outval *string) {
	os.Remove("dbtest.db")
	db, err := Open("dbtest.db")
	if err != nil {
		t.Fatal(err)
	}

	bucket := "test_bucket"

	db.CreateBucket(bucket)

	input, err := json.Marshal(inval)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Put(bucket, "test.key", inval); err != nil {
		t.Fatal(err)
	}

	if err := db.Get(bucket, "test.key", outval); err != nil {
		t.Fatal(err)
	}
	output, err := json.Marshal(outval)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, output) {
		t.Fatal("differences encountered")
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNil(t *testing.T) {
	os.Remove("dbtest.db")
	db, err := Open("dbtest.db")
	if err != nil {
		t.Fatal(err)
	}
	bucket := "test_bucket"

	db.CreateBucket(bucket)

	if err := db.Put(bucket, "key1", "value1"); err != nil {
		t.Fatal(err)
	}
	// can Get() into a nil value
	var data string
	if err := db.Get(bucket, "key1", &data); err != nil {
		t.Fatal(err)
	}
	db.Close()
}

func TestGoroutines(t *testing.T) {
	os.Remove("dbtest.db")
	db, err := Open("dbtest.db")
	if err != nil {
		t.Fatal(err)
	}
	bucket := "test_bucket"

	db.CreateBucket(bucket)

	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			switch rand.Intn(3) {
			case 0:
				if err := db.Put(bucket, "key1", "value1"); err != nil {
					t.Fatal(err)
				}
			case 1:
				var val string
				if err := db.Get(bucket, "key1", &val); err != nil && err != ErrNotFound {
					t.Fatal(err)
				}
			case 2:
				if err := db.Delete(bucket, "key1"); err != nil && err != ErrNotFound {
					t.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkPut(b *testing.B) {
	os.Remove("skv-bench.db")
	db, err := Open("skv-bench.db")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	bucket := "test_bucket"

	db.CreateBucket(bucket)

	for i := 0; i < b.N; i++ {
		if err := db.Put(bucket, fmt.Sprintf("key%d", i), "this.is.a.value"); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	db.Close()
}

func BenchmarkPutGet(b *testing.B) {
	os.Remove("skv-bench.db")
	db, err := Open("skv-bench.db")
	if err != nil {
		b.Fatal(err)
	}
	bucket := "test_bucket"

	db.CreateBucket(bucket)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.Put(bucket, fmt.Sprintf("key%d", i), "this.is.a.value"); err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < b.N; i++ {
		var val string
		if err := db.Get(bucket, fmt.Sprintf("key%d", i), &val); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	db.Close()
}

func BenchmarkPutDelete(b *testing.B) {
	os.Remove("skv-bench.db")
	db, err := Open("skv-bench.db")
	if err != nil {
		b.Fatal(err)
	}

	bucket := "test_bucket"

	db.CreateBucket(bucket)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.Put(bucket, fmt.Sprintf("key%d", i), "this.is.a.value"); err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < b.N; i++ {
		if err := db.Delete(bucket, fmt.Sprintf("key%d", i)); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	db.Close()
}
