# galaxy
hakaton web server

##

* Тематика - космос
* Игра пошаговая
* Игрок управляет флотом, можно создавать новые корабли
* Поле игры разбито на клетки, на клетках генерируется и накапливается энергия

## Ход игры
* Сервер присылает игроку список кораблей и карту для каждого из кораблей
* Игрок делает ход для своих кораблей и отправляет серверу
* Сервер собирает ходы всех игроков
* Сервер обрабатывает ход
* Если два корабля на одной клетке - считается бой
 
## Параметры корабля
* *У*ровень - максимальное количество энергии
* *Р*адар - дальность обзора
* *Д*вигатели - дальность полета
* *О*ружие - надение
* *Б*роня - защита
 
## Бой
Для обоих кораблей считается нанесенный урон по формуле Уа * (Оа - Бб). Нанесенный урон вычитается из энергии корабля. Корабли с отрицательной энергией в конце боя уничтожаются. Уничтоженный корабль оставляет на своём месте энергию.


