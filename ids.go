package loro

// ContainerType singletons. Use these instead of the zero-value struct
// literals (e.g. ContainerTypeMap{}) when building ContainerIdRoot or calling
// any API that takes a ContainerType.
//
//	doc.GetContainer(loro.ContainerIdRoot{Name: "data", ContainerType: loro.MapType})
var (
	MapType         ContainerType = ContainerTypeMap{}
	ListType        ContainerType = ContainerTypeList{}
	TextType        ContainerType = ContainerTypeText{}
	TreeType        ContainerType = ContainerTypeTree{}
	MovableListType ContainerType = ContainerTypeMovableList{}
	CounterType     ContainerType = ContainerTypeCounter{}
)

// StringContainerId is a string that can be used to get maps, lists, etc from
// a document.
type StringContainerId string

func (c StringContainerId) AsContainerId(containerType ContainerType) ContainerId {
	return ContainerIdRoot{
		Name:          string(c),
		ContainerType: containerType,
	}
}

// AsContainerId converts a string to a ContainerIDLike, which can be used to
// get maps, lists, etc from a document.
func AsContainerId(v string) StringContainerId {
	return StringContainerId(v)
}

var _ ContainerIdLike = StringContainerId("")

// MapRoot returns the root ContainerId for a map with the given name. Useful
// for APIs that take a ContainerId directly, such as HasContainer,
// DeleteRootContainer, and GetPathToContainer.
func MapRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: MapType}
}

// ListRoot returns the root ContainerId for a list with the given name. Useful
// for APIs that take a ContainerId directly, such as HasContainer,
// DeleteRootContainer, and GetPathToContainer.
func ListRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: ListType}
}

// TextRoot returns the root ContainerId for a text container with the given
// name. Useful for APIs that take a ContainerId directly, such as
// HasContainer, DeleteRootContainer, and GetPathToContainer.
func TextRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: TextType}
}

// TreeRoot returns the root ContainerId for a tree with the given name. Useful
// for APIs that take a ContainerId directly, such as HasContainer,
// DeleteRootContainer, and GetPathToContainer.
func TreeRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: TreeType}
}

// MovableListRoot returns the root ContainerId for a movable list with the
// given name. Useful for APIs that take a ContainerId directly, such as
// HasContainer, DeleteRootContainer, and GetPathToContainer.
func MovableListRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: MovableListType}
}

// CounterRoot returns the root ContainerId for a counter with the given name.
// Useful for APIs that take a ContainerId directly, such as HasContainer,
// DeleteRootContainer, and GetPathToContainer.
func CounterRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: CounterType}
}
