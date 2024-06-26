package semScript

import (
	genStruct "linux-monitoring-utility/internal/bpfScript/generalStructIPC"
	"errors"
	"os"
	"strconv"
)

var (
	immutablePieces = map[string]string{
		"partFullPath": "",

		"semHeader": "#ifndef BPFTRACE_HAVE_BTF\n#include <linux/sched.h>\n#endif\n\n",
		"semStart":  "tracepoint:syscalls:sys_enter_semget\n{\n\t@semkey[tid] = args.key;\n}\n\ntracepoint:syscalls:sys_exit_semget\n/@semkey[tid]/\n{\n\t@semid[tid] = args.ret;\n}\n\ntracepoint:syscalls:sys_enter_semop,\ntracepoint:syscalls:sys_enter_semtimedop,\n/@semid[tid]/\n{\n\t",
		"semEnd":    "printf(" + `"%` + "x %" + `d",` + "@semkey[tid], @semid[tid]);\n}\n\nEND\n{\n\tclear(@semkey);\n\tclear(@semid);\n\tclear(@full_path_comm);\n}\n",
	}
)

func MakeSemaphoreScript(file *os.File, option []genStruct.OptionStruct, rootInode int) error {
	immutablePieces["partFullPath"] = "$task = (struct task_struct *)curtask;\n\t$part_path = $task->mm->exe_file->f_path.dentry->d_parent;\n\t$i = 0;\n\t@full_path_comm[$i] = $part_path->d_name.name;\n\t$i = 1;\n\twhile ($i != 3000) {\n\t\t$part_path = $part_path->d_parent;\n\t\t@full_path_comm[$i] = $part_path->d_name.name;\n\t\tif ((uint64)$part_path->d_inode->i_ino == " + strconv.Itoa(rootInode) + ") {\n\t\t\tbreak;\n\t\t}\n\t\t$i = $i + 1;\n\t}\n\tprintf(" + `"\n/"` + ");\n\twhile ($i != -1) {\n\t\t$str_ = @full_path_comm[$i];\n\t\tprintf(" + `"%` + `s/"` + ", str($str_));\n\t\t$i = $i - 1;\n\t}\n\tprintf(" + `"%` + `s "` + ",comm);\n\n\t"

	if len(option) == 0 {
		file.WriteString(immutablePieces["semHeader"])
		file.WriteString(immutablePieces["semStart"])
		file.WriteString(immutablePieces["partFullPath"])
		file.WriteString(immutablePieces["semEnd"])
	} else {
		err := errors.New("Options are not available for semaphore.")
		return err
	}
	return nil
}
