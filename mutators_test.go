package loro

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestLoroMap_InsertAny(t *testing.T) {
	doc := NewLoroDoc()
	m := doc.GetMap(AsContainerId("m"))

	must(t, m.InsertAny("s", "hello"))
	must(t, m.InsertAny("i", 7))
	must(t, m.InsertAny("nested", map[string]any{"k": "v"}))

	if v, ok := m.GetString("s"); !ok || v != "hello" {
		t.Fatalf("GetString = (%q, %v)", v, ok)
	}
	if v, ok := m.GetInt64("i"); !ok || v != 7 {
		t.Fatalf("GetInt64 = (%v, %v)", v, ok)
	}
	if mp, ok := m.GetMapOfAny("nested"); !ok || !reflect.DeepEqual(mp, map[string]any{"k": "v"}) {
		t.Fatalf("nested = (%v, %v)", mp, ok)
	}

	if err := m.InsertAny("bad", struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
	if _, ok := m.GetAny("bad"); ok {
		t.Fatal("map should not have been modified on AsValue error")
	}
}

func TestLoroList_InsertAndPushAny(t *testing.T) {
	doc := NewLoroDoc()
	l := doc.GetList(AsContainerId("l"))

	must(t, l.PushAny("x"))
	must(t, l.PushAny(int64(1)))
	must(t, l.InsertAny(0, true))

	if v, ok := l.GetBool(0); !ok || !v {
		t.Fatalf("[0] = (%v, %v)", v, ok)
	}
	if v, ok := l.GetString(1); !ok || v != "x" {
		t.Fatalf("[1] = (%q, %v)", v, ok)
	}
	if v, ok := l.GetInt64(2); !ok || v != 1 {
		t.Fatalf("[2] = (%v, %v)", v, ok)
	}

	if err := l.PushAny(struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
	if err := l.InsertAny(0, struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestLoroMovableList_SetAny(t *testing.T) {
	doc := NewLoroDoc()
	ml := doc.GetMovableList(AsContainerId("ml"))
	must(t, ml.PushAny("a"))
	must(t, ml.PushAny("b"))
	must(t, ml.SetAny(1, int64(42)))

	if v, ok := ml.GetInt64(1); !ok || v != 42 {
		t.Fatalf("[1] after SetAny = (%v, %v)", v, ok)
	}
	if err := ml.SetAny(0, struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestLoroText_MarkAny(t *testing.T) {
	doc := NewLoroDoc()
	text := doc.GetText(AsContainerId("t"))
	must(t, text.Insert(0, "hello world"))

	must(t, text.MarkAny(0, 5, "bold", true))
	if err := text.MarkAny(0, 5, "bad", struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestLoroText_Append(t *testing.T) {
	doc := NewLoroDoc()
	text := doc.GetText(AsContainerId("t"))

	must(t, text.Append("hello"))
	must(t, text.Append(", "))
	must(t, text.Append("world"))

	if got := text.String(); got != "hello, world" {
		t.Fatalf("text = %q, want %q", got, "hello, world")
	}
}

func TestLoroText_Appendf(t *testing.T) {
	doc := NewLoroDoc()
	text := doc.GetText(AsContainerId("t"))

	must(t, text.Appendf("count=%d ", 3))
	must(t, text.Appendf("name=%s", "alice"))

	if got := text.String(); got != "count=3 name=alice" {
		t.Fatalf("text = %q", got)
	}
}

func TestLoroText_Clear(t *testing.T) {
	doc := NewLoroDoc()
	text := doc.GetText(AsContainerId("t"))

	// Clearing an empty text should not error.
	must(t, text.Clear())

	must(t, text.Insert(0, "hello, world"))
	must(t, text.Clear())

	if got := text.String(); got != "" {
		t.Fatalf("text after Clear = %q, want empty", got)
	}
	if !text.IsEmpty() {
		t.Fatal("IsEmpty = false after Clear")
	}
}

func TestLoroText_Write(t *testing.T) {
	doc := NewLoroDoc()
	text := doc.GetText(AsContainerId("t"))

	// Verify it satisfies io.Writer.
	var w io.Writer = text

	n, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatalf("Write returned n=%d, want 5", n)
	}

	if _, err := fmt.Fprintf(w, " %s #%d", "world", 7); err != nil {
		t.Fatal(err)
	}

	if got := text.String(); got != "hello world #7" {
		t.Fatalf("text = %q", got)
	}
}

func TestEphemeralStore_SetAny(t *testing.T) {
	store := NewEphemeralStore(60_000)
	if err := store.SetAny("k", "v"); err != nil {
		t.Fatal(err)
	}
	if err := store.SetAny("bad", struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestAwareness_SetLocalStateAny(t *testing.T) {
	a := NewAwareness(1, 60_000)
	if err := a.SetLocalStateAny("hello"); err != nil {
		t.Fatal(err)
	}
	if err := a.SetLocalStateAny(struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
}
