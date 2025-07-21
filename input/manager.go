package input

import (
	"context"
	"sync"
	"test_gluent_mini/confmanager"
	"test_gluent_mini/data"
	"test_gluent_mini/offset"
)

var m sync.RWMutex
var offsetMap = make(map[string]int64)
var cancelCtx context.Context
var inputChannel map[string]chan data.InputData
var configData confmanager.Config
var offsetChannel chan data.InputData
var cancelMap = make(map[string]context.CancelFunc)

func init() {
	m.Lock()
	defer m.Unlock()
	offsets, err := offset.GetOffsetMap()
	if err != nil {
		panic("Error reading offsets: " + err.Error())
	}
	for file, off := range offsets {
		offsetMap[file] = off // Initialize offset map with existing offsets
	}
}

func Configure(ctx context.Context,
	config confmanager.Config,
	inputChan map[string]chan data.InputData,
	offsetChan chan data.InputData) {

	configData = config
	inputChannel = inputChan
	cancelCtx, _ = context.WithCancel(ctx)
	offsetChannel = offsetChan
}

func ManagingNode() {
	for _, inputConfig := range configData.Inputs {
		input_chan := inputChannel[inputConfig.Name]
		go ManagingFileNode(inputConfig, input_chan)
	}
}
