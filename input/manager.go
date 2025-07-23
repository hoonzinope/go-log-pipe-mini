package input

import (
	"test_gluent_mini/confmanager"
	"test_gluent_mini/offset"
	"test_gluent_mini/shared"
)

var configData confmanager.Config

func init() {
	shared.M.Lock()
	defer shared.M.Unlock()
	offsets, err := offset.GetOffsetMap()
	if err != nil {
		panic("Error reading offsets: " + err.Error())
	}
	for file, off := range offsets {
		shared.OffsetMap[file] = off // Initialize offset map with existing offsets
	}
}

func Configure(config confmanager.Config) {
	configData = config
}

func ManagingNode() {
	for _, inputConfig := range configData.Inputs {
		input_chan := shared.InputChannel[inputConfig.Name]
		go ManagingFileNode(inputConfig, input_chan)
	}
}
