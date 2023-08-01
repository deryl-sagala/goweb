# Welcome To Goweb Docs

## Installation
* `go get https://github.com/DerylDarren/goweb`

## Import To Project
* `import "https://github.com/DerylDarren/goweb"`

## Features
### Simple Hello World
    import "https://github.com/DerylDarren/goweb"

    func helloWorld() {
        web.returnText("Hello World)
    }
    
    web.addRoute("/hello-world", web.wrap(helloWorld))