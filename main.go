package main

func main() {
	a := App{}
	a.Port = ":2083"
	a.Initialize()
	a.Run()
}
