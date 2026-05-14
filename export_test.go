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
