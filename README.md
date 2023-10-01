
[![GoDoc](https://pkg.go.dev/badge/github.com/gin-gonic/gin?status.svg)](https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc)
[![Release](https://img.shields.io/github/release/gin-gonic/gin.svg?style=flat-square)](https://github.com/gin-gonic/gin/releases)

# GO-WEATHER
Простой сервис по выводу информации о погоде.

Написан на языке `Golang-1.20.3` с использованием библиотеки `GIN-1.9.1`.


## START
Сборка проетка:
```bash
go get
```

Запуск:
```bash
go run main.go
```


## Example


### Input:
in `bash`:
```bash
curl "localhost:8080/weather?city=Hamburg"
```


### Output:   

in `service`: 
```
[GIN] 2023/10/01 - 10:18:42 | 404 |       1.548µs |       127.0.0.1 | GET      "/"
[GIN] 2023/10/01 - 10:18:54 | 200 |  500.250509ms |       127.0.0.1 | GET      "/weather?city=Hamburg"
[GIN] 2023/10/01 - 10:19:01 | 200 |  138.697694ms |       127.0.0.1 | GET      "/weather?city=Ufa"
```

in `bash`:

```bash
curl "localhost:3000/weather?city=Ufa"
```
```json
{
    "OK":
    {
        "City":"Ufa",
        "Forecasts":
        [
            {
                "Date":"Sun 00:00",
                "Temperature":"7.1°C"
                },
                {
                    "Date":"Sun 01:00",
                    "Temperature":"7.0°C"
                },
                {
                    "Date":"Sun 02:00",
                    "Temperature":"6.6°C"
                }
            }
        ]
    }
}
```