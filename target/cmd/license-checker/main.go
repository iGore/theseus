package main

import (
	"context"
	"os"
	"theseus-target/license-checker/internal/app"
)

func main() { os.Exit(app.Run(context.Background(), app.DefaultInvocation(os.Args[1:]))) }
