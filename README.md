# Описание
Для практики работы с TCP мы напишем простую игру "Угадай число". Сервер будет запускаться на определенном порту указанном в флаге -port (дефолт 8080),
получать в виде флагов -min и -max (дефолт 0 и 100 соответственно) диапазон в котором он будет загадывать числа (включительно для min и max). После этого он будет ждать подключения клиента.
* Клиент подключается к серверу и отправляет ему значение `message.Start` в виде строки
* Сервер проверяет это значение (другое значение = ошибка) и отправляет клиенту диапазон чисел в виде значения `message.MinMax`
* Клиент получает диапазон и используя бинарный поиск начинает свои попытки угадать число отослав его в виде строки
* Сервер парсит ответ в число и сравнивает с загаданным числом. Если число больше, то отправляет клиенту `message.Lower`, если меньше, то `message.Higher`, если угадал, то `message.Correct`. В случае последнего закрывает канал, но не сам сервер
* Клиент получает ответ и в случае если это `message.Correct` завершает работу, иначе продолжает попытки

# Примечания
* Используйте весь диапазон `[min, max]`
* Используйте пакет `message` в котором указаны все сообщения нашего "протокола" и удобные функции для записи и чтения из коннекшена
* Краште клиент в main если что-то пошло не так
* Не краште сервер если что-то пошло не так. Одно подключение не должно убивать сервер на котором потенциально запущены другие игры
* Логгируйте на клиенте получение диапазона, свою догадку и результат полученный с сервера, а на сервере загаданное число (формат есть в коментариях)
* Для воспроизводимости игры используйте флаг `RAND_SEED`. В его отсутствие используйте `time.Now().UnixNano()`
* Для сканирования диапазона используйте `fmt.Sscanf` с шаблоном `message.MinMax`

# Тест
Запустите `make -s play`. В случае корректной работы вы увидете:
```
guessed -29
min: -100, max: 100
guessing 0
lower
guessing -50
higher
guessing -25
lower
guessing -37
higher
guessing -31
higher
guessing -28
lower
guessing -29
correct
Terminated
make[1]: *** [Makefile:4: run-server] Error 143
```