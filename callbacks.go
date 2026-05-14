package loro

// LocalUpdateCallbackFn adapts a function to the LocalUpdateCallback interface.
//
// Use with LoroDoc.SubscribeLocalUpdate to receive raw update bytes whenever the
// local document changes. Useful for syncing updates to a remote peer or persisting
// them to storage.
//
//	sub := doc.SubscribeLocalUpdate(loro.LocalUpdateCallbackFn(func(update []byte) {
//		_ = conn.Send(update) // send update to a remote peer
//	}))
//	defer sub.Unsubscribe()
type LocalUpdateCallbackFn func(update []byte)

func (fn LocalUpdateCallbackFn) OnLocalUpdate(update []byte) {
	fn(update)
}

// AsLocalUpdateCallback adapts a function to the LocalUpdateCallback interface.
func AsLocalUpdateCallback(fn LocalUpdateCallbackFn) LocalUpdateCallback {
	return fn
}

// SubscriberFn adapts a function to the Subscriber interface.
//
// Use with LoroDoc.Subscribe, LoroDoc.SubscribeRoot, or any container's Subscribe
// method (LoroMap, LoroList, LoroText, etc.) to observe document or container changes.
//
//	sub := doc.SubscribeRoot(loro.SubscriberFn(func(diff loro.DiffEvent) {
//		fmt.Println("changed by:", diff.TriggeredBy)
//		for _, event := range diff.Events {
//			fmt.Println("container:", event.CurrentTarget)
//		}
//	}))
//	defer sub.Unsubscribe()
type SubscriberFn func(diff DiffEvent)

func (fn SubscriberFn) OnDiff(diff DiffEvent) {
	fn(diff)
}

// AsSubscriber adapts a function to the Subscriber interface.
func AsSubscriber(fn SubscriberFn) Subscriber {
	return fn
}

// FirstCommitFromPeerCallbackFn adapts a function to the FirstCommitFromPeerCallback interface.
//
// Use with LoroDoc.SubscribeFirstCommitFromPeer to detect when a new peer starts
// contributing to the document for the first time.
//
//	sub := doc.SubscribeFirstCommitFromPeer(loro.FirstCommitFromPeerCallbackFn(
//		func(payload loro.FirstCommitFromPeerPayload) {
//			fmt.Println("new peer joined:", payload.PeerId)
//		},
//	))
//	defer sub.Unsubscribe()
type FirstCommitFromPeerCallbackFn func(payload FirstCommitFromPeerPayload)

func (fn FirstCommitFromPeerCallbackFn) OnFirstCommitFromPeer(payload FirstCommitFromPeerPayload) {
	fn(payload)
}

// AsFirstCommitFromPeerCallback adapts a function to the FirstCommitFromPeerCallback interface.
func AsFirstCommitFromPeerCallback(fn FirstCommitFromPeerCallbackFn) FirstCommitFromPeerCallback {
	return fn
}

// PreCommitCallbackFn adapts a function to the PreCommitCallback interface.
//
// Use with LoroDoc.SubscribePreCommit to intercept commits before they are finalized.
// The callback receives a payload with change metadata and a ChangeModifier that lets
// you set a commit message or timestamp.
//
//	sub := doc.SubscribePreCommit(loro.PreCommitCallbackFn(
//		func(payload loro.PreCommitCallbackPayload) {
//			payload.Modifier.SetMessage("auto-tagged commit")
//		},
//	))
//	defer sub.Unsubscribe()
type PreCommitCallbackFn func(payload PreCommitCallbackPayload)

func (fn PreCommitCallbackFn) OnPreCommit(payload PreCommitCallbackPayload) {
	fn(payload)
}

// AsPreCommitCallback adapts a function to the PreCommitCallback interface.
func AsPreCommitCallback(fn PreCommitCallbackFn) PreCommitCallback {
	return fn
}

