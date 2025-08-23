package main

import (
	"log"
	"main/src"
)

func main() {
	scriptName, scriptArgs, err := src.PrepareContext()
	if err != nil {
		log.Fatal(err)
	}

	script, ok := src.CurrentContext.Chart.Scripts[scriptName]
	if !ok {
		log.Fatalf("Script %s not found in chart", scriptName)
	}

	scriptContent, ok := script.GetScript(src.CurrentContext.Runner)
	if !ok {
		log.Fatalf("Script %s not found for runner %s", scriptName, src.CurrentContext.Runner)
	}

	err = src.CurrentContext.Runner.Run(scriptName, scriptContent, scriptArgs)
	if err != nil {
		log.Fatal(err)
	}
}
