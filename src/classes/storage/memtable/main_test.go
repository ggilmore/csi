package memtable

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	fuzz "github.com/google/gofuzz"
)

type entry struct {
	Key   []byte
	Value []byte
}

var (
	Foo = entry{Key: []byte("Foo"), Value: []byte("A")}
	Bar = entry{Key: []byte("Bar"), Value: []byte("B")}
	Baz = entry{Key: []byte("Baz"), Value: []byte("C")}
)

type factory func(t *testing.T) DB

func HashIndexFactory(t *testing.T) DB { return NewHashIndex() }

func SkipListFactory(options SkipListOptions) factory {
	return func(t *testing.T) DB {
		return NewSkipList(options)
	}
}
func CombinedSkipAndSSFactory(o CombinedSkipAndSSOptions) factory {
	return func(t *testing.T) DB {
		dir := t.TempDir()
		o.SSTableDir = dir

		combined, err := NewCombinedSkipAndSS(o)
		if err != nil {
			t.Fatalf("NewCombinedSkipAndSS: %s", err)
		}

		return combined
	}
}

func TestDelete(t *testing.T) {
	t.Run("hashtable", func(t *testing.T) {
		testDeletion(t, HashIndexFactory)
	})

	t.Run("skiplist", func(t *testing.T) {
		testDeletion(t, SkipListFactory(SkipListOptions{}))
	})

	t.Run("skiplist combined", func(t *testing.T) {
		testDeletion(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{}))
	})

	t.Run("skiplist combined small sstable size", func(t *testing.T) {
		testDeletion(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{SkipListSizeThresholdBytes: 1}))
	})
}

func TestPutAndGet(t *testing.T) {
	t.Run("hashtable", func(t *testing.T) {
		testPutGet(t, HashIndexFactory)
	})

	t.Run("skiplist", func(t *testing.T) {
		testPutGet(t, SkipListFactory(SkipListOptions{}))
	})

	t.Run("skiplist combined", func(t *testing.T) {
		testPutGet(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{}))
	})

	t.Run("skiplist combined small sstable size", func(t *testing.T) {
		testPutGet(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{SkipListSizeThresholdBytes: 1}))
	})
}

func TestHas(t *testing.T) {
	t.Run("hashtable", func(t *testing.T) {
		testHas(t, HashIndexFactory)
	})

	t.Run("skiplist", func(t *testing.T) {
		testHas(t, SkipListFactory(SkipListOptions{}))
	})

	t.Run("skiplist combined", func(t *testing.T) {
		testHas(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{}))
	})

	t.Run("skiplist combined small sstable size", func(t *testing.T) {
		testHas(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{SkipListSizeThresholdBytes: 1}))
	})
}

func TestRangeScan(t *testing.T) {
	t.Run("hashtable", func(t *testing.T) {
		testRangeScan(t, func(t *testing.T) DB { return NewHashIndex() })
	})

	t.Run("skiplist", func(t *testing.T) {
		testRangeScan(t, func(t *testing.T) DB { return NewSkipList(SkipListOptions{}) })
	})

	t.Run("skiplist combined", func(t *testing.T) {
		testRangeScan(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{}))
	})

	t.Run("skiplist combined small sstable size", func(t *testing.T) {
		testRangeScan(t, CombinedSkipAndSSFactory(CombinedSkipAndSSOptions{SkipListSizeThresholdBytes: 1}))
	})
}

