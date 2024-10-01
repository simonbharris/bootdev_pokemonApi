# This is a repo covering my recent work and learning of golang

This project is following the guidance of [boot.dev](https://www.boot.dev/) project "build a pokedex". 

Project is a REPL program that wraps around the public [Pokemon Api](https://pokeapi.co/)

No code in this project is AI Generated. 

## How to build and run
Assuming you have the [golang development environment](https://go.dev/doc/install) setup 

```sh
go build .\cmd\pokedex\; .\pokedex.exe
```

```sh
Pokedex > help
Usage:
help: Returns help information
exit: exits the pokedex
map: Returns a page of 20 locations, repeated calls return the next page.
mapb: Returns a prior page of 20 locations. repeated calls move return prior pages.
explore: Returns what pokemon can be found in specified area.
        Usage: explore area-name
```
