package loro

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// addSegment appends a transcript-like segment: a map with a text and a
// words list, every value carrying marker.
func addSegment(t *testing.T, segments *LoroList, marker string) *LoroMap {
	t.Helper()
	seg, err := segments.InsertLoroMap(segments.Len())
	must(t, err)
	must(t, seg.InsertAny("speaker", marker+"-speaker"))
	text, err := seg.GetOrCreateLoroText("text")
	must(t, err)
	must(t, text.Insert(0, marker+"-text"))
	words, err := seg.GetOrCreateLoroList("words")
	must(t, err)
	must(t, words.PushAny(marker+"-word"))
	must(t, words.PushAny(2.9899999999999998))
	return seg
}

func mustSnapshot(t *testing.T, doc *LoroDoc) []byte {
	t.Helper()
	b, err := doc.Export(SnapshotMode())
	must(t, err)
	return b
}

// assertAbsent checks the marker is in neither the snapshot bytes nor the
// JSON history.
func assertAbsent(t *testing.T, doc *LoroDoc, marker string) {
	t.Helper()
	if bytes.Contains(mustSnapshot(t, doc), []byte(marker)) {
		t.Errorf("snapshot still contains %q", marker)
	}
	vv := doc.OplogVv()
	defer vv.Destroy()
	empty := NewVersionVector()
	defer empty.Destroy()
	if bytes.Contains([]byte(doc.ExportJsonUpdatesWithoutPeerCompression(empty, vv)), []byte(marker)) {
		t.Errorf("JSON history still contains %q", marker)
	}
}

func assertPresent(t *testing.T, doc *LoroDoc, marker string) {
	t.Helper()
	if !bytes.Contains(mustSnapshot(t, doc), []byte(marker)) {
		t.Errorf("snapshot lost %q", marker)
	}
}

func mustRedact(t *testing.T, doc *LoroDoc, opts RedactOptions) RedactResult {
	t.Helper()
	res, err := doc.RedactDeletedContainers(opts)
	must(t, err)
	if res.Doc != nil {
		t.Cleanup(res.Doc.Destroy)
	}
	return res
}

func TestRedactDeletedContainers_SeveralSpansOnePeer(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	must(t, doc.SetPeerId(7))
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE1")
	addSegment(t, segments, "KEPT")
	addSegment(t, segments, "GONE2")
	// One peer wrote all three, so the deleted ops form two spans with live
	// ops between them.
	must(t, segments.Delete(2, 1))
	must(t, segments.Delete(0, 1))

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	// Per segment: speaker insert, text insert, two words inserts, plus the
	// map ops that create the text and the list.
	if res.RedactedContainers != 6 {
		t.Errorf("RedactedContainers = %d, want 6", res.RedactedContainers)
	}
	assertAbsent(t, res.Doc, "GONE1")
	assertAbsent(t, res.Doc, "GONE2")
	assertPresent(t, res.Doc, "KEPT")
	if got := res.Doc.PeerId(); got != 7 {
		t.Errorf("PeerId = %d, want 7", got)
	}

	seg, ok := res.Doc.GetList(AsContainerId("segments")).GetLoroMap(0)
	if !ok {
		t.Fatal("segment 0 is not a map")
	}
	words, ok := seg.GetLoroList("words")
	if !ok {
		t.Fatal("kept segment lost its words")
	}
	if f, ok := words.GetFloat64(1); !ok || f != 2.9899999999999998 {
		t.Errorf("kept float = (%v, %v), want exact 2.9899999999999998", f, ok)
	}
	assertPresent(t, doc, "GONE1")
}

func TestRedactDeletedContainers_OverwrittenList(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	seg, err := doc.GetList(AsContainerId("segments")).InsertLoroMap(0)
	must(t, err)
	old, err := seg.GetOrCreateLoroList("words")
	must(t, err)
	must(t, old.PushAny("OLDWORD"))
	// Replacing a map value deletes the container it held.
	fresh, err := seg.InsertListContainer("words", NewLoroList())
	must(t, err)
	must(t, fresh.PushAny("NEWWORD"))

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	assertAbsent(t, res.Doc, "OLDWORD")
	assertPresent(t, res.Doc, "NEWWORD")
}

