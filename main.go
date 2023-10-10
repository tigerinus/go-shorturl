package main

import (
	"context"
	"crypto/rand"
	"fmt"
)

const dict = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var nextShort = make(chan string, 1)

var short2long = map[string]string{}

func randomString(n int) string {
    buf := make([]byte, n)
    if _, err := rand.Read(buf); err != nil {
        panic(err)
    }
    
    for i, b := range buf {
        buf[i] = dict[int(b) % len(dict)]
    }

    return string(buf)
}

func generateShort(ctx context.Context) {
    for {
        select {
        case <- ctx.Done():
            break
        default:
            newShort := randomString(7)
            if _, ok := short2long[newShort]; ok {
                continue
            }
             
            nextShort <- newShort
        }
    }
}

func createShortURL(url string) string {
    shortUrl := <- nextShort

    short2long[shortUrl] = url

    return shortUrl
}

func getLongUrl(shortUrl string) string {
    return short2long[shortUrl]
}

func main() {

    backgroundCtx, cancel := context.WithCancel(context.Background())

    defer cancel()

    go generateShort(backgroundCtx)

    shortUrl := createShortURL("google.com")
    fmt.Printf("shortUrl:\t%s\n", shortUrl)

    longUrl := getLongUrl(shortUrl)
    fmt.Printf("longUrl:\t%s\n", longUrl)

}
