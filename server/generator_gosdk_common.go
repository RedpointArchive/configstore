package main

func generateGoCode_common() string {
	return `
func Fnv64a(val string) uint64 {
	hash := fnv.New64a()
	hash.Write(([]byte)(val))
	return hash.Sum64()
}

func Fnv64aPair(a uint64, b uint64) uint64 {
	hash := fnv.New64a()
	if a > b {
		tmp := a
		a = b
		b = tmp
	}
	key := make([]byte, 16)
	binary.LittleEndian.PutUint64(key, a)
	binary.LittleEndian.PutUint64(key[8:], b)
	hash.Write(key)
	return hash.Sum64()
}

func CreateTopLevelKey(partitionId *PartitionId, pathElement *PathElement) *Key {
	return &Key{
		PartitionId: partitionId,
		Path: []*PathElement{pathElement},
	}
}

func CreateIncompleteTopLevelKey(partitionId *PartitionId, kind string) *Key {
	return &Key{
		PartitionId: partitionId,
		Path: []*PathElement{
			&PathElement{
				Kind: kind,
			},
		},
	}
}

func CreateDescendantKey(parent *Key, pathElement *PathElement) *Key {
	newKey := &Key{
		PartitionId: &PartitionId{
			Namespace: parent.PartitionId.Namespace,
		},
		Path: nil,
	}
	for _, elem := range parent.Path {
		switch elem.IdType.(type) {
		case *PathElement_Id:
			newKey.Path = append(newKey.Path, &PathElement{
				Kind: elem.Kind,
				IdType: &PathElement_Id{
					Id: elem.GetId(),
				},
			})
			break
		case *PathElement_Name:
			newKey.Path = append(newKey.Path, &PathElement{
				Kind: elem.Kind,
				IdType: &PathElement_Name{
					Name: elem.GetName(),
				},
			})
			break
		}
	}
	switch pathElement.IdType.(type) {
	case *PathElement_Id:
		newKey.Path = append(newKey.Path, &PathElement{
			Kind: pathElement.Kind,
			IdType: &PathElement_Id{
				Id: pathElement.GetId(),
			},
		})
		break
	case *PathElement_Name:
		newKey.Path = append(newKey.Path, &PathElement{
			Kind: pathElement.Kind,
			IdType: &PathElement_Name{
				Name: pathElement.GetName(),
			},
		})
		break
	}
	return newKey
}

func SerializeTimestamp(ts *timestamp.Timestamp) string {
	if ts == nil {
		return ""
	}

	return ts.String()
}

func SerializeKey(key *Key) string {
	if key == nil {
		return ""
	}

	var elements []string
	for _, pathElement := range key.Path {
		if _, ok := pathElement.IdType.(*PathElement_Id); ok {
			elements = append(elements, fmt.Sprintf("%s:id=%d", pathElement.GetKind(), pathElement.GetId()))
		} else if _, ok := pathElement.IdType.(*PathElement_Name); ok {
			elements = append(elements, fmt.Sprintf("%s:name=%s", pathElement.GetKind(), pathElement.GetName()))
		} else {
			elements = append(elements, fmt.Sprintf("%s:unset", pathElement.GetKind()))
		}
	}
	return fmt.Sprintf("ns=%s|%s", key.PartitionId.Namespace, strings.Join(elements, "|"))
}

func CompareKeys(a *Key, b *Key) bool {
	return SerializeKey(a) == SerializeKey(b)
}
`
}
