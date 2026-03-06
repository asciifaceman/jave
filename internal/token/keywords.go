package token

var keywords = map[string]Kind{
	"outy":        Outy,
	"inny":        Inny,
	"sequence":    Sequence,
	"seq":         Seq,
	"give":        Give,
	"up":          Up,
	"maybe":       Maybe,
	"furthermore": Furthermore,
	"otherwise":   Otherwise,
	"given":       Given,
	"again":       Again,
	"allow":       Allow,
	"install":     Install,
	"from":        From,
	"within":      Within,
	"yee":         Yee,
	"nee":         Nee,
	"exact":       TypeExact,
	"vag":         TypeVag,
	"truther":     TypeTruther,
	"strang":      TypeStrang,
	"nada":        TypeNada,
	"naw":         TypeNaw,
}

// LookupIdentifier maps source text to either a keyword kind or Identifier.
func LookupIdentifier(lit string) Kind {
	if kind, ok := keywords[lit]; ok {
		return kind
	}
	return Identifier
}
