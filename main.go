package main

import app "github.com/AbyssExplorer/Cogent/cmd/cogent"

func main() {
	if err := app.Execute(); err != nil {
		panic(err)
	}
}
