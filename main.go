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

var BuildDate string

func main() {
	const (
		MSG = "%s v%s total execution time was %.3fs"
	)
	start := time.Now()
	global.GitHash = GetGitCommitHash(".")

	structs.AuthAndWorkSeparator = global.AUTHWORKSEPARATOR
	structs.AbbrevLen = global.ABBREVLEN

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
		global.WorkingOnCorpus = "TLG"

		resetdb.DropOneAuthorTable(global.TLGABBREV + cfg.OneGreek)
		resetdb.ResetOneAuthor(global.TLGABBREV + cfg.OneGreek)

		authorbuildpipeline.BuildOneAuthor(cfg.GreekDir, global.WorkingOnCorpus+cfg.OneGreek)
		insert.BuildTrigramIndices()
		idtandbin.MapIdtAuMapOntoMasterAuMap()
		idtandbin.StoreUpdatedMetadata()
		insert.InsertBuildMetadata(global.WorkingOnCorpus, "(mixed build) OneGreek: "+cfg.OneGreek)
	}

	if cfg.OneLatin != "" {
		fmt.Println("One Latin: ", cfg.OneLatin)
		global.WorkingOnCorpus = "LAT"

		resetdb.DropOneAuthorTable(global.LATABBREV + cfg.OneLatin)
		resetdb.ResetOneAuthor(global.LATABBREV + cfg.OneLatin)

		betacode.EarybirdTuples = betacode.GetEarlyBirdTuples()
		authorbuildpipeline.BuildOneAuthor(cfg.LatDir, global.WorkingOnCorpus+cfg.OneLatin)
		insert.BuildTrigramIndices()
		idtandbin.StoreUpdatedMetadata()
		global.WorkingOnCorpus = ""
		betacode.EarybirdTuples = betacode.GetEarlyBirdTuples()
		insert.InsertBuildMetadata(global.WorkingOnCorpus, "(mixed build) OneLatin: "+cfg.OneLatin)
	}

	if cfg.OneIns != "" {
		a := "One Inscription"
		b := "INS"
		c := "oneins"
		onephi(a, b, c, cfg.InsDir, cfg.OneIns)
	}

	if cfg.OneChr != "" {
		a := "One Christian"
		b := "CHR"
		c := "onechr"
		onephi(a, b, c, cfg.ChrDir, cfg.OneChr)
	}

	if cfg.OnePap != "" {
		a := "One Papyrus"
		b := "DDP"
		c := "oneddp"
		onephi(a, b, c, cfg.PapDir, cfg.OnePap)
	}

	if cfg.DoGreek {
		global.HEAD("Greek Corpus")
		global.WorkingOnCorpus = "TLG"
		resetdb.ResetCorpus(global.WorkingOnCorpus)
		authorbuildpipeline.RunCorpusPipeline(cfg.GreekDir, global.WorkingOnCorpus)
	}

	if cfg.DoLatin {
		global.HEAD("Latin Corpus")
		global.WorkingOnCorpus = "LAT"
		resetdb.ResetCorpus(global.WorkingOnCorpus)
		betacode.EarybirdTuples = betacode.GetEarlyBirdTuples()
		authorbuildpipeline.RunCorpusPipeline(cfg.LatDir, global.WorkingOnCorpus)
	}

	if cfg.DoIns {
		global.WorkingOnCorpus = "INS"
		global.HEAD("Inscriptions Corpus")
		resetdb.ResetCorpus(global.WorkingOnCorpus)
		authorbuildpipeline.RunCorpusPipeline(cfg.InsDir, global.WorkingOnCorpus)
	}

	if cfg.DoPap {
		global.WorkingOnCorpus = "DDP"
		global.HEAD("Papyrus Corpus")
		resetdb.ResetCorpus(global.WorkingOnCorpus)
		authorbuildpipeline.RunCorpusPipeline(cfg.PapDir, global.WorkingOnCorpus)
	}

	if cfg.DoChr {
		global.WorkingOnCorpus = "CHR"
		global.HEAD("Christian Corpus")
		resetdb.ResetCorpus(global.WorkingOnCorpus)
		authorbuildpipeline.RunCorpusPipeline(cfg.ChrDir, global.WorkingOnCorpus)
	}

	if cfg.DoGkLx {
		global.HEAD("Greek Lexicon")
		datadir := cfg.GkLxDataLoc
		resetdb.InitializeGreekLexicalSupportTables()
		xmls := lexica.BuildLexDataFileNamesSlice(datadir, ".xml")
		lexica.FanoutGkLexBuilder(xmls, datadir)
		glh := GetGitCommitHash(cfg.GkLxDataLoc)
		insert.InsertBuildMetadata("GLex", "Logeon Commit: "+glh)
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
		resetdb.DropAuthorMultipleTables(global.LATABBREV)
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

func onephi(a string, cp string, c string, fd string, fn string) {
	fmt.Printf("%s: %s\n", a, fn)
	global.WorkingOnCorpus = cp
	resetdb.ResetCorpus(global.WorkingOnCorpus)
	authorbuildpipeline.BuildOneAuthor(fd, global.WorkingOnCorpus+fn)
	insert.BuildTrigramIndices()
	idtandbin.StoreUpdatedMetadata()
	fmt.Printf("WARNING: '%s' is for testing only. You just broke the WHOLE of this corpus.\n", a)
	insert.InsertBuildMetadata(global.WorkingOnCorpus, fmt.Sprintf("(broken corpus) '--%s %s'", c, fn))
}