func TestRedactDeletedContainers_NothingToRedact(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	addSegment(t, doc.GetList(AsContainerId("segments")), "LIVE")

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc != nil || res.RedactedOps != 0 {
		t.Fatalf("got %+v, want no rebuild", res)
	}
}

func TestRedactDeletedContainers_AlreadyRedacted(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	must(t, segments.Delete(0, 1))

	first := mustRedact(t, doc, RedactOptions{})
	if first.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	second := mustRedact(t, first.Doc, RedactOptions{})
	if second.Doc != nil || second.RedactedOps != 0 {
		t.Fatalf("second pass = %+v, want no rebuild", second)
	}
}

func TestRedactDeletedContainers_SyncsWithOriginalHistory(t *testing.T) {
	original := NewLoroDoc()
	defer original.Destroy()
	segments := original.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	addSegment(t, segments, "KEPT")

	// stale has the segment but not its delete.
	stale := NewLoroDoc()
	defer stale.Destroy()
	_, err := stale.Import(mustSnapshot(t, original))
	must(t, err)
	must(t, segments.Delete(0, 1))

	res := mustRedact(t, original, RedactOptions{})
	rebuilt := res.Doc
	if rebuilt == nil {
		t.Fatal("expected a rebuilt doc")
	}

	// An import made only of op IDs the rebuilt doc holds is a no-op, so the
	// content does not come back.
	_, err = rebuilt.Import(mustSnapshot(t, original))
	must(t, err)
	assertAbsent(t, rebuilt, "GONE")

	// A peer that never saw the delete writes into the deleted segment; the
	// next pass redacts it.
	staleSeg, ok := stale.GetList(AsContainerId("segments")).GetLoroMap(0)
	if !ok {
		t.Fatal("stale segment missing")
	}
	must(t, staleSeg.InsertAny("speaker", "LATE"))
	rebuiltVV := rebuilt.OplogVv()
	defer rebuiltVV.Destroy()
	late, err := stale.Export(UpdatesMode(rebuiltVV))
	must(t, err)
	_, err = rebuilt.Import(late)
	must(t, err)
	again := mustRedact(t, rebuilt, RedactOptions{})
	if again.Doc == nil || again.RedactedOps != 1 {
		t.Fatalf("second pass = %+v, want one redacted op", again)
	}
	assertAbsent(t, again.Doc, "LATE")

	// New writes on the rebuilt doc reach a peer holding the original
	// history.
	must(t, again.Doc.GetList(AsContainerId("segments")).PushAny("AFTER"))
	originalVV := original.OplogVv()
	defer originalVV.Destroy()
	update, err := again.Doc.Export(UpdatesMode(originalVV))
	must(t, err)
	_, err = original.Import(update)
	must(t, err)
	if s, ok := original.GetList(AsContainerId("segments")).GetString(1); !ok || s != "AFTER" {
		t.Errorf("original did not receive the new write: (%q, %v)", s, ok)
	}
}

// Reading a root creates it with no op, and the rebuilt doc has no such root.
func TestRedactDeletedContainers_RootCreatedByAReader(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	addSegment(t, segments, "KEPT")
	must(t, segments.Delete(0, 1))
	doc.GetMap(AsContainerId("content")).Destroy()
	doc.GetTree(AsContainerId("sections")).Destroy()

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	assertAbsent(t, res.Doc, "GONE")
	assertPresent(t, res.Doc, "KEPT")
}