// ChangeAncestorsTravelerFn adapts a function to the ChangeAncestorsTraveler interface.
//
// Use with LoroDoc.TravelChangeAncestors to walk the causal history of one or more
// changes. Return true to continue traversal, or false to stop early.
//
//	doc.TravelChangeAncestors(ids, loro.ChangeAncestorsTravelerFn(
//		func(change loro.ChangeMeta) bool {
//			fmt.Println("visiting change from peer:", change.Id.Peer)
//			return true // continue traversal
//		},
//	))
type ChangeAncestorsTravelerFn func(payload ChangeMeta) bool

func (fn ChangeAncestorsTravelerFn) Travel(payload ChangeMeta) bool {
	return fn(payload)
}

// AsChangeAncestorCallback adapts a function to the ChangeAncestorsTraveler interface.
func AsChangeAncestorCallback(fn ChangeAncestorsTravelerFn) ChangeAncestorsTraveler {
	return fn
}

// LocalEphemeralListenerFn adapts a function to the LocalEphemeralListener interface.
//
// Use with EphemeralStore.SubscribeLocalUpdate to receive raw ephemeral update bytes
// whenever the local ephemeral store changes. Useful for broadcasting awareness data
// (e.g. cursor positions) to other peers.
//
//	sub := store.SubscribeLocalUpdate(loro.LocalEphemeralListenerFn(func(update []byte) {
//		_ = conn.Send(update) // broadcast ephemeral update
//	}))
//	defer sub.Unsubscribe()
type LocalEphemeralListenerFn func(update []byte)

func (fn LocalEphemeralListenerFn) OnEphemeralUpdate(update []byte) {
	fn(update)
}

// AsLocalEphemeralListener adapts a function to the LocalEphemeralListener interface.
func AsLocalEphemeralListener(fn LocalEphemeralListenerFn) LocalEphemeralListener {
	return fn
}

// EphemeralSubscriberFn adapts a function to the EphemeralSubscriber interface.
//
// Use with EphemeralStore.Subscribe to observe changes to the ephemeral store,
// such as awareness/presence updates from peers.
//
//	sub := store.Subscribe(loro.EphemeralSubscriberFn(func(event loro.EphemeralStoreEvent) {
//		for _, key := range event.Updated {
//			fmt.Println("peer updated:", key)
//		}
//	}))
//	defer sub.Unsubscribe()
type EphemeralSubscriberFn func(event EphemeralStoreEvent)

func (fn EphemeralSubscriberFn) OnEphemeralEvent(event EphemeralStoreEvent) {
	fn(event)
}

// AsEphemeralSubscriber adapts a function to the EphemeralSubscriber interface.
func AsEphemeralSubscriber(fn EphemeralSubscriberFn) EphemeralSubscriber {
	return fn
}

// JsonPathSubscriberFn adapts a function to the JsonPathSubscriber interface.
//
// Use with LoroDoc.SubscribeJsonpath to be notified when a change might affect
// the result of a JSONPath query. The callback receives no arguments; re-evaluate
// the query to get the updated result.
//
//	sub, err := doc.SubscribeJsonpath("$.users[*].name", loro.JsonPathSubscriberFn(func() {
//		// a change may have affected the query result, re-evaluate it
//		results := doc.JsonpathQuery("$.users[*].name")
//		fmt.Println("names:", results)
//	}))
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer sub.Unsubscribe()
type JsonPathSubscriberFn func()

func (fn JsonPathSubscriberFn) OnJsonpathChanged() {
	fn()
}

// AsJsonPathSubscriber adapts a function to the JsonPathSubscriber interface.
func AsJsonPathSubscriber(fn JsonPathSubscriberFn) JsonPathSubscriber {
	return fn
}

// OnPushFn adapts a function to the OnPush interface.
//
// Use with UndoManager.SetOnPush to be notified when a new item is pushed onto
// the undo or redo stack. Return UndoItemMeta to attach custom metadata (e.g.
// cursor position) that will be available when the item is later popped.
//
//	push := loro.OnPushFn(func(
//		undoOrRedo loro.UndoOrRedo,
//		span loro.CounterSpan,
//		diffEvent *loro.DiffEvent,
//	) loro.UndoItemMeta {
//		return loro.UndoItemMeta{CursorsPos: cursorMap}
//	})
//	undoManager.SetOnPush(&push)
type OnPushFn func(undoOrRedo UndoOrRedo, span CounterSpan, diffEvent *DiffEvent) UndoItemMeta

