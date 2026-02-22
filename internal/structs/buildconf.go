//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

type BuildConfig struct {
	WorkerCount   int
	QuietBuild    bool
	DataBaseName  string
	PGAdminUser   string
	PGAdminPass   string
	HGDBUserName  string
	HGDBUserPass  string
	Reset         bool
	DoGreek       bool
	DoLatin       bool
	DoPap         bool
	DoChr         bool
	DoIns         bool
	DoGkLx        bool
	DoLtLx        bool
	DoGkGr        bool
	DoLtGr        bool
	DoWdCt        bool
	GreekDir      string
	LatDir        string
	PapDir        string
	ChrDir        string
	InsDir        string
	OneGreek      string
	OneLatin      string
	OnePap        string
	OneChr        string
	OneIns        string
	GkLxDataLoc   string
	LtLxDataLoc   string
	GkGrDataLoc   string
	LtGrDataLoc   string
	GkLemFileName string
	LtLemFileName string
	GkAnaFileName string
	LtAnaFileName string
	LtLexFileName string
	Reproducible  bool
	TestRun       bool
}
