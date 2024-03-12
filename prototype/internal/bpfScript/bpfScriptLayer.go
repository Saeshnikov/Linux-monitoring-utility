package bpfScriptLayer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"
)

type Syscall int

const (
	Start Syscall = iota + 1
	Execve
	Execveat
	Open
	Openat
	Openat2
	Open_tree
	Statx
	Fspick
	Name_to_handle_at
	Readlink
	Readlinkat
	End
)

var (
	immutablePieces = map[string]string{
		"BEGIN":  "\n{\n\tprintf(" + `"Tracing file system syscalls... Hit Ctrl-C to end.\n"` + ");\n}\n",
		"END":    "\n{\n\tprint(@filename);\n\tclear(@oldname);\n\tclear(@filename);\n\tclear(@name);\n\tclear(@fd);\n}",
	}

	tmplSimplCom, _ = template.New("SimplCom").Parse("tracepoint:syscalls:sys_enter_{{.SyscallName}}\n{\n\t@filename[str(args.{{.ArgCollected}})] = count();\n}\n")

	tmplToHandleAt, _ = template.New("tmplToHandleAt").Parse("tracepoint:syscalls:sys_enter_{{.SyscallName}}\n{\n\t@name[tid] = args.name;\n}\ntracepoint:syscalls:sys_exit_name_to_handle_at\n/@name[tid]/\n{\n\t$ret = args.ret;\n\t@fd[tid] = $ret >= 0 ? $ret : -1;\n}\ntracepoint:syscalls:sys_enter_open_by_handle_at\n/@fd[tid]/\n{\n\t@filename[str(@name[tid])] = count();\n\tdelete(@fd[tid]);\n}\n")

	tmplToSymlink, _ = template.New("tmplToSymlink").Parse("tracepoint:syscalls:sys_enter_{{.SyscallName}}\n{\n\t@oldname[tid] = args.{{.ArgCollected}};\n}\n")

	tmplSimplComForSymlink, _ = template.New("SimplComForSymlink").Parse("tracepoint:syscalls:sys_enter_{{.SyscallName}}\n/@oldname[tid]/\n{\n\t@filename[str(args.{{.ArgCollected}})] = count();\n}\n")

	tmplToHandleAtForSymlink, _ = template.New("tmplToHandleAtForSymlink").Parse("tracepoint:syscalls:sys_enter_{{.SyscallName}}\n/@oldname[tid]/\n{\n\t@name[tid] = args.name;\n}\ntracepoint:syscalls:sys_exit_name_to_handle_at\n/@name[tid]/\n{\n\t$ret = args.ret;\n\t@fd[tid] = $ret >= 0 ? $ret : -1;\n}\ntracepoint:syscalls:sys_enter_open_by_handle_at\n/@fd[tid]/\n{\n\t@filename[str(@name[tid])] = count();\n\tdelete(@fd[tid]);\n}\n")
)

type SimpleCommand struct {
	SyscallName  string
	ArgCollected string
}

