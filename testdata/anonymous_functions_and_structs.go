package main

import "log"

var fn = func(s string,
	i int,
	x float64) (int, error) {
	return 0, nil
}

type Function func(s string,
	j int,
	y float64) (int, error)

func run(name string,
	fn func(string,
		int) (int, error)) (string, int, error) {

	log.Println(fn(name, 1))
	return "", 10, nil
}

func Test() {
	run("quite a long name",
		func(s string, i int) (int, error) {
			log.Println("foo")
			log.Println("bar")
			log.Println("buzz")

			return 1, nil
		})
}
