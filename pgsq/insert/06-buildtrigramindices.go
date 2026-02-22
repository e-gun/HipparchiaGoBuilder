//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

// BuildTrigramIndices - add an index to an author table to speed up searches
func BuildTrigramIndices() {
	const (
		TMPL1 = `DROP INDEX %s_accented_line_trgm_gist_idx;`
		TMPL3 = `DROP INDEX %s_stripped_line_trgm_gist_idx;`
		TMPL2 = `
			CREATE INDEX %s_accented_line_trgm_gist_idx on public.%s
			  USING gist(accented_line gist_trgm_ops(siglen=16));`
		TMPL4 = `
			CREATE INDEX %s_stripped_line_trgm_gist_idx on public.%s
			  USING gist(stripped_line gist_trgm_ops(siglen=16));`
		TMPL5 = `CREATE INDEX %s_accented_trgm_idx ON public.%s USING gin (accented_line public.gin_trgm_ops);`
		TMPL6 = `CREATE INDEX %s_stripped_trgm_idx ON public.%s USING gin (stripped_line public.gin_trgm_ops);`
	)

	// drop should be superfluous because the index will go when you drop the table it indexes
	// and BuildTrigramIndices() should only run right after CreateAuthorTable() which has the drop

	// it is not yet clear which index style is faster; they seem roughly equivalent in practice
	// but it is definitely slower to build with `gist_trgm_ops(siglen=256)`: add almost 50%
	// and one is generally expected to use GIN

	// note especially that only a relatively small number of large author tables will speed up because of an index

	// https://www.todosdb.com/blog/understanding-postgresqls-advanced-indexing-gin-and-gist/
	// for GiSt we read:
	//	Full-Text Search: While GIN is often preferred for full-text indexing,
	//	GiST can also index text by leveraging specific operator classes,
	//	enabling proximity-based text search and relevance ranking.
	// for GIN we read:
	//	Full-Text Search Efficiency: In full-text search applications, GIN can
	//	rapidly locate documents containing specific terms by leveraging its
	//	inverted index structure. GIN’s ability to handle large datasets of
	//	text and return relevant documents based on keyword searches
	//	significantly enhances search performance.

	// documenting the differences in practice: https://oxilor.com/blog/comparing-indexes-for-text-search-in-postgresql-part-2
	// cutting to the conclusion:
	//	The search speed of GiST and GIN indexes is approximately the same
	//	with a short string length (up to ~250 characters). At the same time,
	//	the GIN index is much smaller that the GiST index (by 3.5 times in
	//	our test). Therefore, it is worth choosing the GIN index for short
	//	strings (e.g. up to 250 characters).

	fmt.Println("Building Indices")

	toinsert := global.TheMasterAuMap.Keys()
	for _, tablename := range toinsert {
		queries := []string{
			//fmt.Sprintf(TMPL1, tablename),
			//fmt.Sprintf(TMPL3, tablename),
			fmt.Sprintf(TMPL5, tablename, tablename),
			fmt.Sprintf(TMPL6, tablename, tablename),
		}

		// fmt.Println("BuildTrigramIndices", tablename)
		dbc.DBCCommandSequence(queries)
	}
}
