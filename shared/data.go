package shared

import (
	"context"
	"sync"
)

var bufferSize = 1000

var OffsetChannel = make(chan InputData, bufferSize)
var Ctx, Cancel = context.WithCancel(context.Background())
var CancelMap = make(map[string]context.CancelFunc)
var InputChannel = make(map[string]chan InputData, bufferSize)
var FilterChannel = make(map[string]chan InputData, bufferSize)
var M sync.RWMutex
var OffsetMap = make(map[string]int64)
var Wg sync.WaitGroup

type InputData struct {
	Raw      string
	Json     map[string]interface{}
	FileName string
	Tag      string
	Offset   int64
}
