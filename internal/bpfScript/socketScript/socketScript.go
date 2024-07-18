package socketScript

import (
	genStruct "linux-monitoring-utility/internal/bpfScript/generalStructIPC"
	"errors"
	"os"
	"strconv"
	"text/template"
)

var (
	immutablePieces = map[string]string{
		"partFullPath": "",

		"sockHeader":       "#ifndef BPFTRACE_HAVE_BTF\n#include <linux/sched.h>\n#include <linux/socket.h>\n#include <net/sock.h>\n#else\n#include <sys/socket.h>\n#endif\n",
		"sockConnectStart": "\ntracepoint:syscalls:sys_enter_connect\n{\n\t$sk = ((struct sockaddr *) args->uservaddr);\n\t@inet_family1[tid] = $sk->sa_family;\n\t@fd1[tid] = args->fd;\n}\n\ntracepoint:syscalls:sys_exit_connect\n/@fd1[tid]/\n{\n\t",
		"sockConnectEnd":   "\t}\n\tdelete(@fd1[tid]);\n\tdelete(@inet_family1[tid]);\n}\n",
		"sockIfConnect":    "$ret = args->ret;\n\tif ($ret == 0) {\n",
		"sockAcceptStart":  "\ntracepoint:syscalls:sys_enter_bind\n{\n\t$sk = ((struct sockaddr *) args->umyaddr);\n\t@inet_family2[tid] = $sk->sa_family;\n\t@fd2[tid] = args->fd;\n}\n\ntracepoint:syscalls:sys_exit_accept\n/@fd2[tid]/\n{\n\t",
		"sockAcceptEnd":    "\t}\n\tdelete(@fd2[tid]);\n\tdelete(@inet_family2[tid]);\n}\n",
		"sockIfAccept":     "$ret = args->ret;\n\tif ($ret > 0) {\n",
	}

	tmplSockProtocol, _ = template.New("SockProtocol").Parse("\t\tif (@inet_family{{.Number}}[tid] == AF_{{.Protocol}}) {\n\t\t\tprintf(" + `"%` + `s %` + `s %` + `d"` + ", " + `"{{.TypeSyscall}}"` + ", " + `"{{.Protocol}}"` + ", @fd{{.Number}}[tid]);\n\t\t}\n")
	tmplSockEnd, _      = template.New("SockEnd").Parse("\nEND\n{\n\t{{ range $index, $element := .}}{{ if $index }}; {{ end }}{{$element}}{{ end }}\n}\n")
)

type socketType struct {
	TypeSyscall string
	Protocol    string
	Number      string
}

type socket_Syscall int
type socket_Protocol int

const (
	Connect socket_Syscall = iota + 1
	Accept

	Unix socket_Protocol = iota
	Inet
	Inet6
)

func (syscall socket_Syscall) String() string {
	return [...]string{"connect", "accept"}[syscall-1]
}

func (protocol socket_Protocol) String() string {
	return [...]string{"UNIX", "INET", "INET6"}[protocol-2]
}