func TestSkipListFuzz(t *testing.T) {
	db := NewSkipList(SkipListOptions{})

	entries := generateNUniqueEntries(t, 10000)

	for _, v := range entries {
		err := db.Put(v.Key, v.Value)
		if err != nil {
			t.Fatalf("unexpeceted error when putting value %s: %s", []byte(v.Key), err)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	for _, v := range entries {
		got, err := db.Get(v.Key)
		if err != nil {
			t.Fatalf("unexpeceted error when getting value %s: %s", []byte(v.Key), err)
		}

		if diff := cmp.Diff(v.Value, got); diff != "" {
			log.Println(entries)
			t.Errorf("unexpected diff in values for key %q (-want +got):\n%s", v.Key, diff)
		}

	}
}

func testDeletion(t *testing.T, factory func(t *testing.T) DB) {

	Foo := entry{Key: []byte("Foo"), Value: []byte("A")}
	BarOld := entry{Key: []byte("Bar"), Value: []byte("B-old")}
	BarNew := entry{Key: []byte("Bar"), Value: []byte("B-new")}
	Baz := entry{Key: []byte("Baz"), Value: []byte("C")}

	db := factory(t)

	// load all values
	for _, e := range []entry{Foo, BarOld, BarNew, Baz} {
		err := db.Put(e.Key, e.Value)
		if err != nil {
			t.Fatalf("unexpected error when putting key %q with value %q: %s", e.Key, e.Value, err)
		}
	}

	// ensure that we can retrieve the values that we put in
	for _, e := range []entry{Foo, BarNew, Baz} {
		v, err := db.Get(e.Key)
		if err != nil {
			t.Fatalf("unexpected error when getting key %q: %s", e.Key, err)
		}

		if diff := cmp.Diff(e.Value, v); diff != "" {
			t.Fatalf("unexpected diff in values (-want +got):\n%s", diff)
		}
	}

	// delete one value
	err := db.Delete(BarNew.Key)
	if err != nil {
		t.Fatalf("unexpected error when deleting key %q: %s", Bar.Key, err)
	}

	// ensure that we can still retrieve the values that we didn't delete
	for _, e := range []entry{Foo, Baz} {
		v, err := db.Get(e.Key)
		if err != nil {
			t.Fatalf("unexpected error when getting key %q: %s", e.Key, err)
		}

		if diff := cmp.Diff(e.Value, v); diff != "" {
			t.Fatalf("unexpected diff in values for key %q (-want +got):\n%s", e.Key, diff)
		}
	}

	// ensure that the key we deleted can't be retrieved
	_, err = db.Get(BarNew.Key)
	if !errors.Is(err, KeyNotFound) {
		t.Errorf("expected KeyNotFound error for deleted key %q, got: %s", Bar.Key, err)
	}
}

func testPutGet(t *testing.T, factory factory) {
	tests := []struct {
		name   string
		values []entry
	}{
		{
			name: "multiple",
			values: []entry{
				Foo,
				Bar,
				Baz,
			},
		},
		{
			name:   "fuzz",
			values: generateNUniqueEntries(t, 50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := factory(t)

			for _, v := range tt.values {
				err := db.Put(v.Key, v.Value)
				if err != nil {
					t.Fatalf("unexpeceted error when putting value %s: %s", []byte(v.Key), err)
				}
			}

			for _, v := range tt.values {
				got, err := db.Get(v.Key)
				if err != nil {
					t.Fatalf("unexpeceted error when getting value %s: %s", []byte(v.Key), err)
				}

				if diff := cmp.Diff(v.Value, got); diff != "" {
					t.Errorf("unexpected diff in values (-want +got):\n%s", diff)
				}

			}

		})
	}

	t.Run("should return key not found error for unknown key", func(t *testing.T) {
		db := factory(t)
		err := db.Put(Foo.Key, Foo.Value)
		if err != nil {
			t.Errorf("unexpected error when putting value %q, %s", "Foo", err)
		}

		_, err = db.Get([]byte("wut"))
		if !errors.Is(err, KeyNotFound) {
			t.Errorf("expected KeyNotFound error, got: %s", err)
		}
	})
	t.Run("should return the latest version of a key if it was inserted multiple times", func(t *testing.T) {
		fooV1 := entry{Key: []byte("foo"), Value: []byte("v1")}
		fooV2 := entry{Key: []byte("foo"), Value: []byte("v2")}
		fooV3 := entry{Key: []byte("foo"), Value: []byte("v3")}
		entries := []entry{
			fooV1, Bar, fooV2, Baz, fooV3,
		}

		db := factory(t)
		for _, e := range entries {
			err := db.Put(e.Key, e.Value)
			if err != nil {
				t.Fatalf("unexpected error when putting (%q, %q): %s", e.Key, e.Value, err)
			}
		}

		for _, e := range []entry{Bar, Baz, fooV3} {
			value, err := db.Get(e.Key)
			if err != nil {
				t.Fatalf("unexpected error when retrieving value for %q: %s", e.Key, err)
			}
			if diff := cmp.Diff(e.Value, value); diff != "" {
				t.Errorf("unexpected diff in value (-want +got):\n%s", diff)
			}
		}

	})
}

func testHas(t *testing.T, factory factory) {
	db := factory(t)

	for _, e := range []entry{Foo, Bar, Baz} {
		contains, err := db.Has(e.Key)
		if err != nil {
			t.Fatalf("unexpected error when checking to see if db contains %q: %s", e.Key, err)
		}

		if contains {
			t.Fatalf("db contains %q when it hasn't been loaded yet", e.Key)
		}
	}

	for _, e := range []entry{Foo, Baz} {
		err := db.Put(e.Key, e.Value)
		if err != nil {
			t.Fatalf("unexpeceted error when loading %q: %s", e.Key, err)
		}
	}

	for _, e := range []entry{Foo, Baz} {
		contains, err := db.Has(e.Key)
		if err != nil {
			t.Fatalf("unexpected error when checking to see if db contains %q: %s", e.Key, err)
		}

		if !contains {
			t.Errorf("db doesn't contain %q even though it has been loaded", e.Key)
		}
	}

	contains, err := db.Has(Bar.Key)
	if err != nil {
		t.Fatalf("unexpected error when checking to see if db contains %q: %s", Bar.Key, err)
	}

	if contains {
		t.Errorf("db contains %q when it was never loaded", Bar.Key)
	}
}

func testRangeScan(t *testing.T, factory factory) {
	One := entry{Key: []byte("1"), Value: []byte("A")}
	Two := entry{Key: []byte("2"), Value: []byte("B")}
	Three := entry{Key: []byte("3"), Value: []byte("C")}
	// skip 4
	FiveOld := entry{Key: []byte("5"), Value: []byte("D-old")}
	FiveNew := entry{Key: []byte("5"), Value: []byte("D-new")}
	Six := entry{Key: []byte("6"), Value: []byte("E")}
	// skip 7
	Eight := entry{Key: []byte("8"), Value: []byte("F")}
	Nine := entry{Key: []byte("9"), Value: []byte("G")}

	db := factory(t)

	for _, e := range []entry{One, Two, Three, FiveOld, Six, Eight, Nine, FiveNew} {
		err := db.Put(e.Key, e.Value)
		if err != nil {
			t.Fatalf("unexpected error when loading key %q: %s", e.Key, err)
		}
	}

	iter, err := db.RangeScan(One.Key, Nine.Key)
	if err != nil {
		t.Fatalf("unexpected error when range scanning [%q, %q): %s", One.Key, Nine.Key, err)
	}

	for _, e := range []entry{One, Two, Three, FiveNew, Six, Eight} {
		more := iter.Next()
		if !more {
			t.Fatalf("iterator stopped before yielding key %q", e.Key)
		}

		if iter.Error() != nil {
			t.Fatalf("unexpected error in iterator: %s", err)
		}

		key := iter.Key()
		if diff := cmp.Diff(e.Key, key); diff != "" {
			t.Errorf("unexpected diff in key (-want +got):\n%s", diff)
		}

		value := iter.Value()
		if diff := cmp.Diff(e.Value, value); diff != "" {
			t.Errorf("unexpected diff in value (-want +got):\n%s", diff)
		}
	}

	more := iter.Next()
	if more {
		t.Fatalf("iterator didn't stop after yielding what should have been the final key (%q)", Eight.Key)
	}
}

func TestSStable(t *testing.T) {
	entries := generateNUniqueEntries(t, 5000)

	o := SkipListOptions{}
	skip := NewSkipList(o)

	for _, e := range entries {
		err := skip.Put(e.Key, e.Value)
		if err != nil {
			t.Fatalf("loading skiplist with (%q, %q): %s", e.Key, e.Value, err)
		}
	}

	file, err := os.CreateTemp(t.TempDir(), "sstable")
	if err != nil {
		t.Fatalf("creating temp file for sstable: %s", err)
	}
	defer file.Close()

	err = skip.flushSSTable(file)
	if err != nil {
		t.Fatalf("flusing skip list to sstable file %q: %s", file.Name(), err)
	}

	sstable, err := OpenSStable(file)
	if err != nil {
		t.Fatalf("opening sstable file %q: %s", file.Name(), err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	for _, e := range entries {
		value, err := sstable.Get(e.Key)
		if err != nil {
			t.Fatalf("getting sstable value for key %q: %s", e.Key, err)
		}

		if diff := cmp.Diff(e.Value, value, cmpopts.EquateEmpty()); diff != "" {
			t.Fatalf("unexpected diff in values for key %q (-want +got):\n%s", e.Key, diff)
		}
	}
}

func generateNUniqueEntries(t *testing.T, n int) []entry {
	t.Helper()

	f := fuzz.New().NilChance(0)

	uniqueEntries := make(map[string]entry)

	for len(uniqueEntries) < n {
		var e entry
		f.Fuzz(&e)
		uniqueEntries[string(e.Key)] = e
	}

	var out []entry
	for _, e := range uniqueEntries {
		out = append(out, e)
	}

	return out
}
