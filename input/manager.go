package input

import (
	"test_gluent_mini/confmanager"
	"test_gluent_mini/shared"
)

var configData confmanager.Config

func Configure(config confmanager.Config) {
	configData = config
}

func ManagingNode() {
	for _, inputConfig := range configData.Inputs {
		input_chan := shared.InputChannel[inputConfig.Name]
		go ManagingFileNode(inputConfig, input_chan)
	}
}
