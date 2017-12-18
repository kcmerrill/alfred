package hook

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestHook(t *testing.T) {
	// register an add function
	Register("add", func(a, b, c int, r *int) {
		*r = a + b + c
	})
	// shove our results somewhere
	var result int
	Trigger("add", 1, 1, 1, &result)
	if result != 3 {
		log.Fatalf("result should equal 3")
	}
}
func TestRaceCondition(t *testing.T) {
	Register("race", func() {
		<-time.After(200 * time.Millisecond)
		return
	})
	// sweet ... now
	// lets try to trigger a bunch of things
	var wg sync.WaitGroup
	for x := 0; x < 100; x++ {
		wg.Add(1)
		go func() {
			Trigger("race")
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestPriorityOrder(t *testing.T) {
	Register("wordp", 3, func(word *string) {
		*word += "c"
	})
	Register("wordp", 1, func(word *string) {
		*word += "a"
	})
	Register("wordp", 2, func(word *string) {
		*word += "b"
	})

	wordp := ""
	Trigger("wordp", &wordp)
	if wordp != "abc" {
		log.Fatalf("abc was expected")
	}
}

func TestOrder(t *testing.T) {
	Register("word", func(word *string) {
		*word += "c"
	})
	Register("word", func(word *string) {
		*word += "a"
	})
	Register("word", func(word *string) {
		*word += "b"
	})
	word := ""
	Filter("word", &word)
	if word != "cab" {
		log.Fatalf("cab was expected")
	}
}

func TestPlugin(t *testing.T) {
	Register("extra-text-to-word", "python plugin/python.py")

	p := "hi"
	Filter("extra-text-to-word", &p)

	if p != "hi-from-plugin" {
		t.Fatalf("Unable to call plugin")
	}
}
