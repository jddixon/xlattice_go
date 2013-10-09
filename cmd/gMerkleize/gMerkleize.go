package main

// xlattice_go/cmd/gMerkleize/gMerkleize.go

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	xm "github.com/jddixon/xlattice_go/util/merkletree"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var _ = errors.New

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const ()

var (
	// these need to be referenced as pointers
	hashOutput    = flag.Bool("x", false, "output top level hash")
	inDir         = flag.String("i", "", "path to directory being scanned")
	justShow      = flag.Bool("j", false, "display option settings and exit")
	outPath       = flag.String("o", "", "write serialzed merkletree here")
	showTimestamp = flag.Bool("t", false, "output UTC timestamp")
	showTree      = flag.Bool("m", false, "output the merkletree")
	showVersion   = flag.Bool("V", false, "output package version info")
	testing       = flag.Bool("T", false, "test run")
	usingSHA1     = flag.Bool("1", false, "test run")
	verbose       = flag.Bool("v", false, "be talkative")
)

// ------------------------------------------------------------------
type exclusions []string

var excludePats exclusions

func (ex *exclusions) Set(value string) (err error) {
	*ex = append(*ex, value)
	return
}
func (ex *exclusions) String() string {
	return fmt.Sprint(*ex)
}

// ------------------------------------------------------------------
type matches []string

var matchPats matches

func (ma *matches) Set(value string) (err error) {
	*ma = append(*ma, value)
	return
}
func (ma *matches) String() string {
	return fmt.Sprint(*ma)
}

// ------------------------------------------------------------------

func init() {
	flag.Var(&excludePats, "X",
		"list of patterns, file names patterns to be excluded")
	flag.Var(&matchPats, "P",
		"list of patterns, file name patterns to be matched")
}

func main() {
	var (
		dirName, inPath string
		err             error
	)

	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////
	if !*justShow && *inDir == "" {
		err = errors.New("no inDir specified")
	} else if *inDir != "" {
		parts := strings.Split(*inDir, "/")
		partsCount := len(parts)
		if partsCount == 1 {
			inPath = "."
			dirName = *inDir
		} else {
			dirName = parts[partsCount-1] // the last one
			parts = parts[:partsCount-1]
			inPath = strings.Join(parts, "/") // the earlier parts
		}
	}
	if *testing {
	}
	// SANITY CHECKS ////////////////////////////////////////////////
	if err == nil {
	}
	// DISPLAY OPTIONS //////////////////////////////////////////////
	if err == nil && *verbose || *justShow {
		fmt.Printf("excludePats 	= %s\n", excludePats)
		fmt.Printf("hashOutput  	= %v\n", *hashOutput)
		fmt.Printf("inDir       	= %v\n", *inDir)
		fmt.Printf("justShow    	= %v\n", *justShow)
		fmt.Printf("matchPats   	= %s\n", matchPats)
		fmt.Printf("outPath     	= %v\n", *outPath)
		fmt.Printf("showTimestamp   = %v\n", *showTimestamp)
		fmt.Printf("showTree		= %v\n", *showTree)
		fmt.Printf("showVersion 	= %v\n", *showVersion)
		fmt.Printf("testing     	= %v\n", *testing)
		fmt.Printf("usingSHA1       = %v\n", *usingSHA1)
		fmt.Printf("verbose     	= %v\n", *verbose)

		fmt.Println("DEBUG:")
		fmt.Printf("dirName         = %v\n", dirName)
		fmt.Printf("inPath          = %v\n", inPath)
	}
	// DO IT ////////////////////////////////////////////////////////
	if err == nil && !*justShow {
		var (
			doc *xm.MerkleDoc
			ss  []string
		)
		fullPath := path.Join(inPath, dirName)
		doc, err = xm.CreateMerkleDocFromFileSystem(fullPath, *usingSHA1,
			excludePats, matchPats)
		if err == nil {
			tree := doc.GetTree()
			if *showTree {
				tree.ToStrings("", &ss)
			}
			if *hashOutput {
				treeHashAsHex := hex.EncodeToString(tree.GetHash())
				ss = append(ss, treeHashAsHex)
			}
			output := []byte(strings.Join(ss, "\r\n") + "\r\n")
			if *outPath == "" {
				os.Stdout.Write(output)
			} else {
				parts := strings.Split(*outPath, "/")
				partCount := len(parts)
				if partCount > 1 {
					// we need to be sure the directory exists
					parts := parts[:partCount-1]
					outDir := strings.Join(parts, "/")
					err = os.MkdirAll(outDir, 0740)
				}
				if err == nil {
					err = ioutil.WriteFile(*outPath, output, 0644)
				}
			}
		}
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(-1)
	}

}
