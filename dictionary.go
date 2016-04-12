package jiebago

import (
	"math"
	"sync"

	"github.com/linkedin-inc/jiebago/dictionary"
)

// A Dictionary represents a thread-safe dictionary used for word segmentation.
type Dictionary struct {
	total, logTotal float64
	freqMap         map[string]float64
	posMap          map[string]string
	sync.RWMutex
}

// Load loads all tokens from given channel
func (d *Dictionary) Load(ch <-chan dictionary.Token) {
	d.Lock()
	for token := range ch {
		d.addToken(token)
	}
	d.Unlock()
	d.updateLogTotal()
}

// AddToken adds one token
func (d *Dictionary) AddToken(token dictionary.Token) {
	d.Lock()
	d.addToken(token)
	d.Unlock()
	d.updateLogTotal()
}

func (d *Dictionary) addToken(token dictionary.Token) {
	d.freqMap[token.Text()] = token.Frequency()
	if len(token.Pos()) > 0 {
		d.posMap[token.Text()] = token.Pos()
	}
	d.total += token.Frequency()
	runes := []rune(token.Text())
	n := len(runes)
	for i := 0; i < n; i++ { //TODO: n-1?
		frag := string(runes[:i+1])
		if _, ok := d.freqMap[frag]; !ok {
			d.freqMap[frag] = 0.0
		}
	}
}

func (d *Dictionary) updateLogTotal() {
	d.logTotal = math.Log(d.total)
}

// Frequency returns the frequency and existence of give word
func (d *Dictionary) Frequency(key string) (float64, bool) {
	d.RLock()
	freq, ok := d.freqMap[key]
	d.RUnlock()
	return freq, ok
}

// Pos returns the pos and existence of give word
func (d *Dictionary) Pos(key string) (string, bool) {
	d.RLock()
	pos, ok := d.posMap[key]
	d.RUnlock()
	return pos, ok
}

func (d *Dictionary) loadDictionary(fileName string) error {
	return dictionary.LoadDictionary(d, fileName)
}
