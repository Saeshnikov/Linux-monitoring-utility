package bpfScript

import (
	"bytes"
	"log"
	"os"
)

var (
	immutablePieces = map[string]string{
		"BEGIN":  "\n{\n\tprintf(" + `"Tracing file system syscalls... Hit Ctrl-C to end.\n"` + ");\n}\n",
		"END":    "\n{\n\tprint(@filename);\n\tclear(@oldname);\n\tclear(@filename);\n\tclear(@name);\n\tclear(@fd);\n}",
		"filter": "\n/@oldname[tid]/",
	}

	mainPieces = map[string]string{
		"tracepoint:syscalls:sys_enter_execve":            "\n{\n\t@filename[str(args.filename)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_execveat":          "\n{\n\t@filename[str(args.filename)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_open":              "\n{\n\t@filename[str(args.filename)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_openat":            "\n{\n\t@filename[str(args.filename)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_openat2":           "\n{\n\t@filename[str(args.filename)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_open_tree":         "\n{\n\t@filename[str(args.filename)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_statx":             "\n{\n\t@filename[str(args.filename)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_fspick":            "\n{\n\t@filename[str(args.path)] = count();\n}\n",
		"tracepoint:syscalls:sys_enter_name_to_handle_at": "\n{\n\t@name[tid] = args.name;\n}\ntracepoint:syscalls:sys_exit_name_to_handle_at\n/@name[tid]/\n{\n\t$ret = args.ret;\n\t@fd[tid] = $ret >= 0 ? $ret : -1;\n}\ntracepoint:syscalls:sys_enter_open_by_handle_at\n/@fd[tid]/\n{\n\t@filename[str(@name[tid])] = count();\n\tdelete(@fd[tid]);\n}\n",
	}

	symlinkPieces = map[string]string{
		"tracepoint:syscalls:sys_enter_readlink":   "\n{\n\t@oldname[tid] = args.path;\n}\n",
		"tracepoint:syscalls:sys_enter_readlinkat": "\n{\n\t@oldname[tid] = args.pathname;\n}\n",
	}

	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func GenerateBpfScript(commands []string) *os.File {
	path := "./script.bt"
	createFile(path)

	var file, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer file.Close()

	var script bytes.Buffer

	script.WriteString("BEGIN" + immutablePieces["BEGIN"])
	if len(commands) == 0 {
		script.WriteString(makeDefaultScript())
	} else {
		temp := makeSpecificScript(commands)
		if len(temp) == 0 {
			errorLog.Fatal(err)
		} else {
			script.WriteString(temp)
		}
	}
	script.WriteString("END" + immutablePieces["END"])

	_, err = file.WriteString(script.String())
	if err != nil {
		errorLog.Fatal(err)
	}

	//fmt.Println("==> file successful")
	return file
}

func createFile(path string) {
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			errorLog.Fatal(err)
		}
		defer file.Close()
	}
	//fmt.Println("==> file created successfully", path)
}

func makeDefaultScript() string {
	var scriptBody bytes.Buffer

	for nameCom, code := range mainPieces {
		scriptBody.WriteString(nameCom + code)
	}

	for symlinkCom, part := range symlinkPieces {
		scriptBody.WriteString(symlinkCom + part)
	}
	for nameCom, code := range mainPieces {
		scriptBody.WriteString(nameCom + immutablePieces["filter"] + code)
	}

	//fmt.Println("==> makeDefaultScript")
	return scriptBody.String()
}

func makeSpecificScript(commands []string) string {
	var scriptBody bytes.Buffer

	for _, com := range commands {
		key := "tracepoint:syscalls:sys_enter_" + com
		if nameCom, found := mainPieces[key]; found {
			scriptBody.WriteString(key + nameCom)
		}
		if key == "tracepoint:syscalls:sys_enter_open_by_handle_at" {
			tmp := "tracepoint:syscalls:sys_enter_name_to_handle_at"
			scriptBody.WriteString(tmp + mainPieces[tmp])
		}
		if symlinkCom, found := symlinkPieces[key]; found {
			scriptBody.WriteString(key + symlinkCom)
			for nameCom, code := range mainPieces {
				scriptBody.WriteString(nameCom + immutablePieces["filter"] + code)
			}
		}
	}

	//fmt.Println("==> makeSpecificScript")
	return scriptBody.String()
}
