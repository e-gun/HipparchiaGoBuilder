## HipparchiaGoBuilder
### v.0.9.2 : release

the python builder is old and confusing and hard to maintain

this is supposed to be faster, cleaner, and easier to maintain

improvements
1. more legible inscriptions and papyri
2. improved lexica
3. better dating
4. better/smarter word counts
5. better testing infrastructure
6. better handling of metadata

all major segments implemented; at the debugging and polishing stage

*caveats*:

only works with `HipparchiaGoServer v2.0.0+`

current status:
1. will build a good version of everything
2. very good for TLG, LAT
3. good for dictionaries, and grammar
4. good for INS, DDP, and CHR

items on the todo list:
1. check / improve dating for all but TLG (c. 90% PHI has a non-empty date)

```
./HipparchiaGoBuilder -h
HipparchiaGoBuilder 0.8.0b
Built: 2025-04-29@14:01:45
Git: 8c4d36f0

HipparchiaGoBuilder
	-h            help
	-v            show version
	-cc {int}     # of multi-core workers to dispatch (default is 'runtime.NumCPU() - 1')
	-rp           reproducible author/word names for INS, etc builds: SLOW (ensured via eliminating multi-core workers)
	-q            quiet build: suppress most trivial error messages

	-00           reset whole database (probably a good idea before a big build unless you know why it is not)
	-allcorp      build TLG, LAT, INS, DDP, and CHR (but no grammar, dictionaries or wordcounts)
	-allbutwc     build TLG, LAT, INS, DDP, and CHR along with grammar, and dictionaries (but no wordcounts)
	-all          build TLG, LAT, INS, DDP, and CHR along with grammar, dictionaries and wordcounts (+ reset DB before start)
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
    
```

```
% cloc --not-match-d="data" . 
     199 text files.
     194 unique files.                                          
      20 files ignored.

github.com/AlDanial/cloc v 2.04  T=0.15 s (1288.2 files/s, 160311.8 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                             179           2233           2862          16569
Python                           3            256            379            841
XML                              4              0              0            460
Text                             3             57              0            263
Markdown                         1             27              0             91
JSON                             2              0              0             83
Bourne Shell                     2              6              0             15
-------------------------------------------------------------------------------
SUM:                           194           2579           3241          18322
-------------------------------------------------------------------------------

``` 

substantially faster builds (8x+ speedup; even more for word counts: 50x?)

```
./HipparchiaGoBuilder -all

HipparchiaGoBuilder v0.8.2b total execution time was 298.811s

(683s to build w/ 8 cores on Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz)

```

and faster search results in HGS

```
Sought »λέγεϲθαι« within 1 lines of »ἄγεϲθαι«
Searched 7,462 works and found 2 passages (0.77s)
Sorted by author name 

vs

Sought »λέγεϲθαι« within 1 lines of »ἄγεϲθαι«
Searched 7,461 works and found 2 passages (1.58s)
Sorted by author name 
```