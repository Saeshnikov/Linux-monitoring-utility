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
