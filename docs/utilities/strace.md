# Утилита strace 
– инструмент для «трассировки системных вызовов». 

Данное средство широко применяется в отладке развертывания (имеется ввиду взаимодействие программы и среды/операционной системы)

В операционной системе и используемых в ней программах иногда возникают ошибки, причину которых очень сложно понять, анализируя файлы 
журналов и сообщения об ошибках. Для таких ситуаций в Linux есть strace. За процессом работы любой из программ можно проследить, наблюдая системные вызовы, которые использует программа.

С помощью системных вызовов можно понять, к каким файлам обращается программа, какие сетевые порты она использует, какие ресурсы ей нужны, а также какие ошибки возвращает ей система. 
<br>Команда ***strace*** показывает все системные вызовы программы, которые та отправляет к системе во время выполнения, а также их параметры и результат выполнения. 
Но при необходимости можно подключиться и к уже запущенному процессу.
<br>В самом простом варианте ***strace*** запускает переданную команду с её аргументами и выводит в стандартный поток ошибок все системные вызовы команды. 

## Опции ##
Опции утилиты, с помощью которых можно управлять её поведением:

|Ключь|Описание|
|:---:|---|
|-i |- выводить указатель на инструкцию во время выполнения системного вызова;|
|-k |- выводить стек вызовов для отслеживаемого процесса после каждого системного вызова;|
|-o |- выводить всю информацию о системных вызовах не в стандартный поток ошибок, а в файл;|
|-q |- не выводить сообщения о подключении о отключении от процесса;|
|-qq |- не выводить сообщения о завершении работы процесса;|
|-r |- выводить временную метку для каждого системного вызова;|
|-s |- указать максимальный размер выводимой строки, по умолчанию 32;|
|-t |- выводить время суток для каждого вызова;|
|-tt |- добавить микросекунды;|
|-ttt |- добавить микросекунды и количество секунд после начала эпохи Unix;|
|-T |- выводить длительность выполнения системного вызова;|
|-x |- выводить все не ASCI-строки в шестнадцатеричном виде;|
|-xx |- выводить все строки в шестнадцатеричном виде;|
|-y |- выводить пути для файловых дескрипторов;|
|-yy |- выводить информацию о протоколе для файловых дескрипторов;|
|-c |- подсчитывать количество ошибок, вызовов и время выполнения для каждого системного вызова;|
|-O |- добавить определённое количество микросекунд к счетчику времени для каждого вызова;|
|-S |- сортировать информацию выводимую при опции -c. Доступны поля time, calls, name и nothing. По умолчанию используется time;|
|-w |- суммировать время между началом и завершением системного вызова;|
|-e |- позволяет отфильтровать только нужные системные вызовы или события;|
|-P |- отслеживать только системные вызовы, которые касаются указанного пути;|
|-v |- позволяет выводить дополнительную информацию, такую как версии окружения, статистику и так далее;|
|-b |- если указанный системный вызов обнаружен, трассировка прекращается;|
|-f |- отслеживать также дочерние процессы, если они будут созданы;|
|-ff |- если задана опция -o, то для каждого дочернего процесса будет создан отдельный файл с именем имя_файла.pid.|
|-I |- позволяет блокировать реакцию на нажатия Ctrl+C и Ctrl+Z;|
|-E |- добавляет переменную окружения для запускаемой программы;|
|-p |- указывает pid процесса, к которому следует подключиться;|
|-u |- запустить программу, от имени указанного пользователя.|

Для эффективной работы с утилитой необходимо знать основные системные вызовы. Условно их можно разделить на пять типов: Системные вызовы управления процессами, 
Системные вызовы управления файлами, Системные вызовы управления устройствами, Системные вызовы управления сетью и Системные вызовы системной информации. 
В рамках работы нас будут интересовать следующие:
*	Системные вызовы управления процессами − эти системные вызовы используются для управления процессами, такими как запуск новых, остановка существующих и ожидание их завершения. Fork(), exec(), wait() и exit() - все это примеры системных вызовов управления процессами ().
*	Системные вызовы управления файлами − эти вызовы системы используются для открытия, чтения, записи и закрытия документов, а также для их создания, переименования и удаления. Некоторые системные вызовы управления файлами () - это open(), read(), write(), close(), mkdir() и rmdir().
*	Системные вызовы управления сетью − Эти системные вызовы используются для управления сетевыми ресурсами, такими как подключение и отключение от сетей, отправка и получение данных по сетям и разрешение сетевых адресов. Socket(), connect(), send() и recv() являются примерами системных вызовов сетевого управления ().

## Синтаксис вывода ##
```
имя_системного_вызова (параметр1, параметр2) = результат сообщение
```
Имя системного вызова указывает, какой именно вызов использовала программа. Для большинства вызовов характерно то, 
что им нужно передавать параметры, имена файлов, данные и так далее. Эти параметры передаются в скобках. Далее идет знак равенства и результат выполнения. 
Если всё прошло успешно, то здесь будет ноль или положительное число. Если же возвращается отрицательное значение, делаем вывод, что произошла ошибка. В таком случае выводится сообщение.

