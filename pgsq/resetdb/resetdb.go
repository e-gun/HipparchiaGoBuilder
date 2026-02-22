//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import "github.com/e-gun/HipparchiaGoBuilder/pgsq"

func ResetDatabase() {
	adminpass := RequestPostgresAdminPW()
	InitializeHDB(adminpass, pgsq.DEFAULTPSQLPASS)
	InitializeAuthorTable()
	InitializeWorksTable()
	InitializeAllGreekSupportTables()
	InitializeAllLatinSupportTables()
	InitializeBuildMetadataTable()
}
