//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package main

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/idtandbin"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/resetdb"
	"github.com/e-gun/HipparchiaGoBuilder/support/grammar"
	"github.com/e-gun/HipparchiaGoBuilder/support/lexica"
	"github.com/e-gun/HipparchiaGoBuilder/support/wordcounts"
	"os/exec"
	"strings"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

// these next variables should be injected at build time: 'go build -ldflags "-X main.GitCommit=$GIT_COMMIT"', etc

var BuildDate string

func main() {
	const (
		MSG = "%s v%s total execution time was %.3fs"
	)
	start := time.Now()
	global.GitHash = GetGitCommitHash(".")

	structs.AuthAndWorkSeparator = global.AUTHWORKSEPARATOR

	// defer profile.Start().Stop()
	// go tool pprof --pdf ./HipparchiaGoBuilder ./default.pgo > CPUProfile.pdf
	// defer profile.Start(profile.MemProfile).Stop()

	cfg := readconfigfile()
	cfg = readcommandline(cfg)

	ReadyDBConnection()

	if cfg.Reset {
		global.HEAD("Reset DB")
		resetdb.ResetDatabase()
	}

	if cfg.OneGreek != "" {
		fmt.Println("One Greek: ", cfg.OneGreek)
		authorbuildpipeline.BuildOneAuthor(cfg.GreekDir, "TLG"+cfg.OneGreek)
		insert.BuildTrigramIndices()
		idtandbin.MapIdtAuMapOntoMasterAuMap()
		idtandbin.StoreUpdatedMetadata()
	}

	if cfg.OneLatin != "" {
		fmt.Println("One Latin: ", cfg.OneLatin)
		global.WorkingOnCorpus = "LAT"
		betacode.EarybirdTuples = betacode.GetEarlyBirdTuples()
		authorbuildpipeline.BuildOneAuthor(cfg.LatDir, "LAT"+cfg.OneLatin)
		insert.BuildTrigramIndices()
		idtandbin.StoreUpdatedMetadata()
		global.WorkingOnCorpus = ""
		betacode.EarybirdTuples = betacode.GetEarlyBirdTuples()
	}

	if cfg.OneIns != "" {
		fmt.Println("One Inscription: ", cfg.OneIns)
		authorbuildpipeline.BuildOneAuthor(cfg.InsDir, "INS"+cfg.OneIns)
		idtandbin.StoreUpdatedMetadata()
		insert.BuildTrigramIndices()
		fmt.Println("WARNING: One Inscription is for testing only. You just broke ALL of the inscriptions.")
	}

	if cfg.OneChr != "" {
		fmt.Println("One Christian: ", cfg.OneChr)
		authorbuildpipeline.BuildOneAuthor(cfg.ChrDir, "CHR"+cfg.OneChr)
		insert.BuildTrigramIndices()
		idtandbin.StoreUpdatedMetadata()
		fmt.Println("WARNING: One Christian is for testing only. You just broke ALL of the Christians.")
	}

	if cfg.OnePap != "" {
		fmt.Println("One Papyrus: ", cfg.OnePap)
		authorbuildpipeline.BuildOneAuthor(cfg.PapDir, "DDP"+cfg.OnePap)
		insert.BuildTrigramIndices()
		idtandbin.StoreUpdatedMetadata()
		fmt.Println("WARNING: One Papyrus is for testing only. You just broke ALL of the papyri.")
	}

	if cfg.DoGreek {
		global.WorkingOnCorpus = "TLG"
		global.HEAD("Greek Corpus")
		authorbuildpipeline.RunCorpusPipeline(cfg.GreekDir, "TLG")
	}

	if cfg.DoLatin {
		global.HEAD("Latin Corpus")
		global.WorkingOnCorpus = "LAT"
		betacode.EarybirdTuples = betacode.GetEarlyBirdTuples()
		authorbuildpipeline.RunCorpusPipeline(cfg.LatDir, "LAT")
	}

	if cfg.DoIns {
		global.WorkingOnCorpus = "INS"
		global.HEAD("Inscriptions Corpus")
		authorbuildpipeline.RunCorpusPipeline(cfg.InsDir, "INS")
	}

	if cfg.DoPap {
		global.WorkingOnCorpus = "DDP"
		global.HEAD("Papyrus Corpus")
		authorbuildpipeline.RunCorpusPipeline(cfg.PapDir, "DDP")
	}

	if cfg.DoChr {
		global.WorkingOnCorpus = "CHR"
		global.HEAD("Christian Corpus")
		authorbuildpipeline.RunCorpusPipeline(cfg.ChrDir, "CHR")
	}

	if cfg.DoGkLx {
		global.HEAD("Greek Lexicon")
		datadir := cfg.GkLxDataLoc
		resetdb.InitializeGreekLexicalSupportTables()
		xmls := lexica.BuildLexDataFileNamesSlice(datadir, ".xml")
		lexica.FanoutGkLexBuilder(xmls, datadir)
		glh := GetGitCommitHash(cfg.GkLxDataLoc)
		insert.InsertBuildMetadata("GLex", "Commit: "+glh)
	}

	if cfg.DoGkGr {
		global.HEAD("Greek Grammar")
		datadir := cfg.GkGrDataLoc
		glem := cfg.GkLemFileName
		ganf := cfg.GkAnaFileName
		grammar.BuildAndLoadGreekGrammar(datadir, glem, ganf)
		insert.InsertBuildMetadata("GGram", "")
	}

	if cfg.DoLtLx {
		global.HEAD("Latin Lexicon")
		datadir := cfg.LtLxDataLoc
		datafile := cfg.LtLexFileName
		resetdb.InitializeLatinLexicalSupportTables()
		lexica.FanoutLatinLexBuilder(datadir, datafile)
		insert.InsertBuildMetadata("LLex", cfg.LtLexFileName)
	}

	if cfg.DoLtGr {
		global.HEAD("Latin Grammar")
		datadir := cfg.LtGrDataLoc
		glem := cfg.LtLemFileName
		ganf := cfg.LtAnaFileName
		grammar.BuildAndLoadLatinGrammar(datadir, glem, ganf)
		insert.InsertBuildMetadata("LGram", "")
	}

	if cfg.DoWdCt {
		global.HEAD("Word Counts")
		wordcounts.HeadwordLookupMap = wordcounts.PrepareAnalysisLookupMap()
		wordcounts.DoAllWordcounts()
	}

	if cfg.TestRun {
		global.HEAD("TestRun")
		idtandbin.LoadLatinCanon("/Users/erik/Development/go/src/github.com/e-gun/HipparchiaGoBuilder/data/LAT/")
		// idtandbin.LoadLatinCanonIDT("/Users/erik/Development/go/src/github.com/e-gun/HipparchiaGoBuilder/data/LAT/")

	}

	d := fmt.Sprintf(MSG, global.NAME, global.VERSION, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func ReadyDBConnection() {
	pl := dbc.PostgresLogin{
		Host:   pgsq.DEFAULTPSQLHOST,
		Port:   pgsq.DEFAULTPSQLPORT,
		User:   global.Config.HGDBUserName,
		Pass:   global.Config.HGDBUserPass,
		DBName: pgsq.DEFAULTPSQLDB,
	}

	dbc.SQLPool = dbc.FillDBConnectionPool(pl)
}

func GetGitCommitHash(repopath string) string {
	const (
		LENGTH = 8
	)
	// Get the latest commit hash
	gitcommand := exec.Command("git", "rev-parse", "HEAD")
	gitcommand.Dir = repopath
	hashbytes, err := gitcommand.CombinedOutput()
	if err != nil {
		fmt.Println("Error getting git commit hash")
		return "(unknown)"
	}
	cleanhash := strings.TrimSpace(string(hashbytes))

	trimphash := cleanhash[:LENGTH]
	return trimphash
}