// A root emptied by its own ops is on the rebuilt doc, and a reader's root is
// only on the source.
func TestRedactDeletedContainers_RootEmptiedByItsOps(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	sections := doc.GetTree(AsContainerId("sections"))
	node, err := sections.Create(TreeParentIdRoot{})
	must(t, err)
	meta, err := sections.GetMeta(node)
	must(t, err)
	must(t, meta.InsertAny("title", "GONE"))
	must(t, sections.Delete(node))
	doc.GetMap(AsContainerId("content")).Destroy()

	res := mustRedact(t, doc, RedactOptions{IncludeRevivable: true})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	assertAbsent(t, res.Doc, "GONE")
}

func TestRedactDeletedContainers_OverwrittenMapValue(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	seg, err := segments.InsertLoroMap(0)
	must(t, err)
	must(t, seg.InsertAny("text", "OLDTEXT"))
	must(t, seg.InsertAny("speaker", "SPEAKER"))
	doc.Commit()
	must(t, seg.InsertAny("text", "NEWTEXT"))

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil || res.RedactedOps != 1 || res.RedactedContainers != 0 {
		t.Fatalf("got %+v, want one redacted op in no deleted container", res)
	}
	assertAbsent(t, res.Doc, "OLDTEXT")
	rebuiltSeg, ok := res.Doc.GetList(AsContainerId("segments")).GetLoroMap(0)
	if !ok {
		t.Fatal("rebuilt segment missing")
	}
	if s, ok := rebuiltSeg.GetString("text"); !ok || s != "NEWTEXT" {
		t.Errorf("text = (%q, %v), want NEWTEXT", s, ok)
	}
	if s, ok := rebuiltSeg.GetString("speaker"); !ok || s != "SPEAKER" {
		t.Errorf("speaker = (%q, %v), want SPEAKER", s, ok)
	}

	second := mustRedact(t, res.Doc, RedactOptions{})
	if second.Doc != nil || second.RedactedOps != 0 {
		t.Fatalf("second pass = %+v, want no rebuild", second)
	}
}

func TestRedactDeletedContainers_DeletedMapKey(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	m := doc.GetMap(AsContainerId("meta"))
	must(t, m.InsertAny("summary", "OLDSUMMARY"))
	must(t, m.InsertAny("kept", "KEPT"))
	doc.Commit()
	must(t, m.Delete("summary"))

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	assertAbsent(t, res.Doc, "OLDSUMMARY")
	assertPresent(t, res.Doc, "KEPT")
}

// An op's lamport is its change's plus its offset in the change, and ops of
// other containers in the change move that offset.
func TestRedactDeletedContainers_OverwriteInsideOneChange(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	words := doc.GetList(AsContainerId("words"))
	must(t, words.PushAny("one"))
	must(t, words.PushAny("two"))
	m := doc.GetMap(AsContainerId("meta"))
	must(t, m.InsertAny("text", "OLDTEXT"))
	must(t, words.PushAny("three"))
	must(t, m.InsertAny("text", "NEWTEXT"))

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil || res.RedactedOps != 1 {
		t.Fatalf("got %+v, want one redacted op", res)
	}
	assertAbsent(t, res.Doc, "OLDTEXT")
	if s, ok := res.Doc.GetMap(AsContainerId("meta")).GetString("text"); !ok || s != "NEWTEXT" {
		t.Errorf("text = (%q, %v), want NEWTEXT", s, ok)
	}
}

// Two writes at one lamport: the higher peer holds the key on every copy, and
// the other value is the one to go, whichever way the peers are numbered.
func TestRedactDeletedContainers_ConcurrentWritesTieOnPeer(t *testing.T) {
	for _, tc := range []struct {
		name          string
		peerA, peerB  uint64
		winner, loser string
	}{
		{"second peer higher", 1, 2, "FROMB", "FROMA"},
		{"first peer higher", 2, 1, "FROMA", "FROMB"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := NewLoroDoc()
			defer a.Destroy()
			must(t, a.SetPeerId(tc.peerA))
			b := NewLoroDoc()
			defer b.Destroy()
			must(t, b.SetPeerId(tc.peerB))
			must(t, a.GetMap(AsContainerId("meta")).InsertAny("text", "FROMA"))
			must(t, b.GetMap(AsContainerId("meta")).InsertAny("text", "FROMB"))
			_, err := a.Import(mustSnapshot(t, b))
			must(t, err)

			res := mustRedact(t, a, RedactOptions{})
			if res.Doc == nil || res.RedactedOps != 1 {
				t.Fatalf("got %+v, want one redacted op", res)
			}
			assertAbsent(t, res.Doc, tc.loser)
			if s, ok := res.Doc.GetMap(AsContainerId("meta")).GetString("text"); !ok || s != tc.winner {
				t.Errorf("text = (%q, %v), want %q", s, ok, tc.winner)
			}
		})
	}
}

