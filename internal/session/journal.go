package session

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrNoChanges                = errors.New("no session changes recorded")
	ErrLatestChangeIrreversible = errors.New("latest session change is not reversible")
)

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

func (j *Journal) UndoCandidate() (Change, error) {
	if len(j.changes) == 0 {
		return Change{}, ErrNoChanges
	}

	latest := j.changes[len(j.changes)-1]
	if !latest.Reversible {
		return Change{}, ErrLatestChangeIrreversible
	}
	return latest, nil
}

func (j *Journal) CommitUndo() error {
	if _, err := j.UndoCandidate(); err != nil {
		return err
	}
	j.changes = j.changes[:len(j.changes)-1]
	return nil
}

func (j *Journal) Diff() string {
	var diff strings.Builder
	for _, change := range j.changes {
		if change.Before == change.After {
			continue
		}

		if diff.Len() > 0 {
			diff.WriteByte('\n')
		}
		resource := safeDiffResourceLabel(change.Resource)
		diff.WriteString("--- " + resource + " (before)\n")
		diff.WriteString("+++ " + resource + " (after)\n")
		diff.WriteString("@@\n")
		writeDiffLines(&diff, "-", change.Before)
		writeDiffLines(&diff, "+", change.After)
	}
	return diff.String()
}

func safeDiffResourceLabel(resource string) string {
	for _, r := range resource {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || r == '\u2028' || r == '\u2029' {
			quoted := strconv.Quote(resource)
			return quoted[1 : len(quoted)-1]
		}
	}
	return resource
}

func writeDiffLines(dst *strings.Builder, prefix, content string) {
	if content == "" {
		return
	}
	for _, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {
		dst.WriteString(prefix)
		dst.WriteString(line)
		dst.WriteByte('\n')
	}
}
