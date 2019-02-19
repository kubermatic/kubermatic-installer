package dns

type RecordKind int

const (
	RecordKindA RecordKind = iota
	RecordKindCNAME
)

type Record struct {
	Name   string
	Target string
	Kind   RecordKind
}