// A delta from a peer that has not seen the winner arrives after the pass. It
// loses on the redacted copy as it does on the original, and the next pass
// removes it.
func TestRedactDeletedContainers_LateConcurrentWriteLosesToo(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	must(t, doc.SetPeerId(5))
	m := doc.GetMap(AsContainerId("meta"))
	must(t, m.InsertAny("text", "FIRST"))
	doc.Commit()

	stale := NewLoroDoc()
	defer stale.Destroy()
	must(t, stale.SetPeerId(1))
	_, err := stale.Import(mustSnapshot(t, doc))
	must(t, err)

	must(t, m.InsertAny("text", "SECOND"))
	doc.Commit()
	must(t, m.InsertAny("text", "THIRD"))
	doc.Commit()
	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}

	must(t, stale.GetMap(AsContainerId("meta")).InsertAny("text", "LATE"))
	for _, target := range []*LoroDoc{res.Doc, doc} {
		vv := target.OplogVv()
		late, err := stale.Export(UpdatesMode(vv))
		vv.Destroy()
		must(t, err)
		_, err = target.Import(late)
		must(t, err)
		if s, ok := target.GetMap(AsContainerId("meta")).GetString("text"); !ok || s != "THIRD" {
			t.Errorf("text = (%q, %v), want THIRD", s, ok)
		}
	}

	again := mustRedact(t, res.Doc, RedactOptions{})
	if again.Doc == nil || again.RedactedOps != 1 {
		t.Fatalf("second pass = %+v, want one redacted op", again)
	}
	assertAbsent(t, again.Doc, "LATE")
	assertPresent(t, again.Doc, "THIRD")
}

// Pins Loro's guard against peer ID reuse: an import that brings a new op
// is checked against every op ID it shares with the doc, and a redacted op
// never matches its original. Peers must exchange deltas from the version
// vector, not whole histories, with a redacted doc.
func TestRedactDeletedContainers_OverlappingImportWithNewOpsIsRefused(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	stale := NewLoroDoc()
	defer stale.Destroy()
	_, err := stale.Import(mustSnapshot(t, doc))
	must(t, err)
	must(t, segments.Delete(0, 1))
	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	must(t, stale.GetList(AsContainerId("other")).PushAny("NEW"))

	if _, err := res.Doc.Import(mustSnapshot(t, stale)); err == nil || !strings.Contains(err.Error(), "UsedOpId") {
		t.Fatalf("snapshot import err = %v, want UsedOpId", err)
	}
	vv := res.Doc.OplogVv()
	defer vv.Destroy()
	delta, err := stale.Export(UpdatesMode(vv))
	must(t, err)
	_, err = res.Doc.Import(delta)
	must(t, err)
}

func TestRedactDeletedContainers_TreeMetaIsRevivable(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	tree := doc.GetTree(AsContainerId("sections"))
	node, err := tree.Create(TreeParentIdRoot{})
	must(t, err)
	meta, err := tree.GetMeta(node)
	must(t, err)
	must(t, meta.InsertAny("title", "SECTITLE"))
	must(t, tree.Delete(node))

	skipped := mustRedact(t, doc, RedactOptions{})
	if skipped.Doc != nil || skipped.SkippedRevivableContainers != 1 {
		t.Fatalf("default options = %+v, want the meta map skipped", skipped)
	}

	res := mustRedact(t, doc, RedactOptions{IncludeRevivable: true})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	assertAbsent(t, res.Doc, "SECTITLE")
}