func (fn OnPushFn) OnPush(undoOrRedo UndoOrRedo, span CounterSpan, diffEvent *DiffEvent) UndoItemMeta {
	return fn(undoOrRedo, span, diffEvent)
}

// AsOnPush adapts a function to the OnPush interface.
func AsOnPush(fn OnPushFn) OnPush {
	return fn
}

// OnPopFn adapts a function to the OnPop interface.
//
// Use with UndoManager.SetOnPop to be notified when an item is popped from the
// undo or redo stack (i.e. when an undo or redo is performed). The callback
// receives the metadata that was attached when the item was pushed.
//
//	pop := loro.OnPopFn(func(
//		undoOrRedo loro.UndoOrRedo,
//		span loro.CounterSpan,
//		undoMeta loro.UndoItemMeta,
//	) {
//		// restore cursor position from the saved metadata
//		restoreCursors(undoMeta.CursorsPos)
//	})
//	undoManager.SetOnPop(&pop)
type OnPopFn func(undoOrRedo UndoOrRedo, span CounterSpan, undoMeta UndoItemMeta)

func (fn OnPopFn) OnPop(undoOrRedo UndoOrRedo, span CounterSpan, undoMeta UndoItemMeta) {
	fn(undoOrRedo, span, undoMeta)
}

// AsOnPop adapts a function to the OnPop interface.
func AsOnPop(fn OnPopFn) OnPop {
	return fn
}

// --- Convenience methods that accept plain functions directly on structs ---

// SubscribeLocalUpdateFn is a convenience wrapper around LoroDoc.SubscribeLocalUpdate
// that accepts a plain function instead of a LocalUpdateCallback interface.
//
//	sub := doc.SubscribeLocalUpdateFn(func(update []byte) {
//		_ = conn.Send(update)
//	})
//	defer sub.Unsubscribe()
func (doc *LoroDoc) SubscribeLocalUpdateFn(fn func(update []byte)) *Subscription {
	return doc.SubscribeLocalUpdate(LocalUpdateCallbackFn(fn))
}

// SubscribeFn is a convenience wrapper around LoroDoc.Subscribe
// that accepts a plain function instead of a Subscriber interface.
//
//	sub := doc.SubscribeFn(containerId, func(diff loro.DiffEvent) {
//		fmt.Println("container changed:", diff.TriggeredBy)
//	})
//	defer sub.Unsubscribe()
func (doc *LoroDoc) SubscribeFn(containerId ContainerId, fn func(diff DiffEvent)) *Subscription {
	return doc.Subscribe(containerId, SubscriberFn(fn))
}

// SubscribeRootFn is a convenience wrapper around LoroDoc.SubscribeRoot
// that accepts a plain function instead of a Subscriber interface.
//
//	sub := doc.SubscribeRootFn(func(diff loro.DiffEvent) {
//		fmt.Println("document changed:", diff.TriggeredBy)
//	})
//	defer sub.Unsubscribe()
func (doc *LoroDoc) SubscribeRootFn(fn func(diff DiffEvent)) *Subscription {
	return doc.SubscribeRoot(SubscriberFn(fn))
}

// SubscribeFirstCommitFromPeerFn is a convenience wrapper around
// LoroDoc.SubscribeFirstCommitFromPeer that accepts a plain function instead
// of a FirstCommitFromPeerCallback interface.
//
//	sub := doc.SubscribeFirstCommitFromPeerFn(func(payload loro.FirstCommitFromPeerPayload) {
//		fmt.Println("new peer:", payload.PeerId)
//	})
//	defer sub.Unsubscribe()
func (doc *LoroDoc) SubscribeFirstCommitFromPeerFn(fn func(payload FirstCommitFromPeerPayload)) *Subscription {
	return doc.SubscribeFirstCommitFromPeer(FirstCommitFromPeerCallbackFn(fn))
}

