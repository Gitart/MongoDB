

Запуск в командной строке и запустить MongoDB:

```
mongo
```

Просмотр баз данных

```
show dbs
```

База данных `test`. 
В ней мы создадим коллекцию `shows` и поместим туда несколько записей в формате JSON. 
Создадание пременных

```
a = { title:"Arrested Development", airdate:"November 2, 2003", network:"FOX" }
b = { title:"Stella", airdate:"June 28, 2005", network:"Comedy Central" }
c = { title:"Modern Family", airdate:"September 23, 2009", network:"ABC" }

db.shows.save(a)
db.shows.save(b)
db.shows.save(c)
```

В результате выполнения последних команд, 
коллекция shows будет заполнена тремя JSON строками. 
Для просмотра всех коллекций можно воспользоваться командой `show collections`.

Просмотр данных коллекции - методо `find()`:

```
db.shows.find()
```

