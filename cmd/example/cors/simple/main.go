package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title></title>
        <link href="css/style.css" rel="stylesheet">
    </head>
    <body>
        <h1>Simple CORS</h1>
        <div id="output"></div>
        <script>
            document.addEventListener('DOMContentLoaded', function() {
                // fetch 发起跨域请求
                fetch('http://localhost:4000/v1/healthcheck')
                    .then(function (response) {
                        response.text().then(function (text) {
                            document.getElementById('output').innerHTML = text;
                        });
                    }, function (error) {
                        document.getElementById('output').innerHTML = error;
                    });
            });
        </script>
    
    </body>
</html>
`

func main() {
    addr := flag.String("addr", ":9000", "HTTP network address")
    flag.Parse()

    log.Printf("Server starting on %s", *addr)

    err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(html))
    }))
    log.Fatal(err)
}
