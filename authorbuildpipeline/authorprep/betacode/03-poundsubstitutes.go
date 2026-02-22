//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

// ReplacePoundSigns() - convert #1 into "ϟ", etc.

var (
	poundsouter  = regexp.MustCompile(`#\d{1,4}`)
	poundsubsmap = map[int]string{
		// u'\u",
		// it is noisy to build with the 'undocumenteds' commented out
		// but, if you want to get the codes as HTML, this is what you need to do...
		1:  "\u03df",
		2:  "\u03da", // stigma; supposed to be able to handle '*//2' and not just '//2'
		3:  "\u03d9", // koppa; supposed to be able to handle '*//3' and not just '//3'
		4:  "\u03d9", // koppa, variant
		5:  "\u03e1", // sampi; supposed to be able to handle '*//5' and not just '//5'
		6:  "\u2e0f",
		7:  "<hb-sp-unk>(＃7)</hb-sp-unk>", // idiosyncratic
		8:  "\u2e10",
		9:  "\u0301",
		10: "\u03fd",
		11: "\u03ff",
		12: "\u2014",
		13: "\u203b",
		14: "\u2e16",
		15: "⟩", // officially: "\u003e", greater than sign, diple
		16: "\u03fe",
		// 17: "002f",  // careful: '/' is dangerous
		17:  "／", // fulwidth solidus instead
		18:  "⟨", // officially: "\u003c", less than sign, reversed diple
		19:  "\u0300",
		20:  "𐅵",
		21:  "𐅵",
		22:  "\u0375",
		23:  "\u03d9",
		24:  "𐅵",
		25:  "𐅶",
		26:  "\u2e0f",
		27:  "𐄂", // 'check mark'; non tlg; and so - AEGEAN CHECK MARK; Unicode: U+10102, UTF-8: F0 90 84 82
		28:  "<hmu_mark_deleting_entry />␥",
		29:  "\u00b7",                       // middle dot: ·
		30:  "<hb-sp-unk>(＃30)</hb-sp-unk>", // idiosyncratic
		31:  "<hb-sp-unk>(＃31)</hb-sp-unk>", // idiosyncratic
		48:  "<hb-sp-unk>(＃48)</hb-sp-unk>",
		50:  "<hb-sp-unk>(＃50)</hb-sp-unk>",
		51:  "\u00b7", // middle dot: ·
		52:  "\u205a",
		53:  "\u205d",
		54:  "<hb-sp-unk>(＃54)</hb-sp-unk>", // appears in INS0010
		55:  "\u2059",
		56:  "∣", // 'dividers of other forms'; not a helpful description: trying u2223 for now
		57:  "﹤", // small variant;  as per http://noapplet.epigraphy.packhum.org/text/3292?&bookid=5&location=7
		58:  "﹥", // cf //57
		59:  "\u03fd",
		60:  "\u0399",
		61:  "𐅂",
		62:  "𐅃",
		63:  "\u0394",
		64:  "𐅄",
		65:  "\u0397",
		66:  "𐅅",
		67:  "\u03a7",
		68:  "𐅆",
		69:  "\u039c",
		70:  "\u002e",
		71:  "\u00b7",
		72:  "\u02d9",
		73:  "\u205a",
		74:  "\u205d",
		75:  "\u002e",
		80:  "\u0308",
		81:  "＇", // fullwidth apostrophe instead of the dangerous simple apostrophe
		82:  "\u02ca",
		83:  "\u02cb",
		84:  "\u1fc0",
		85:  "\u02bd",
		86:  "\u02bc",
		87:  "\u0394\u0345", // 'ΔΕ'
		90:  "\u2014",
		92:  "<hb-sp-unk>(＃92)</hb-sp-unk>",
		99:  "<hb-sp-unk>(＃99)</hb-sp-unk>",
		100: "𐆆",
		101: "𐅻",           // trouble with the four character unicode codes: uh oh
		102: "𐆂<6\u03c56>", // upsilon supposed to be superscript too: add betacode for that <6...6>
		103: "\u039b\u0338",
		104: "𐆂<6\u03bf6>",                   // the omicron is supposed to be superscript too: add betacode for that <6...6>
		105: "<hb-sp-unk>(＃105)</hb-sp-unk>", // idiosyncratic
		106: "𐆄",
		107: "<hb-sp-unk>(＃107)</hb-sp-unk>", // idiosyncratic
		108: "<hb-sp-unk>(＃108)</hb-sp-unk>", // idiosyncratic
		109: "𐆂<6\u03bf6>",                   // the omicron is supposed to be superscript too: add betacode for that <6...6>
		110: "<11α>11<10\u0375>10",           // need to do the combining accent second, right?
		111: "𐆂<6\u03b56>",
		112: "𐆈", // 𐆈- GREEK GRAMMA SIGN; Unicode: U+10188, UTF-8: F0 90 86 88
		113: "𐅼",
		114: "𐅀",
		115: "𐆉",
		116: "\u2053",
		117: "𐆃", // 𐆃GREEK LITRA SIGN; Unicode: U+10183, UTF-8: F0 90 86 83
		118: "\u03bb\u0338",
		119: "𐅽",
		121: "\u03be\u0338",
		122: "𐅽",
		123: "𐅼",
		124: "<hb-sp-unk>(＃124)</hb-sp-unk>", // idiosyncratic
		125: "𐆂<6\u03c56>",                   // the upsilon is supposed to be superscript too: add betacode for that <6...6>
		126: "<hb-sp-unk>(＃126)</hb-sp-unk>", // idiosyncratic
		127: "\u039b\u0325",                  // 'ΛΕ'
		128: "\u03fc",
		129: "\u039b\u0325", // 'ΛΕ'
		130: "𐆊",
		131: "𐅷",
		132: "\u03b2\u0388", // 'βΈ'
		133: "\u0393<6\u03b26>",
		134: "\u0393<6\u03b26>", // the beta is supposed to be superscript too: add betacode for that <6...6>
		135: "\u02d9",
		136: "\u03a3",           // capital sigma: stater
		137: "\u0393<6\u03b26>", //the beta is supposed to be superscript: add betacode for that <6...6>
		150: "\u221e",
		151: "\u2014",
		152: "\u205a\u2014",
		153: "\u2026\u0305",
		154: "\u2c80",
		155: "\u2014\u0323",
		156: "\u2310",
		157: "<hb-sp-unk>(＃157)</hb-sp-unk>", // idiosyncratic
		158: "\u2237\u0336",
		159: "\u2237\u0344",
		160: "\u007e\u0323",
		161: "𐅵",
		162: "\u25a1",
		163: "\u00b6",
		165: "\u00d7",
		166: "\u2a5a",
		167: "\u039c\u039c",       // supposed to stack this too
		168: "\u039c\u039c\u039c", // supposed to stack this too
		169: "𐅵",
		171: "𐅵",
		172: "𐅵",
		200: "\u2643",
		201: "\u25a1",
		202: "\u264f",
		203: "\u264d",
		204: "\u2640",
		205: "\u2650",
		206: "\u2644",
		207: "\u2609",
		208: "\u263f",
		209: "\u263e",
		210: "\u2642",
		211: "\u2651",
		212: "\u264c",
		213: "\u2648",
		214: "\u264e",
		215: "\u264a",
		216: "\u264b",
		217: "\u2653",
		218: "\u2652",
		219: "\u2649",
		220: "♃",
		221: "\u263d",
		222: "\u260c",
		223: "\u2605",
		240: "𐅷", // 𐅷 GREEK TWO THIRDS SIGN; Unicode: U+10177, UTF-8: F0 90 85 B7
		241: "\u260b",
		242: "\u2651",
		243: "<hb-sp-unk>(＃243)</hb-sp-unk>", // idiosyncratic
		244: "\u264c",
		246: "<hb-sp-unk>(＃246)</hb-sp-unk>", // idiosyncratic
		247: "<hb-sp-unk>(＃247)</hb-sp-unk>", // idiosyncratic
		300: "\u2e0e",                        // but supposed to be just upper half of a coronis
		302: "<hb-sp-unk>(＃302)</hb-sp-unk>", // idiosyncratic
		303: "›",
		304: "\u2e0e", // but supposed to be just part of a coronis
		305: "\u2e0e",
		306: "\u2e0f", // but supposed to be a double paragraphos
		307: "\u2e0e", // but supposed to be just part of a coronis
		308: "\u2e0e", // but supposed to be just part of a coronis
		310: "\u2e0e",
		311: "\u2e0e", // but supposed to be just lower half of a coronis
		312: "\u2e0e", // but supposed to be just upper half of a coronis
		313: "\u2e0e",
		314: "<hb-sp-unk>(＃314)</hb-sp-unk>", // idiosyncratic
		315: "\u2e0e",
		316: "<hb-sp-unk>(＃316)</hb-sp-unk>", // deprecated: no further info
		317: "<hb-document_cancelled_with_slashes />⑊⑊⑊⑊⑊⑊⑊⑊",
		318: "<hb-line_filled_with_cross-strokes />⧷⧷⧷⧷⧷⧷⧷⧷",
		319: "\u25cf",
		320: "\u2629",
		321: "\u2629",
		322: "\u2627",
		323: "﹥", // 'greater-than sign -> line filler' says the instructions; small version instead of the markup version (uFE65, cf uFE64)
		324: "<hb-filler_stroke_to_margin>(filler stroke to margin)</hb-filler_stroke_to_margin>",
		325: "<hb-large_single_X>✕</hmu_large_single_X>",
		326: "<hb-pattern_of_Xs>✕✕✕✕</hmu_pattern_of_Xs>",
		327: "<hb-tachygraphic_marks>(tachygraphic marks)</hb-tachygraphic_marks>",
		329: "<hb-monogram />",
		330: "<hb-drawing>(drawing)</hb-drawing>",
		331: "<hb-wavy_line_as_divider />〜〜〜〜〜",
		332: "<hb-impression_of_stamp_on_papyrus>⦻</hb-impression_of_stamp_on_papyrus>",
		333: "<hb-text_enclosed_in_box_or_circle />",
		334: "<hb-text_enclosed_in_brackets />",
		335: "<span class=\"strikethrough\">N</span>",
		336: "<hum_redundant_s-type_sign />",
		337: "<hb-seal_attached_to_papyrus>❊</hb-seal_attached_to_papyrus>",
		400: "ͱ", // heta; supposed to know how to do capital too: *//400
		401: "ϳ", // yot; supposed to know how to do capital too: *//401
		451: "\u0283",
		452: "\u2310",
		453: "\u2e11",
		454: "\u2e10",
		456: "\u2e0e",
		457: "<hb-sp-unk>(＃457)</hb-sp-unk>",
		458: "\u0387",
		459: "\u00b7",
		460: "\u2014",
		461: "\u007c",
		465: "\u2627",
		466: "<hb-sp-unk>(＃466)</hb-sp-unk>",
		467: "\u2192",
		468: "\u2e0e",
		476: "\u0283",
		// 486: "<hb-undoc-pound betacodeval="486">⊚</hb-undoc-pound>",
		500: "⋮",
		// 500: "<hb-undoc-pound betacodeval="500">⊚</hb-undoc-pound>",
		501: "π<6ιθ6>", // abbreviation for πιθανόν: added own betacode - <6...6>
		502: "🜚",       // listed as idiosyncratic; but looks like 'alchemical symbol for gold': U+1F71A
		503: "ΡΠ",      // but supposed to be on top of one another
		504: "\u2e0e",
		505: "\u205c",
		506: "\u2e15",
		507: "\u2e14",
		508: "\u203b",
		509: "\u0305\u0311",
		510: "ε/π", // but supposed to be stacked
		511: "ι/κ", // but supposed to be stacked
		512: "\u03fd",
		513: "<hb-sp-unk>(＃513)</hb-sp-unk>",
		514: "<hb-sp-unk>(＃514)</hb-sp-unk>",
		515: "𐆅",
		516: "\u0394\u0345",
		517: "𐆅",
		518: "𐅹",
		519: "\u2191",
		520: "\u2629",
		521: "<hb-sp-unk>(＃521)</hb-sp-unk>",
		522: "<span class=\"90degreerotate\">Η</span>", // markup <rotate> 0397
		523: "\u2e13",
		524: "\u2297",
		526: "\u2190",
		527: "\u02c6",
		528: "\u03bb\u032d",
		529: "\u204b",
		530: "<hb-sp-unk>(＃530)</hb-sp-unk>",
		531: "\u035c",
		532: "\u2e12",
		533: "\u03da",
		534: "\u0302",
		535: "<hb-sp-unk>(＃535)</hb-sp-unk>",
		536: "<hb-sp-unk>(＃536)</hb-sp-unk>",
		537: "<hb-sp-unk>(＃537)</hb-sp-unk>",
		538: "<hb-sp-unk>(＃538)</hb-sp-unk>",
		540: "<hb-sp-unk>(＃540)</hb-sp-unk>",
		541: "<hb-sp-unk>(＃541)</hb-sp-unk>",
		542: "\u03a1\u0336",
		543: "<hb-sp-unk>(＃543)</hb-sp-unk>",
		544: "\u2058",
		545: "<hb-sp-unk>(＃545)</hb-sp-unk>",
		546: "<hb-sp-unk>(＃546)</hb-sp-unk>",
		547: "<hb-sp-unk>(＃547)</hb-sp-unk>",
		548: "<hb-sp-supsc>‖̴</hb-sp-supsc>", // 2016+0334
		549: "<hb-sp-unk>(＃549)</hb-sp-unk>",
		550: "\u003a\u003a\u2e2e",
		551: "\u25cc",
		552: "<hb-sp-unk>(＃552)</hb-sp-unk>",
		553: "<hb-sp-unk>(＃553)</hb-sp-unk>",
		554: "<hb-sp-unk>(＃554)</hb-sp-unk>",
		555: "<hb-sp-unk>(＃555)</hb-sp-unk>",
		556: "\u2629",
		557: "<hb-sp-unk>(＃557)</hb-sp-unk>",
		558: "<hb-sp-unk>(＃558)</hb-sp-unk>",
		559: "<hb-sp-unk>(＃559)</hb-sp-unk>",
		561: "\u2191",
		562: "\u0305",
		563: "⏗",
		564: "⏘",
		565: "⏙",
		// GREEK INSTRUMENTAL NOTATIONS
		566: "𝈱", // 32
		567: "𝈓",
		568: "𝈳",
		569: "𝈶", // 40
		570: "\u03f9",
		571: "𐅃",
		572: "𝈩", // 𝈩GREEK INSTRUMENTAL NOTATION SYMBOL-19; Unicode: U+1D229, UTF-8: F0 9D 88 A9
		573: "𝈒",
		574: "\u0393",
		575: "𝈕",
		576: "𝈖",
		577: "\u03a6",
		578: "\u03a1",
		579: "\u039c",
		580: "\u0399",
		581: "\u0398",
		582: "𝈍",
		583: "\u039d",
		584: "\u2127",
		585: "\u0396",
		586: "𝈸", // 43
		587: "\u0395",
		588: "𝈈", // Vocal //9' Instrum //44
		589: "𝈿", // 𝈿GREEK INSTRUMENTAL NOTATION SYMBOL-52; Unicode: U+1D23F, UTF-8: F0 9D 88 BF
		590: "𝈿",
		591: "𝈛",
		592: "𝉀",
		593: "\u039b",
		594: "<hb-sp-unk>(＃594)</hb-sp-unk>", // 'rare' and so has no unicode representation
		595: "<hb-sp-unk>(＃595)</hb-sp-unk>", // 'rare' and so has no unicode representation
		596: "<hb-sp-unk>(＃596)</hb-sp-unk>", // 'rare' and so has no unicode representation
		597: "<hb-sp-unk>(＃597)</hb-sp-unk>", // 'rare' and so has no unicode representation
		598: "\u0394",
		599: "𝈔",
		600: "𝈨", // Instrum //18
		// 601: "", // 'rare' and so has no unicode representation
		602: "𝈷",
		603: "\u03a0",
		604: "𝈦",                             // 𝈦GREEK INSTRUMENTAL NOTATION SYMBOL-14; Unicode: U+1D226, UTF-8: F0 9D 88 A6
		605: "<hb-sp-unk>(＃605)</hb-sp-unk>", // 'rare' and so has no unicode representation
		606: "<hb-sp-unk>(＃606)</hb-sp-unk>", // 'rare' and so has no unicode representation
		607: "<hb-sp-unk>(＃607)</hb-sp-unk>", // 'rare' and so has no unicode representation
		608: "<hb-sp-unk>(＃608)</hb-sp-unk>", // 'rare' and so has no unicode representation
		609: "<hb-sp-unk>(＃609)</hb-sp-unk>", // 'rare' and so has no unicode representation
		610: "<hb-sp-unk>(＃610)</hb-sp-unk>", // 'rare' and so has no unicode representation
		611: "<hb-sp-unk>(＃611)</hb-sp-unk>", // 'rare' and so has no unicode representation
		612: "<hb-sp-unk>(＃612)</hb-sp-unk>", // 'rare' and so has no unicode representation
		613: "<hb-sp-unk>(＃613)</hb-sp-unk>", // 'rare' and so has no unicode representation
		614: "<hb-sp-unk>(＃614)</hb-sp-unk>", // 'rare' and so has no unicode representation
		615: "𝈰",                             // 𝈰GREEK INSTRUMENTAL NOTATION SYMBOL-30; Unicode: U+1D230, UTF-8: F0 9D 88 B0
		616: "𝈞",
		617: "Ω",
		618: "<hb-sp-unk>(＃618)</hb-sp-unk>", // 'rare' and so has no unicode representation
		619: "λ",
		620: "<hb-sp-unk>(＃620)</hb-sp-unk>", // 'rare' and so has no unicode representation
		621: "𝈅",
		622: "𝈁",
		623: "\u2127",
		624: "\u03fd",
		625: "<hb-sp-unk>(＃625)</hb-sp-unk>", // 'rare' and so has no unicode representation
		626: "<hb-sp-unk>(＃626)</hb-sp-unk>", // 'rare' and so has no unicode representation
		627: "𝈗",
		628: "Ο",
		629: "Ξ",
		630: "Δ",
		631: "\u039a",
		632: "𝈎",
		633: "𝈲",
		634: "𝈹",
		635: "𝈝", // 𝈝GREEK INSTRUMENTAL NOTATION SYMBOL-1; Unicode: U+1D21D, UTF-8: F0 9D 88 9D
		636: "𝈃",
		637: "𝈆",
		638: "𝈉",
		639: "𝈌",
		640: "𝈑",
		641: "Ω",
		642: "Η",
		643: "𝈝",
		644: "𝈟",
		645: "𝈡",
		646: "𝈥",
		647: "𝈬",
		648: "𝈵",
		649: "𝈋",
		650: "𝈏",
		651: "\u03a7",
		652: "\u03a4",
		653: "𝈙",
		654: "𝈜",
		655: "𝈂",
		656: "𝈤",
		657: "𝈮",
		658: "𝈾",
		659: "𝉁",
		660: "\u0391",
		661: "\u0392",
		662: "\u03a5",
		663: "\u03a8",
		664: "𝈺",
		665: "𝈴", // 𝈴GREEK INSTRUMENTAL NOTATION SYMBOL-38; Unicode: U+1D234, UTF-8: F0 9D 88 B4
		666: "𝈯", // 𝈯GREEK INSTRUMENTAL NOTATION SYMBOL-29; Unicode: U+1D22F, UTF-8: F0 9D 88 AF
		667: "𝈭",
		668: "𝈐",
		669: "𝈊",
		670: "𝈇",
		671: "𝈛",
		672: "𝈘",
		673: "𝈣",
		674: "𝈢",
		675: "𝉀",
		676: "𝈽",
		677: "μ",
		678: "𝈠",
		679: "𝈄",
		681: "<hb-sp-unk>(＃681)</hb-sp-unk>",
		682: "<hb-sp-unk>(＃682)</hb-sp-unk>",
		683: "\u2733",
		684: "𝈪",
		// 685: "", // 'rare' and so has no unicode representation
		// 686: "", // 'rare' and so has no unicode representation
		// 687: "", // 'rare' and so has no unicode representation
		688: "\u03bc\u030a",
		689: "𐅵",
		690: "\u27d8",
		691: "\u27c0",
		692: "\u27c1",
		693: "<hb-sp-unk>(＃693)</hb-sp-unk>",
		694: "𝈼",
		695: "—", // em-dash 2014
		696: "𝈧",
		697: "𝉅",
		700: "\u205e",
		701: "<hb-sp-unk>(＃701)</hb-sp-unk>",
		702: "<hb-sp-unk>(＃702)</hb-sp-unk>",
		703: "\u25cb\u25cb\u25cb",
		704: "\u2014\u0307",
		705: "<hb-sp-unk>(＃705)</hb-sp-unk>",
		706: "<hb-sp-unk>(＃706)</hb-sp-unk>",
		707: "<hb-sp-unk>(＃707)</hb-sp-unk>",
		708: "<hb-sp-unk>(＃708)</hb-sp-unk>",
		709: "\u223b",
		710: "\u039a\u0336",
		711: "\u03fb",
		741: "<hb-sp-unk>(＃741)</hb-sp-unk>",
		751: "\u0661", // arabic-indic digits...
		752: "\u0662",
		753: "\u0663",
		754: "\u0664",
		755: "\u0665",
		756: "\u0666",
		757: "\u0667",
		758: "\u0668",
		759: "\u0669",
		760: "\u0660",
		800: "\u2733",
		801: "𐅁",
		802: "𐅀",
		803: "\u03a7",
		804: "／", // fulwidth solidus instead
		805: "\u03a4",
		806: "\u039A",
		807: "𐅦",
		808: "𐅈",
		809: "𐅃", // http://noapplet.epigraphy.packhum.org/text/335386?&bookid=859&location=16
		811: "\u03a4",
		812: "𐅈",
		813: "𐅉",
		814: "𐅊",
		815: "𐅋",
		816: "𐅌",
		817: "𐅍",
		818: "𐅎",
		821: "\u03a3",
		822: "𐅟",
		823: "𐅐",
		824: "𐅑",
		825: "𐅒",
		826: "𐅓",
		827: "𐅔",
		829: "𐅕",
		830: "𐅇",
		831: "𐅇",
		832: "𐅖",
		833: "\u039c",
		834: "𐅗",
		835: "\u03a7",
		836: "\u03a3",
		837: "\u03a4",
		838: "𐅃",
		839: "𐅁",
		840: "\u007c\u007c",
		841: "\u007c\u007c\u007c",
		842: "\u00b7",
		843: "𐅛",
		844: "\u205d",
		845: "𐅘",
		846: "𐄐",
		847: "𐅞",
		848: "𐄒",
		// 850: [not known by PHI]
		850: "<hb-sp-unk>(＃850)</hb-sp-unk>",
		853: "\u0399",
		862: "\u0394",
		863: "𐅄",
		865: "𐅅",
		866: "\u03a7",
		867: "𐅆",
		870: "Ε",
		// PHI will show you glyphs, but 'private use area' means that they are feeding them to you
		// it seems that there is no official support for these characters
		// 875: "", // see http://noapplet.epigraphy.packhum.org/text/247092?&bookid=489&location=1689; private use area
		874: "<hb-sp-unk>(＃874)</hb-sp-unk>", // u'\ue022", // private use area
		875: "<hb-sp-unk>(＃875)</hb-sp-unk>", // u'\ue022", // private use area
		876: "<hb-sp-unk>(＃876)</hb-sp-unk>", // u'\ue023", // private use area
		877: "<hb-sp-unk>(＃877)</hb-sp-unk>", // u'\ue024", // inferred; private use area
		878: "<hb-sp-unk>(＃878)</hb-sp-unk>", // u'\ue025", // inferred; private use area
		879: "<hb-sp-unk>(＃879)</hb-sp-unk>", // u'\ue026",  // inferred; private use area
		880: "<hb-sp-unk>(＃880)</hb-sp-unk>", // u'\ue027",  // inferred; private use area
		881: "<hb-sp-unk>(＃881)</hb-sp-unk>", // u'\ue028",  // inferred; private use area
		882: "<hb-sp-unk>(＃882)</hb-sp-unk>", // u'\ue029",  // inferred; private use area
		883: "<hb-sp-unk>(＃883)</hb-sp-unk>", // u'\ue02a", // private use area
		897: "<hb-sp-unk>(＃897)</hb-sp-unk>",
		// 898: [not known by PHI]
		898: "<hb-sp-unk>(＃898)</hb-sp-unk>",
		899: "<hmu_unknown_numeral>",
		900: "○", // '❦' (?!) is what you can see at http://noapplet.epigraphy.packhum.org/text/232427?&bookid=396&location=7
		// 901: [not known by PHI]
		901: "<hb-sp-unk>(＃901)</hb-sp-unk>",
		// 921: [not known by PHI]
		921: "<hb-sp-unk>(＃921)</hb-sp-unk>",
		922: "𝈨",
		923: "<hb-sp-unk>(＃923)</hb-sp-unk>",
		924: "<hb-sp-unk>(＃924)</hb-sp-unk>",
		925: "𝈗",
		926: "𝈫",
		927: "W",
		928: "𝈋",
		929: "𝈔",
		930: "<hb-sp-unk>(＃930)</hb-sp-unk>",
		932: "\u2733",
		933: "<hb-sp-unk>(＃933)</hb-sp-unk>",
		934: "<hb-sp-unk>(＃934)</hb-sp-unk>",
		// 936: [not known by PHI]
		936: "<hb-sp-unk>(＃936)</hb-sp-unk>",
		937: "<hmu_miscellaneous_illustrations>",
		// 938: "", // http://noapplet.epigraphy.packhum.org/text/260647?&bookid=509&location=1035; private use area?
		938: "Ƨ", // 01a7
		939: "~",
		940: "<hb-sp-unk>(＃940)</hb-sp-unk>",
		943: "﹥", // PHI
		945: "<hb-sp-unk>(＃945)</hb-sp-unk>",
		946: "<hb-sp-unk>(＃946)</hb-sp-unk>",
		// 947: [not known by PHI]
		947: "<hb-sp-unk>(＃947)</hb-sp-unk>",
		// 948: [not known by PHI]
		948: "<hb-sp-unk>(＃948)</hb-sp-unk>",
		949: "—", // http://noapplet.epigraphy.packhum.org/text/251612?&bookid=491&location=1689
		950: "<hb-sp-unk>(＃950)</hb-sp-unk>",
		961: "<hmu_line_on_stone_stops_but_edition_continues_line />",
		971: "<hb-sp-unk>(＃971)</hb-sp-unk>",
		972: "<hb-sp-unk>(＃972)</hb-sp-unk>",
		// 973: [not known by PHI]
		973: "<hb-sp-unk>(＃973)</hb-sp-unk>",
		975: "<hb-sp-unk>(＃975)</hb-sp-unk>",
		977: "§", // Caria (Stratonikeia), 8 2, line 12; http://noapplet.epigraphy.packhum.org/text/262496?&bookid=526&location=1035
		// 990: "<hb-undoc-pound betacodeval="990">⊚</hb-undoc-pound>",
		// 981: [not known by PHI]
		981: "<hb-sp-unk>(＃981)</hb-sp-unk>",
		// 982: [not known by PHI]
		982:  "<hb-sp-unk>(＃982)</hb-sp-unk>",
		1000: "𐅼",
		1001: "𐅽",
		1002: "𐅾",
		1003: "𐅿",
		1004: "𐆀",
		1009: "", // http://noapplet.epigraphy.packhum.org/text/247092?&bookid=489&location=1689
		// a huge run of undocumented poundsigns in the inscriptions: this only scratches the surface
		// packhum.org has representations of many of them
		// see especially: http://noapplet.epigraphy.packhum.org/text/260603?&bookid=509&location=1035
		1012: "\ue036",
		1023: "ηʹ", // http://noapplet.epigraphy.packhum.org/text/247092?&bookid=489&location=1689
		1024: "ΛΒ", // Greek Capital Letter and Greek Capital Letter Beta
		// 1045: [not known by PHI]
		1045: "<hb-sp-unk>(＃1045)</hb-sp-unk>",
		1053: "<hb-sp-unk>(＃1053)</hb-sp-unk>",
		// 1057: "", // http://noapplet.epigraphy.packhum.org/text/258019?&bookid=493&location=1035; private use area?
		1056: "<hb-sp-unk>(＃1056)</hb-sp-unk>",
		1057: "<hb-sp-unk>(＃1057)</hb-sp-unk>",
		1059: "<hb-sp-unk>(＃1059)</hb-sp-unk>",
		1061: "γʹ",
		1062: "δʹ",
		1063: "εʹ",
		1064: "ϛʹ",
		1065: "ζʹ",
		1067: "θʹ",
		1068: "ιʹ",
		1069: "κʹ",
		1070: "λʹ",
		1071: "μʹ",
		1072: "νʹ",
		1073: "ξʹ",
		1074: "οʹ",
		1075: "πʹ",
		1077: "ρʹ",
		1078: "σʹ",
		1079: "τʹ",
		1080: "υʹ",
		1082: "χʹ",
		1084: "ωʹ", // Caria (Tralles), 243: line 16; http://noapplet.epigraphy.packhum.org/text/263093?&bookid=531&location=1035
		1085: " ϡʹ",
		1086: "͵α",
		1087: "͵β",
		1100: "\u2183",
		1101: "IS",
		1102: "H",
		1103: "\u0323\u0313",
		1104: "S\u0038", // deprecated, use &S%162$
		1105: "\u004d\u030a",
		1106: "<hb-sp-unk>(＃1106)</hb-sp-unk>",
		1107: "\u0053\u0335\u0053\u0336",
		1108: "\u0058\u0036",
		1109: "\u003d",
		1110: "\u002d",
		1111: "\u00b0",
		1112: "<hb-sp-unk>(＃1112)</hb-sp-unk>",
		1113: "<hb-sp-unk>(＃1113)</hb-sp-unk>",
		1114: "𝈁",
		1115: "\u007c",
		1116: "\u01a7",
		1117: "\u005a",
		1118: "<hb-sp-unk>(＃1118)</hb-sp-unk>",
		1119: "\u0110",
		1120: "<hb-sp-unk>(＃1120)</hb-sp-unk>",
		1121: "\u005a",
		1122: "<hb-sp-unk>(＃1122)</hb-sp-unk>",
		1123: "<hb-sp-unk>(＃1123)</hb-sp-unk>",
		1124: "\u211e",
		1125: "<hb-sp-unk>(＃1125)</hb-sp-unk>",
		1126: "\u004f",
		1127: "\u0076\u0338",
		1128: "\u0049\u0336\u0049\u0336\u0053\u0336",
		1129: "\u005a\u0336",
		1130: "＼", // fullwidth reverse solidus (vs just reverse)
		1131: "\u005c\u005c",
		1132: "\u005c\u0336",
		1133: "\u005c\u0336\u005c\u0336",
		1134: "<hb-sp-unk>(＃1134)</hb-sp-unk>",
		1135: "\u002f\u002f",
		1136: "\u2112",
		1221: "\u0131",
		1222: "\u0130",
		1314: "\u006e\u030a",
		1315: "ΜΡ", // but supposed to be on top of one another
		1316: "\u0292",
		1317: "\u02d9\002f\u002f\u002e",
		1318: "\u223b",
		1320: "\u0375\u0311",
		1321: "🜚", // listed as idiosyncratic; but looks like 'alchemical symbol for gold': U+1F71A
		1322: "\u2644",
		1323: "\u03b6\u0337\u03c2\u0300",
		1324: "\u03b8\u03c2\u0302",
		1326: "<hb-sp-unk>(＃1326)</hb-sp-unk>",
		1327: "<hb-sp-unk>(＃1327)</hb-sp-unk>",
		1328: "<hb-sp-unk>(＃1328)</hb-sp-unk>",
		1334: "<hb-sp-unk>(＃1334)</hb-sp-unk>",
		1335: "／／", // fulwidth solidus instead
		1336: "<hmu_unsupported_hebrew_character>□</hmu_unsupported_hebrew_character>",
		1337: "﹥", // supposed to be 003e, ie simple angle bracket ; this is fe65
		1338: "𐅾",
		1341: "<hb-sp-unk>(＃1341)</hb-sp-unk>",
		1500: "\u03b3\u030a",
		1501: "<hb-sp-unk>(＃1501)</hb-sp-unk>",
		1502: "\u03a7\u0374",
		1503: "<hb-sp-unk>(＃1503)</hb-sp-unk>",
		1504: "<hb-sp-unk>(＃1504)</hb-sp-unk>",
		1505: "<hb-unknown_abbreviation betacodeval=\"1505\">◦</hmu_unknown_abbreviation>",
		1506: "\u0300\u0306",
		1509: "πληθ",                     // supposed to be a symbol
		1510: "Α\u0338<6\u0304ν\u002f>6", // A%162<6E%26N%3>6 [!]
		1511: "π<hb-sp-supsc>ε:`</hb-sp-supsc>",
		// 1806: [not known by PHI]
		1806: "<hb-sp-unk>(＃1806)</hb-sp-unk>",
	}
)