// SubscribePreCommitFn is a convenience wrapper around LoroDoc.SubscribePreCommit
// that accepts a plain function instead of a PreCommitCallback interface.
//
//	sub := doc.SubscribePreCommitFn(func(payload loro.PreCommitCallbackPayload) {
//		payload.Modifier.SetMessage("auto-tagged")
//	})
//	defer sub.Unsubscribe()
func (doc *LoroDoc) SubscribePreCommitFn(fn func(payload PreCommitCallbackPayload)) *Subscription {
	return doc.SubscribePreCommit(PreCommitCallbackFn(fn))
}

// SubscribeJsonpathFn is a convenience wrapper around LoroDoc.SubscribeJsonpath
// that accepts a plain function instead of a JsonPathSubscriber interface.
//
//	sub, err := doc.SubscribeJsonpathFn("$.users[*].name", func() {
//		results := doc.JsonpathQuery("$.users[*].name")
//		fmt.Println("names:", results)
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer sub.Unsubscribe()
func (doc *LoroDoc) SubscribeJsonpathFn(path string, fn func()) (*Subscription, error) {
	return doc.SubscribeJsonpath(path, JsonPathSubscriberFn(fn))
}

// TravelChangeAncestorsFn is a convenience wrapper around LoroDoc.TravelChangeAncestors
// that accepts a plain function instead of a ChangeAncestorsTraveler interface.
// Return true to continue traversal, or false to stop early.
//
//	doc.TravelChangeAncestorsFn(ids, func(change loro.ChangeMeta) bool {
//		fmt.Println("peer:", change.Id.Peer)
//		return true
//	})
func (doc *LoroDoc) TravelChangeAncestorsFn(ids []Id, fn func(change ChangeMeta) bool) error {
	return doc.TravelChangeAncestors(ids, ChangeAncestorsTravelerFn(fn))
}

// SubscribeFn is a convenience wrapper around LoroCounter.Subscribe
// that accepts a plain function instead of a Subscriber interface.
// Returns nil if the container is detached.
//
//	sub := counter.SubscribeFn(func(diff loro.DiffEvent) {
//		fmt.Println("counter changed")
//	})
//	if sub != nil {
//		defer sub.Unsubscribe()
//	}
func (c *LoroCounter) SubscribeFn(fn func(diff DiffEvent)) *Subscription {
	result := c.Subscribe(SubscriberFn(fn))
	if result == nil {
		return nil
	}
	return *result
}

// SubscribeFn is a convenience wrapper around LoroList.Subscribe
// that accepts a plain function instead of a Subscriber interface.
// Returns nil if the container is detached.
//
//	sub := list.SubscribeFn(func(diff loro.DiffEvent) {
//		fmt.Println("list changed")
//	})
//	if sub != nil {
//		defer sub.Unsubscribe()
//	}
func (l *LoroList) SubscribeFn(fn func(diff DiffEvent)) *Subscription {
	result := l.Subscribe(SubscriberFn(fn))
	if result == nil {
		return nil
	}
	return *result
}

// SubscribeFn is a convenience wrapper around LoroMap.Subscribe
// that accepts a plain function instead of a Subscriber interface.
// Returns nil if the container is detached.
//
//	sub := m.SubscribeFn(func(diff loro.DiffEvent) {
//		fmt.Println("map changed")
//	})
//	if sub != nil {
//		defer sub.Unsubscribe()
//	}
func (m *LoroMap) SubscribeFn(fn func(diff DiffEvent)) *Subscription {
	result := m.Subscribe(SubscriberFn(fn))
	if result == nil {
		return nil
	}
	return *result
}