// Pins the caveat on IncludeRevivable: a concurrent move from a peer that
// missed the delete brings the node back with redacted metadata.
func TestRedactDeletedContainers_RevivedTreeNodeShowsNull(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	tree := doc.GetTree(AsContainerId("sections"))
	parent, err := tree.Create(TreeParentIdRoot{})
	must(t, err)
	node, err := tree.Create(TreeParentIdRoot{})
	must(t, err)
	meta, err := tree.GetMeta(node)
	must(t, err)
	must(t, meta.InsertAny("title", "SECTITLE"))

	stale := NewLoroDoc()
	defer stale.Destroy()
	_, err = stale.Import(mustSnapshot(t, doc))
	must(t, err)
	must(t, tree.Delete(node))

	res := mustRedact(t, doc, RedactOptions{IncludeRevivable: true})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	// The move is authored after the delete was imported elsewhere, so it
	// wins on lamport once both sides see it.
	staleTree := stale.GetTree(AsContainerId("sections"))
	must(t, staleTree.Mov(node, TreeParentNode(parent)))
	must(t, staleTree.Mov(node, TreeParentIdRoot{}))
	for _, target := range []*LoroDoc{res.Doc, doc} {
		vv := target.OplogVv()
		move, err := stale.Export(UpdatesMode(vv))
		vv.Destroy()
		must(t, err)
		_, err = target.Import(move)
		must(t, err)
	}

	revived, err := res.Doc.GetTree(AsContainerId("sections")).GetMeta(node)
	must(t, err)
	if !revived.IsExplicitlyNil("title") {
		v, _ := revived.GetAny("title")
		t.Errorf("revived title on the rebuilt doc = %v, want null", v)
	}
	original, err := tree.GetMeta(node)
	must(t, err)
	if s, ok := original.GetString("title"); !ok || s != "SECTITLE" {
		t.Errorf("revived title on the original = (%q, %v), want SECTITLE", s, ok)
	}
}

func TestRedactDeletedContainers_KeepsPendingOps(t *testing.T) {
	base := NewLoroDoc()
	defer base.Destroy()
	must(t, base.GetList(AsContainerId("pending")).PushAny("base"))

	writer := NewLoroDoc()
	defer writer.Destroy()
	_, err := writer.Import(mustSnapshot(t, base))
	must(t, err)
	baseVV := writer.OplogVv()
	defer baseVV.Destroy()
	must(t, writer.GetList(AsContainerId("pending")).PushAny("dependent"))
	dependent, err := writer.Export(UpdatesMode(baseVV))
	must(t, err)

	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	must(t, segments.Delete(0, 1))
	status, err := doc.Import(dependent)
	must(t, err)
	if status.Pending == nil {
		t.Fatal("expected the dependent update to be pending")
	}

	res := mustRedact(t, doc, RedactOptions{})
	if res.Doc == nil {
		t.Fatal("expected a rebuilt doc")
	}
	_, err = res.Doc.Import(mustSnapshot(t, base))
	must(t, err)
	if s, ok := res.Doc.GetList(AsContainerId("pending")).GetString(1); !ok || s != "dependent" {
		t.Errorf("pending op lost in the rebuild: (%q, %v)", s, ok)
	}
}

func TestRedactDeletedContainers_Shallow(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	must(t, segments.Delete(0, 1))
	frontiers := doc.OplogFrontiers()
	defer frontiers.Destroy()
	shallowBytes, err := doc.Export(ShallowSnapshotMode(frontiers))
	must(t, err)
	shallow := NewLoroDoc()
	defer shallow.Destroy()
	_, err = shallow.Import(shallowBytes)
	must(t, err)

	if _, err := shallow.RedactDeletedContainers(RedactOptions{}); !errors.Is(err, ErrRedactShallow) {
		t.Fatalf("err = %v, want ErrRedactShallow", err)
	}
}