func ReplacePoundSigns(ttc string) string {
	// Purge # markup
	// note that this has to run before the check for "<6" because we will insert some of that superscript markup
	// it is a shame that you cannot set a match group to the fnc but instead have to send a whole match...

	// see also the notes at cleanembeddedvalue() in hexrunner.go

	ttc = poundsouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return poundsubstitutes(match)
	})

	return ttc
}

func poundsubstitutes(match string) string {
	// Format Additional Characters
	// #26-49 Reserved for Greek documentary papyri (common miscellaneous characters)
	// #50-69 Reserved for Greek inscriptions (punctuation)
	// #70-99 Reserved for Greek documentary papyri (punctuation)
	// #300-399 Reserved for Greek documentary papyri (miscellaneous characters)
	// #800-899 Reserved for Greek inscriptions (numeric)
	// #900-999 Reserved for Greek inscriptions (miscellaneous)
	// #1000-1099 Reserved for Greek inscriptions (numeric, continued)
	// #1100-1199 Reserved for Latin numeric
	// #1200-1299 Reserved for Founding Fathers Project

	const (
		UHP = `<hgb-build-error>＃%s<hgb-build-error>`
	)
	// match looks like "#14": drop that initial character
	// the real question is whether this is the fastest way to do it...
	// convert to rune, drop first, convert back is an alternative
	match = strings.TrimPrefix(match, "#")

	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := poundsubsmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("poundsubstitutes()\t%s", match))
	}
	return substitute

}
