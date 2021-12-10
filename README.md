# Todo-list - сервис на Go

### В процессе работы использованы:

- Вёрстка по стандарту HTML5, стили CSS3(Flexbox, Grid), "резиновая", адаптивная верстка, БЭМ
- JavaScript ES6
- Webpack, GIT

### Функционал проекта:

- всплывающие подсказки при наведении на задачу
- в Backlog попадают задачи, за которыми не закреплен исполнитель
- прокручивание недели вперед и назад
- Drag and Drop для задач из Backlog
- возможность перетаскивания задач как на конкретную дату в строке, так и на исполнителя
- подсветка просроченных задач: задачи без отметки о завершении становятся красного цвета начиная с момента просрочки по сегодняшний день
- возможность поиска задач в Backlog
- текущая дата выделена серым цветом

### Что планируется улучшить:

- использование дополнительного свойства - затрачиваемое время на задачу
- реализовать всплывающие подсказки на графике, указывающие на дату завершения задачи
- работа с LocalStorage
- использовать SCSS

### Развёртывание:

- Клонировать репозиторий командой

```
git clone https://github.com/nikolaevfo/elma-test-work.git
```

- Установить сторонние библиотеки командой npm i
- Запустить вертуальный сервер командой npm run dev
- Для финальной сборки введите команду npm run build
- Для работы с проектом у вас должен быть установлен Node.js - https://nodejs.org/en/download/
