package fsorwScript

import (
	"errors"
	"os"
	"strconv"
	"text/template"
)

var (
	immutablePieces = map[string]string{
		"partFullPath": "",

		"fsorwHeader":    "#ifndef BPFTRACE_HAVE_BTF\n#include <linux/sched.h>\n#endif\n\n",
		"filter":         "\n/@oldname[tid]/",
		"openToHandleAt": "\ntracepoint:syscalls:sys_exit_name_to_handle_at\n/@name[tid]/\n{\n\t$ret = args.ret;\n\t@fdHandle[tid] = $ret >= 0 ? $ret : -1;\n}\ntracepoint:syscalls:sys_enter_open_by_handle_at\n/@fdHandle[tid]/\n{\n\t@filename[tid] = @name[tid];\n\tdelete(@fdHandle[tid]);\n}\n",
		"openEnd":        "\n/@filename[tid]/\n{\n\t$ret = args.ret;\n\t@fd[tid] = $ret >= 0 ? $ret : -1;\n}\n",
	}

	tmplOpenSyscallEnter, _  = template.New("OpenSyscallExit").Parse("\ntracepoint:syscalls:sys_enter_{{.SyscallName}}")
	tmplOpenSyscallBody, _   = template.New("OpenSyscallBody").Parse("\n{\n\t@{{.ReceivingArray}}[tid] = args.{{.CollectingArgs}};\n}\n")
	tmplFsorwCommandStart, _ = template.New("FsorwCommandStart").Parse("\ntracepoint:syscalls:sys_exit_{{.SyscallName}},\ntracepoint:syscalls:sys_exit_{{.SyscallName}}v,\ntracepoint:syscalls:sys_exit_p{{.SyscallName}}v\n/@fd[tid]/\n{\n\t$ret = args.ret;\n\t$nbyte = $ret >= 0 ? $ret : -1;\n\t$nothing = " + `0` + ";\n\n\t")
	tmplFsorwCommandEnd, _   = template.New("FsorwCommandEnd").Parse("printf(" + `"%` + "d %" + "s %" + "d %" + `d"` + ", @fd[tid], str(@filename[tid]), ${{.ReadByte}}, ${{.WriteByte}});\n\n\tdelete(@filename[tid]);\n\tdelete(@fd[tid]);\n}\n")
	tmplFsorwEnd, _          = template.New("FsorwEnd").Parse("\nEND\n{\n\t{{ range $index, $element := .}}{{ if $index }} {{ end }}{{$element}}{{ end }}\n}\n")

	tmplOpenSyscallExit, _ = template.New("OpenSyscallExit").Parse("\ntracepoint:syscalls:sys_exit_{{.SyscallName}},")
)

type open_Syscall int

const (
	StartOpen open_Syscall = iota + 1
	Execve
	Execveat
	Open
	Openat
	Openat2
	Open_tree
	Statx
	Fspick
	Open_by_handle_at
	Readlink
	Readlinkat
	EndOpen
)

type openCommand struct {
	SyscallName    string
	ReceivingArray string
	CollectingArgs string
}

type fsorwCommand struct {
	SyscallName string
	ReadByte    string
	WriteByte   string
}

type fsorw_Syscall int

const (
	Read fsorw_Syscall = iota + 1
	Write
)

func (s open_Syscall) String() string {
	return [...]string{"start", "execve", "execveat", "open", "openat", "openat2", "open_tree", "statx", "fspick", "open_by_handle_at", "readlink", "readlinkat", "end"}[s-1]
}

func (syscall fsorw_Syscall) String() string {
	return [...]string{"read", "write"}[syscall-1]
}

func isValidOpenCom(commad string) bool {
	for nameCom := StartOpen + 1; nameCom < EndOpen; nameCom++ {
		tempStr := nameCom.String()
		switch tempStr == commad {
		case true:
			return true
		}
	}
	return false
}

