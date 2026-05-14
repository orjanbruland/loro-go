package loro

import (
	"testing"
)

// exportImport round-trips a doc through the given ExportMode and applies it
// to a fresh doc, returning the destination doc for inspection.
func exportImport(t *testing.T, src *LoroDoc, mode ExportMode) *LoroDoc {
	t.Helper()
	bytes, err := src.Export(mode)
	if err != nil {
		t.Fatalf("Export error: %v", err)
	}
	if len(bytes) == 0 {
		t.Fatal("Export returned empty bytes")
	}
	dst := NewLoroDoc()
	if _, err := dst.Import(bytes); err != nil {
		t.Fatalf("Import error: %v", err)
	}
	return dst
}

func TestExportMode_Snapshot(t *testing.T) {
	src := NewLoroDoc()
	m := src.GetMap(AsContainerId("m"))
	must(t, m.InsertAny("k", "v"))

	dst := exportImport(t, src, SnapshotMode())
	got, ok := dst.GetMap(AsContainerId("m")).GetString("k")
	if !ok || got != "v" {
		t.Fatalf("snapshot import lost data: (%q, %v)", got, ok)
	}
}

func TestExportMode_Updates(t *testing.T) {
	src := NewLoroDoc()
	m := src.GetMap(AsContainerId("m"))
	must(t, m.InsertAny("k", int64(1)))

	dst := exportImport(t, src, UpdatesMode(NewVersionVector()))
	got, ok := dst.GetMap(AsContainerId("m")).GetInt64("k")
	if !ok || got != 1 {
		t.Fatalf("updates import lost data: (%v, %v)", got, ok)
	}
}

func TestExportMode_Updates_DeltaAgainstPeer(t *testing.T) {
	a := NewLoroDoc()
	am := a.GetMap(AsContainerId("m"))
	must(t, am.InsertAny("a", int64(1)))
	a.Commit()

	b := NewLoroDoc()
	full, err := a.Export(SnapshotMode())
	if err != nil {
		t.Fatalf("snapshot export error: %v", err)
	}
	if _, err := b.Import(full); err != nil {
		t.Fatalf("initial import error: %v", err)
	}

	must(t, am.InsertAny("b", int64(2)))
	a.Commit()

	delta, err := a.Export(UpdatesMode(b.StateVv()))
	if err != nil {
		t.Fatalf("delta export error: %v", err)
	}
	if len(delta) >= len(full) {
		t.Fatalf("delta (%d bytes) should be smaller than full snapshot (%d bytes)", len(delta), len(full))
	}
	if _, err := b.Import(delta); err != nil {
		t.Fatalf("delta import error: %v", err)
	}

	bm := b.GetMap(AsContainerId("m"))
	if v, ok := bm.GetInt64("a"); !ok || v != 1 {
		t.Fatalf("a after delta sync = (%v, %v), want 1", v, ok)
	}
	if v, ok := bm.GetInt64("b"); !ok || v != 2 {
		t.Fatalf("b after delta sync = (%v, %v), want 2", v, ok)
	}
}

func TestConcurrentEdits_Converge(t *testing.T) {
	a := NewLoroDoc()
	b := NewLoroDoc()

	must(t, a.GetMap(AsContainerId("m")).InsertAny("from-a", int64(1)))
	must(t, b.GetMap(AsContainerId("m")).InsertAny("from-b", int64(2)))

	aUpdates, err := a.Export(UpdatesMode(b.StateVv()))
	if err != nil {
		t.Fatalf("a export error: %v", err)
	}
	bUpdates, err := b.Export(UpdatesMode(a.StateVv()))
	if err != nil {
		t.Fatalf("b export error: %v", err)
	}
	if _, err := b.Import(aUpdates); err != nil {
		t.Fatalf("b import error: %v", err)
	}
	if _, err := a.Import(bUpdates); err != nil {
		t.Fatalf("a import error: %v", err)
	}

	for _, d := range []*LoroDoc{a, b} {
		m := d.GetMap(AsContainerId("m"))
		if v, ok := m.GetInt64("from-a"); !ok || v != 1 {
			t.Fatalf("from-a = (%v, %v), want 1", v, ok)
		}
		if v, ok := m.GetInt64("from-b"); !ok || v != 2 {
			t.Fatalf("from-b = (%v, %v), want 2", v, ok)
		}
	}
}

func TestExportMode_StateOnly_Nil(t *testing.T) {
	src := NewLoroDoc()
	m := src.GetMap(AsContainerId("m"))
	must(t, m.InsertAny("k", "v"))

	dst := exportImport(t, src, StateOnlyMode(nil))
	got, ok := dst.GetMap(AsContainerId("m")).GetString("k")
	if !ok || got != "v" {
		t.Fatalf("state-only import lost data: (%q, %v)", got, ok)
	}
}

func TestExportMode_SnapshotAt(t *testing.T) {
	src := NewLoroDoc()
	m := src.GetMap(AsContainerId("m"))
	must(t, m.InsertAny("a", int64(1)))
	src.Commit()
	frontiers := src.StateFrontiers()
	must(t, m.InsertAny("b", int64(2)))
	src.Commit()

	dst := exportImport(t, src, SnapshotAtMode(frontiers))
	dm := dst.GetMap(AsContainerId("m"))
	if v, ok := dm.GetInt64("a"); !ok || v != 1 {
		t.Fatalf("a at frontiers = (%v, %v), want 1", v, ok)
	}
	if _, ok := dm.GetInt64("b"); ok {
		t.Fatal("b should not be present at earlier frontiers")
	}
}
