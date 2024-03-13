| Пункт мониторинга | Утилита для получения информации| Актуальность | Обоснование выбора утилиты |
| :---: | --- | --- | --- |
| Неименованные каналы | bpftrace-сценарий [pipes.bt](#pipes.bt) <br> В качестве ключевой будет выводиться следующая информация: <br>- pid процесса <br>- команда <br>- файловые дескрипторы <br>- ошибка открытия канала <br>- флаги режима открытия файла `16ссч` <br>- файл | Неименованные каналы являются традиционным средством взаимодействия между связными процессами в ОС UNIX. <br>Характеристики, которые собираются bpftrace-сценарием необходимы для последующего составления пар взаимодействующих процессов.| Отслеживание системных вызовов - единственный способ проследить ipc по каналам. <br>Bpftrace-сценарии наиболее удобный инструмент, так как позволяет сразу отфильтровать необходимую информацию из большого потока приходящих системных вызовов и предоставить её в читаемом формате, удобном для последующей обработки. |  
| Именованные каналы | bpftrace-сценарий [named_pipes.bt](#named_pipes.bt) <br> В качестве ключевой будет выводиться следующая информация: <br>- pid процесса <br>- команда <br>- файловый дескриптор <br>- ошибка открытия канала <br>- флаги режима открытия файла `16ссч` <br>- файл | Именованные каналы являются вляются традиционным средством взаимодействия и синхронизации произвольных процессов в ОС UNIX. <br>Характеристики, которые собираются bpftrace-сценарием необходимы для последующего составления пар взаимодействующих процессов.| Отслеживание системных вызовов - единственный способ проследить ipc по каналам. <br>Bpftrace-сценарии наиболее удобный инструмент, так как позволяет сразу отфильтровать необходимую информацию из большого потока приходящих системных вызовов и предоставить её в читаемом формате, удобном для последующей обработки. |  
| Семафоры (System V) |  bpftrace-сценарий [semaphores.bt](#semaphores.bt) <br> В качестве ключевой будет выводиться следующая информация: <br>- pid процесса <br>- команда <br>- ключь группы семафоров <br>- id семафора  <br>- кол-во семафоров в наборе | Семафоры обеспечивают возможность синхронизации процессов при доступе к совместно используемым ресурсам, что актуально в рамках рассмотрения межпроцессорного взаимодействия. | Для мониторинга ipc с семафорами был выбран вариант с bpftrace-сценариями, так как используя другие инструменты для получения информации по семафорам и процессам, которые их используют, требуется постоянно вызывать утилиту с разным набором аргуметов. К тому же предоставляемая информация оказывается избыточной. Bpftrace-сценарии позволяют избежать всего этого.|
| Семафоры (Posix) | `ls -l /dev/shm` <br> `lsof /dev/shm/<sem_name>` | Семафоры обеспечивают возможность синхронизации процессов при доступе к совместно используемым ресурсам, что актуально в рамках рассмотрения межпроцессорного взаимодействия. <br>Семафоры POSIX более новые и сейчас гораздо чаще используются в работе. | Представленный способ мониторинга не является достаточно хорошим, постараюсь найти варианты, как улучшить существующий или найти замену. |

# Неименованные каналы (Pipes) 
Используются для связанных процессов, т.е. ими могут пользоваться только создавший их процесс и его потомки.

### Tracepiont входа в системный вызов: ###

![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/bf612938-2745-4ef8-88ac-ab6b18cf3890)

fildes: указывает на массив файловых дескрипторов, где первый элемент указывает на конец чтения, а второй на сторону записи.

### Tracepiont выхода из системного вызова: ###

![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/4c553260-dd1f-41be-ae14-ea97711e1ce8)

ret: возвращаемое значение

## pipes.bt ##
```
 BEGIN
{
	printf("Pipes... Hit Ctrl-C to end.\n");
	printf("%-6s %-16s %4s %3s %-16s %-5s %-5s\n", "PID", "COMM", "FD", "ERR", "FLAGS", "PIPE1", "PIPE2");
}

tracepoint:syscalls:sys_enter_pipe,
tracepoint:syscalls:sys_enter_pipe2
{
	@pipename1[tid] = args.fildes[1]; ???????????
	@pipename2[tid] = args.fildes[2];
}

tracepoint:syscalls:sys_enter_open,
tracepoint:syscalls:sys_enter_openat
/@pipename1[tid] || @pipename2[tid]/
{
	@flag[tid] = args.flags;
}

tracepoint:syscalls:sys_exit_open,
tracepoint:syscalls:sys_exit_openat
/@pipename1[tid] || @pipename2[tid]/
{
	$ret = args.ret;
	$fd = $ret >= 0 ? $ret : -1;
	$errno = $ret >= 0 ? 0 : - $ret;

	printf("%-6d %-16s %4d %3d %-16x %-5d %-5d\n", pid, comm, $fd, $errno, @flag[tid], @pipename1[tid], @pipename2[tid]);
	delete(@pipename1[tid]);
	delete(@pipename2[tid]);
	delete(@flag[tid]);
}

END
{
    	clear(@pipename1);
    	clear(@pipename2);
	clear(@flag);
}

```
```
Example output:
localhost:/home/anna # bpftrace /home/anna/Desktop/bpftrace/pipes.bt
Attaching 8 probes...
Pipes... Hit Ctrl-C to end.
PID    COMM               FD ERR FLAGS            PIPE1 PIPE2
1591   systemd            22   0 80101            -1    0    
5081   bash               -1   2 80000            32563 0    
5074   konsole            15   0 80000            32563 1817418240
5081   bash                3   0 241              32766 37   
5081   bash                3   0 0                32766 3    
5081   bash                3   0 0                -1    5081 
2130   plasmashell        20   0 90800            32693 0    
5109   bash                3   0 241              32764 37   
5109   bash               -1   2 0                32764 3    
5109   bash                3   0 0                -1    5105

```

# Именованные каналы 
Служат для общения и синхронизации произвольных процессов, знающих имя данного канала и имеющих соответствующие права доступа.

Требуемые для отслеживания системные вызовы с передаваемыми аргументами: 

![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/d9e8f5bc-4782-40df-9125-35f84f23c053)

filename: узел файловой системы
<br>mode: определяет тип файла и режим файла 
<br>dev: информация об устройстве
```
mknod(FIFO_FILE, S_IFIFO|0640, 0)
это означает чтение и запись (4 + 2 = 6) для владельца, чтение (4) для группы и отсутствие разрешений (0) для других
```
> !Важно! bpftrace не поддерживает вывод в 8ССч, а именно в этом формате представляются данные, поэтому необходимо делать перевод сомастоятельно для удобного чтения.

## named_pipes.bt ##
```
BEGIN
{
	printf("Named Pipes... Hit Ctrl-C to end.\n");
	printf("%-6s %-16s %4s %3s %-16s %s\n", "PID", "COMM", "FD", "ERR", "FLAGS", "PIPE");
}

tracepoint:syscalls:sys_enter_mknod,
tracepoint:syscalls:sys_enter_mknodat
{
	@pipename[tid] = args.filename;
}

tracepoint:syscalls:sys_enter_open,
tracepoint:syscalls:sys_enter_openat
/@pipename[tid]/
{
	@flag[tid] = args.flags;
}

tracepoint:syscalls:sys_exit_open,
tracepoint:syscalls:sys_exit_openat
/@pipename[tid]/
{
	$ret = args.ret;
	$fd = $ret >= 0 ? $ret : -1;
	$errno = $ret >= 0 ? 0 : - $ret;

	printf("%-6d %-16s %4d %3d %-16x %s\n", pid, comm, $fd, $errno, @flag[tid], str(@pipename[tid]));
	delete(@pipename[tid]);
	delete(@flag[tid]);
}

END
{
  clear(@pipename);
	clear(@flag);
}
```
```
Example output:
# ./named_pipe.bt
Attaching 3 probes...
Named Pipe... Hit Ctrl-C to end.
PID    COMM               FD ERR FLAGS PATH
2440   snmp-pass           4   0       /proc/cpuinfo
2440   snmp-pass           4   0       /proc/stat
25706  ls                  3   0       /etc/ld.so.cache
25706  ls                  3   0       /lib/x86_64-linux-gnu/libselinux.so.1

localhost:/home/anna # bpftrace /home/anna/Desktop/bpftrace/namedpipe.bt
Attaching 8 probes...
Named Pipes... Hit Ctrl-C to end.
PID    COMM               FD ERR FLAGS            PIPE
833    systemd-logind     22   0                  /run/systemd/inhibit/7.ref
4210   mkfifoo             3   0                  /tmp/my_named_pipe1
^C
[3]+  Stopped                 bpftrace /home/anna/Desktop/bpftrace/namedpipe.bt

```

# Семафоры (Semaphres)
## System V Semaphres ##
Cредства, обеспечивающие возможность синхронизации процессов при доступе к совместно используемым ресурсам, например, к разделяемой памяти.

Существует три варианта отследить работу с семафорами:
1. Псевдофайловая система /proc/sysvipc/sem
   ![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/d1dbf8b0-bcda-4d00-9c5c-b16d66f6eacd)

2. Утилита `ipcs -s`. Ключь `-i` с указанием semid позволяет увидеть pid процесса, который последним работал с семафором.
![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/2b7cb0bf-26a6-42d3-b416-414d59c7c9cf)

3. `Bpftrace`-сценарии
	Семафор в ОС UNIX состоит из следующих элементов:
	*	значение семафора;
	*	идентификатор процесса, который хронологически последним работал с семафором;
	*	число процессов, ожидающих увеличения значения семафора;
	*	число процессов, ожидающих нулевого значения семафора.
	Для работы с семафорами поддерживаются три системных вызова:
	*	`semget` для создания и получения доступа к набору семафоров;
	*	`semop` для манипулирования значениями семафоров (это именно тот системный вызов, который 		позволяет процессам синхронизоваться на основе использования семафоров);
	*	`semctl` для выполнения разнообразных управляющих операций над набором семафоров.

	Требуемые для отслеживания системные вызовы с передаваемыми аргументами: 
	![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/fdfdfc7b-f130-4452-b142-0d648e782f2e)
	
## semaphores.bt ##
```
BEGIN
{
	printf("Semaphores... Hit Ctrl-C to end.\n");
	printf("%s %s %s %s\n", "PID", "COMM", "KEY", "SEMID", "NSEM");
}

tracepoint:syscalls:sys_enter_semget
{
	@semkey[tid] = args.key;
	@nsems[tid] = args.nsem;
}

tracepoint:syscalls:sys_exit_semget
/@semkey[tid]/
{
	@semid[tid] = args.ret;
}

tracepoint:syscalls:sys_exit_semop
/@semid[tid]/
{
	printf("%d %s %x %d %d\n", pid, comm,  @semkey[tid], @semid[tid], @nsems[tid]);
	delete(@pipename[tid]);
	delete(@flag[tid]);
	delete(@nsems[tid]);
}

END
{
    	clear(@semkey);
	clear(@semid);
	clear(@nsems);
}
```
```
Example output:
# ./semaphores.bt
Attaching 4 probes...
Semaphores... Hit Ctrl-C to end.
PID COMM KEY SEMID NSEM
3143 lsof 16 12345 1

```

## POSIX Semaphres ##
***POSIX Named Semaphore calls***
```
#создание или открытие семафора
sem_t *sem_open (const char *name, int oflag);
sem_t *sem_open (const char *name, int oflag, mode_t mode, unsigned int value);

#увеличивают или уменьшают значение семафора
int sem_post (sem_t *sem);
int sem_wait (sem_t *sem);

#получает текущее значение семафора
int sem_getvalue (sem_t *sem, int *sval);

#закрывает семафор
sem_close(sem_t * sem);

#удаляет имя семафора и помечает его как удаленное, когда все процессы закрывают семафор
int sem_unlink (const char *name);
```
***POSIX Unnamed Semaphore calls***
```
#создание или открытие семафора
int sem_init (sem_t *sem, int pshared, unsigned int value);

#закрывает семафор
int sem_destroy (sem_t *sem);
```

Мониторинг именованных семафоров

В данном случае поиск по доступным вызовам в bpftrace, к сожалению, не дал никакого результата, поэтому написать bpftrace-сценарии для отслеживания такого типа ipc не получится. 
<br> Единственным вариантом для отслеживания POSIX семафоров является `/dev/shm` в сочетании с `lsof`
```
$ ls -l /dev/shm
total 4
-rw-------. 1 alice alice 32 Jan 13 14:17 sem.my_named_semaphore

$ lsof /dev/shm/sem.my_named_semaphore
COMMAND    PID   USER  FD   TYPE DEVICE SIZE/OFF NODE NAME
sem_posix 4923  alice DEL    REG   0,22            10 /dev/shm/sem.my_named_semaphore
```

Существенным недостатоком данного метода сбора информации является невозможность обеспечить непрерывное получение данных, что может привести к потере информации.
