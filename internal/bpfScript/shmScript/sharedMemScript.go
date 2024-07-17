package shmScript

import (
	genStruct "linux-monitoring-utility/internal/bpfScript/generalStructIPC"
	"errors"
	"os"
	"strconv"
)

var (
	immutablePieces = map[string]string{
		"partFullPath": "",

		"shmHeader":     "#ifndef BPFTRACE_HAVE_BTF\n#include <linux/sched.h>\n#include <linux/limits.h>\n#endif\n\n",
		"shmSystVStart": "tracepoint:syscalls:sys_enter_shmget\n{\n\t@shmkey[tid] = args->key;\n}\n\ntracepoint:syscalls:sys_exit_shmget\n/@shmkey[tid]/\n{\n\t@shmid[tid] = args->ret;\n}\n\ntracepoint:syscalls:sys_exit_shmat\n/@shmid[tid]/\n{\n\t",
		"shmSystVEnd":   "$type = " + `"systemV"` + ";\n\t" + "printf(" + `"%` + `x %` + `d %` + `s"` + ", @shmkey[tid], @shmid[tid], $type);\n}\n\n",
		"shmPosixStart": "tracepoint:syscalls:sys_enter_openat,\n{\n\t$ret = strcontains(str(args->filename), " + `"dev/shm/"` + ");\n\tif ($ret == 1) {\n\t\t@filename[tid] = args->filename;\n\t}\n}\n\ntracepoint:syscalls:sys_exit_openat,\n/@filename[tid]/\n{\n\tif (args->ret > 0) {\n\t",
		"shmPosixEnd":   "@posshmid[tid] = args->ret;\n\t@posshmid[tid] = args->ret;\n\tprintf(" + `"%` + `s "` + ", str(@filename[tid]));\n\t}\n}\n\ntracepoint:syscalls:sys_enter_mmap\n/@posshmid[tid]/\n{\n\t$type = " + `"posix"` + ";\n\tprintf(" + `"%` + `d %` + `s"` + ", @posshmid[tid], $type);\n}\n\n",
		"shmENDStrart":  "END\n{\n\t",
		"shmENDsysV":    "clear(@shmkey);\n\tclear(@shmid);\n\t",
		"shmENDposix":   "clear(@posshmid);\n\tclear(@full_path_comm);\n\tclear(@filename);\n\t",
		"shmENDend":     "\n}\n",
	}
)

type standard int

const (
	SystemV standard = iota + 1
	Posix
)

func (norm standard) String() string {
	return [...]string{"systemV", "posix"}[norm-1]
}

func MakeSharedMemScript(file *os.File, option []genStruct.OptionStruct, rootInode int) error {
	immutablePieces["partFullPath"] = "$task = (struct task_struct *)curtask;\n\t$part_path = $task->mm->exe_file->f_path.dentry->d_parent;\n\t$i = 0;\n\t@full_path_comm[$i] = $part_path->d_name.name;\n\t$i = 1;\n\twhile ($i != 3000) {\n\t\t$part_path = $part_path->d_parent;\n\t\t@full_path_comm[$i] = $part_path->d_name.name;\n\t\tif ((uint64)$part_path->d_inode->i_ino == " + strconv.Itoa(rootInode) + ") {\n\t\t\tbreak;\n\t\t}\n\t\t$i = $i + 1;\n\t}\n\tprintf(" + `"\n/"` + ");\n\twhile ($i != -1) {\n\t\t$str_ = @full_path_comm[$i];\n\t\tprintf(" + `"%` + `s/"` + ", str($str_));\n\t\t$i = $i - 1;\n\t}\n\tprintf(" + `"%` + `s "` + ",comm);\n\n\t"

		var (
		sysV  bool
		posix bool
	)

	if len(option) == 0 {
		file.WriteString(immutablePieces["shmHeader"])
		file.WriteString(immutablePieces["shmSystVStart"])
		file.WriteString(immutablePieces["partFullPath"])
		file.WriteString(immutablePieces["shmSystVEnd"])
		file.WriteString(immutablePieces["shmPosixStart"])
		file.WriteString(immutablePieces["partFullPath"])
		file.WriteString(immutablePieces["shmPosixEnd"])
		file.WriteString(immutablePieces["shmENDStrart"])
		file.WriteString(immutablePieces["shmENDsysV"])
		file.WriteString(immutablePieces["shmENDposix"])
		file.WriteString(immutablePieces["shmENDend"])
	} else {
		for _, opt := range option {
			switch opt.OptionType {
			case "standards":
				file.WriteString(immutablePieces["shmHeader"])
				for _, value := range opt.Options {
					switch value {
					case SystemV.String():
						sysV = true
						file.WriteString(immutablePieces["shmSystVStart"])
						file.WriteString(immutablePieces["partFullPath"])
						file.WriteString(immutablePieces["shmSystVEnd"])
					case Posix.String():
						posix = true
						file.WriteString(immutablePieces["shmPosixStart"])
						file.WriteString(immutablePieces["partFullPath"])
						file.WriteString(immutablePieces["shmPosixEnd"])
					default:
						err := errors.New("Standards for sharedMem is not valid.")
						return err
					}
				}
				file.WriteString(immutablePieces["shmENDStrart"])
				if sysV && posix {
					file.WriteString(immutablePieces["shmENDsysV"])
					file.WriteString(immutablePieces["shmENDposix"])
				} else if sysV {
					file.WriteString(immutablePieces["shmENDsysV"])
				} else if posix {
					file.WriteString(immutablePieces["shmENDposix"])
				}
				file.WriteString(immutablePieces["shmENDend"])
			default:
				err := errors.New("Type option for sharedMem is not valid.")
				return err
			}
		}
	}
	return nil
}
