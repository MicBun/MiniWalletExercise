package main

import "miniWalletExercise/route"

func main() {
	r := route.SetupRouter()
	r.Run(":80")
}