func GenerateBpfScript(commands []string) error {
	path := "./script.bt"
	err := createFile(path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("BEGIN" + immutablePieces["BEGIN"])
	if err != nil {
		return err
	}
	if len(commands) == 0 {
		err = makeDefaultScript(file)
		if err != nil {
			return err
		}
	} else {
		err = makeSpecificScript(file, commands)
		if err != nil {
			return err
		}
	}
	_, err = file.WriteString("END" + immutablePieces["END"])
	if err != nil {
		return err
	}

	//fmt.Println("==> file successful")
	return err
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
	//fmt.Println("==> file created successfully", path)
	return nil
}

func (s Syscall) String() string {
	return [...]string{"start", "execve", "execveat", "open", "openat", "openat2", "open_tree", "statx", "fspick", "open_by_handle_at", "readlink", "readlinkat", "end"}[s-1]
}

func isValid(commad string) bool {
	for nameCom := Start + 1; nameCom < End; nameCom++ {
		tempStr := nameCom.String()
		switch tempStr == commad {
		case true:
			return true
		}
	}
	return false
}

func makeDefaultScript(file *os.File) error {
	for nameCom := Start + 1; nameCom < End; nameCom++ {
		mainPieces := SimpleCommand{nameCom.String(), "filename"}
		fspickPiecess := SimpleCommand{nameCom.String(), "path"}
		tempStr := nameCom.String()
		switch tempStr {
		case "fspick":
			err := tmplSimplCom.Execute(file, fspickPiecess)
			if err != nil {
				return err
			}
		case "open_by_handle_at":
			mainPieces := SimpleCommand{"name_to_handle_at", "filename"}
			err := tmplToHandleAt.Execute(file, mainPieces)
			if err != nil {
				return err
			}
		case "readlink", "readlinkat":
			if tempStr == "readlink" {
				symlinkPiecess := SimpleCommand{nameCom.String(), "path"}
				err := tmplToSymlink.Execute(file, symlinkPiecess)
				if err != nil {
					return err
				}
			}
			if tempStr == "readlinkat" {
				symlinkPiecess := SimpleCommand{nameCom.String(), "pathname"}
				err := tmplToSymlink.Execute(file, symlinkPiecess)
				if err != nil {
					return err
				}
			}
			for symlink := Start + 1; symlink < End; symlink++ {
				mainPieces = SimpleCommand{symlink.String(), "filename"}
				fspickPiecess = SimpleCommand{symlink.String(), "path"}
				tempStrSymlink := symlink.String()
				switch tempStrSymlink {
				case "fspick":
					err := tmplSimplComForSymlink.Execute(file, fspickPiecess)
					if err != nil {
						return err
					}
				case "open_by_handle_at":
					mainPieces := SimpleCommand{"name_to_handle_at", "filename"}
					err := tmplToHandleAtForSymlink.Execute(file, mainPieces)
					if err != nil {
						return err
					}
				case "readlink", "readlinkat":
				default:
					err := tmplSimplComForSymlink.Execute(file, mainPieces)
					if err != nil {
						return err
					}
				}
			}
		default:
			err := tmplSimplCom.Execute(file, mainPieces)
			if err != nil {
				return err
			}
		}
	}

	//fmt.Println("==> makeDefaultScript")
	return nil
}

func makeSpecificScript(file *os.File, commands []string) error {

	for _, com := range commands {
		if isValid(com) {
			mainPieces := SimpleCommand{com, "filename"}
			fspickPiecess := SimpleCommand{com, "path"}
			switch com {
			case "fspick":
				err := tmplSimplCom.Execute(file, fspickPiecess)
				if err != nil {
					return err
				}
			case "open_by_handle_at":
				mainPieces = SimpleCommand{"name_to_handle_at", "filename"}
				err := tmplToHandleAt.Execute(file, mainPieces)
				if err != nil {
					return err
				}
			case "readlink", "readlinkat":
				if com == "readlink" {
					symlinkPiecess := SimpleCommand{com, "path"}
					err := tmplToSymlink.Execute(file, symlinkPiecess)
					if err != nil {
						return err
					}
				}
				if com == "readlinkat" {
					symlinkPiecess := SimpleCommand{com, "pathname"}
					err := tmplToSymlink.Execute(file, symlinkPiecess)
					if err != nil {
						return err
					}
				}
				for symlink := Start + 1; symlink < End; symlink++ {
					mainPieces = SimpleCommand{symlink.String(), "filename"}
					fspickPiecess = SimpleCommand{symlink.String(), "path"}
					tempStrSymlink := symlink.String()
					switch tempStrSymlink {
					case "fspick":
						err := tmplSimplComForSymlink.Execute(file, fspickPiecess)
						if err != nil {
							return err
						}
					case "open_by_handle_at":
						mainPieces = SimpleCommand{"name_to_handle_at", "filename"}
						err := tmplToHandleAtForSymlink.Execute(file, mainPieces)
						if err != nil {
							return err
						}
					case "readlink", "readlinkat":
					default:
						err := tmplSimplComForSymlink.Execute(file, mainPieces)
						if err != nil {
							return err
						}
					}
				}
			default:
				tmplSimplCom.Execute(file, mainPieces)
			}
		} else {
			err := errors.New("The system call is not valid.")
			return err
		}
	}

	//fmt.Println("==> makeSpecificScript")
	return nil
}
