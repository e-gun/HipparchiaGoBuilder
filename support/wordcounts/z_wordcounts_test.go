package wordcounts

import (
	"bytes"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline"
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"os"
	"strings"
	"testing"
)

// two short ones:
// TLG2460	 0.051s
// TLG0379	 0.048s

func ReadyDBConnection() {
	pl := dbc.PostgresLogin{
		Host:   pgsq.DEFAULTPSQLHOST,
		Port:   pgsq.DEFAULTPSQLPORT,
		User:   pgsq.DEFAULTPSQLUSER,
		Pass:   pgsq.DEFAULTPSQLPASS,
		DBName: pgsq.DEFAULTPSQLDB,
	}

	dbc.SQLPool = dbc.FillDBConnectionPool(pl)
}

func TestFanoutCorpusCount(t *testing.T) {
	// this pseudo-test is useful to figure out which characters need stripping in BuildStrippedAndAccentedlines()
	// at the head of the index file you will see all sorts of random garbage; BuildStrippedAndAccentedlines() is
	// where you can clean it

	const (
		OUTTMPL = "%s %d\n"
	)
	// Create a byte buffer to accumulate the strings
	var buffer bytes.Buffer

	ReadyDBConnection()

	auu := authorbuildpipeline.BuildAuthorsSlice("../../data/TLG/", "TLG")
	auu = authorbuildpipeline.SortWorkpileByFilesize(auu, "../../data/TLG/")
	var testauthors []string
	for _, a := range auu {
		testauthors = append(testauthors, strings.ReplaceAll(a, "TLG", "gr"))
	}

	//testauthors := []string{"gr1444"}
	//testauthors = FindAuthorsToGrab("gr")
	fmt.Println("to count:", len(testauthors))
	results := FanoutCorpusCount(false, testauthors)
	rkeys := generic.PolytonicSort(generic.StringMapKeysIntoSlice(results))
	for _, k := range rkeys {
		buffer.WriteString(fmt.Sprintf(OUTTMPL, k, results[k]))
	}
	file, err := os.Create("wordcounts.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the buffered data to the file
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func TestParsedCountOfAllEras(t *testing.T) {
	x := ParsedCountOfAllEras()
	fmt.Println(len(x))
}

func TestCalculateCorpusWeights(t *testing.T) {
	ReadyDBConnection()
	CalculateUnparsedWordcountWeights()
}

func TestCalculateGenreWeights(t *testing.T) {
	ReadyDBConnection()
	CalculateParsedWordcountTotals()
}

// gr2460
//  	ΑΡΑΒΙΚΗ ΑΡΧΑΙΟΛΟΓΙΑ.
//t1.1 	E LIBRO SECVNDO.
//1.1 	   Steph. Byz﹕ Ἀταφηνοὶ, ἔθνοϲ μέγα Ἀραβίαϲ,
//	περὶ οὗ Γλαῦκοϲ ἐν δευτέρᾳ.
//	   Δούμαθα, πόλιϲ Ἀραβίαϲ. Ὁ πολίτηϲ Δουμαθη-
//	νόϲ, ὡϲ Γλαῦκοϲ ἐν βʹ Ἀραβικῆϲ ἀρχαιολογίαϲ.
//	   Ἔρθα, πόλιϲ Παρθίαϲ ἐπὶ τῷ Εὐφράτῃ. Τὸ ἐθνι-
//	κὸν Ἐρθηνόϲ, ὡϲ Γλαῦκοϲ ἐν Ἀραβικῶν δευτέρῳ.
//	   Νέγλα, πολίχνιον Ἀραβίαϲ. Γλαῦκοϲ δευτέρῳ
//	Ἀραβικῆϲ ἀρχαιολογίαϲ. Τὸ ἐθνικὸν Νέγλιοϲ ἢ Νεχλίτηϲ
//	τῷ ἔθει τῆϲ χώραϲ.
//1.10 	   Ὄμανα, πόλιϲ τῆϲ εὐδαίμονοϲ Ἀραβίαϲ. Γλαῦ-
//	κοϲ δευτέρῳ Ἀρ. ἀρχ. Τὸ ἐθνικὸν Ὀμανεύϲ.
//	   Εὐαληνοὶ, ἔθνοϲ περὶ οὗ φηϲὶ Γλαῦκοϲ ἐν δευτέρῳ
//	Περὶ Ἀραβίαϲ.
//t2.1 	E LIBRO TERTIO.
//2.1 	   Idem﹕ Ἀΐλανον, πόλιϲ Ἀραβίαϲ, ἧϲ ὁ πολίτηϲ
//	Ἀϊλανίτηϲ. Τινὲϲ δὲ κόλπον Αἴλα (﹖) φαϲί. Γλαῦκοϲ δὲ
//	κώμην αὐτὴν λέγει ἐν Ἀραβικῶν τρίτῳ· «Τὰ πρὸϲ ἕω
//	τῆϲ Αΐλαϲ.» Ὁ πολίτηϲ Αἰλανίτηϲ.
//	   Βαϲιννοὶ, Ἀραβικὸν ἔθνοϲ. Γλαῦκοϲ ἐν τρίτῳ
//	Ἀραβικῆϲ ἀρχαιολογίαϲ.
//t3.1 	E LIBRO QVARTO.
//3.1 	   Idem﹕ Γάδδα, χωρίον Ἀραβίαϲ. Γλαῦκοϲ ἐν τε-
//	τάρτῳ. Καὶ θηλυκῶϲ καὶ οὐδετέρωϲ. Τὸ ἐθνικὸν Γαδ-
//	δηνόϲ.
//	   Idem﹕ Χαράκμωβα, πόλιϲ τῆϲ νῦν τρίτηϲ Πα-
//	λαιϲτίνηϲ ... ὁ πολίτηϲ .. Χαρακμωβηνόϲ ... Γλαῦκοϲ ἐν
//	Ἀραβικῆϲ ἀρχαιολογίαϲ τετάρτῳ· «Ἡϲύχαζον δ’ ἐν
//	τούτοιϲ Χαρακμωβηνοί.»
//t4.1 	E LIBRIS INCERTIS.
//4.1 	   Idem﹕ Ϲαλμηνοὶ, ἔθνοϲ νομαδικὸν, ὡϲ Γλαῦκοϲ
//	ἐν ﹡ Ἀραβικῆϲ ἀρχαιολογίαϲ.
//	   Γέα, πόλιϲ πληϲίον Πετρῶν ἐν Ἀραβίᾳ, ὡϲ Γλαῦ-
//	κοϲ ἐν Ἀραβικῇ ἀρχαιολογίᾳ.
//	   Ἀρίνδηλα, πόλιϲ Παλαιϲτίνηϲ, Γλαῦκοϲ δὲ κώ-
//	μην αὐτὴν καλεῖ.

// gr0379
//  	Κύκλωψ ἢ Γαλάτεια
//2.69 	ἄνδρα δὲ τὸν Κυθέρηθεν ὃν ἐθρέψαντο τιθῆναι
//	   Βάκχου καὶ λωτοῦ πιϲτότατον ταμίην
//2.70 	Μοῦϲαι παιδευθέντα Φιλόξενον, οἷα τιναχθεὶϲ
//2.70 	   Ὀρτυγίηι ταύτηϲ ἦλθε διὰ πτόλεωϲ
//	γιγνώϲκειϲ, ἀίουϲα μέγαν πόθον ὃν Γαλατείη
//	   αὐτοῖϲ μηλείοιϲ θήκαθ’ ὑπὸ προγόνοιϲ.
//6,1.1 	θρεττανελό
//	ἀλλ’ εἶα τέκεα θαμίν’ ἐπαναβοῶντεϲ
//7.1 	πήραν ἔχοντα λάχανά τ’ ἄγρια δροϲερὰ
//8.1 	ὦ καλλιπρόϲωπε χρυϲεοβόϲτρυχε [Γαλάτεια]
//	χαριτόφωνε θάλοϲ Ἐρώτων
//9.1 	Μούϲαιϲ εὐφώνοιϲ ἰωμένη τὸν ἔρωτα
//10.1 	ἔθυϲαϲ, ἀντιθύϲηι.
//11.1 	οἵωι μ’ ὁ δαίμων τέρατι ϲυγκαθεῖρξεν·
//15.1 	Γάμε θεῶν λαμπρότατε
//16.1 	αὐτοὶ γὰρ διὰ Παρναϲϲοῦ
//	χρυϲορόφων Νυμφέων εἴϲω
//	θαλάμων
//17.1 	ϲτρεπταίγλαν
//18.1 	εὐρείταϲ οἶνοϲ πάμφωνοϲ.
//19.1 	⟨τὸν⟩ ἀρκεϲίγυιον
//20.1 	ϲυμβαλοῦμαί τι μέλοϲ ὑμῖν εἰϲ ἔρωτα
//	[2] Philoxenus Lyr., Tituli
//1.tit 	Γενεαλογία τῶν Αἰακιδῶν
//12.tit 	Κωμαϲτήϲ﹖
//13.tit 	Μυϲοί﹖
//14.tit 	Ϲύροϲ﹖