func TestRedactDeletedContainers_Detached(t *testing.T) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	frontiers := doc.OplogFrontiers()
	defer frontiers.Destroy()
	must(t, segments.Delete(0, 1))
	must(t, doc.Checkout(frontiers))

	if _, err := doc.RedactDeletedContainers(RedactOptions{}); !errors.Is(err, ErrRedactDetached) {
		t.Fatalf("err = %v, want ErrRedactDetached", err)
	}
}

// Pins the handle contract the Rust side relies on: it consumes the handle
// incrementPointer clones, and returns one Lift takes ownership of. A wrong
// refcount shows up as a crash or a use-after-free under -race here.
func TestRedactDeletedContainers_HandleOwnership(t *testing.T) {
	doc := NewLoroDoc()
	segments := doc.GetList(AsContainerId("segments"))
	addSegment(t, segments, "GONE")
	must(t, segments.Delete(0, 1))

	var rebuilt []*LoroDoc
	for range 50 {
		res, err := doc.RedactDeletedContainers(RedactOptions{})
		must(t, err)
		rebuilt = append(rebuilt, res.Doc)
	}
	// The source is still usable after being passed fifty times.
	must(t, segments.PushAny("still-alive"))
	for _, r := range rebuilt {
		if r.GetList(AsContainerId("segments")).Len() != 0 {
			t.Fatal("rebuilt doc changed with the source")
		}
		r.Destroy()
	}
	doc.Destroy()
}

// About an hour of transcript, every segment then deleted, as a wipe does.
func BenchmarkRedactDeletedContainers_Wipe(b *testing.B) {
	doc := NewLoroDoc()
	defer doc.Destroy()
	segments := doc.GetList(AsContainerId("segments"))
	for i := range 1500 {
		seg, _ := segments.InsertLoroMap(uint32(i))
		_ = seg.InsertAny("speaker", "speaker")
		text, _ := seg.GetOrCreateLoroText("text")
		_ = text.Insert(0, "the quick brown fox jumps over the lazy dog and keeps on talking for a while")
		words, _ := seg.GetOrCreateLoroList("words")
		for w := range 15 {
			_ = words.PushAny(map[string]any{"w": "word", "start": float64(i) + float64(w)*0.1, "end": float64(i) + float64(w)*0.1 + 0.07})
		}
	}
	_ = segments.Delete(0, segments.Len())
	b.ResetTimer()
	for b.Loop() {
		res, err := doc.RedactDeletedContainers(RedactOptions{})
		if err != nil {
			b.Fatal(err)
		}
		res.Doc.Destroy()
	}
}

// The JSON path parses floats with serde_json, which returns the same double
// only with its float_roundtrip feature (enabled by loro-internal).
func TestJsonUpdates_FloatRoundTrip(t *testing.T) {
	values := []float64{2.9899999999999998, 0.44999999999999996}
	doc := NewLoroDoc()
	defer doc.Destroy()
	list := doc.GetList(AsContainerId("floats"))
	for _, v := range values {
		must(t, list.PushAny(v))
	}

	empty := NewVersionVector()
	defer empty.Destroy()
	vv := doc.OplogVv()
	defer vv.Destroy()
	vr := NewVersionRange()
	defer vr.Destroy()
	redacted, err := doc.RedactJsonUpdates(doc.ExportJsonUpdatesWithoutPeerCompression(empty, vv), vr)
	must(t, err)
	rebuilt := NewLoroDoc()
	defer rebuilt.Destroy()
	_, err = rebuilt.ImportJsonUpdates(redacted)
	must(t, err)

	got := rebuilt.GetList(AsContainerId("floats"))
	for i, want := range values {
		if f, ok := got.GetFloat64(uint32(i)); !ok || f != want {
			t.Errorf("float %d = (%v, %v), want exact %v", i, f, ok, want)
		}
	}
}
