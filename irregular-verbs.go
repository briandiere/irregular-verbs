package main

import (
    "fmt"
    "log"
    "os"
    "io/ioutil"
    "encoding/json"
    "math/rand"
    "time"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
)

const indexVerb = 0
const indexPreterit = 1
const indexPastPart = 2
const indexTradFr = 3

var errorColor = walk.RGB(80, 0, 0)
var validColor = walk.RGB(0, 80, 0)

type IrregularVerbs []IrregularVerb

type IrregularVerb struct {
    Verb         string `json:"base_verbale"`
    Preterit     string `json:"preterit"`
    PastParticip string `json:"participe_passe"`
    TranslateFr  string `json:"traduction_fr"`
    Actif        bool   `json:"actif"`
}

func IsVerbActif(verb IrregularVerb) bool {
    return verb.Actif
}

func Filter(verbs []IrregularVerb, test func(IrregularVerb) bool) (ret []IrregularVerb) {
    verbsValid  := make([]IrregularVerb, 0)
    for _, verb := range verbs {
        if test(verb) {
            verbsValid = append(verbsValid, verb)
        }
    }
    return verbsValid
}

func RecupVerbToLearn() (ret []IrregularVerb) {
	// Open our jsonFile
	jsonFile, err := os.Open("./irregular-verbs.json")
	// if we os.Open returns an error then handle it
	if err != nil {
	    fmt.Println(err)
	} else {
		fmt.Println("Lectures du fichier OK")
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var verbs IrregularVerbs
	json.Unmarshal(byteValue, &verbs)

	verbsToLearn := Filter(verbs, IsVerbActif)
	fmt.Println(len(verbsToLearn), "verbes à connaître sur un total de", len(verbs), "verbes.")

	//	for i := 0; i < len(verbs) && i < 1; i++ {
	//	}

	return verbsToLearn
}

func main() {
	mw := &MyMainWindow{model: NewIrregularVerbModel(RecupVerbToLearn())}

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Révisions des verbes irréguliers",
		MinSize:  Size{100, 50},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "Verbe:",},
			mw.InitTextEdit(mw.verbTE, indexVerb, mw.model.indexRandom, mw.model.items[mw.model.indexCourant].Verb),
			Label{Text: "Prétérit:",},
			mw.InitTextEdit(mw.preteritTE, indexPreterit, mw.model.indexRandom, mw.model.items[mw.model.indexCourant].Preterit),
			Label{Text: "Participe passé:",},
			mw.InitTextEdit(mw.pastPartTE, indexPastPart, mw.model.indexRandom, mw.model.items[mw.model.indexCourant].PastParticip),
			Label{Text: "Traduction:",},
			mw.InitTextEdit(mw.traducFrTE, indexTradFr, mw.model.indexRandom, mw.model.items[mw.model.indexCourant].TranslateFr),
			PushButton{
				AssignTo: &mw.validButton, 
				Text: "Valider",
				Visible: true,
				OnClicked: mw.ValiderResultats,
			},
			PushButton{
				AssignTo: &mw.nextButton, 
				Text: "Suivant",
				Visible: false,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

type MyMainWindow struct {
	*walk.MainWindow
	model         *IrregularVerbModel
	verbTE        *walk.TextEdit
	preteritTE    *walk.TextEdit
	pastPartTE    *walk.TextEdit
	traducFrTE    *walk.TextEdit
	validButton   *walk.PushButton
	nextButton    *walk.PushButton
}

type IrregularVerbModel struct {
	items         []IrregularVerb
	indexCourant  int
	indexRandom   int
	nbTotal       int
	nbRepValide   int
}

func NewIrregularVerbModel(verbs []IrregularVerb) *IrregularVerbModel {
	rand.Seed(time.Now().UnixNano())
	return &IrregularVerbModel {
		items: verbs,
		indexCourant: rand.Intn(len(verbs)),
		indexRandom: rand.Intn(4),
		nbTotal: len(verbs) * 3,
		nbRepValide: 0,
	}
}

func (mw *MyMainWindow) InitTextEdit(assignTo *walk.TextEdit, iConst int, iRandom int, value string) (ret TextEdit) {
	var textTE = ""
	var reaOnlyTE = iRandom == iConst

	if reaOnlyTE {
		textTE = value
	}

	return TextEdit{
		AssignTo: &assignTo,
		ReadOnly: reaOnlyTE,
		Text: textTE,
	}
}

func (mw *MyMainWindow) ValiderResultats() {
	item := &mw.model.items[mw.model.indexCourant]
	VerifierChamp(mw.verbTE, item.Verb)
	VerifierChamp(mw.preteritTE, item.Preterit)
	VerifierChamp(mw.pastPartTE, item.PastParticip)
	VerifierChamp(mw.traducFrTE, item.TranslateFr)
}

func VerifierChamp(champ *walk.TextEdit, expectedValue string) () {
	if (! champ.ReadOnly()) {
		champ.SetReadOnly(true)
		if (strings.ToLower(champ.Text()) == strings.ToLower(expectedValue)) {
			champ.SetTextColor(errorColor)
		} else {
			champ.SetTextColor(validColor)
			champ.AppendText(" --> " + expectedValue)
		}
	}
	return
}

