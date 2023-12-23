ProcFS — специальная файловая система, используемая в Linux, позволяет получить доступ к информации из ядра о системных процессах.

ВАЖНО: содержание папки /proc может изменяться в зависимости от системы

****
### Информация о конкретном процессе (/proc/PID):
| file | content |
| ---- | ------- |
clear_refs(--) | Clears page referenced bits shown in smaps output 
cmdline         |Command line arguments
cwd             |Link to the current working directory
environ(-)      |Values of environment variables
exe             |Link to the executable of this process
fd              |Directory, which contains all file descriptors
maps            |Memory maps to executables and library files
mem(--)         |Memory held by this process
root            |Link to the root directory of this process
stat            |Process status
statm           |Process memory status information
status          |Process status in human readable form
wchan           |Present with CONFIG_KALLSYMS=y: it shows the kernel function symbol the task is blocked in - or "0" if not blocked.
pagemap         |Page table
stack           |Report full stack trace, enable via CONFIG_STACKTRACE
smaps           |An extension based on maps, showing the memory consumption of each mapping and flags associated with it
smaps_rollup    |Accumulated smaps stats for all mappings of the process.  This can be derived from smaps, but is faster and more convenient
numa_maps       |An extension based on maps, showing the memory locality and binding policy as well as mem usage (in pages) of each mapping.

****
### /status 

| field | content |
| ----- | ------- |
Name                              |filename of the executable
Umask                             |file mode creation mask
State                             |state (R is running, S is sleeping, D is sleeping in an uninterruptible wait, Z is zombie, T is traced or stopped)
Tgid                              |thread group ID
Ngid                              |NUMA group ID (0 if none)
Pid                               |process id
PPid                              |process id of the parent process
TracerPid                         |PID of process tracing this process (0 if not)
Uid                               |Real, effective, saved set, and  file system UIDs
Gid                               |Real, effective, saved set, and  file system GIDs
FDSize                            |number of file descriptor slots currently allocated
Groups                            |supplementary group list
NStgid                            |supplementary group list
NStgid                            |descendant namespace thread group ID hierarchy
NSpid                             |descendant namespace process ID hierarchy
NSpgid                            |descendant namespace process group ID hierarchy 
NSsid                             |descendant namespace session ID hierarchy
VmPeak                            |peak virtual memory size
VmSize                            |total program size
VmLck                             |locked memory size
VmPin                             |pinned memory size
VmHWM                             |peak resident set size ("high water mark")
VmRSS                             |size of memory portions. It contains the three  following parts (VmRSS = RssAnon + RssFile + RssShmem)
RssAnon                           |size of resident anonymous memory
RssFile                           |size of resident file mappings
RssShmem                          |size of resident shmem memory (includes SysV shm, mapping of tmpfs and shared anonymous mappings)
VmData                            |size of private data segments
VmStk                             |size of stack segments
VmExe                             |size of text segment
VmLib                             |size of shared library code
VmPTE                             |size of page table entries
VmSwap                            |amount of swap used by anonymous private data (shmem swap usage is not included)
HugetlbPages                      |size of hugetlb memory portions
CoreDumping                       |process's memory is currently being dumped (killing the process may lead to a corrupted core)
THP_enabled                       |process is allowed to use THP (returns 0 when PR_SET_THP_DISABLE is set on the process)
Threads                           |number of threads
SigQ                              |number of signals queued/max. number for queue
SigPnd                            |bitmap of pending signals for the threa
ShdPnd                            |bitmap of shared pending signals for the process
SigBlk                            |bitmap of blocked signals
SigIgn                            |bitmap of ignored signals
SigCgt                            |bitmap of caught signals
CapInh                            |bitmap of inheritable capabilities
CapPrm                            |bitmap of permitted capabilities
CapEff                            |bitmap of effective capabilities
CapBnd                            |bitmap of capabilities bounding set
CapAmb                            |bitmap of ambient capabilities
NoNewPrivs                        |no_new_privs, like prctl(PR_GET_NO_NEW_PRIV, ...)
Seccomp                           | seccomp mode, like prctl(PR_GET_SECCOMP, ...)
Speculation_Store_Bypass          | speculative store bypass mitigation status
Cpus_allowed                      | mask of CPUs on which this process may run
Cpus_allowed_list                 | Same as previous, but in "list format"
Mems_allowed                      | mask of memory nodes allowed to this process
Mems_allowed_list                 | Same as previous, but in "list format"
voluntary_ctxt_switches           | number of voluntary context switches
nonvoluntary_ctxt_switches        |number of non voluntary context switches

