//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lat

import (
	"fmt"
	"testing"
)

func TestConvertLatinDiacriticals(t *testing.T) {
	ttc := `u=V=i+e/o\`
	result := ConvertLatinDiacriticals(ttc)
	want := "ûÛïéò"
	if result != want {
		t.Errorf("ConvertLatinDiacriticals(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestRewritegreekfontshift(t *testing.T) {
	ttc := `L. Aelium magistrum suum in $E)TUMOLOGI/A|<hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	result := rewritegreekfontshift(ttc)
	want := `L. Aelium magistrum suum in <hb-fs-g-normal>E)TUMOLOGI/A|</hb-fs-g-normal><hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	if result != want {
		t.Errorf("rewritegreekfontshift(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `$KU/BON<hb-fs-l-normal>, quid </hb-fs-l-normal>GRAMMH/N<hb-fs-l-normal>; quibusque ista omnia Latinis uoca-</hb-fs-l-normal> █`
	result = rewritegreekfontshift(ttc)
	want = `<hb-fs-g-normal>KU/BON</hb-fs-g-normal><hb-fs-l-normal>, quid </hb-fs-l-normal><hb-fs-g-normal>GRAMMH/N</hb-fs-g-normal><hb-fs-l-normal>; quibusque ista omnia Latinis uoca-</hb-fs-l-normal> █`
	if result != want {
		t.Errorf("rewritegreekfontshift(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `<hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext />$EU)FHMEI=N XRH\ KA)CI/STASQAI TOI=S H(METE/ROISI XOROI=SIN, █`
	result = rewritegreekfontshift(ttc)
	want = `<hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-fs-g-normal>EU)FHMEI=N XRH\ KA)CI/STASQAI TOI=S H(METE/ROISI XOROI=SIN,</hb-fs-g-normal> █`
	if result != want {
		t.Errorf("rewritegreekfontshift(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestApplybetacodeconversion(t *testing.T) {
	ttc := `L. Aelium magistrum suum in <hb-fs-g-normal>E)TUMOLOGI/A|</hb-fs-g-normal><hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	result := applybetacodeconversion1(ttc)
	want := `L. Aelium magistrum suum in <hb-fs-g-normal>ἐτυμολογίᾳ</hb-fs-g-normal><hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	if result != want {
		t.Errorf("applybetacodeconversion(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestLatinCleanup(t *testing.T) {
	ttc := `{&7Pa.&} censeo. {&7Ph.&} sed heus tu. {&7Pa.&} quid vis? {&7Ph.&} censen posse me obfirmare et `
	result := LatinCleanup(ttc)
	want := `{<hb-fs-l-smallcapitals>Pa.</hb-fs-l-smallcapitals><hb-fs-l-normal>} censeo. {</hb-fs-l-normal><hb-fs-l-smallcapitals>Ph.</hb-fs-l-smallcapitals>} sed heus tu. {<hb-fs-l-smallcapitals>Pa.</hb-fs-l-smallcapitals><hb-fs-l-normal>} quid vis? {</hb-fs-l-normal><hb-fs-l-smallcapitals>Ph.</hb-fs-l-smallcapitals>} censen posse me obfirmare et `
	if result != want {
		t.Errorf("LatinCleanup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestLatinCleanup2(t *testing.T) {
	// ttc := ` § █⑧⓪ <hb-tabbedtext /> █⑨⑧ █⑧ⓓ █⑧⑧ █⑧ⓔ 1 $⟨QA/L⟩ASSA[N] PEIRATEUOME/NHN U(PO\ A)POSTATW=N DOU/-█⑧⓪ $LWN [EI)⟨RH/N]EUSA⟩: ⟨E)C⟩ ⟨W(=N⟩ ⟨TREI=S⟩ ⟨POU⟩ ⟨MURIA/D⟩AS TOI=S █⑧⓪ $DE[SPO/TAI]S EI)S KO/LASIN PARE/DWKA. <hb-fs-l-normal>2 </hb-fs-l-normal>W)/MO⟨SEN⟩ █⑧⓪ $[⟨EI)S⟩ ⟨TOU\S⟩ ⟨E)MOU\]S⟩ ⟨LO/GOUS⟩ A(/PASA H( *)ITALI/A E(KOU=SA KA)-█⑧⓪ $[ME\ POLE/MOU], W(=[I] E)P＇ *)AKTI/⟨W⟩I ⟨E)NE[I/]KHSA⟩, ⟨H(GEMO/NA⟩ ⟨E)C⟩H|-█⑧⓪ $[TH/SATO: W)/]MOSAN EI)S TOU\S [AU)TOU\]S LO/GOUS E)⟨PA[R]⟩-█⑧⓪ $⟨XE[I/AI⟩ ⟨*GALA]TI/A⟩ ⟨*(ISPANI/A⟩ ⟨*LI⟩BU/H *SI[KELI/A *SAR]DW/. <hb-fs-l-normal>3 </hb-fs-l-normal>OI( U(P＇ E)-█⑧⓪ $M[AI=S SHME/AIS TO/]TE STRATEU[SA/ME⟨NOI⟩ ⟨H)=SAN⟩ ⟨SUNKLHTI]⟩-█⑧⓪ $⟨[KOI\⟩ ⟨PLE⟩I/OUS E(PT]A[KOSI/]WN: [E)]N [AU)TOI=S OI(\ H)\ PRO/TERON H)\] █⑧⓪ [$METE/PEI⟨TA⟩] ⟨E)G[E/NON]TO⟩ [⟨U(/P]A[TOI⟩ ⟨A)/XRI⟩ ⟨E)⟩K]E[I/]N[HS TH=S H(]ME/-█⑧⓪ $[RAS E)N H(=I TAU=TA GE/GRAPTA]I O)[GDOH/⟨KO]NTA⟩ ⟨TRE[I=]S⟩, ⟨I(ER[EI=]S █⑨⓪ $⟨PRO/S1⟩POU E(KATO\N E(BDOMH/[K]ONTA. █⑧⓪ <hb-tabbedtext /> █ⓐ⓪ █⑨ⓕ █ⓓ⑥ █ⓕⓕ █⑧⑧ █⑧⑨ 1 omniu⟨m⟩ ⟨pr⟩o/v[inciarum populi Romani], quibus fi⟨niti⟩mae fuer⟨unt⟩ █⑧⓪ gente/s qu⟨ae⟩ ⟨n[o/n⟩ ⟨p﹖⟩arerent imperio nos]tro, fines auxi. 2 G⟨all⟩ias et Hispa-█⑧⓪ ⟨nia/⟩s pro/v⟨i⟨n⟩cia/[s⟩, ⟨i﹖⟩tem Germaniam qua inclu]dit O/ceanus a Ga/dibus ad o/sti-█⑧⓪ um Albis f⟨lu/m[in⟩is pacavi. § 3 Alpes a re]gio/ne ea/, quae proxima est Ha-█⑧⓪ dria/n⟨o/⟩ ⟨mar⟩i/, [ad Tuscum pacari fec]i nulli/ genti/ bello per iniu/riam █⑧⓪ inla/to. § 4 cla[ssis m⟨ea⟩ ⟨p⟩er Oceanum] ab o/stio/ Rhe/ni ad so/lis orientis re-█⑧⓪ gionem usque ad fi[nes Cimbroru]m navigavit, [§] quo/ neque terra neque █⑧⓪ mari quisquam Romanus ante id tempus adi/t, § Cimbrique et Charydes █ⓕⓔ █⓪ █ⓔⓕ █⑧⓪ █ⓑ⓪ █ⓑ① █ⓑ② █ⓑ⓪ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ① █ⓕⓕ █ⓔⓕ █⑧② █ⓓ② █ⓒ⑦ █ⓕⓕ █ⓔⓕ █⑧③ █ⓒⓒ █ⓔ① █ⓕ④ █ⓔ⑨ █ⓔⓔ █ⓐ⓪ █ⓓⓑ █ⓔⓕ █ⓕ④ █ⓔ⑧ █ⓔ⑤ █ⓕ② █ⓓⓓ █ⓕⓕ █ⓐ⑧ █⑨ⓐ █⑨ⓕ █ⓓ⑥ █ⓕⓕ █⑧⑧ █⑨① █ⓔⓕ █ⓔ① █ⓒ⑦ █ⓔ① █ⓔⓒ █ⓔ① █ⓕ④ █ⓔ⑨ █ⓔ① █ⓕⓕ █ⓔⓕ █ⓔ② █ⓒ① █ⓔⓔ █ⓔⓑ █ⓕ⑨ █ⓕ② █ⓔ① █ⓕⓕ █ⓔⓕ █ⓔ③ █ⓕⓕ █ⓔⓐ █ⓔ④ █⑨③ █ⓐ⓪ █ⓔ① █ⓔ③ █ⓕⓕ █ⓔⓕ █ⓔ⑤ █ⓕⓕ et Semnones et eiusdem tractu/s alii/ Germa/no/rum popu[l]i per lega/to/s amici-█⑧⓪`
	// ttc := `e/r[um exe]mpla imi-█⑧⓪ tan⟨da⟩ ⟨p⟩os[teris tradidi]. █⑧⓪ <hb-tabbedtext /> █⑨④ █⑧⑧ █⑧⑧ 1 $TW=N [PAT]RIKI/WN TO\N A)RIQMO\N EU)/CHSA <hb-sp-rectified_form>H)U/CHSA</hb-sp-rectified_form> PE/MPTON █⑧⓪ $U(/PAT[OS E)PIT]AGH=I TOU= TE DH/MOU KAI\ TH=S SUNKLH/-█⑧⓪ $TOU. <hb-fs-l-normal>§ 2 </hb-fs-l-normal>[TH\N SU/]NKLHTON TRI\S E)PE/LECA. § E(/KTON U(/PA-█⑧⓪ $TOS TH\N A)P[O]TEI/MHSIN TOU= DH/MOU SUNA/RXON-█⑧⓪ $[T]A E)/XWN *MA=RKON *)AGRI/PPAN E)/LABON, H(/TIS A)PO-█⑧⓪ $[TEI/MH]SIS META\ [DU\O KAI\] TESSARAKOSTO\N E)NIAU-█⑧⓪ $TO\N [S]UNE[K]LEI/SQH. E)N H(=I A)POTEIMH/SEI *(RWMAI/WN █⑧⓪ $E)TEI[MH/S]A[NTO] KEFALAI\ TETRAKO/[SIAI E(]CH/KON-█⑧⓪ $TA MU[RIA/DES KAI\ TRISXI/LIAI. <hb-fs-l-normal>3 </hb-fs-l-normal>EI)=TA DEU/TERON U(]PATI-█⑧⓪ $KH=I E)C[OUSI/AI MO/NOS, *GAI/+WI *KHNSWRI/NWI KAI\] █⑧⓪ $*GAI/+WI [*)ASINI/WI U(PA/TOIS, TH\N A)POTEI/MHSIN E)/LABON]: █⑧⓪ $E)N [H(=I] A)P[OTEIMH/SEI E)TEIMH/SANTO *(RWMAI/]-█⑧⓪ $WN TET[RAKO/SIAI EI)/KOSI TREI=S MURIA/DES KAI\ T]RI[S]-█⑧⓪ $XI/LIOI. <hb-fs-l-normal>4  </hb-fs-l-normal>K[AI\ TRI/TON U(PATIKH=I E)COUSI/AI TA\S A)POTEIMH/]-█⑧⓪ $SE[I]S E)/LA[BO]N, [E)/XW]N [SUNA/RXONTA *TIBE/RION] █⑧⓪ $*KAI/SARA, TO\N UI(O/N MO[U, *SE/CTWI *POMPHI/WI KAI\] █⑨⓪ $*SE/CTWI *)APPOULHI/WI U(PA/TOIS: E)N H(=I A)POTEIMH/SEI █⑧⓪ $E)TEIMH/SANTO *(RWMAI/WN TETRAKO/SIAI E)NENH/KONTA █⑧⓪ $TREI=S MURIA/DES KAI\ E(PTAKISXEI/LIOI. <hb-fs-l-normal>§ 5 </hb-fs-l-normal>EI)SAGAGW\N KAI-█⑧⓪ $NOU\S NO/MOUS POLLA\ H)/DH TW=N A)RXAI/WN E)QW=N KA-█⑧⓪ $TALUO/MENA DIWRQWSA/MHN KAI\ AU)TO\S POLLW=N █⑧⓪ $PRAGMA/TWN MEI/MHMA E)MAUTO\N TOI=S METE/PEI-█⑧⓪ $TA PARE/DWKA. █⑧⓪ <hb-tabbedtext /> █ⓐ⓪ █⑨ⓕ █ⓒ⑨ █ⓒ⑨ █ⓕⓕ █⑧⑧ █⑧ⓕ 1 [⟨vo/ta⟩ ⟨p﹖⟩ro valetudine mea susc⟨i﹖p﹖i﹖⟩ ⟨p﹖⟩er cons]ule/s et sacerdotes qu[in⟨t﹖o﹖⟩] █⑧⓪ [⟨qu⟩oque anno senatus decrevit. ex iis] votis s[ae]pe fe⟨cer⟩unt ⟨vi/vo⟩ █⑧⓪ [⟨m⟩e ludos aliquotiens sacerdotu]m quatt⟨uo⟩r am⟨pliss⟩ima col⟨le/⟩-█⑧⓪ [gia, aliquotiens consules.`
	ttc := ` P. Su⟨lpicio⟩ C. Va⟨lgi⟩o consulibu[s]. § █⑧⓪ <hb-tabbedtext /> █⑨⑤ █⑧⑧ █⑨⓪ 1 $TO\ O)/N[OM]A/ MOU SUNKLH/TOU DO/GMATI E)NPERIELH/-█⑧⓪ $FQH EI)[S TOU\]S *SALI/WN U(/MNOUS, KAI\ I(/NA I(ERO\S W)=I █⑧⓪ $DIA\ [BI/O]U [T]E TH\N DHMARXIKH\N E)/XWI E)COUSI/AN, █⑧⓪ $NO/[MWI E)K]URW/QH. <hb-fs-l-normal>§ 2 </hb-fs-l-normal>A)RXIERWSU/NHN, H(\N O( PATH/R █⑧⓪ $[M]OU [E)SX]H/KEI TOU= DH/MOU MOI KATAFE/RONTOS █⑧⓪ $EI)S TO\N TOU= ZW=NTOS TO/PON, OU) PROSEDECA/-█⑧⓪ $M[H]N. [H(\]N A)RXIERATEI/AN META/ TINAS E)NIAUTOU\S, █⑨⓪ $A)POQANO/NTOS TOU= PROKATEILHFO/TOS AU)-█⑧⓪ $TH\N E)N POLEITIKAI=S TARAXAI=S, A)NEI/LHFA, EI)S █⑧⓪ $TA\ E)MA\ A)RXAIRE/SIA E)C O(/LHS TH=S *)ITALI/AS TOSOU/-█⑧⓪ $TOU PLH/QOUS SUNELHLUQO/TOS, O(/SON OU)DEI\S █⑧⓪ $E)/NPROS1⟨QEN⟩ ⟨I(STO/RHS1＇⟩ <hb-fs-l-normal><hb-sp-alternative_reading>Apoll. </hb-fs-l-normal>I(STO/RHSEN</hb-sp-alternative_reading> ⟨E)PI\⟩ ⟨*(RW/MHS⟩ G⟨EGONE/NAI⟩, ⟨*PO⟩-█⑧⓪ $⟨PLI/W⟩I ⟨*SOULPIKI/W⟩I ⟨KAI\⟩ ⟨*GAI/+WI⟩ ⟨*OU)ALGI/WI⟩ ⟨U(PA/TOIS⟩. █⑧⓪ <hb-tabbedtext /> █ⓐ⓪ █⑨ⓕ █ⓒ⑨ █ⓒ⑨ █ⓕⓕ █⑧⑧ █⑨ⓓ [⟨aram⟩ Fortunae ⟨R﹖⟩educis a⟨nte⟩ ⟨ae⟩]de/s Honoris et Virtutis ad portam █⑧⓪ [⟨Cap⟩enam pro ⟨red⟩itu me⟨o/⟩ ⟨se]na/⟩tus consacravit, in qua ponti-█⑧⓪ [fices et ⟨vir﹖⟩gines Ve⟨stal⟩es anni]v⟨er⟩sa/rium sacrificium facere █⑧⓪ [iussit eo ⟨d﹖i﹖⟩e, quo co⟨n﹖sul﹖⟩ibus Q. Luc]retio et [M. Vi⟨nic﹖⟩i]o in urbem ex █⑧⓪ [Syria redieram, et diem Augustali]a ex [c]o[gnomine ⟨nos⟩t]ro appellavit.`
	result := LatinCleanup(ttc)
	//want := ``
	//if result != want {
	//	t.Errorf("LatinCleanup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	//}
	fmt.Println(result)
}
