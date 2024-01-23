# Файловая система (FS)
| Пункт мониторинга | Утилита для получения информации| Актуальность | Обоснование выбора утилиты |
| :---: | --- | --- | --- |
| Вывод информации о том, какие файлы используются теми или иными процессами | bpftrace-сценарий [fsorw.bt](#fsorw.bt) <br> В качестве ключевой будет выводиться следующая информация: <br>- pid процесса <br>- команда <br>- файловые дескрипторы <br>- флаги режима открытия файла `16ссч` <br>- файл <br>- кол-во считанных байт <br>- кол-во записанных байт|В рамках проекта самый высокий приоритет отдается отслеживанию различного рода взаимодействий как между процессами, так и с операционной системой, поэтому необходимо знать, какие процессы в какие файлы пишут и из каких читают. | Среди возможных вариантов рассмотривалась так же утилита `lsof`, однако выводимая через неё информация содержит избыточные элементы и подлежит фильтрации, помимо этого кол-во прочитанных/записанных байт придётся отслеживать дополнительно. <br> Обращение к информации из системных вызовов с помощью `bpftrace`-сценарии - наиболее удобный инструмент сбора данного пункта мониторинга. |
| Общее кол-во записанных и прочитанных байт процессом | `bpftrace -e 'tracepoint:syscalls:sys_exit_write /args->ret/ { @[comm] = sum(args->ret); }'` <br> Информация выводится в формате: [command]: {кол-во байт} | Общая информация о том, как много процесс потребляет ресурсов может быть актуальна в рамках мониторинга и сбора статистики для финального отчета. | Однострочный `bpftrace`-сценарий - простое и удобное средство для сбора статистики, список которой, сразу отсортирован в порядке возрастания.|
| Список из 10-20 файлов, которые чаще всего используются процессами | bpftrace-сценарий [topfiles.bt](#topfiles.bt) <br> Выводит файл и кол-во обращений к нему. | Актуально в рамках рассмотрения нагрузки на файловую ситсему и сбра статистки по взаимодействиям между ней и процессами. | Единственный и самый простой способ выводить файлы, отсорированные по возрастанию колличества обращений к ним. |

В рамках рассмотрения IPC, FS будет нас интересовать только с точки зрения взаимодействий с ней различных процессов, а именно: 
   *	какие файлы открывает процесс;
   *	для чего он их открывает (записи или чтения);
   *	кол-во считанных/записанных байт;
     
### Требуемые для отслеживания вызовы с передаваемыми аргументами ###
Системные вызовы, связанные с открытием файла:

![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/810543f6-a657-4e51-bc00-ac10927b0581)

Второй аргумент - это режим открытия файла, представляющий собой один или несколько флагов открытия, объединенных оператором побитового ИЛИ.
### Флаги режима открытия файла ###
| Флаг | Описание |
|:---:| --- |
|O_RDONLY	|Только чтение (0)|
|O_WRONLY	|Только запись (1)|
|O_RDWR	|Чтение и запись (2)|
|O_CREAT	|Создать файл, если не существует|
|O_TRUNC	|Стереть файл, если существует|
|O_APPEND	|Дописывать в конец|
|O_EXCL	|Выдать ошибку, если файл существует при использовании O_CREAT|

Системные вызовы возвращают файловый дескриптор, по которому можно обращаться к файлу.

Системные вызовы, связанные с чтением/записью из/в файл(а):

![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/bd3bd86b-091e-4ddd-a363-c0480a681f1b)
![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/269b7dfc-5650-4c47-b710-77337e0b14e2)

Системный вызов `readv()` работает также как `read`, но считывает несколько буферов.
<br>Системный вызов `writev()` работает также как `write`, но записывает несколько буферов.
<br>В системном вызове `preadv()` объединены возможности `readv()` и `pread`. Он выполняет ту же задачу что и `readv()`, но имеет четвёртый аргумент `offset`, задающий файловое смещение, по которому нужно выполнить операцию чтения.

При успешном выполнении системных вызовов возвращается количество считанных/записанных байт. В случае ошибки возвращается -1.
>Для успешного выполнения не считается ошибкой передача меньшего количества байт чем запрошено.

## fsorw.bt ##
```
BEGIN
{
	printf("Tracing file system syscalls... Hit Ctrl-C to end.\n");
	printf("%-6s %-16s %4s %-6s %-40s %-4s %-4s\n", "PID", "COMM", "FD", "FLAGS", "PATH", "R", "W");
}

tracepoint:syscalls:sys_enter_open,
tracepoint:syscalls:sys_enter_openat
{
	@filename[tid] = args.filename;
	@flags[tid] = args.flags;
}

tracepoint:syscalls:sys_exit_open,
tracepoint:syscalls:sys_exit_openat
/@filename[tid]/
{
	$ret = args.ret;
	@fd[tid] = $ret >= 0 ? $ret : -1;
}

tracepoint:syscalls:sys_exit_read,
tracepoint:syscalls:sys_exit_readv,
tracepoint:syscalls:sys_exit_preadv
/@fd[tid]/
{
  $ret = args.ret;
	$nbyte = $ret >= 0 ? $ret : -1;
	$nothing = "-";
	printf("%-6d %-16s %4d %-6x %-40s %-4d %-4s\n", pid, comm, @fd[tid], @flags[tid], str(@filename[tid]), $ret, $nothing);

  delete(@filename[tid]);
	delete(@flags[tid]);
	delete(@fd[tid]);
}

tracepoint:syscalls:sys_exit_write,
tracepoint:syscalls:sys_exit_writev,
tracepoint:syscalls:sys_exit_pwritev
/@fd[tid]/
{
  $ret = args.ret;
	$nbyte = $ret >= 0 ? $ret : -1;
	$nothing = "-";
	printf("%-6d %-16s %4d %-6x %-40s %-4s %-4d\n", pid, comm, @fd[tid], @flags[tid], str(@filename[tid]), $nothing, $ret);

  delete(@filename[tid]);
	delete(@flags[tid]);
	delete(@fd[tid]);
}

END
{
	clear(@filename);
	clear(@flags);
	clear(@fd);
}
```

```
Example output:
localhost:/home/anna # bpftrace /home/anna/Desktop/bpftrace/fsorw.bt
Attaching 12 probes...
Tracing file system syscalls... Hit Ctrl-C to end.
PID    COMM               FD FLAGS  PATH                                     R    W     
1235   X                  67 0      /usr/bin/VBoxClient                      -    32  
1465   VBoxClient          8 0      /run/user/1000/xauth_cctNJU              111  -   
1235   X                  67 0      /usr/bin/VBoxClient                      31   -   
1834   plasmashell        -1 90800                                           -    8   
1318   master             91 801    public/pickup                            -    1   
1320   pickup             11 90800  ���������.localdomai              -    12  
1465   VBoxClient          8 0      /run/user/1000/xauth_cctNJU              111  -   
1235   X                  67 0      /usr/bin/VBoxClient                      31   -   
2267   kate               30 80241  /home/anna/Desktop/bpftrace/fsorw.bt     -    1290
2267   kate               29 80000  /home/anna/Desktop/bpftrace/fsorw.bt     1290 -   
2267   kate               29 80000  /home/anna/Desktop/bpftrace/fsorw.bt     1290 -   
2267   kate               29 80000  /home/anna/Desktop/bpftrace/fsorw.bt     1290 -   
1465   VBoxClient          8 0      /run/user/1000/xauth_cctNJU              111  -   
1235   X                  67 0      /usr/bin/VBoxClient                      31   -   
2171   KIO::WorkerThre    17 90800  /home/anna/Desktop/bpftrace              -    8   
...
```

## One-Liners bpftrace-сценарии ##
Статистика по системному вызову `write` \ `read`

```
# bpftrace -e 'tracepoint:syscalls:sys_exit_write /args->ret/ { @[comm] = sum(args->ret); }'
Attaching 1 probe...
^C

@[kactivitymanage]: 16
@[rtkit-daemon]: 16
@[wireplumber]: 32
@[Qt bearer threa]: 56
@[plasmashell]: 80
@[kwin_x11]: 96
@[auditd]: 120
@[InputThread]: 197
@[bash]: 436
@[konsole]: 2409
@[QXcbEventQueue]: 8680
@[top]: 15355
@[lsof]: 6296840
```

## topfiles.bt ##
```
BEGIN
{
	printf("Tracing file system syscalls... Hit Ctrl-C to end.\n");
}

tracepoint:syscalls:sys_enter_open,
tracepoint:syscalls:sys_enter_openat
{
	@filename[str(args.filename)] = count();
}


END
{
	printf("\nTop 20 files :\n");
	print(@filename, 20);
	clear(@filename);
}
```

```
Example output:
ocalhost:/home/anna # bpftrace /home/anna/Desktop/bpftrace/topfiles.bt
Attaching 4 probes...
Tracing file system syscalls... Hit Ctrl-C to end.
^C
Top 20 files :
@filename[/home/anna/.local/share/RecentDocuments/bpftrace[5].desktop]: 23
@filename[/proc/1651/cmdline]: 25
@filename[/home/anna/.local/share/RecentDocuments/fsopen.bt[3].desktop]: 26
@filename[/proc/sys/kernel/random/boot_id]: 27
@filename[/var/lib/dbus/machine-id]: 27
@filename[/home/anna/.local/share/RecentDocuments/fsorw.bt[3].desktop]: 28
@filename[/run/user/1000/xauth_jrgNiT]: 33
@filename[/home/anna/.local/share/RecentDocuments/bpftrace[2].desktop]: 34
@filename[/home/anna/.local/share/RecentDocuments/bpftrace[4].desktop]: 34
@filename[/home/anna/.local/share/RecentDocuments/topfiles.bt.desktop]: 37
@filename[/home/anna/.local/share/RecentDocuments/fsorw.bt[2].desktop]: 38
@filename[/home/anna/.local/share/RecentDocuments/fsorw.bt.desktop]: 40
@filename[/home/anna/.local/share/RecentDocuments/lsof.c.desktop]: 41
@filename[/sys/fs/cgroup/unified/user.slice/user-1000.slice/user@1000.ser]: 44
@filename[/home/anna/.local/share/RecentDocuments/bpftrace[3].desktop]: 44
@filename[/run/mount/utab]: 46
@filename[/home/anna/.local/share/RecentDocuments/bpftrace.desktop]: 46
@filename[/proc/self/mountinfo]: 51
@filename[/home/anna/.local/share/RecentDocuments]: 57
@filename[/etc/ld.so.cache]: 77
```
