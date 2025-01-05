# Go Home

## Домен для Go пакетов

[gohome.4gophers.ru](https://gohome.4gophers.ru/)

Очень легко устанавливать Go пакеты с гитхаба. Но все сложнее если библиотека лежит на [gitflic.ru](https://gitflic.ru). Тут обычный go get не сработает.

Для этого есть Go Home. С его помощью пакеты можно устанавливать через go get.

Для начала нужно создать пакет, название которого будет начинаться с [gohome.4gophers.ru](https://gohome.4gophers.ru/):

```
go mod init gohome.4gophers.ru/kovardin/example
```

После этого создать репозиторий на gitflic.ru:

```
https://gitflic.ru/project/kovardin/example
```

Теперь вы можете указывать в своем коде зависимость

```go
package main

import (
	"fmt"

	"gohome.4gophers.ru/kovardin/example"
)

func main() {
	fmt.Println(example.Hello("Artem"))
}
```

Go разрулит зависимости через meta теги на этом сайте и установит зависимость.