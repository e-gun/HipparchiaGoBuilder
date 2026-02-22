//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
)

const (
	CONFIGFILENAME   = "hgb-config.json"
	DEFAULTGKDIR     = "./data/TLG/"
	DEFAULTLTDIR     = "./data/LAT/"
	DEFAULTINSDIR    = "./data/PHI/"
	DEFAULTPAPDIR    = "./data/PHI/"
	DEFAULTCHRDIR    = "./data/PHI/"
	DEFAULTGKLEXDIR  = "./data/LSJLogeion/"
	DEFAULTLATEXDIR  = "./data/LatinLexicon/"
	DEFAULTGKGRAMDIR = "./data/Grammar/"
	DEFAULTGKLMFN    = "greek-lemmata.txt"
	DEFAULTLTLMFN    = "latin-lemmata.txt"
	DEFAULTGKANFN    = "greek-analyses.txt"
	DEFAULTLTANFN    = "latin-analyses.txt"
	DEFAULTLTLXFN    = "latin-lexicon_1999.04.0059.xml"
)

func readcommandline(cfg structs.BuildConfig) structs.BuildConfig {
	const (
		HELP = `
HipparchiaGoBuilder
	-h            help
	-v            show version
	-cc {int}     # of multi-core workers to dispatch (default is 'runtime.NumCPU() - 1')
	-rp           reproducible author/word names for INS, etc builds: SLOW (ensured via eliminating multi-core workers)
	-q            quiet build: suppress most trivial error messages
	
	-00           reset whole database (probably a good idea before a big build unless you know why it is not)
	-all          build TLG, LAT, INS, DDP, and CHR along with grammar, dictionaries and wordcounts (+ reset DB before start)
	-allcorp      build TLG, LAT, INS, DDP, and CHR (but no grammar, dictionaries or wordcounts)
	-allbutwc     build TLG, LAT, INS, DDP, and CHR along with grammar, and dictionaries (but no wordcounts)
	-chr          build all Christian inscriptions
	-ddp          build all DDP texts
	-ggr          build Greek grammar
	-glx          build Greek lexicon
	-grk          build all TLG authors
	-ins          build all Greek Inscriptions
	-lat          build all LAT authors
	-lgr          build Latin grammar
	-llx          build Latin lexicon
	-onechr NNNN  build one CHR file [DANGER: debugging only; will *break* authors and works tables; also true of ins and pap]
	-onegrk NNNN  build one TLG author, e.g 0012
	-oneins NNNN  build one INS file
	-onelat NNNN  build one Latin author
	-oneddp NNNN  build one DDP file
	-wc           build wordcounts

to cut to the chase and just do it all:
    ./HipparchiaGoBuilder -all

to build Greek: 
    ./HipparchiaGoBuilder -ggr -glx -grk

`
	)

	args := os.Args[1:len(os.Args)]
	for i, a := range args {
		switch a {
		case "-h", "--help":
			fmt.Println(global.NAME + " " + global.VERSION)
			fmt.Println("Built:", BuildDate)
			fmt.Println("Git:", global.GitHash)
			fmt.Println(HELP)
			os.Exit(1)
		case "-v", "--version":
			fmt.Println(global.NAME + " " + global.VERSION)
			fmt.Println("Git:", global.GitHash)
			fmt.Println("Built:", BuildDate)
			os.Exit(1)
		case "-all":
			cfg.Reset = true
			cfg.DoChr = true
			cfg.DoPap = true
			cfg.DoGreek = true
			cfg.DoIns = true
			cfg.DoLatin = true
			cfg.DoGkLx = true
			cfg.DoLtLx = true
			cfg.DoGkGr = true
			cfg.DoLtGr = true
			cfg.DoWdCt = true
		case "-allcorp":
			cfg.DoChr = true
			cfg.DoPap = true
			cfg.DoGreek = true
			cfg.DoIns = true
			cfg.DoLatin = true
		case "-allbutwc":
			cfg.DoChr = true
			cfg.DoPap = true
			cfg.DoGreek = true
			cfg.DoIns = true
			cfg.DoLatin = true
			cfg.DoGkLx = true
			cfg.DoLtLx = true
			cfg.DoGkGr = true
			cfg.DoLtGr = true
			cfg.DoWdCt = true
		case "-00":
			cfg.Reset = true
		case "-chr":
			cfg.DoChr = true
		case "-ddp":
			cfg.DoPap = true
		case "-grk":
			cfg.DoGreek = true
		case "-ins":
			cfg.DoIns = true
		case "-lat":
			cfg.DoLatin = true
		case "-onechr":
			cfg.OneChr = args[i+1]
		case "-oneddp":
			cfg.OnePap = args[i+1]
		case "-onegrk":
			cfg.OneGreek = args[i+1]
		case "-oneins":
			cfg.OneIns = args[i+1]
		case "-onelat":
			cfg.OneLatin = args[i+1]
		case "-glx":
			cfg.DoGkLx = true
		case "-ggr":
			cfg.DoGkGr = true
		case "-lgr":
			cfg.DoLtGr = true
		case "-llx":
			cfg.DoLtLx = true
		case "-q":
			cfg.QuietBuild = true
		case "-rp":
			cfg.Reproducible = true
		case "-test":
			cfg.TestRun = true
		case "-wc":
			cfg.DoWdCt = true
		case "-cc":
			cores, e := strconv.Atoi(args[i+1])
			if e != nil {
				panic(e)
			}
			cfg.WorkerCount = cores
		default:
			//
		}
	}

	if cfg.WorkerCount < 1 {
		cfg.WorkerCount = 1
	}

	global.Config = cfg
	return cfg
}

