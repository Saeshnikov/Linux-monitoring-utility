# Утилита bpftrace (Dtrace)
Утилита DTrace обеспечивает динамическую трассировку, которая представляет собой возможность инструментирования запущенного ядра операционной системы.

DTrace позволяет связывать действия, такие как сбор или печать трассировок стека, аргументов функций, временных меток и статистических агрегатов, 
с зондами, которые могут быть событиями времени выполнения или местоположениями исходного кода. Она позволяет изучать поведение пользовательских программ 
и операционной системы, понимать, как работает система, отслеживать проблемы с производительностью и находить причины отклоняющегося поведения.

В Linux есть полноценный аналог DTrace под названием bpftrace.
<br>bpftrace — это высокоуровневый язык трассировки для Linux enhanced Berkeley Packet Filter (eBPF), доступный в последних версиях ядер Linux (4.x). 
bpftrace использует LLVM в качестве серверной части для компиляции сценариев в байт-код BPF и использует BCC для взаимодействия с системой Linux BPF, а также существующую трассировку Linux. 

## Возможности использования bpftrace ##
bpftrace позволяет выяснить, что происходит внутри системного или библиотечного вызова, имеет возможность не просто составить список вызовов, но и, 
например, собрать статистику по определённому поведению, а также трассировать несколько процессов и сопоставить данные из нескольких источников. 
Имеется доступ к различной контекстной информации, такой как текущий PID, трассировка стека, время, аргументы вызовов, возвращаемые значения, и т.д.

Возможности: динамическая трассировка ядра (kprobes), динамическая трассировка на уровне пользователя (uprobes) и точки трассировки. 
Язык трассировки bpf вдохновлен awk и C, а также предшествующими трассировщиками, такими как DTrace и SystemTap.

Внутри bpfTrace работает так: пишем bpf-программу, которая парсится, конвертируется в C, потом обрабатывается через Clang, который генерирует bpf-байт-код, после этого программа загружается в ядро.

## Экосистема трассировки (bpftrace Probe Types) ##
![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/f2a6b0fa-7a1d-4f31-8d86-eebda4f2f93f)

События для инструментирования:
1.	Специальные события: `BEGIN`, `END` (начало и конец выполнения bpftrace)
2.	События, основанные на kprobes (динамическая трассировка):
    *	`kprobe, kretprobe` (запускается в момент входа в функцию, запускается в момент выхода из функции соответственно) - используют технологию ядра
      Linux для обеспечения динамической трассировки функций ядра. При этом код возврата функции содержится в специальной встроенной переменной retval.

     	С его помощью можно трассировать не только сами системные вызовы, но и то, что происходит внутри них (потому что точки входа системных вызовов вызывают другие внутренние функции).
     	Также можно использовать kprobes для трассировки событий ядра, не являющихся системными вызовами, например, «буферизированные данные записываются на диск»,
     	«TCP-пакет посылается по сети» или «в данный момент происходит переключение контекста».
    *	`uprobe, uretprobe` - инструментирование программ, работающих в пространстве пользователя. Внедряют технологию ядра Linux для обеспечения динамической трассировки функций пользовательского уровня

> В теории, strace может быть реализован с помощью kprobes, а ltrace — с помощью uprobes.

3.  События, основанные на tracepoints(статическая трассировка): `tracepoint`
    * Точки трассировки (tracepoints) ядра позволяют трассировать нестандартные события, определённые разработчиками ядра. Эти события находятся
      не на уровне вызовов функций. Для создания таких точек разработчики ядра вручную размещают макрос TRACE_EVENT в коде ядра.

> У обоих источников есть и плюсы, и минусы. Kprobes работает «автоматически», т.к. не требует от разработчиков ядра ручной разметки кода.
> Но события kprobe могут произвольно меняться от одной версии ядра к другой, потому что функции постоянно изменяются — добавляются, удаляются, переименовываются.
>
> Точки трассировки ядра, как правило, более стабильны с течением времени и могут предоставлять полезную контекстную информацию, которая может быть недоступна в случае использования kprobes.
> Используя kprobes, можно получить доступ к аргументам вызовов функций. Но с помощью точек трассировки можно получить любую информацию, которую разработчик ядра решит вручную описать.

4.	События, основанные на BPF trampolines: `kfunc, kretfunc`
5.	Статическая отладка в пространстве пользователя: `usdt` (позволяет добавить в программу статических точек останова в момент компиляции.)
6.	События, основанные на подсистеме perf:  `software`, `hardware`, `profile`, `interval`, `watchpoint`
    *	software: События программного обеспечения ядра
    *	hardware: События на уровне процессора
    *	interval: интервальное событие
    *	profile: интервальное событие для профилирования (сколько раз в секунду)

Некоторые типы зондов допускают подстановочные знаки для соответствия нескольким зондам, например, `probe:vfs_*`.

