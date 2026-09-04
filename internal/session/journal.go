package session

type Change struct {
	Tool       string
	Resource   string
	Before     string
	After      string
	Reversible bool
}

type Journal struct {
	changes []Change
}

func NewJournal() *Journal {
	return &Journal{}
}

func (j *Journal) Record(change Change) {
	j.changes = append(j.changes, change)
}

func (j *Journal) Changes() []Change {
	return append([]Change(nil), j.changes...)
}
