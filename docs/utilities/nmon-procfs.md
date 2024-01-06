## Сравнительный анализ nmon и procfs
### CPU
По процессору nmon берет информацию из procfs /proc/stat
#### Вывод в режиме записи:
- User% (/proc/stat user)
- Sys% (/proc/stat system)
- Wait% (/proc/stat iowait)
- Idle% (/proc/stat idle)
- Steal% (/proc/stat steal (since Linux 2.6.11))

### Memory
По памяти nmon берет информацию из procfs /proc/meminfo
#### Вывод в режиме записи:
- Memtotal
- hightotal
- lowtotal
- swaptotal
- memfree
- highfree
- lowfree
- swapfree
- memshared
- cached
- active
- buffers
- swapcached
- inactive

### Disks
По дискам nmon берет информацию из procfs /proc/diskstats и преобразует ее
Вывод в режиме записи:
- DISKBUSY time spent doing I/Os (ms) / интервал
 Процент времени, в течение которого диск активен.
- DISKREAD delta reads merged / интервал
Общее количество операций чтения с диска в КБ в секунду.
- DISKWRITE delta writes completed / интервал
Общее количество операций записи на диск в КБ в секунду.
- DISKXFER delta( reads completed successfully + writes completed) / интервал
Количество трансферов в секунду.
- DISKBSIZE (delta sectors read + delta sectors written)/ delta ( reads completed successfully + writes completed)
Общее количество дисковых блоков, прочитанных и записанных за интервал.
### Kernel
По ядру nmon берет информацию из procfs /proc/stat и преобразует ее
- Runnable (/proc/stat procs_running)
Число работоспособных процессов, готовых к запуску в секунду в runочереди. Очередь runподдерживается планировщиком процессов и содержит список потоков, готовых к отправке.
- Swap-in (/proc/stat procs_blocked)
Длина очереди swapв секунду, что означает количество готовых процессов, ожидающих выгрузки в секунду. Очередь подкачки содержит список процессов, которые готовы к запуску, но заменены запущенными в данный момент процессами.
- Pswitch (/proc/stat delta ctxt / интервал)
Количество переключений контекста процесса в секунду.
- Fork (/proc/stat delta processes / интервал)
Количество forkсистемных вызовов, выполняемых в секунду.

### VM
По виртуальной памяти nmon берет информацию из procfs /proc/vmstat и преобразует ее
Вывод в режиме записи:
- Paging and Virtual Memory
- nr_dirty
- nr_writeback
- nr_unstable
- nr_page_table_pages
- nr_mapped
- nr_slab_reclaimable
- pgpgin
- pgpgout
- pswpin
- pswpout
- pgfree
- pgactivate
- pgdeactivate
- pgfault
- pgmajfault
- pginodesteal
- slabs_scanned
- kswapd_steal
- kswapd_inodesteal
- pageoutrun
- allocstall
- pgrotated
- pgalloc_high
- pgalloc_normal
- pgalloc_dma
- pgrefill_high
- pgrefill_normal
- pgrefill_dma
- pgsteal_high
- pgsteal_normal
- pgsteal_dma
- pgscan_kswapd_high
- pgscan_kswapd_normal
- pgscan_kswapd_dma-
- pgscan_direct_high
- pgscan_direct_normal
- pgscan_direct_dma
### Network
Вывод в режиме записи:
По сетям nmon берет информацию из procfs /proc/net/dev и преобразует ее
NET,Network I/O, [имя устройства], lo-read-KB/s,eth0-read-KB/s,lo-write-KB/s,eth0-write-KB/s
NETPACKET,Network Packets, [имя устройства], lo-read/s,eth0-read/s,lo-write/s,eth0-write/s
### JFS 
По ядру nmon берет информацию из /etc/mtab (хотя можно и из procfs /proc/mounts) и преобразует ее
Вывод в режиме записи:
- JFS Filespace
- %Used

## Вывод
Исходя из анализа nmon, можно сделать вывод, что все метрики nmon собираются по procfs, следовательно nmon можно полностью исключить из-за его чрезмерной простоты и ненадобности.
