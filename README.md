## Возможности контекста в языке Go
Этот проект посвящён исследыванию стандартной библиотеки `context` 
и возможностей для её применения.
Контекст обычно создаётся в какой-то родительской функции и передаётся в 
дочернюю. Он имеет смысл при параллельном выполнении. Например, когда 
дочерняя функция - отдельная горутина или вообще отдельный сервис, к которому
посылается запрос.

### Создание контекстов (уровень родителя)
Контексты создаются в вызывающей функции/процессоре и бывают разные 
по функционалу. Ниже примеры их создания:
#### 1. `context.Background() Context`
Это просто пустой контекст. Он используется на самом высоком уровне 
(то есть для передачи данных из `main` или обработчике высокого уровня).

#### 2. `context.TODO() Context`
Аналогичен предыдущему, но используется как заглушка, если пока что мы не 
уверены, какой контекст использовать. Это полезно, например, при выявлении
ошибок через CI/CD.
***
При многократной вложенности обычно контекст может получать дополнительный
функционал. Для этого используется наследование контекста, т.е. к родительскому
контексту добавляются дополнительные свойства. Следующие функции создают 
копию родительского контекста с добавлением различных функций:

#### 3. `context.WithValue(parent Context, key, val interface{}) (Context, CancelFunc)`
По сути, создаёт производный от родительского контекст, где к родительскому
добавлена передаваемая пара ключ-значение.

#### 4. `context.WithCancel(parent Context) (Context, CancelFunc)`
Возвращает производный контекст и **функцию отмены**. Когда выполняется 
функция отмены, все функции с данным контекстом должны завершиться. 
Корректно работать с функцией отмены в той функции, где она была создана.

#### 5. `context.WithDeadline(parent Context, d time.Time) (Context, CancelFunc)`
Аналогично п.4, кроме функции отмены, позволяет установить 
дедлайн для отмены контекста. 

#### 6. `context.WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)`
Здесь задаётся не дедлайн, а таймаут отмены. Т.е. использование `WithDeadline`
задаёт конкретное время прерывания, а `WithTimeout` - время с текущего момента
до прерывания. В остальном аналогично п.5.
