package loro

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
	return ContainerIdRoot{Name: name, ContainerType: ContainerTypeMap{}}
}

// ListRoot returns the root ContainerId for a list with the given name. Useful
// for APIs that take a ContainerId directly, such as HasContainer,
// DeleteRootContainer, and GetPathToContainer.
func ListRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: ContainerTypeList{}}
}

// TextRoot returns the root ContainerId for a text container with the given
// name. Useful for APIs that take a ContainerId directly, such as
// HasContainer, DeleteRootContainer, and GetPathToContainer.
func TextRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: ContainerTypeText{}}
}

// TreeRoot returns the root ContainerId for a tree with the given name. Useful
// for APIs that take a ContainerId directly, such as HasContainer,
// DeleteRootContainer, and GetPathToContainer.
func TreeRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: ContainerTypeTree{}}
}

// MovableListRoot returns the root ContainerId for a movable list with the
// given name. Useful for APIs that take a ContainerId directly, such as
// HasContainer, DeleteRootContainer, and GetPathToContainer.
func MovableListRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: ContainerTypeMovableList{}}
}

// CounterRoot returns the root ContainerId for a counter with the given name.
// Useful for APIs that take a ContainerId directly, such as HasContainer,
// DeleteRootContainer, and GetPathToContainer.
func CounterRoot(name string) ContainerId {
	return ContainerIdRoot{Name: name, ContainerType: ContainerTypeCounter{}}
}
