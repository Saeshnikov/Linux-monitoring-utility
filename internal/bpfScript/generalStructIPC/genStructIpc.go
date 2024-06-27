package generalStructIPC

type OptionStruct struct {
	OptionType string
	Options    []string
}

type IpcStruct struct {
	IpcType string
	Option  []OptionStruct
}
