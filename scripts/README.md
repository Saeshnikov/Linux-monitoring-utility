# testing_overall_work.sh
Данный bash-скрипт предназначен для проверки корректности работы `prototype`.

## Описание ##
Общее: 
* Пользователю на систему устанавливается пакет с текстовым редакотором `emacs`
* Ожидаемые результаты:
  <br>После завершения работы `prototype`, в результирующих файлах не должны находиться пакеты, файл которых открываются в скрипте, но должен быть установленный пакет `emacs-27.2-150400.3.6.1.x86-64`.
          
Скрипт состоит из трех частей:
1. Первая (получение пакетов):
      - Получает количество пакетов, установленных на систему.
      - Определяет % пакетов от общего числа, которые не будут открываться bash-скриптом (20%).
      - Получает список все установленных на системе пакетов.
2. Вторая (проверка неиспользуемого пакета):
      - Скачивает из интернета и устанавливает пакет `emacs` (данный пакет не будет расположен в списке пакетов из предыдущего этапа, которые могут быть откыты скриптом).
      - Запускает ожидание на 10 минут, в течении которого необходимо запустить прототип.
3. Третья (проверка используемых пакетов):
      - Генерирует рандомное число в диапозоне, верхней границей которого является общее число пакетов на системе и берет соответствующий пакет из списка, полученного на первом этапе.
      - Использует `rpm -ql` для определения файлов, принадлежащих пакету и берет первый из них.
      - Открывает файл пакета, который точно есть на системе.
  
## Использование ##
```
$ chmod +x testing_overall_work.sh
$ ./testing_overall_work.sh
```
`!` Требования: 
<br>        - запускаемый `prototype` должен работать не меньше, чем работает `bash-скрипт`
<br>        - должен быть доступ к интернету для скачивания неиспользуемого пакета