// SubscribeFn is a convenience wrapper around LoroMovableList.Subscribe
// that accepts a plain function instead of a Subscriber interface.
// Returns nil if the container is detached.
//
//	sub := ml.SubscribeFn(func(diff loro.DiffEvent) {
//		fmt.Println("movable list changed")
//	})
//	if sub != nil {
//		defer sub.Unsubscribe()
//	}
func (ml *LoroMovableList) SubscribeFn(fn func(diff DiffEvent)) *Subscription {
	result := ml.Subscribe(SubscriberFn(fn))
	if result == nil {
		return nil
	}
	return *result
}

// SubscribeFn is a convenience wrapper around LoroText.Subscribe
// that accepts a plain function instead of a Subscriber interface.
// Returns nil if the container is detached.
//
//	sub := text.SubscribeFn(func(diff loro.DiffEvent) {
//		fmt.Println("text changed")
//	})
//	if sub != nil {
//		defer sub.Unsubscribe()
//	}
func (t *LoroText) SubscribeFn(fn func(diff DiffEvent)) *Subscription {
	result := t.Subscribe(SubscriberFn(fn))
	if result == nil {
		return nil
	}
	return *result
}

// SubscribeFn is a convenience wrapper around LoroTree.Subscribe
// that accepts a plain function instead of a Subscriber interface.
// Returns nil if the container is detached.
//
//	sub := tree.SubscribeFn(func(diff loro.DiffEvent) {
//		fmt.Println("tree changed")
//	})
//	if sub != nil {
//		defer sub.Unsubscribe()
//	}
func (tr *LoroTree) SubscribeFn(fn func(diff DiffEvent)) *Subscription {
	result := tr.Subscribe(SubscriberFn(fn))
	if result == nil {
		return nil
	}
	return *result
}

// SubscribeFn is a convenience wrapper around EphemeralStore.Subscribe
// that accepts a plain function instead of an EphemeralSubscriber interface.
//
//	sub := store.SubscribeFn(func(event loro.EphemeralStoreEvent) {
//		fmt.Println("ephemeral keys updated:", event.Updated)
//	})
//	defer sub.Unsubscribe()
func (s *EphemeralStore) SubscribeFn(fn func(event EphemeralStoreEvent)) *Subscription {
	return s.Subscribe(EphemeralSubscriberFn(fn))
}

// SubscribeLocalUpdateFn is a convenience wrapper around EphemeralStore.SubscribeLocalUpdate
// that accepts a plain function instead of a LocalEphemeralListener interface.
//
//	sub := store.SubscribeLocalUpdateFn(func(update []byte) {
//		_ = conn.Send(update) // broadcast ephemeral update
//	})
//	defer sub.Unsubscribe()
func (s *EphemeralStore) SubscribeLocalUpdateFn(fn func(update []byte)) *Subscription {
	return s.SubscribeLocalUpdate(LocalEphemeralListenerFn(fn))
}

// SetOnPushFn is a convenience wrapper around UndoManager.SetOnPush
// that accepts a plain function instead of an OnPush interface.
//
//	undoManager.SetOnPushFn(func(
//		undoOrRedo loro.UndoOrRedo,
//		span loro.CounterSpan,
//		diffEvent *loro.DiffEvent,
//	) loro.UndoItemMeta {
//		return loro.UndoItemMeta{CursorsPos: cursorMap}
//	})
func (um *UndoManager) SetOnPushFn(fn func(undoOrRedo UndoOrRedo, span CounterSpan, diffEvent *DiffEvent) UndoItemMeta) {
	var push OnPush = OnPushFn(fn)
	um.SetOnPush(&push)
}

// SetOnPopFn is a convenience wrapper around UndoManager.SetOnPop
// that accepts a plain function instead of an OnPop interface.
//
//	undoManager.SetOnPopFn(func(
//		undoOrRedo loro.UndoOrRedo,
//		span loro.CounterSpan,
//		undoMeta loro.UndoItemMeta,
//	) {
//		restoreCursors(undoMeta.CursorsPos)
//	})
func (um *UndoManager) SetOnPopFn(fn func(undoOrRedo UndoOrRedo, span CounterSpan, undoMeta UndoItemMeta)) {
	var pop OnPop = OnPopFn(fn)
	um.SetOnPop(&pop)
}
