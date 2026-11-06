package main

import (
	"fmt"
	"os"

	"github.com/ranas-mukminov/retrace/internal/index"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "index":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: retrace index <directory>")
			os.Exit(1)
		}
		if err := index.IndexDir(os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "index:", err)
			os.Exit(1)
		}
	case "search":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: retrace search <query>")
			os.Exit(1)
		}
		results, err := index.Search(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "search:", err)
			os.Exit(1)
		}
		for _, r := range results {
			fmt.Printf("%s\n  %s\n", r.Path, r.Snippet)
		}
		if len(results) == 0 {
			fmt.Println("No matches.")
		}
	case "status":
		n, err := index.Count()
		if err != nil {
			fmt.Fprintln(os.Stderr, "status:", err)
			os.Exit(1)
		}
		fmt.Printf("%d documents indexed\n", n)
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `retrace — where did I see that?

Usage:
  retrace index <directory>   index text files
  retrace search <query>      search your index
  retrace status              show indexed document count`)
}