## Предоставляемые возможности ##
1. Фильтрация системных вызовов.
   <br>Утилита выводит слишком много данных, которые зачастую нас не интересуют. С помощью опции ***-e*** можно применять различные фильтры для более удобного поиска проблемы.
   Например: отобразить только вызовы stat, передав в опцию ***-e*** такой параметр trace=stat:
   ```
   strace -e trace=stat {process}
   ```
   Кроме непосредственно системных вызовов, в качестве параметра для trace можно передавать и такие значения:
   *	file - все системные вызовы, которые касаются файлов;
   *	process - управление процессами;
   *	network - сетевые системные вызовы;
   *	signal - системные вызовы, что касаются сигналов;
   *	ipc - системные вызовы IPC;
   *	desc - управление дескрипторами файлов;
   *	memory - работа с памятью программы.
2. Подключение к запущенной программе.
   <br>Если программа, которую нам надо отследить, уже запущена, то можно подключиться к ней по ее идентификатору PID.  
   ```
   strace -o {file} -p {process PID}
   ```
3. Пути к файлам вместо дескрипторов
   <br>По умолчанию отображается только файловый дескриптор, чтобы узнать имена конкретных файлов, к которым обращается программа используем флаг ***-y***:
   ```
   strace -y -e trace=write,read -o {file} {process}
   ```
4. Фильтрация по пути
   <br>Для вывода всех системных вызовов, связанных только с определенным, можно выполнить фильтрацию по нему с помощью опции ***-P***.
   ```
   strace -P {file} {process}
   ```
5. Стистика системных вызовов
   <br>С помощью опции ***-с*** можно собрать статистику для системных вызовов, которые использует программа.
   ```
   strace -с {process}
   ```
   Во время работы утилита ничего выводить не будет. Результат будет рассчитан и выведен после завершения отладки. В выводе будут:
   *	time - процент времени от общего времени выполнения системных вызовов;
   *	seconds - общее количество секунд, затраченное на выполнение системных вызовов этого типа;
   *	calls - количество обращений к вызову;
   *	errors - количество ошибок;
   *	syscall - имя системного вызова.
   Для получения информации в режиме реального времени, используют опцию ***-C***.
6. Можно убирать системные вызовы — например, связанные с выделением и освобождением памяти(***/!***):
    ```
   strace -e trace=\!brk,mmap,mprotect,munmap -o {file} {process}
   ```
7. Отслеживание дочерних процессов
   <br>По умолчанию выводятся только системные вызовы родительского процесса. Отслеживать дерево процессов целиком помогает флаг ***-f***,
   с которым strace отслеживает системные вызовы в процессах-потомках. К каждой строке вывода при этом добавляется PID процесса, делающего системный вывод.
   ```
   strace -f -e trace=%process -o {file} {process}
   ```
   В этом контексте может пригодиться фильтрация по группам системных вызовов.
8. Многопоточные программы
   <br>Отслеживание информации в многопоточных программах. Флаг ***-f***, как и в случае с обычными процессами, добавит в начало каждой строки PID процесса.

   Речь идёт не об идентификаторе потока в смысле реализации стандарта POSIX Threads, а о номере, используемом планировщиком задач в Linux.
   С точки зрения последнего нет никаких процессов и потоков — есть задачи, которые надо распределить по доступным ядрам машины.
   <br>При работе в несколько потоков системных вызовов становится слишком много. Имеет смысл ограничиться только управлением процессами и системным вызовом write:
   ```
   strace -f -e trace=%process,write -o {file} {process}
   ```
   Флаг ***-ff*** если задана опция ***-o***, то разделяет все результаты трассировки по отдельным файлам с именем file_name.pid)
   
   Отличия будут хорошо просматриваться в параметрах для системного вызова clone (для создания нового потока (в случае многопоточных программ) и создании нового процесса (для дочерних процессов))
9. Межсетевое взаимодействие
   <br>Работа с системными вызовами управления сетью дает возможность получения ряда полезной информации
   ```
   strace -e trace=socket,connect,send,recv -o {file} {process}
   ```

   socket(AF_UNIX, SOCK_STREAM, 0) =sockfd;
   <br>параметры: 
   * Домен соединения (локальное соединение, IPv4 протоколы Интернет, IPv6 протоколы Интернет, IPX, устройство для взаимодействия с ядром и тд)
   * Семантика коммуникации (SOCK_STREAM двусторонний надежный и последовательный поток байтов, SOCK_DGRAM поддерживает датаграммы, и тд)

   connect() - устанавливает соединение с конечной точкой внешней сети.
   Путь, по которому сокет связан с дескриптором.

   Понимание того, что процесс работает с чем-то через сокет по определенному пути. Таким образом через, например, утилиту ss, можно найти процесс, привязанный к этому сокету с другой стороны.

   Утилита позволяет установить использует ли процесс данное компьютерное соединение для передачи данных и/или получения и что передается.

### Стоит обратить внимание ###
Если выполняется системный вызов, в то время как другой вызов вызывается из другого потока/процесса, то strace попытается сохранить порядок этих событий 
и пометить текущий вызов как незавершённый. Когда вызов вернется, он будет помечен как возобновлённый.