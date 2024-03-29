Утилита **glances**.

Glances — это кроссплатформенный инструмент мониторинга системы на основе командной строки curses, написанный на языке Python, который использует библиотеку psutil для получения информации из системы. Он также может работать в режиме клиент/сервер. Удаленный мониторинг может осуществляться через терминал, веб-интерфейс или API (XMLRPC и RESTful).

С помощью glances мы можем отслеживать процессор, среднюю загрузку, память, сетевые интерфейсы, дисковый ввод-вывод, процессы и использование пространства файловой системы.

Особенности glances:

- Информация о процессоре (связанные с пользователем приложения, основные программы системы и неактивные программы.
- Информация об общей памяти, включая оперативную память, подкачку, свободную память и т. д.
- Средняя загрузка процессора за последние 1 минуту, 5 минут и 15 минут.
- Скорость загрузки сетевых подключений.
- Общее количество процессов, активных, спящих процессов и т.д.
- Сведения о скорости дискового ввода-вывода (чтения или записи), относящиеся к диску.
- Использование диска подключенными в данный момент устройствами.
- top-овые процессы с указанием их использования процессора/памяти, имен и местоположения приложения.
- Выделяет красным цветом процессы, которые потребляют больше всего системных ресурсов.

Основы glances.

Если запустить команду glancesбез аргументов командной строки или опций, начнется мониторинг локального компьютера. Чтобы закончить мониторинг, нужно нажать 'q' или 'ESC'.

Варианты вывода:

- Отображать необработанную (RAW) статистику (Python) непосредственно в стандартный вывод:

**glances --stdout cpu.user,mem.used,load**

- Отображать в формате CSV благодаря опции stdout-csv:

**glances --stdout-csv now,cpu.user,mem.used,load**

- Отображать в формате JSON благодаря опции stdout-json (атрибут не поддерживается в этом режиме):

**glances --stdout-json cpu,mem**

Он будет отображать по одной строке для каждой характеристики при каждом обновлении.

Режим Клиент/сервер:

Чтобы удаленно контролировать машину, называемую server, с другой машины, называемой client, запустить на сервере:

**server$ glances -s**

инаклиенте:

**client$ glances -c @server**

где @server - IP-адрес или имя хоста сервера.

В режиме сервера можно задать bind адрес с помощью -B ADDRESS, а прослушивающий TCP-порт - с помощью -p PORT. В режиме клиента можно задать TCP-порт сервера с помощью -p PORT. Адрес привязки по умолчанию - 0.0.0.0 (Glances будет прослушивать все доступные сетевые интерфейсы), а TCP-порт - 61209.

В режиме клиент/сервер ограничения устанавливаются на стороне сервера.

Режим веб-сервера:

Чтобы удаленно контролировать машину, называемую server, с любого устройства с веб-браузером, запустить сервер с параметром -w:

**server$ glances -w**

затем на клиенте ввести следующий URL-адрес в веб-браузере:

**http://@server:61208**

Чтобы изменить частоту обновления страницы, нужно добавить период в секундах в конце URL.

Опции:

- **-h, --help** : отобразить сводку опций и выйти.
- **-V, --version**:вывод информации о версии и выход.
- **-d, --debug** : включить режим отладки.
- **-C CONF\_FILE, --config CONF \_FILE**: путь к конфигурационному файлу.
- **--modules-list** : отображение списка модулей (плагинов и экспорта) и выход.
- **--disable-plugin PLUGIN** : отключить PLUGIN (список, разделенный запятыми).
- **--enable-plugin PLUGIN** : включить PLUGIN (список, разделенный запятыми).
- **--stdout PLUGINS\_STATS** :

вывод статистики в стандартный вывод (список плагинов/plugins.attribute через запятую).

- **--export EXPORT** : включить модуль EXPORT (список, разделенный запятыми).
- **--export-csv-file EXPORT \_ CSV\_FILE** : путь к файлу для экспортера CSV.
- **--export-json-file EXPORT\_JSON\_FILE** : путькфайлудляэкспортера JSON.
- **--disable-process**:

отключить модуль обработки (уменьшите потребление Glances CPU).

- **--disable-webui** : отключить веб-интерфейс (только RESTful API будет отвечать).
- **--light, --enable-light** :

легкий режим для пользовательского интерфейса Curses (отключить все, кроме верхнего меню).

- **-0, --disable-irix** :

загрузка процессора задачей будет разделена на общее количество CPU.

- **-1, --percpu** : запустить Glances в режиме для каждого процессора.
- **-2, --disable-left-sidebar** : отключить сетевые, дисковые I/O, FS и сенсорные модули.
- **-3, --disable-quicklook** : отключить модуль быстрого просмотра.
- **-4, --full-quicklook** : отключите все, кроме быстрого просмотра и загрузки.
- **-5, --disable-top** : отключить верхнее меню (QuickLook, CPU, MEM, SWAP и LOAD).
- **-6, --meangpu** : запустить Glances в режиме mean GPU.
- **--enable-history** : включить режим истории.
- **--enable-process-extended** : включить расширенную статистику в главном процессе.
- **-c CLIENT, --client CLIENT** :

подключитесь к серверу Glances по адресу IPv4/IPv6, имени хоста или hostname:port.

- **-s, --server**: запустить Glances в режиме сервера.
- **--browser** : запустить клиентский браузер (список серверов).
- **--disable-autodiscover** : отключить функцию автообнаружения.
- **-p PORT, --port PORT** : определите TCP-порт клиента/сервера [по умолчанию: 61209].
- **-B BIND\_ADDRESS, --bind BIND\_ADDRESS** :

привязать сервер к заданному IPv4/IPv6 адресу или имени хоста.

- **--username** : определите имя пользователя клиента/сервера.
- **--password** : определите пароль клиента/сервера.
- **--snmp-community SNMP\_COMMUNITY** : сообщество SNMP.
- **--snmp-port SNMP\_PORT** : SNMP порт.
- **--snmp-version SNMP\_VERSION** : SNMP версия (1, 2сили 3).
- **--snmp-user SNMP\_USER** : имя пользователя SNMP (только для SNMPv3).
- **--snmp-auth SNMP\_AUTH** : ключ аутентификации SNMP (только для SNMPv3).
- **--snmp-force** : принудительный режим SNMP.
- **-t TIME, --time TIME** :

установите время обновления в секундах [по умолчанию: 3 секунды].

- **-w, --webserver** :

запустить Glances в режиме веб-сервера (требуется библиотека bottle).

- **--cached-time CACHED\_TIME** :

установите время кэширования сервера [по умолчанию: 1 сек].

- **--open-web-browser** :

попробовать открыть веб-интерфейс в веб-браузере по умолчанию **.**

- **-q, --quiet** : не отображать интерфейс curses.
- **-f PROCESS\_FILTER, --process-filter PROCESS\_FILTER** :

установите шаблон фильтра процесса (regular expression).

- **--process-short-name** : принудительное краткое имя для имени процесса.
- **--hide-kernel-threads** :

скрыть потоки ядра в списке процессов (недоступно в Windows) **.**

- **-b, --byte** : отображать скорость работы сети в байтах в секунду.
- **--diskio-show-ramfs** : показать RAM FS в плагине DiskIO.
- **--diskio-iops** :

показывать количество операций ввода-вывода в секунду в плагине DiskIO.

- **--fahrenheit** : отобразить температуру в градусах Фаренгейта.
- **--fs-free-space** : отображать FS свободное пространство вместо используемого.
- **--theme-white** : оптимизировать цвета отображения для белого фона.
- **--disable-check-update** :отключить онлайн проверку версии Glances.

Описание приложения:

Заголовок:

В заголовке указаны имя хоста, название операционной системы, версия выпуска, IP-адреса архитектуры платформы (частные и общедоступные) и время безотказной работы системы. Кроме того, в GNU/Linux также отображается версия ядра.

Quick Look:

Плагин quicklook отображается только на широком экране и предлагает просмотр панели для процессора и памяти (виртуальной и подкачки). В интерфейсе Curses/terminal также можно переключиться с bar на sparkline, используя горячую клавишу "S" или параметр командной строки –sparkline (в системе требуется библиотека sparklines Python lib).

CPU:

Статистика процессора отображается в процентах или значениях и за настроенное время обновления. Общая загрузка процессора отображается в первой строке.

Описание статистики процессора:

- user: процент времени, затраченного в пользовательском пространстве. Процессорное время пользователя — это время, затраченное процессором на выполнение кода вашей программы (или кода в библиотеках).
- system: процент времени, затраченного в пространстве ядра. Системное процессорное время — это время, затраченное на выполнение кода в ядре операционной системы.
- idle: процент использования ЦП любой программой. Если ЦП выполнил все задачи, он находится в режиме ожидания.
- nice (\*nix): процент времени, занятого процессами пользовательского уровня, с положительным значением nice. Время, затраченное процессором на запуск пользовательских процессов, которые были улучшены.
- irq (Linux, \*BSD): процент времени, затраченного на обслуживание/обработку аппаратных/программных прерываний.
- iowait (Linux): процент времени, затрачиваемого процессором на ожидание завершения операций ввода-вывода.
- steal (Linux): процент времени, в течение которого виртуальный процессор ожидает реального процессора, в то время как гипервизор обслуживает другой виртуальный процессор.
- ctx\_sw: количество переключений контекста (добровольных + непроизвольных) в секунду. Переключение контекста — это процедура, которой следует центральный процессор компьютера для перехода от одной задачи (или процесса) к другой, гарантируя, что задачи не конфликтуют.
- inter: количество прерываний в секунду.
- sw\_inter: количество программных прерываний в секунду.
- syscal: количество системных вызовов в секунду. Не отображается в Linux (всегда 0).
- dpc: (Windows): время, затраченное на обслуживание отложенных вызовов процедур.

GPU:

Статистика графического процессора отображается в процентах от значения и для настроенного времени обновления. Отображается: общее использование GPU, потребление памяти, температура (версии 3.1.4 или выше)

Память:

Glances использует два столбца: один для оперативной памяти и один для подкачки.

Описание статистики:

- percent: процент использования.
- total: общая доступная физическая память.
- used: используемая память, рассчитывается по-разному в зависимости от платформы и предназначена только для информационных целей. Вычисляется следующим образом: используемая память = полностью свободна (при этом свободна = доступно + буферы + кэширование)
- free: память, которая не используется вообще (обнуляется) и которая легко доступна; это не отражает фактическую доступную память.
- active (UNIX): память, используемая в данный момент или использовавшаяся совсем недавно, и поэтому она находится в оперативной памяти.
- inactive (UNIX): память, помеченная как неиспользуемая.
- buffers (Linux, BSD): кэш для таких вещей, как метаданные файловой системы.
- cached (Linux, BSD): кэш для различных целей.

Дополнительная статистика доступна через API:

- available: фактический объем доступной памяти, который может быть мгновенно предоставлен процессам, запрашивающим больше памяти в байтах; он вычисляется путем суммирования различных значений памяти в зависимости от платформы (например, свободная + буферы + кэшированная в Linux) и предполагается использовать для мониторинга фактического использования памяти в кроссплатформенной среде.
- wired (BSD, macOS): память, которая помечена как всегда остающаяся в оперативной памяти. Она никогда не перемещается на диск.
- shared (BSD): память, к которой могут одновременно обращаться несколько процессов.

Загрузка:

Вкратце, это средняя сумма количества процессов, ожидающих в очереди выполнения, плюс количество процессов, выполняющихся в данный момент в течение периодов времени 1, 5 и 15 минут.

Сеть:

Glances отображает скорость передачи данных сетевого интерфейса. Устройство настраивается динамически (бит/с, кбит/с, Мбит/с и т.д.). Если обнаружена скорость интерфейса (не во всех системах), применяются пороговые значения по умолчанию (70% для осторожности, 80% для предупреждения и 90% для критичности).

Соединения:

Этот плагин отображает расширенную информацию о сетевых подключениях.

Состояния:

- Listen: все порты, созданные сервером и ожидающие подключения клиента
- Initialized: Все состояния при инициализации соединения (сумма SYN\_SENT и SYN\_RECEIVED)
- Established: Все установленные соединения между клиентом и сервером
- Terminated: Все состояния при завершении соединения (FIN\_WAIT1, CLOSE\_WAIT, LAST\_ACK, FIN\_WAIT2, TIME\_WAIT и CLOSE)
- Tracked: Текущее количество и максимальное подключение к Netfilter tracker (nf\_conntrack\_count/nf\_conntrack\_max)

Wi-Fi:

Glances отображает названия точек доступа Wi-Fi и качество сигнала. Если Glances запущен от имени root, отображаются все доступные точки доступа.

Порты:

Этот плагин предназначен для предоставления списка хостов /портов и URL-адреса для сканирования.

Disk I/O:

Glances отображает пропускную способность дискового ввода-вывода. Модуль адаптируется динамически.

Можно отобразить:

- байт в секунду (поведение по умолчанию / Байт/с, Кбайт/с, Мбайт/с и т.д.)
- запросов в секунду (с использованием опции –diskio-iops или горячей клавиши B)

Файловая система:

Glances отображает используемое и общее пространство на диске файловой системы. Устройство адаптируется динамически.

Папки:

Плагин folders позволяет пользователю с помощью файла конфигурации отслеживать размер предопределенного списка папок.

Список процессов:

Представление процесса состоит из 3 частей:

- Сводка процессов
- Список отслеживаемых процессов (необязательно, только в автономном режиме)
- Расширенная статистика для выбранного процесса (необязательно)
- Список процессов

Отображается строка сводки процессов:

- Общее количество задач/процессов (псевдонимы как общее количество в Glances API)
- Количество потоков
- Количество запущенных задач/процессов
- Количество спящих задач/процессов
- Другое количество задач/процессов (не находящихся в запущенном или спящем состоянии)
- Ключ сортировки для списка процессов

Можно отфильтровать список процессов, используя клавишу ENTER.

Контейнеры:

Glances может отслеживать контейнеры Docker или Podman. Glances использует API контейнеров через библиотеки docker-py и podman-py.

Процесс мониторинга приложений:

Благодаря Glances и его модулю AMP можно добавить специальный мониторинг к запущенным процессам. Усилители определены в конфигурационном файле Glances.