package loro

import "testing"

func TestRootHelpers_ContainerTypes(t *testing.T) {
	cases := []struct {
		name string
		id   ContainerId
		want ContainerType
	}{
		{"MapRoot", MapRoot("m"), ContainerTypeMap{}},
		{"ListRoot", ListRoot("l"), ContainerTypeList{}},
		{"TextRoot", TextRoot("t"), ContainerTypeText{}},
		{"TreeRoot", TreeRoot("tr"), ContainerTypeTree{}},
		{"MovableListRoot", MovableListRoot("ml"), ContainerTypeMovableList{}},
		{"CounterRoot", CounterRoot("c"), ContainerTypeCounter{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, ok := tc.id.(ContainerIdRoot)
			if !ok {
				t.Fatalf("%s did not return ContainerIdRoot, got %T", tc.name, tc.id)
			}
			if got, want := root.ContainerType, tc.want; got != want {
				t.Fatalf("ContainerType = %T, want %T", got, want)
			}
		})
	}
}

func TestRootHelpers_HasContainer(t *testing.T) {
	doc := NewLoroDoc()
	_ = doc.GetMap(AsContainerId("m"))
	_ = doc.GetList(AsContainerId("l"))
	_ = doc.GetText(AsContainerId("t"))
	_ = doc.GetMovableList(AsContainerId("ml"))
	_ = doc.GetCounter(AsContainerId("c"))
	_ = doc.GetTree(AsContainerId("tr"))

	if !doc.HasContainer(MapRoot("m")) {
		t.Fatal("HasContainer(MapRoot) should be true")
	}
	if !doc.HasContainer(ListRoot("l")) {
		t.Fatal("HasContainer(ListRoot) should be true")
	}
	if !doc.HasContainer(TextRoot("t")) {
		t.Fatal("HasContainer(TextRoot) should be true")
	}
	if !doc.HasContainer(MovableListRoot("ml")) {
		t.Fatal("HasContainer(MovableListRoot) should be true")
	}
	if !doc.HasContainer(CounterRoot("c")) {
		t.Fatal("HasContainer(CounterRoot) should be true")
	}
	if !doc.HasContainer(TreeRoot("tr")) {
		t.Fatal("HasContainer(TreeRoot) should be true")
	}
}
