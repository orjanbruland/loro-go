package loro

// Constructors for ExportMode variants, for use with LoroDoc.Export. The
// underlying ExportMode* structs remain available; these helpers shave off
// the struct-literal boilerplate and make the variants discoverable via
// package-level autocomplete.

// SnapshotMode returns an ExportMode that exports a full snapshot of the
// document, including its complete history.
//
//	bytes, err := doc.Export(loro.SnapshotMode())
func SnapshotMode() ExportMode {
	return ExportModeSnapshot{}
}

// UpdatesMode returns an ExportMode that exports updates not present in the
// given version vector. Pass the remote peer's VersionVector to produce a
// delta suitable for incremental sync.
//
//	bytes, err := doc.Export(loro.UpdatesMode(remoteVV))
func UpdatesMode(from *VersionVector) ExportMode {
	return ExportModeUpdates{From: from}
}

// UpdatesInRangeMode returns an ExportMode that exports the updates within
// the given id spans.
func UpdatesInRangeMode(spans []IdSpan) ExportMode {
	return ExportModeUpdatesInRange{Spans: spans}
}

// ShallowSnapshotMode returns an ExportMode that exports a snapshot trimmed
// at the given frontiers, discarding earlier history.
func ShallowSnapshotMode(frontiers *Frontiers) ExportMode {
	return ExportModeShallowSnapshot{Frontiers: frontiers}
}

// StateOnlyMode returns an ExportMode that exports only the state at the
// given frontiers, without history. Pass nil to export the current state.
func StateOnlyMode(frontiers *Frontiers) ExportMode {
	if frontiers == nil {
		return ExportModeStateOnly{Frontiers: nil}
	}
	return ExportModeStateOnly{Frontiers: &frontiers}
}

// SnapshotAtMode returns an ExportMode that exports a snapshot of the
// document as it existed at the given frontiers.
func SnapshotAtMode(frontiers *Frontiers) ExportMode {
	return ExportModeSnapshotAt{Frontiers: frontiers}
}
