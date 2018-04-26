package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func usage() {
	log.Printf(`Usage: %v [in.dat] [out.json]
	Default input/output files are "mod-settings.dat" and "mod-settings.json" in the current directory.
	"-" cad be used for any of arguments to use stdin/stdout.`, os.Args[0])
}

func main() {
	log.SetFlags(0)
	// Default output path
	outPath := "mod-settings.json"
	var in io.Reader

	switch len(os.Args) {
	case 1:
		file, err := os.Open("mod-settings.dat")
		switch {
		case err == nil:

		case os.IsNotExist(err):
			usage()
			os.Exit(1)

		default:
			log.Fatal(err)
		}

		defer file.Close()
		in = file

	case 3:
		outPath = os.Args[2]
		fallthrough

	case 2:
		switch os.Args[1] {
		case "-h", "--help":
			usage()
			return
		case "-":
			in = os.Stdin
		default:
			file, err := os.Open(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			in = file
		}

	default:
		usage()
		os.Exit(1)
	}

	var modData FModData
	modData.Decode(in)
	bytes, _ := json.MarshalIndent(modData.Data, "", "  ")

	if outPath == "-" {
		log.Println(string(bytes))
	} else {
		out, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}
		_, err = out.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		err = out.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
