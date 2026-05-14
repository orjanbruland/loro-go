package loro

import (
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