****
### /maps

| field | content |
| ----- | ------- |
address| the address space in the process that it occupies
perms|r = read </br> w = write </br> x = execute </br> s = shared </br> p = private (copy on write)
offset  | the offset into the mapping
dev   | device (major:minor)
inode      |the inode  on that device
pathname | name associated file for this mapping </br> If the mapping is not associated with a file: </br> [heap]= the heap of the program </br> [stack] =the stack of the main process <\br> [vdso]= the "virtual dynamic shared object",the kernel system call handler or if empty, the mapping is anonymous.

****

### Общая информация:

| file | content |
| ---- | ------- |
apm | Advanced power management info  
buddyinfo|Kernel memory allocator information
bus |Directory containing bus specific information
cmdline|Kernel command line 
cpuinfo|Info about the CPU
devices|Available devices (block and character) 
dma | Used DMS channels 
filesystems|Supported filesystems
driver|Various drivers grouped here, currently rtc
execdomains|Execdomains, related to security
fb|Frame Buffer devices
fs|File system parameters, currently nfs/exports
ide|Directory containing info about the IDE subsystem
interrupts|Interrupt usage 
iomem|Memory map
ioports|I/O port usage 
irq|Masks for irq to cpu affinity
isapnp |ISA PnP (Plug&Play) Info
kcore|Kernel core image (can be ELF or A.OUT(deprecated in 2.4))
kmsg|Kernel messages 
ksyms|Kernel symbol table  
loadavg|Load average of last 1, 5 & 15 minutes 
locks|Kernel locks
meminfo|Memory info
misc |Miscellaneous 
modules|List of loaded modules 
mounts|Mounted filesystems  
net |Networking info
pagetypeinfo|Additional page allocator information
partitions|Table of partitions known to the system 
pci|Deprecated info of PCI bus (new way → /proc/bus/pci/, decoupled by lspci
rtc|Real time clock 
scsi |SCSI info
slabinfo    |Slab pool info  
softirqs|softirq usage
stat|Overall statistics 
swaps |Swap space utilization 
sysvipc|Info of SysVIPC Resources (msg, sem, shm)
tty|Info of tty drivers
uptime|Wall clock since boot, combined idle time of all cpus
version|Kernel version
video|bttv info of video resources
vmallocinfo|Show vmalloced areas

****

### /net

Network info:

| file | content |
| ---- | ------- |
arp|Kernel  ARP table
dev|network devices with statistics 
dev_mcast |the Layer2 multicast groups a device is listening too (interface index, label, number of references, number of bound addresses). 
dev_stat|network device status 
ip_fwchains|Firewall chain linkage
ip_fwnames|Firewall chain names 
ip_masq |Directory containing the masquerading tables  
ip_masquerade|Major masquerading table
netstat |Network statistics 
raw |raw device statistics 
route|Kernel routing table 
rpc|Directory containing rpc info
rt_cache |Routing cache 
snmp |SNMP data 
sockstat|Socket statistics 
tcp |TCP  sockets 
udp|UDP sockets
unix|UNIX domain sockets 
wireless|Wireless interface data (Wavelan etc)   
igmp |IP multicast addresses, which this host joined  
psched |Global packet scheduler parameters.  
netlink|List of PF_NETLINK sockets  
ip_mr_vifs|List of multicast virtual interfaces  
ip_mr_cache|List of multicast routing cache 
