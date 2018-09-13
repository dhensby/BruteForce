package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/ngirot/BruteForce/bruteforce"
	"github.com/ngirot/BruteForce/bruteforce/conf"
	"github.com/ngirot/BruteForce/bruteforce/hashs"
	"os"
)

func main() {
	var bench = flag.Bool("benchmark", false, "Launch a benchmark")
	var value = flag.String("value", "", "The value to be tested")
	var alphabet = flag.String("alphabet", "", "The file containing all characters")
	var dictionary = flag.String("dictionary", "", "The file containing all words to be tested")
	var hashType = flag.String("type", "sha256", "The hash type")
	var salt = flag.String("salt", "", "The salt added to the end of the generated word")
	var test = flag.Bool("test", false, "Test value to use")
	flag.Parse()

	if *bench {
		var types = hashs.AllHasherTypes()
		for _, t := range types {
			var hasherCreator, _ = hashs.HasherCreator(t)
			var timeOneCpu = bruteforce.BenchHasherOneCpu(hasherCreator)
			fmt.Printf("One CPU (%s) hasher: %d kh/s\n", t, timeOneCpu/1000)
			var timeMultiCpu = bruteforce.BenchHasherMultiCpu(hasherCreator)
			fmt.Printf("Multi CPU (%s) hasher: %d kh/s\n", t, timeMultiCpu/1000)
		}

		fmt.Printf("One CPU word generator: %d kw/s\n", bruteforce.BenchBruter()/1000)
		os.Exit(0)
	}

	if *value == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *test {
		var hasherCreator, _ = hashs.HasherCreator(*hashType)
		var hashBytes = hasherCreator().Hash(*value);
		//fmt.Printf("%d\n", binary.LittleEndian.Uint32(hashBytes))
		var hash = hex.EncodeToString(hashBytes)
		fmt.Printf("Hash: %s \nLength: %d\n", hash, len(hash))
		os.Exit(0)
	}

	fmt.Printf("Start brute-forcing (%s)...\n", *hashType)

	var hashConf = conf.NewHash(*value, *hashType)
	var wordConf = conf.NewWordConf(*dictionary, *alphabet, *salt)

	var chrono = bruteforce.NewChrono()
	chrono.Start()
	if result, error := bruteforce.Launch(hashConf, wordConf); error == nil {
		chrono.End()

		if result != "" {
			fmt.Printf("\rFound: %s in %d s\n", result, chrono.DurationInRoundedSeconds())
		} else {
			fmt.Printf("\rNot found\n")
		}
	} else {
		fmt.Printf("Hasher %s invalid: %q\n", *hashType, error)
	}
}
