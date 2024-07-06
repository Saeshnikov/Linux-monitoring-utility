package generalStructIPC

type OptionStruct struct {
	OptionType string `yaml:"optionType"`
	Options    []string `yaml:"options"`
}

type IpcStruct struct {
	IpcType string `yaml:"ipcType"`
	Enable  bool `yaml:"enable"`
	Option  []OptionStruct `yaml:"option"`
}
