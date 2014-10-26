package main

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Limiter struct {
	plugin    interface{}
	limit     int
	isPercent bool

	currentRPS  int
	currentTime int64
}

func parseLimitOptions(options string) (limit int, isPercent bool) {
	if strings.Contains(options, "%") {
		limit, _ = strconv.Atoi(strings.Split(options, "%")[0])
		isPercent = true
	} else {
		limit, _ = strconv.Atoi(options)
		isPercent = false
	}

	return
}

func NewLimiter(plugin interface{}, options string) io.ReadWriter {
	l := new(Limiter)
	l.limit, l.isPercent = parseLimitOptions(options)
	l.plugin = plugin
	l.currentTime = time.Now().UnixNano()

	return l
}

func (l *Limiter) isLimited() bool {
	if l.isPercent {
		return l.limit <= rand.Intn(100)
	}

	if (time.Now().UnixNano() - l.currentTime) > time.Second.Nanoseconds() {
		l.currentTime = time.Now().UnixNano()
		l.currentRPS = 0
	}

	if l.currentRPS >= l.limit {
		return true
	}

	l.currentRPS++

	return false
}

func (l *Limiter) Write(data []byte) (n int, err error) {
	if l.isLimited() {
		return 0, nil
	}

	n, err = l.plugin.(io.Writer).Write(data)

	return
}

func (l *Limiter) Read(data []byte) (n int, err error) {
	if l.isLimited() {
		return 0, nil
	}

	n, err = l.plugin.(io.Reader).Read(data)

	return
}

func (l *Limiter) String() string {
	return fmt.Sprintf("Limiting %s to: %d (isPercent: %b)", l.plugin, l.limit, l.isPercent)
}
