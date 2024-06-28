package generalStructIPC

type OptionStruct struct {
	OptionType string `yaml:"optionType"`
	Options    []string `yaml:"options"`
}

type IpcStruct struct {
	IpcType string `yaml:"ipcType"`
	Option  []OptionStruct `yaml:"option"`
}
