package loro

import "fmt"

// Convenience wrappers around Insert / Push / Set / Mark / SetLocalState that
// accept plain `any` values instead of LoroValueLike. Each method runs the
// argument through AsValue, so the supported types match AsValue's contract
// (string, bool, float64, int64, nil, []any, map[string]any).

// InsertAny inserts a plain Go value at key. The value is converted via
// AsValue; passing an unsupported type returns an error without modifying
// the map.
func (m *LoroMap) InsertAny(key string, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return m.Insert(key, val)
}

// InsertAny inserts a plain Go value at pos. The value is converted via
// AsValue; passing an unsupported type returns an error without modifying
// the list.
func (l *LoroList) InsertAny(pos uint32, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return l.Insert(pos, val)
}

// PushAny appends a plain Go value. The value is converted via AsValue;
// passing an unsupported type returns an error without modifying the list.
func (l *LoroList) PushAny(v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return l.Push(val)
}

// InsertAny inserts a plain Go value at pos. The value is converted via
// AsValue; passing an unsupported type returns an error without modifying
// the list.
func (ml *LoroMovableList) InsertAny(pos uint32, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return ml.Insert(pos, val)
}

// PushAny appends a plain Go value. The value is converted via AsValue;
// passing an unsupported type returns an error without modifying the list.
func (ml *LoroMovableList) PushAny(v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return ml.Push(val)
}

// SetAny replaces the value at pos with a plain Go value. The value is
// converted via AsValue; passing an unsupported type returns an error
// without modifying the list.
func (ml *LoroMovableList) SetAny(pos uint32, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return ml.Set(pos, val)
}

// SetAny sets the ephemeral key to a plain Go value. The value is converted
// via AsValue; passing an unsupported type returns an error without
// modifying the store.
func (s *EphemeralStore) SetAny(key string, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	s.Set(key, val)
	return nil
}

// SetLocalStateAny sets the local awareness state to a plain Go value. The
// value is converted via AsValue; passing an unsupported type returns an
// error without modifying the awareness state.
func (a *Awareness) SetLocalStateAny(v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	a.SetLocalState(val)
	return nil
}

// MarkAny applies a mark with a plain Go value over [from, to). The value is
// converted via AsValue; passing an unsupported type returns an error without
// applying the mark.
func (t *LoroText) MarkAny(from uint32, to uint32, key string, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return t.Mark(from, to, key, val)
}

// MarkUtf16Any applies a mark with a plain Go value over a UTF-16 range. The
// value is converted via AsValue; passing an unsupported type returns an
// error without applying the mark.
func (t *LoroText) MarkUtf16Any(from uint32, to uint32, key string, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return t.MarkUtf16(from, to, key, val)
}

// MarkUtf8Any applies a mark with a plain Go value over a UTF-8 range. The
// value is converted via AsValue; passing an unsupported type returns an
// error without applying the mark.
func (t *LoroText) MarkUtf8Any(from uint32, to uint32, key string, v any) error {
	val, err := AsValue(v)
	if err != nil {
		return err
	}
	return t.MarkUtf8(from, to, key, val)
}

// Append appends s to the end of the text. It is an idiomatic alias for
// PushStr.
func (t *LoroText) Append(s string) error {
	return t.PushStr(s)
}

// Appendf formats according to a format specifier and appends the result to
// the end of the text.
func (t *LoroText) Appendf(format string, args ...any) error {
	return t.PushStr(fmt.Sprintf(format, args...))
}

// Clear removes all content from the text.
func (t *LoroText) Clear() error {
	n := t.LenUnicode()
	if n == 0 {
		return nil
	}
	return t.Delete(0, n)
}

// Write appends p to the end of the text, implementing io.Writer. The bytes
// are interpreted as UTF-8; passing invalid UTF-8 will return an error from
// the underlying append.
func (t *LoroText) Write(p []byte) (int, error) {
	if err := t.PushStr(string(p)); err != nil {
		return 0, err
	}
	return len(p), nil
}
