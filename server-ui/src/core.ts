import { Key, PathElement, PartitionId } from "./api/meta_pb";

export function g<T>(v: T | undefined): T {
  return v as T;
}

export function serializeKey(key: Key): string {
  return `ns=${g(key.getPartitionid()).getNamespace()}|${key
    .getPathList()
    .map(pe => {
      if (pe.getIdtypeCase() === PathElement.IdtypeCase.ID) {
        return `${pe.getKind()}:id=${pe.getId()}`;
      } else if (pe.getIdtypeCase() === PathElement.IdtypeCase.NAME) {
        return `${pe.getKind()}:name=${pe.getName()}`;
      } else {
        return `${pe.getKind()}:unset`;
      }
    })
    .join("|")}`;
}

export function deserializeKey(keyString: string): Key {
  const components = keyString.split("|");
  const ns = components[0].substr(3);
  const partitionId = new PartitionId();
  partitionId.setNamespace(ns);
  const key = new Key();
  key.setPartitionid(partitionId);
  for (let i = 1; i < components.length; i++) {
    const subcomponent = components[i].split(":", 2);
    if (subcomponent[1].startsWith("id=")) {
      const pathElement = new PathElement();
      pathElement.setKind(subcomponent[0]);
      pathElement.setId(parseInt(subcomponent[1].substr(3)));
      key.addPath(pathElement);
    } else if (subcomponent[1].startsWith("name=")) {
      const pathElement = new PathElement();
      pathElement.setKind(subcomponent[0]);
      pathElement.setName(subcomponent[1].substr(5));
      key.addPath(pathElement);
    } else if (subcomponent[1] == "unset") {
      const pathElement = new PathElement();
      pathElement.setKind(subcomponent[0]);
      key.addPath(pathElement);
    }
  }
  return key;
}
