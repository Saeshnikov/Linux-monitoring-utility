package namedPipeScript

import (
	"errors"
	"os"
	"strconv"
	"text/template"
)

var (
	immutablePieces = map[string]string{
		"partFullPath": "",

		"pipeHeader":    "#ifndef BPFTRACE_HAVE_BTF\n#include <linux/sched.h>\n#endif\n\n",
		"pipeStart":     "tracepoint:syscalls:sys_enter_mknod,\ntracepoint:syscalls:sys_enter_mknodat\n{\n\tif ((args.mode & 0170000) == 0010000) {\n\t\t@pipename[tid] = args.filename;\n\t}\n}\n",
		"pipeOpenStart": "\n/@pipename[tid]/\n{\n\t$ret = args.ret;\n\t$fd = $ret >= 0 ? $ret : -1;\n\t$errno = $ret >= 0 ? 0 : - $ret;\n\n\t",
		"pipeOpenEnd":   "printf(" + `"%` + `d %` + `s"` + ", $fd, str(@pipename[tid]));\n}\n\n",
		"pipeEnd":       "END\n{\n\tclear(@pipename);\n\tclear(@full_path_comm);\n}\n",
	}

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

func (s open_Syscall) String() string {
	return [...]string{"start", "execve", "execveat", "open", "openat", "openat2", "open_tree", "statx", "fspick", "open_by_handle_at", "readlink", "readlinkat", "end"}[s-1]
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

func MakeNamedPipeScript(file *os.File, option map[string][]string, rootInode int) error {
	immutablePieces["partFullPath"] = "$task = (struct task_struct *)curtask;\n\t$part_path = $task->mm->exe_file->f_path.dentry->d_parent;\n\t$i = 0;\n\t@full_path_comm[$i] = $part_path->d_name.name;\n\t$i = 1;\n\twhile ($i != 3000) {\n\t\t$part_path = $part_path->d_parent;\n\t\t@full_path_comm[$i] = $part_path->d_name.name;\n\t\tif ((uint64)$part_path->d_inode->i_ino == " + strconv.Itoa(rootInode) + ") {\n\t\t\tbreak;\n\t\t}\n\t\t$i = $i + 1;\n\t}\n\tprintf(" + `"\n/"` + ");\n\twhile ($i != -1) {\n\t\t$str_ = @full_path_comm[$i];\n\t\tprintf(" + `"%` + `s/"` + ", str($str_));\n\t\t$i = $i - 1;\n\t}\n\tprintf(" + `"%` + `s "` + ",comm);\n\n\t"

	file.WriteString(immutablePieces["pipeHeader"])
	file.WriteString(immutablePieces["pipeStart"])

	if len(option) == 0 {
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
	} else {
		for typeOpt, opt := range option {
			switch typeOpt {
			case "openSyscall":
				for _, value := range opt {
					if isValidOpenCom(value) {
						switch value {
						case Readlink.String(), Readlinkat.String():
							err := errors.New("SyscallOpen for pipe is not valid.")
							return err
						default:
							openCommPiecess := openCommand{value, "", ""}
							err := tmplOpenSyscallExit.Execute(file, openCommPiecess)
							if err != nil {
								return err
							}
						}
					} else {
						err := errors.New("SyscallOpen for pipe is not valid.")
						return err
					}
				}
			default:
				err := errors.New("Type option for pipe is not valid.")
				return err
			}
		}
	}

	file.WriteString(immutablePieces["pipeOpenStart"])
	file.WriteString(immutablePieces["partFullPath"])
	file.WriteString(immutablePieces["pipeOpenEnd"])
	file.WriteString(immutablePieces["pipeEnd"])

	return nil
}