func MakeFsorwScript(file *os.File, option map[string][]string, rootInode int) error {
	immutablePieces["partFullPath"] = "$task = (struct task_struct *)curtask;\n\t$part_path = $task->mm->exe_file->f_path.dentry->d_parent;\n\t$i = 0;\n\t@full_path_comm[$i] = $part_path->d_name.name;\n\t$i = 1;\n\twhile ($i != 3000) {\n\t\t$part_path = $part_path->d_parent;\n\t\t@full_path_comm[$i] = $part_path->d_name.name;\n\t\tif ((uint64)$part_path->d_inode->i_ino == " + strconv.Itoa(rootInode) + ") {\n\t\t\tbreak;\n\t\t}\n\t\t$i = $i + 1;\n\t}\n\tprintf(" + `"\n/"` + ");\n\twhile ($i != -1) {\n\t\t$str_ = @full_path_comm[$i];\n\t\tprintf(" + `"%` + `s/"` + ", str($str_));\n\t\t$i = $i - 1;\n\t}\n\tprintf(" + `"%` + `s "` + ",comm);\n\n\t"

	endPiecess := make([]string, 0, 6)
	immutableEndPiecess := []string{"clear(@filename);", "clear(@fd);", "clear(@full_path_comm);"}
	openToHandleAtEndPiecess := []string{"clear(@name);", "clear(@fdHandle);"}
	symlinkEndPiecess := "clear(@oldname);"

	file.WriteString(immutablePieces["fsorwHeader"])

	if len(option) == 0 {
		arrFsorwSyscall := []string{"read", "write"}
		for nameCom := StartOpen + 1; nameCom < EndOpen; nameCom++ {
			tempStr := nameCom.String()
			switch tempStr {
			case Readlink.String(), Readlinkat.String():
				openCommPiecess := openCommand{}
				if tempStr == Readlink.String() {
					openCommPiecess = openCommand{tempStr, "oldname", "path"}
				} else {
					openCommPiecess = openCommand{tempStr, "oldname", "pathname"}
				}
				err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
				if err1 != nil {
					return err1
				}
				err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
				if err2 != nil {
					return err2
				}
				for nameComLink := StartOpen + 1; nameComLink < EndOpen; nameComLink++ {
					tempStrLink := nameComLink.String()
					switch tempStrLink {
					case Readlink.String(), Readlinkat.String():
					case Fspick.String():
						openCommPiecess = openCommand{tempStrLink, "filename", "path"}
						err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
						if err1 != nil {
							return err1
						}
						file.WriteString(immutablePieces["filter"])
						err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
						if err2 != nil {
							return err2
						}
					case Open_by_handle_at.String():
						openCommPiecess = openCommand{"name_to_handle_at", "name", "name"}
						err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
						if err1 != nil {
							return err1
						}
						file.WriteString(immutablePieces["filter"])
						err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
						if err2 != nil {
							return err2
						}
						file.WriteString(immutablePieces["openToHandleAt"])
					default:
						openCommPiecess = openCommand{tempStrLink, "filename", "filename"}
						err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
						if err1 != nil {
							return err1
						}
						file.WriteString(immutablePieces["filter"])
						err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
						if err2 != nil {
							return err2
						}
					}
				}
			case Fspick.String():
				openCommPiecess := openCommand{tempStr, "filename", "path"}
				err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
				if err1 != nil {
					return err1
				}
				err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
				if err2 != nil {
					return err2
				}
			case Open_by_handle_at.String():
				openCommPiecess := openCommand{"name_to_handle_at", "name", "name"}
				err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
				if err1 != nil {
					return err1
				}
				err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
				if err2 != nil {
					return err2
				}
				file.WriteString(immutablePieces["openToHandleAt"])
			default:
				openCommPiecess := openCommand{tempStr, "filename", "filename"}
				err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
				if err1 != nil {
					return err1
				}
				err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
				if err2 != nil {
					return err2
				}
			}
		}
		for nameCom := StartOpen + 1; nameCom < EndOpen; nameCom++ {
			tempStr := nameCom.String()
			switch tempStr {
			case Readlink.String(), Readlinkat.String():
			default:
				openCommPiecess := openCommand{tempStr, "", ""}
				err := tmplOpenSyscallExit.Execute(file, openCommPiecess)
				if err != nil {
					return err
				}
			}
		}
		file.WriteString(immutablePieces["openEnd"])

		for _, value := range arrFsorwSyscall {
			fsorwCommandPiecess := fsorwCommand{}
			if value == Read.String() {
				fsorwCommandPiecess = fsorwCommand{value, "ret", "nothing"}
			} else {
				fsorwCommandPiecess = fsorwCommand{value, "nothing", "ret"}
			}
			err1 := tmplFsorwCommandStart.Execute(file, fsorwCommandPiecess)
			if err1 != nil {
				return err1
			}
			file.WriteString(immutablePieces["partFullPath"])
			err2 := tmplFsorwCommandEnd.Execute(file, fsorwCommandPiecess)
			if err2 != nil {
				return err2
			}
		}

		endPiecess = append(endPiecess, immutableEndPiecess...)
		endPiecess = append(endPiecess, openToHandleAtEndPiecess...)
		endPiecess = append(endPiecess, symlinkEndPiecess)
		err := tmplFsorwEnd.Execute(file, endPiecess)
		if err != nil {
			return err
		}

	} else {
		var (
			isFsorwSyscall, isOpenSyscall   bool
			arrFsorwSyscall, arrOpenSyscall []string
		)
		for typeOpt, opt := range option {
			switch typeOpt {
			case "fsorwSyscall":
				isFsorwSyscall = true
				for _, value := range opt {
					arrFsorwSyscall = append(arrFsorwSyscall, value)
				}
			case "openSyscall":
				isOpenSyscall = true
				for _, value := range opt {
					arrOpenSyscall = append(arrOpenSyscall, value)
				}
			default:
				err := errors.New("Type option fsorw is not valid.")
				return err
			}
		}
		if !isOpenSyscall {
			for nameCom := StartOpen + 1; nameCom < EndOpen; nameCom++ {
				tempStr := nameCom.String()
				switch tempStr {
				case Readlink.String(), Readlinkat.String():
					openCommPiecess := openCommand{}
					if tempStr == Readlink.String() {
						openCommPiecess = openCommand{tempStr, "oldname", "path"}
					} else {
						openCommPiecess = openCommand{tempStr, "oldname", "pathname"}
					}
					err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
					if err1 != nil {
						return err1
					}
					err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
					if err2 != nil {
						return err2
					}
					for nameComLink := StartOpen + 1; nameComLink < EndOpen; nameComLink++ {
						tempStrLink := nameComLink.String()
						switch tempStrLink {
						case Readlink.String(), Readlinkat.String():
						case Fspick.String():
							openCommPiecess = openCommand{tempStrLink, "filename", "path"}
							err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
							if err1 != nil {
								return err1
							}
							file.WriteString(immutablePieces["filter"])
							err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
							if err2 != nil {
								return err2
							}
						case Open_by_handle_at.String():
							openCommPiecess = openCommand{"name_to_handle_at", "name", "name"}
							err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
							if err1 != nil {
								return err1
							}
							file.WriteString(immutablePieces["filter"])
							err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
							if err2 != nil {
								return err2
							}
							file.WriteString(immutablePieces["openToHandleAt"])
						default:
							openCommPiecess = openCommand{tempStrLink, "filename", "filename"}
							err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
							if err1 != nil {
								return err1
							}
							file.WriteString(immutablePieces["filter"])
							err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
							if err2 != nil {
								return err2
							}
						}
					}
				case Fspick.String():
					openCommPiecess := openCommand{tempStr, "filename", "path"}
					err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
					if err1 != nil {
						return err1
					}
					err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
					if err2 != nil {
						return err2
					}
				case Open_by_handle_at.String():
					openCommPiecess := openCommand{"name_to_handle_at", "name", "name"}
					err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
					if err1 != nil {
						return err1
					}
					err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
					if err2 != nil {
						return err2
					}
					file.WriteString(immutablePieces["openToHandleAt"])
				default:
					openCommPiecess := openCommand{tempStr, "filename", "filename"}
					err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
					if err1 != nil {
						return err1
					}
					err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
					if err2 != nil {
						return err2
					}
				}
			}
			for nameCom := StartOpen + 1; nameCom < EndOpen; nameCom++ {
				tempStr := nameCom.String()
				switch tempStr {
				case Readlink.String(), Readlinkat.String():
				default:
					openCommPiecess := openCommand{tempStr, "", ""}
					err := tmplOpenSyscallExit.Execute(file, openCommPiecess)
					if err != nil {
						return err
					}
				}
			}
			file.WriteString(immutablePieces["openEnd"])

			for _, value := range arrFsorwSyscall {
				fsorwCommandPiecess := fsorwCommand{}
				switch value {
				case Read.String():
					fsorwCommandPiecess = fsorwCommand{value, "ret", "nothing"}
				case Write.String():
					fsorwCommandPiecess = fsorwCommand{value, "nothing", "ret"}
				default:
					err := errors.New("Type fsorw syscall is not valid.")
					return err
				}
				err1 := tmplFsorwCommandStart.Execute(file, fsorwCommandPiecess)
				if err1 != nil {
					return err1
				}
				file.WriteString(immutablePieces["partFullPath"])
				err2 := tmplFsorwCommandEnd.Execute(file, fsorwCommandPiecess)
				if err2 != nil {
					return err2
				}
			}
			endPiecess = append(endPiecess, immutableEndPiecess...)
			endPiecess = append(endPiecess, openToHandleAtEndPiecess...)
			endPiecess = append(endPiecess, symlinkEndPiecess)
			err := tmplFsorwEnd.Execute(file, endPiecess)
			if err != nil {
				return err
			}
		} else {
			for _, value := range arrOpenSyscall {
				if isValidOpenCom(value) {
					switch value {
					case Readlink.String(), Readlinkat.String():
						openCommPiecess := openCommand{}
						if value == Readlink.String() {
							openCommPiecess = openCommand{value, "oldname", "path"}
						} else {
							openCommPiecess = openCommand{value, "oldname", "pathname"}
						}
						err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
						if err1 != nil {
							return err1
						}
						err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
						if err2 != nil {
							return err2
						}
						for _, nameComLink := range arrOpenSyscall {
							switch nameComLink {
							case Readlink.String(), Readlinkat.String():
							case Fspick.String():
								openCommPiecess = openCommand{nameComLink, "filename", "path"}
								err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
								if err1 != nil {
									return err1
								}
								file.WriteString(immutablePieces["filter"])
								err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
								if err2 != nil {
									return err2
								}
							case Open_by_handle_at.String():
								openCommPiecess = openCommand{"name_to_handle_at", "name", "name"}
								err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
								if err1 != nil {
									return err1
								}
								file.WriteString(immutablePieces["filter"])
								err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
								if err2 != nil {
									return err2
								}
								file.WriteString(immutablePieces["openToHandleAt"])
							default:
								openCommPiecess = openCommand{nameComLink, "filename", "filename"}
								err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
								if err1 != nil {
									return err1
								}
								file.WriteString(immutablePieces["filter"])
								err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
								if err2 != nil {
									return err2
								}
							}
						}
						endPiecess = append(endPiecess, symlinkEndPiecess)
					case Fspick.String():
						openCommPiecess := openCommand{value, "filename", "path"}
						err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
						if err1 != nil {
							return err1
						}
						err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
						if err2 != nil {
							return err2
						}
					case Open_by_handle_at.String():
						openCommPiecess := openCommand{"name_to_handle_at", "name", "name"}
						err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
						if err1 != nil {
							return err1
						}
						err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
						if err2 != nil {
							return err2
						}
						file.WriteString(immutablePieces["openToHandleAt"])
						endPiecess = append(endPiecess, openToHandleAtEndPiecess...)
					default:
						openCommPiecess := openCommand{value, "filename", "filename"}
						err1 := tmplOpenSyscallEnter.Execute(file, openCommPiecess)
						if err1 != nil {
							return err1
						}
						err2 := tmplOpenSyscallBody.Execute(file, openCommPiecess)
						if err2 != nil {
							return err2
						}
					}
				} else {
					err := errors.New("SyscallOpen for fsorw is not valid.")
					return err
				}
			}
			if !isFsorwSyscall {
				arrFsorwSyscall = []string{"read", "write"}
				for _, value := range arrFsorwSyscall {
					fsorwCommandPiecess := fsorwCommand{}
					if value == Read.String() {
						fsorwCommandPiecess = fsorwCommand{value, "ret", "nothing"}
					} else {
						fsorwCommandPiecess = fsorwCommand{value, "nothing", "ret"}
					}
					err1 := tmplFsorwCommandStart.Execute(file, fsorwCommandPiecess)
					if err1 != nil {
						return err1
					}
					file.WriteString(immutablePieces["partFullPath"])
					err2 := tmplFsorwCommandEnd.Execute(file, fsorwCommandPiecess)
					if err2 != nil {
						return err2
					}
				}
			} else {
				for _, value := range arrFsorwSyscall {
					fsorwCommandPiecess := fsorwCommand{}
					switch value {
					case Read.String():
						fsorwCommandPiecess = fsorwCommand{value, "ret", "nothing"}
					case Write.String():
						fsorwCommandPiecess = fsorwCommand{value, "nothing", "ret"}
					default:
						err := errors.New("Type fsorw syscall is not valid.")
						return err
					}
					err1 := tmplFsorwCommandStart.Execute(file, fsorwCommandPiecess)
					if err1 != nil {
						return err1
					}
					file.WriteString(immutablePieces["partFullPath"])
					err2 := tmplFsorwCommandEnd.Execute(file, fsorwCommandPiecess)
					if err2 != nil {
						return err2
					}
				}
			}
			err := tmplFsorwEnd.Execute(file, endPiecess)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
