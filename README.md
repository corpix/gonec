[![GitHub issues](https://img.shields.io/github/issues/covrom/gonec.svg)](https://github.com/covrom/gonec/issues) [![Travis](https://travis-ci.org/covrom/gonec.svg?branch=master)](https://github.com/covrom/gonec/releases)

# ГОНЕЦ (gonec)
## Интерпретатор и платформа создания микросервисов на 1С-подобном языке

[![Gonec Logo](/gonec.png)](https://github.com/covrom/gonec/releases)

[![Download](/button_down.png)](https://github.com/covrom/gonec/releases)
[![Docs](/button_doc.png)](https://github.com/covrom/gonec/wiki)
[![Demo site](/button_play.png)](https://gonec.herokuapp.com/)

[Чат в гиттере](https://gitter.im/gonec/Lobby)

## Цели

Интерпретатор создан для решения программистами 1С множества задач, связанных с высокопроизводительными распределенными вычислениями, создания вэб-сервисов и вэб-порталов для работы тысяч пользователей, работы с высокоэффективными базами данных с использованием синтаксиса языка, похожего, но не ограниченного возможностями языка 1С.

Включив такой интерпретатор в свое решение, Вы можете предоставить высокий уровень сервиса для своих клиентов, который обгонит решения не только ваших конкурентов на рынке 1С, но и конкурентных платформ в enterprise.

Интерпретатор разрабатывается “от простого к сложному”. На начальных этапах будет включена базовая функциональность многопоточных вычислений и сетевых сервисов. В перспективе планируется организация работы с различными базами данных и визуализация управляемых форм, созданных в конфигураторе.

Еще никогда не были так просто доступны программистам 1С возможности:
* Создать микросервис с произвольным сетевым протоколом, развернуть его на linux, в docker контейнере или кластере kubernetes
* Выполнить сложную многопоточную вычислительную задачу для десятков тысяч подключающихся пользователей за миллисекунды
* Взаимодействовать с пользователем через web-браузер с минимальным трафиком
* Сохранять и получать данные с максимально доступной скоростью в key-value базах данных

## Описание синтаксиса языка и примеры использования интерпретатора

[Документация находится здесь](https://github.com/covrom/gonec/wiki)

## Почему синтаксис похож на 1С?

Синтаксис 1С знаком и удобен сотням тысяч программистов в России и СНГ, а в перспективе и зарубежом. Это позволяет создавать решения, которые могут поддерживаться любыми программистами 1С, и которые не будут требовать дополнительной квалификации.

## Какие основные отличия от языка 1С?
Язык интерпретатора поддерживает синтаксис языка 1С, но дополнительно к этому имеет возможности, унаследованные от синтаксиса языка Go (Golang) и Javascript:
* Многопоточное программирование: создание и работа с параллельно выполняемыми функциями и каналами передачи информации между ними (полный аналог chan и go из языка Го)
* Поддержка срезов массивов и строк, как в языке Python (высокоскоростная реализация на слайсах Go)
* Поддержка множества возвращаемых значений из функций и множественного присваивания между переменными (а,б=б,а) как в языке Go
* Возможность указания структурных литералов и содержимого массивов прямо в исходном коде (как в Go)
* Передача функций как параметров (функциональное программирование)

## Архитектура языка и платформы

Архитектура платформы
![System architecture](/architecture.png)

Архитектура языка Гонец
![Language architecture](/langarch.png)

## Масштабируемость языка и платформы
Язык Гонец расширяется путем изменения правил синтаксиса в формате YACC, а так же написания собственных высокоэффективных библиотек структур и функций на Го, которые могут быть доступны как объекты метаданных в языке Гонец.

В системных библиотеках языка могут создаваться объекты метаданных, при импорте соответствующей библиотеки. Данные объекты метаданных являются функциональными структурными типами, функциональность которых скомпилирована на языке Го. Объекты таких метаданных имеют встроенные методы для работы с ними. Добавление нового типа с новой функциональностью не представляет никакой сложности для программистов на Го. Компиляция усовершенствованной версии интерпретатора выполняется одной командой `go build .` в папке с исходными текстами, и занимает всего несколько секунд.

Интерпретатор может и запускать микросервисы, и сам выступать в роли такого микросервиса.

Посмотреть на использование интерпретатора в роли микросервиса можно по [ссылке](https://gonec.herokuapp.com/) выше.
В этой реализации в интерпретатор встроена простая система запуска кода через обычный браузер, которая работает на технологии ajax, общающаяся с микросервисом сессий исполнения кода интерпретатором.

## Каковы отличия в метаданных и среде исполнения?
На первом этапе разработки планируется:
* поддержка стандартной библиотеки Go в части создания вэб-сервисов с html-шаблонами
* поддержка и работа с postgresql, включая hstore
* мултиплатформенность: выполнение кода на любой платформе (Windows, Linux, MacOs)
* выполнение в легковесных ОС alpinelinux
* запуск на встраиваемых устройствах с низким энергопотреблением класса "умный дом"/"интернет вещей" (например, Raspberry Pi)
* запуск в контейнерах docker

## Какова производительность интерпретатора?
Производительность ожидается сравнимой или выше, чем у интерпретатора языка Python.
Скорость интерпретации кода соответствует скорости программ на Go и скорости работы библиотек, написанных на Go.

## Какая технологическая основа используется в интерпретаторе?
Интерптетатор реализован на языке Go путем адаптации исходных кодов интерпретатора языка anko (https://github.com/mattn/anko).
Интерпретатор использует собственную виртуальную машину, также написанную на языке Go, а значит, имеет отличную производительность и стабильность.

## Какой статус разработки интерпретатора?
Интерпретатор находится в стадии разработки стандартной библиотеки.
Первая версия уже доступна к тестированию!
