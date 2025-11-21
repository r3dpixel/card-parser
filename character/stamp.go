package character

type Spec string // chara card spec

const (
	SpecV2 Spec = "chara_card_v2"
	SpecV3 Spec = "chara_card_v3"
)

type Version string // chara card version

const (
	V2 Version = "2.0"
	V3 Version = "3.0"
)

type Revision int // chara card revision

const (
	RevisionV2 Revision = 2
	RevisionV3 Revision = 3
)

// Stamp structure of a mapping from revision to spec/version
type Stamp struct {
	Spec     Spec
	Version  Version
	Revision Revision
}

// Stamps mappings from revision to spec/versions
var Stamps = map[Revision]Stamp{
	RevisionV2: {Spec: SpecV2, Version: V2, Revision: RevisionV2},
	RevisionV3: {Spec: SpecV3, Version: V3, Revision: RevisionV3},
}