![image](https://github.com/Saeshnikov/Linux-monitoring-utility/assets/121693400/b613c7f7-503e-43b4-865d-873126fc2496)
<br>*Основные источники трассируемых событий в Linux*

## Опции ##
### Использование: ###
    bpf trace [options] filename
> Программы, сохраненные в виде файлов, часто называются скриптами, и их можно запустить, указав их имя файла с расширением файла .bt, сокращение от bpf trace, но расширение игнорируется.

    bpf trace [options] -e 'program'

### Опции ###

| Ключь | Описание |
| :---: | --- |
|-B MODE           | Режим буферизации вывода ('line', 'full', or 'none')|
|    -d            | отладочная информация для сухого запуска|
|    -dd           | подробная информация об отладке при запуске|
|    -e 'program'  | выполняет эту программу|
|    -h            | показывает это справочное сообщение|
|    -I DIR        | добавляет указанный каталог в путь поиска включаемых файлов|
|    -l [search]   | список зондов|
|    -p PID        | включить USDT-зонды или выполнить поиск uprobes/uretprobes в адресном пространстве PID|
|    -v            | подробные сообщения|

## Структура программ bpftrace ##
Программы состоят из списка блоков вида:

    probe[,probe,...] /filter/ { action }

Например,
 ```
# bpftrace -e 'kprobe:do_sys_open { printf("opening: %s\n", str(arg1)); }'
Attaching 1 probe...
opening: /proc/cpuinfo
opening: /proc/stat
opening: /proc/diskstats
opening: /proc/stat
opening: /proc/vmstat
[...]
```
`<filter>` является опциональным и используется для фильтрации событий,

Например, представленная ниже программа будет передавать привет, только если запускается на CPU 0.
```
# bpftrace -e 'p:s:1 /cpu == 0/ { printf("Привет с CPU%d\n", cpu); }'
Attaching 1 probe...
Привет с CPU0
Привет с CPU0
^C
```

Поддерживает: ++, -- , ?, if-else statements, unroll, array[]

## Примеры использования ##
*	Зонды из библиотек tracepoint и kprobe могут быть перечислены с помощью `-l`.
```
# bpftrace -l | more
tracepoint:xfs:xfs_attr_list_sf
tracepoint:xfs:xfs_attr_list_sf_all
tracepoint:xfs:xfs_attr_list_leaf
tracepoint:xfs:xfs_attr_list_leaf_end
[...]
# bpftrace -l | wc -l
46260
```
* Параметр -v при перечислении точек трассировки покажет их аргументы для использования из встроенного args. Например:
```
# bpftrace -lv tracepoint:syscalls:sys_enter_open
tracepoint:syscalls:sys_enter_open
    int __syscall_nr;
    const char * filename;
    int flags;
    umode_t mode;
```
* Также можно найти события по шаблону. 
```
# bpftrace -l '*kill_all*'
kprobe:rfkill_alloc
kprobe:kill_all
kprobe:btrfs_kill_all_delayed_nodes
```
* Найти события только среди tracepoints:
```
# bpftrace -l 't:*kill*'
tracepoint:cfg80211:rdev_rfkill_poll
tracepoint:syscalls:sys_enter_kill
tracepoint:syscalls:sys_exit_kill
```
>Для каждого системного вызова X определены две точки останова:
>`tracepoint:syscalls:sys_enter_X`
><br>`tracepoint:syscalls:sys_exit_X`
*	Наблюдение в реальном времени, какие файлы открываются какими процессами
```
sudo bpftrace -e 'tracepoint:syscalls:sys_enter_open,
  tracepoint:syscalls:sys_enter_openat {
    printf("%s %s\n", comm, str(args->filename));
  }'
```
* Какие функции доступны для трассировки:
```
# это для конкретного процесса
sudo bpftrace -p 20419 -l | grep 'uprobe:'
```
* Какие функции доступны для конкретного исполняемого файла:
```
sudo bpftrace -l 'uprobe:/home/eax/pginstall/bin/postgres:*'
```

### Готовые решения для bpftrace ###
> bpftrace идет с рядом готовых утилит

| Название | Описание |
| :---: | --- |
|tools/bashreadline.bt: | печать введенных команд bash в масштабах всей системы |
|tools/biosnoop.bt: | инструмент трассировки блочного ввода-вывода, показывающий задержку ввода-вывода |
|tools/cpuwalk.bt: | пример того, какие процессоры выполняют данные процессы |
|tools/dcsnoop.bt:  | поиск в кэше записей каталогов трассировки (dcache) | 
|tools/execsnoop.bt: | трассировка новых процессов с помощью системных вызовов exec(); |
|tools/gethostlatency.bt: | Показывает задержку для вызовов  getaddrinfo/gethostbyname, путем определения имени удаленного хоста, поиск которого был медленным и насколько |
|tools/naptime.bt: | Показывать добровольные вызовы режима сна |
|tools/opensnoop.bt: | Отслеживание системных вызовов open() с отображением имен файлов к которым они обращаются |
|tools/pidpersec.bt: | Подсчитывать новые процессы (через fork) |
|tools/statsnoop.bt: | Отслеживание системных вызовов stat() для общей отладки |
|tools/syscount.bt: | Подсчет системных вызовов |
|tools/tcpaccept.bt: | Отслеживание пассивных соединений TCP (accept()) |
|tools/tcpconnect.bt: | Отслеживание активных соединений TCP (connect()) |
|tools/tcpdrop.bt: | Этот инструмент отслеживает TCP-пакеты или сегменты, которые были отброшены ядром, и показывает подробную информацию из заголовков IP и TCP, состояния сокета и трассировки стека ядра |
|tools/tcplife.bt: | Отслеживание продолжительности сеанса TCP с указанием сведений о подключении |

***Пример использования `syscount`***
```
syscount подсчитывает системные вызовы, и выводит топ 10 индефикаторов системных вызовов с их колличеством и
топ 10 процессов с колличеством системных вызовов,которые они совершают

# ./syscount.bt
Attaching 3 probes...
Counting syscalls... Hit Ctrl-C to end.
^C
Top 10 syscalls IDs:
@syscall[6]: 36862
@syscall[21]: 42189
@syscall[13]: 44532
@syscall[12]: 58456
@syscall[9]: 82113
@syscall[8]: 95575
@syscall[5]: 147658
@syscall[3]: 163269
@syscall[2]: 270801
@syscall[4]: 326333

Top 10 processes:
@process[rm]: 14360
@process[tail]: 16011
@process[objtool]: 20767
@process[fixdep]: 28489
@process[as]: 48982
@process[gcc]: 90652
@process[command-not-fou]: 172874
@process[sh]: 270515
@process[cc1]: 482888
@process[make]: 1404065
```
> Приведенный выше вывод был отслежен во время сборки ядра Linux, и имя процесса
с наибольшим количеством системных вызовов было "make" с 1 404 065 системными вызовами во время трассировки. 
Наивысший идентификатор системного вызова был равен 4, что cоответствует stat().

***Пример использования `tcpconnect`***
```
# ./tcpconnect.bt
TIME     PID      COMM             SADDR          SPORT  DADDR          DPORT
00:36:45 1798396  agent            127.0.0.1      5001   10.229.20.82   56114
00:36:45 1798396  curl             127.0.0.1      10255  10.229.20.82   56606
00:36:45 3949059  nginx            127.0.0.1      8000   127.0.0.1      37780
```

> Этот вывод показывает три подключения, одно от процесса "agent", одно от
"curl" и одно от "nginx". В выходных данных указаны версия IP,
адрес источника, порт исходного сокета, адрес назначения и порт назначения-получателя.
> <br>Это отслеживает попытки
подключения: возможно, они завершились неудачей.

***Пример использования `tcplife`***
```
Этот инструмент показывает продолжительность сеансов TCP, включая статистику пропускной способности,
и для повышения эффективности измеряет только изменения состояния TCP (а не все пакеты).
Например:

# ./tcplife.bt
PID   COMM       LADDR           LPORT RADDR           RPORT TX_KB RX_KB MS
20976 ssh        127.0.0.1       56766 127.0.0.1       22         6 10584 3059
20977 sshd       127.0.0.1       22    127.0.0.1       56766  10584     6 3059
14519 monitord   127.0.0.1       44832 127.0.0.1       44444      0     0 0
4496  Chrome_IOT 7f00:6:5ea7::a00:0 42846 0:0:bb01::      443        0     3 12441
4496  Chrome_IOT 7f00:6:5aa7::a00:0 42842 0:0:bb01::      443        0     3 12436
4496  Chrome_IOT 7f00:6:62a7::a00:0 42850 0:0:bb01::      443        0     3 12436
4496  Chrome_IOT 7f00:6:5ca7::a00:0 42844 0:0:bb01::      443        0     3 12442
4496  Chrome_IOT 7f00:6:60a7::a00:0 42848 0:0:bb01::      443        0     3 12436
4496  Chrome_IOT 10.0.0.65       33342 54.241.2.241    443        0     3 10717
4496  Chrome_IOT 10.0.0.65       33350 54.241.2.241    443        0     3 10711
4496  Chrome_IOT 10.0.0.65       33352 54.241.2.241    443        0     3 10712
14519 monitord   127.0.0.1       44832 127.0.0.1       44444      0     0 0
```
> Вывод начинается с ssh-соединения localhost, поэтому видны обе конечные
точки: процесс ssh (PID 20976), который получил 10584 Кбайт, и процесс sshd
(PID 20977), который передал 10584 Кбайт. Этот сеанс длился 3059
миллисекунд. Также можно просмотреть другие сеансы, включая подключения по протоколу IPv6.