func MakeSocketScript(file *os.File, option []genStruct.OptionStruct, rootInode int) error {
	immutablePieces["partFullPath"] = "$task = (struct task_struct *)curtask;\n\t$part_path = $task->mm->exe_file->f_path.dentry->d_parent;\n\t$i = 0;\n\t@full_path_comm[$i] = $part_path->d_name.name;\n\t$i = 1;\n\twhile ($i != 3000) {\n\t\t$part_path = $part_path->d_parent;\n\t\t@full_path_comm[$i] = $part_path->d_name.name;\n\t\tif ((uint64)$part_path->d_inode->i_ino == " + strconv.Itoa(rootInode) + ") {\n\t\t\tbreak;\n\t\t}\n\t\t$i = $i + 1;\n\t}\n\tprintf(" + `"\n/"` + ");\n\twhile ($i != -1) {\n\t\t$str_ = @full_path_comm[$i];\n\t\tprintf(" + `"%` + `s/"` + ", str($str_));\n\t\t$i = $i - 1;\n\t}\n\tprintf(" + `"%` + `s "` + ",comm);\n\n\t"

	var (
		isTypeSyscall, isProtocol bool
		arrTypeSyscall            []string
		arrProtocol               []string
		clear_list                []string
		protocol_list             []string
	)

	file.WriteString(immutablePieces["sockHeader"])
	if len(option) == 0 {
		protocol_list = []string{Unix.String(), Inet.String(), Inet6.String()}

		file.WriteString(immutablePieces["sockConnectStart"])
		file.WriteString(immutablePieces["partFullPath"])
		file.WriteString(immutablePieces["sockIfConnect"])
		for _, value := range protocol_list {
			tmpl := socketType{"C", value, "1"}
			err := tmplSockProtocol.Execute(file, tmpl)
			if err != nil {
				return err
			}
		}
		file.WriteString(immutablePieces["sockConnectEnd"])

		file.WriteString(immutablePieces["sockAcceptStart"])
		file.WriteString(immutablePieces["partFullPath"])
		file.WriteString(immutablePieces["sockIfAccept"])
		for _, value := range protocol_list {
			tmpl := socketType{"A", value, "2"}
			err := tmplSockProtocol.Execute(file, tmpl)
			if err != nil {
				return err
			}
		}
		file.WriteString(immutablePieces["sockAcceptEnd"])

		clear_list = []string{"clear(@fd1)", "clear(@fd2)", "clear(@inet_family1)", "clear(@inet_family2)", "clear(@full_path_comm);"}
		err := tmplSockEnd.Execute(file, clear_list)
		if err != nil {
			return err
		}

	} else {
		for _, opt := range option {
			switch opt.OptionType {
			case "sockSyscall":
				isTypeSyscall = true
				for _, value := range opt.Options {
					arrTypeSyscall = append(arrTypeSyscall, value)
				}
			case "protocol":
				isProtocol = true
				for _, value := range opt.Options {
					arrProtocol = append(arrProtocol, value)
				}
			default:
				err := errors.New("Type option socket is not valid.")
				return err
			}
		}
		if !isTypeSyscall {
			file.WriteString(immutablePieces["sockConnectStart"])
			file.WriteString(immutablePieces["partFullPath"])
			file.WriteString(immutablePieces["sockIfConnect"])
			for _, value := range arrProtocol {
				switch value {
				case Unix.String():
					tmpl := socketType{"C", value, "1"}
					err := tmplSockProtocol.Execute(file, tmpl)
					if err != nil {
						return err
					}
				case Inet.String():
					tmpl := socketType{"C", value, "1"}
					err := tmplSockProtocol.Execute(file, tmpl)
					if err != nil {
						return err
					}
				case Inet6.String():
					tmpl := socketType{"C", value, "1"}
					err := tmplSockProtocol.Execute(file, tmpl)
					if err != nil {
						return err
					}
				default:
					err := errors.New("Type protocol socket is not valid.")
					return err
				}
			}
			file.WriteString(immutablePieces["sockConnectEnd"])

			file.WriteString(immutablePieces["sockAcceptStart"])
			file.WriteString(immutablePieces["partFullPath"])
			file.WriteString(immutablePieces["sockIfAccept"])
			for _, value := range arrProtocol {
				switch value {
				case Unix.String():
					tmpl := socketType{"A", value, "2"}
					err := tmplSockProtocol.Execute(file, tmpl)
					if err != nil {
						return err
					}
				case Inet.String():
					tmpl := socketType{"A", value, "2"}
					err := tmplSockProtocol.Execute(file, tmpl)
					if err != nil {
						return err
					}
				case Inet6.String():
					tmpl := socketType{"A", value, "2"}
					err := tmplSockProtocol.Execute(file, tmpl)
					if err != nil {
						return err
					}
				default:
					err := errors.New("Type protocol socket is not valid.")
					return err
				}
			}
			file.WriteString(immutablePieces["sockAcceptEnd"])

			clear_list = []string{"clear(@fd1)", "clear(@fd2)", "clear(@inet_family1)", "clear(@inet_family2)", "clear(@full_path_comm);"}
			err := tmplSockEnd.Execute(file, clear_list)
			if err != nil {
				return err
			}
		} else {
			for _, value := range arrTypeSyscall {
				switch value {
				case Connect.String():
					file.WriteString(immutablePieces["sockConnectStart"])
					file.WriteString(immutablePieces["partFullPath"])
					file.WriteString(immutablePieces["sockIfConnect"])
					if !isProtocol {
						protocol_list = []string{Unix.String(), Inet.String(), Inet6.String()}
						for _, value := range protocol_list {
							tmpl := socketType{"C", value, "1"}
							err := tmplSockProtocol.Execute(file, tmpl)
							if err != nil {
								return err
							}
						}
					} else {
						for _, value := range arrProtocol {
							switch value {
							case Unix.String():
								tmpl := socketType{"C", value, "1"}
								err := tmplSockProtocol.Execute(file, tmpl)
								if err != nil {
									return err
								}
							case Inet.String():
								tmpl := socketType{"C", value, "1"}
								err := tmplSockProtocol.Execute(file, tmpl)
								if err != nil {
									return err
								}
							case Inet6.String():
								tmpl := socketType{"C", value, "1"}
								err := tmplSockProtocol.Execute(file, tmpl)
								if err != nil {
									return err
								}
							default:
								err := errors.New("Type protocol socket is not valid.")
								return err
							}
						}
					}
					file.WriteString(immutablePieces["sockConnectEnd"])
					clear_list = append(clear_list, "clear(@fd1)")
					clear_list = append(clear_list, "clear(@inet_family1)")

				case Accept.String():
					file.WriteString(immutablePieces["sockAcceptStart"])
					file.WriteString(immutablePieces["partFullPath"])
					file.WriteString(immutablePieces["sockIfAccept"])
					if !isProtocol {
						protocol_list = []string{Unix.String(), Inet.String(), Inet6.String()}
						for _, value := range protocol_list {
							tmpl := socketType{"A", value, "2"}
							err := tmplSockProtocol.Execute(file, tmpl)
							if err != nil {
								return err
							}
						}
					} else {
						for _, value := range arrProtocol {
							switch value {
							case Unix.String():
								tmpl := socketType{"A", value, "2"}
								err := tmplSockProtocol.Execute(file, tmpl)
								if err != nil {
									return err
								}
							case Inet.String():
								tmpl := socketType{"A", value, "2"}
								err := tmplSockProtocol.Execute(file, tmpl)
								if err != nil {
									return err
								}
							case Inet6.String():
								tmpl := socketType{"A", value, "2"}
								err := tmplSockProtocol.Execute(file, tmpl)
								if err != nil {
									return err
								}
							default:
								err := errors.New("Type protocol socket is not valid.")
								return err
							}
						}
					}
					file.WriteString(immutablePieces["sockAcceptEnd"])
					clear_list = append(clear_list, "clear(@fd2)")
					clear_list = append(clear_list, "clear(@inet_family2)")

				default:
					err := errors.New("Type syscall socket is not valid.")
					return err
				}
			}
			clear_list = append(clear_list, "clear(@full_path_comm);")
			err := tmplSockEnd.Execute(file, clear_list)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
