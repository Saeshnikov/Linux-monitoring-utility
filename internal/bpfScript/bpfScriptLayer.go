package bpfScriptLayer

import (
	"errors"
	"fmt"
	fs "linux-monitoring-utility/internal/bpfScript/fsorwScript"
	pipe "linux-monitoring-utility/internal/bpfScript/namedPipeScript"
	sem "linux-monitoring-utility/internal/bpfScript/semScript"
	sock "linux-monitoring-utility/internal/bpfScript/socketScript"
	"os"
)

type IPC int

const (
	StartIpc IPC = iota + 1
	Fsorw
	Socket
	NamedPipe
	Semaphore
	EndIpc
)

func (ipc IPC) String() string {
	return [...]string{"start", "fsorw", "socket", "namedpipe", "semaphore", "end"}[ipc-1]
}

func isValidIpc(ipc string) bool {
	for nameIpc := StartIpc + 1; nameIpc < EndIpc; nameIpc++ {
		tempStr := nameIpc.String()
		switch tempStr == ipc {
		case true:
			return true
		}
	}
	return false
}

func GenerateBpfScript(ipc map[string]map[string][]string, dirPath string, inode int) ([]*os.File, error) {
	rootInode := inode

	var returnFilesArr []*os.File

	for ipcType, option := range ipc {
		if isValidIpc(ipcType) {
			switch ipcType {
			case Socket.String():
				path, err := checkDir(dirPath, ipcType)
				if err != nil {
					return nil, err
				}
				err = createFile(path)
				if err != nil {
					return nil, err
				}
				file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
				if err != nil {
					return nil, err
				}
				defer file.Close()

				err = sock.MakeSocketScript(file, option, rootInode)
				if err != nil {
					return nil, err
				}
				returnFilesArr = append(returnFilesArr, file)
			case NamedPipe.String():
				path, err := checkDir(dirPath, ipcType)
				if err != nil {
					return nil, err
				}
				err = createFile(path)
				if err != nil {
					return nil, err
				}
				file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
				if err != nil {
					return nil, err
				}
				defer file.Close()

				err = pipe.MakeNamedPipeScript(file, option, rootInode)
				if err != nil {
					return nil, err
				}
				returnFilesArr = append(returnFilesArr, file)
			case Fsorw.String():
				path, err := checkDir(dirPath, ipcType)
				if err != nil {
					return nil, err
				}
				err = createFile(path)
				if err != nil {
					return nil, err
				}
				file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
				if err != nil {
					return nil, err
				}
				defer file.Close()

				err = fs.MakeFsorwScript(file, option, rootInode)
				if err != nil {
					return nil, err
				}
				returnFilesArr = append(returnFilesArr, file)
			case Semaphore.String():
				path, err := checkDir(dirPath, ipcType)
				if err != nil {
					return nil, err
				}
				err = createFile(path)
				if err != nil {
					return nil, err
				}
				file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
				if err != nil {
					return nil, err
				}
				defer file.Close()

				err = sem.MakeSemaphoreScript(file, option, rootInode)
				if err != nil {
					return nil, err
				}
				returnFilesArr = append(returnFilesArr, file)
			}
		} else {
			err := errors.New("The ipc is not valid.")
			return nil, err
		}
	}
	return returnFilesArr, nil
}

func checkDir(dirPath string, nameFile string) (string, error) {
	path := "./" + nameFile + ".bt"
	if len(dirPath) != 0 {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			err := errors.New("Путь '" + dirPath + "' не может быть создан.")
			return "", err
		}
		path = dirPath + "/" + nameFile + ".bt"
		return path, nil
	}
	return path, nil
}

func createFile(path string) error {
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	fmt.Println("==> file created successfully", path)
	return nil
}