func writeblankconfig() {
	cfg := builddefaultconfig()
	content, err := json.MarshalIndent(cfg, "  ", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(CONFIGFILENAME, content, 0644)
	if err != nil {
		panic(err)
	}
}

func readconfigfile() structs.BuildConfig {
	loadedcfg, e := os.Open(CONFIGFILENAME)
	if e != nil {
		fmt.Printf("Error opening %s: %s\n", CONFIGFILENAME, e)
		fmt.Printf("Building a new blank configuration and saving it to %s.\n", CONFIGFILENAME)
		writeblankconfig()
		return builddefaultconfig()
	}

	decoderc := json.NewDecoder(loadedcfg)
	cfg := structs.BuildConfig{}
	err := decoderc.Decode(&cfg)
	_ = loadedcfg.Close()
	if err != nil {
		panic(err)
	}
	return cfg
}

func builddefaultconfig() structs.BuildConfig {
	return structs.BuildConfig{
		WorkerCount:   runtime.NumCPU() - 1,
		DataBaseName:  pgsq.DEFAULTPSQLDB,
		QuietBuild:    false,
		PGAdminUser:   "",
		PGAdminPass:   "",
		HGDBUserName:  pgsq.DEFAULTPSQLUSER,
		HGDBUserPass:  pgsq.DEFAULTPSQLPASS,
		Reset:         false,
		DoGreek:       false,
		DoLatin:       false,
		DoPap:         false,
		DoChr:         false,
		DoIns:         false,
		DoGkLx:        false,
		DoLtLx:        false,
		DoGkGr:        false,
		DoLtGr:        false,
		DoWdCt:        false,
		GreekDir:      DEFAULTGKDIR,
		LatDir:        DEFAULTLTDIR,
		InsDir:        DEFAULTPAPDIR,
		PapDir:        DEFAULTINSDIR,
		ChrDir:        DEFAULTCHRDIR,
		GkLxDataLoc:   DEFAULTGKLEXDIR,
		LtLxDataLoc:   DEFAULTLATEXDIR,
		GkGrDataLoc:   DEFAULTGKGRAMDIR,
		LtGrDataLoc:   DEFAULTGKGRAMDIR,
		GkLemFileName: DEFAULTGKLMFN,
		LtLemFileName: DEFAULTLTLMFN,
		GkAnaFileName: DEFAULTGKANFN,
		LtAnaFileName: DEFAULTLTANFN,
		LtLexFileName: DEFAULTLTLXFN,
		Reproducible:  false,
		TestRun:       false,
	}
}
