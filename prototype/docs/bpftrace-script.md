# Bpftrace-script
Обращается к системным вызовам, которые относятся к файлам и директориям. <br>Выводит все файлы, к которым обращалась система за время работы и колличесвто обращений к ним.
```
BEGIN
{
	printf("Tracing file system syscalls... Hit Ctrl-C to end.\n");
}

//SYMLINK

tracepoint:syscalls:sys_enter_readlink
{
	@oldname[tid] = args.path;
}
tracepoint:syscalls:sys_enter_readlinkat
{
	@oldname[tid] = args.pathname;
}

tracepoint:syscalls:sys_enter_open,
tracepoint:syscalls:sys_enter_openat,
tracepoint:syscalls:sys_enter_execve,
tracepoint:syscalls:sys_enter_execveat,
tracepoint:syscalls:sys_enter_openat2,
tracepoint:syscalls:sys_enter_open_tree,
tracepoint:syscalls:sys_enter_statx
/@oldname[tid]/
{
	@filename[str(@oldname[tid])] = count();
	@filename[str(args.filename)] = count();
}

tracepoint:syscalls:sys_enter_fspick
/@oldname[tid]/
{
	@filename[str(@oldname[tid])] = count();
	@filename[str(args.path)] = count();
}

tracepoint:syscalls:sys_enter_name_to_handle_at
/@oldname[tid]/
{
	@name[tid] = args.name;
}


//ALL

tracepoint:syscalls:sys_enter_name_to_handle_at
{
	@name[tid] = args.name;
}
tracepoint:syscalls:sys_exit_name_to_handle_at
/@name[tid]/
{
	$ret = args.ret;
	@fd[tid] = $ret >= 0 ? $ret : -1;
}
tracepoint:syscalls:sys_enter_open_by_handle_at
/@fd[tid]/
{
	@filename[str(@name[tid])] = count();
	delete(@fd[tid]);
}

tracepoint:syscalls:sys_enter_open,
tracepoint:syscalls:sys_enter_openat,
tracepoint:syscalls:sys_enter_execve,
tracepoint:syscalls:sys_enter_execveat,
tracepoint:syscalls:sys_enter_openat2,
tracepoint:syscalls:sys_enter_open_tree,
tracepoint:syscalls:sys_enter_statx
{
	@filename[str(args.filename)] = count();
}

tracepoint:syscalls:sys_enter_fspick
{
	@filename[str(args.path)] = count();
}

END
{
	print(@filename);
	clear(@oldname);
	clear(@filename);
	clear(@name);
	clear(@fd);
}
```

```
Example output:
localhost:/home/anna # bpftrace /home/anna/Desktop/bpftrace/allfiles.bt
Attaching 22 probes...
Tracing file system syscalls... Hit Ctrl-C to end.
^C
All files :
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
...
```

### Краткое описание системных вызовов ###
Были просмотрены все `tracepoint:syscalls:` доступные через `bpftrace` и выбраны те, которые могут быть использованы.

| t:s:sys_enter_ | Описание | Аргументы | Вопросы |
|:---:| --- | --- | --- |
| chdir | Изменяет текущий рабочий каталог вызвавшего процесса на каталог, указанный в path. | Const char* filename | ```-``` |
| fchdir | Идентичен chdir, разница в том, что каталог указывается в виде открытого файлового дескриптора. | Unsigned int fd | ```-``` |
| chroot | Изменить корневой каталог на тот, что задан аргументом filename. | Const char* filename | ```-``` |
| Execve <br> Execveat | Выполнять программу, заданную параметром filename. | (int fd)<br>Const char* filename<br>Const char* const* argv<br>Const char* const* envp<br>(int flags)|  EFAULT (filename указывает за пределы доступного адресного пространства)<br>ENOEXEC (исполняемый файл в неизвестном формате или же встречены ошибки, препятствующие его выполнению)<br>EACCES (не прав на выполнение/поиск в одном из каталогов по пути filename)<br>ENOENT (filename не сущетсвует)<br>ENOTDIR (компонент пути не является каталогом)|
| fspick | Используется для открытия/выбора или отсоединения файлов из файловой системы. Обычно используется для проверки, существует ли конкретный файл в опр. каталоге. | Int dfd<br>Const char* path<br>Unsigned int flags|  ```!``` обратить внимание на сочетание dfd и path в процессе обработки результатов|
| Link<br>Linkat | Создает новую “жесткую” ссылку на существующий файл. | (int olddfd)<br>Const char* oldname<br>(int newdfd)<br>Const char* newname<br>(int flags)| ```-``` Можно перемещать, переименовывать и удалять файл без вреда ссылке -> нет необходимости отслеживать её создание. |
| Open<br>Openat<br>[Openat2] | Открытие файла. | (int dfd)<br>Const char* filename<br>Int flags<br>Umode_t mode| EACCES<br>EFAULT<br>ELOOP (filename был символической ссылкой, и были указаны флаги<br>O_NOFOLLOW, но не O_PATH)<br>ENAMETOOLONG (путь был слишком длинным)<br>ENOENT (O_CREAT не задан, и именованный файл не существует/ компонент каталога в pathname не существует или являетсявисячей символической ссылкой)|
| open_by_handle_at | Открытие файла через описатель (name_to_handle_at). | Int mountdirfd<br>Struct file_handle* handle<br>Int flags|  |
| open_tree | Можно использовать для открытия файла относительно файлового дескриптора открытого каталога. | Int dfd<br>Const char* filename<br>Unsigned flags |  |
| Readlink<br>Readlinkat | Считывает значение символьной ссылки (помещает содержимое ссылки в buf). | (int dfd)<br>Const char* path(name)<br>Char * buf<br>Int bufsiz| Используется для проверки того, что файл по ссылке откроется. |
| Rename<br>Renameat<br>[Renameat2] | Изменяет имя или расположение файла. | (int olddfd)<br>Const char* oldname<br>(int newdfd)<br>Const char* newname<br>[unsigned int flags] | ```-``` |
| splice | Перемещает данные между двумя файловыми дескрипторами, не выполняя при этом копирования между адресным пространством пользователя и ядра (один обязательно канал). | Int fd_in<br>Loff_t* off_in<br>Int fd_out<br>Loff_t* off_out<br>Size_t len<br>Unsigned int flags | ```-``` |
| statx | Возвращает информацию о файле, записывая ее в буфер. | Int dfd<br>Const char* filename<br>Unsigned flags<br>Unsigned int mask<br>Struct statx * buffer| Информация о файле может считаться за использование этого файла |
| Symlink<br>Symlinkat| Создает новую “символьную” ссылку на существующий файл. | Const char* oldname<br>(int newdfd)<br>Const char* newname | Символьная ссылка после удаления исходного файла становится недействительными, следовательно, хоть и прямых обращений к файлу-оригиналу нет, мы не можем утверждать, что он нам не нужен. Исходя из этого данный системный вызов надо отслеживать и добавлять в общий список как oldname, так и newname. При этом (чтобы избежать того, что ссылка создана, но нигде не используется) отфильтруем и будем добавлять в результат только то, используется|
