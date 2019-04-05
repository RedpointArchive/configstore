import { Key, PathElement, PartitionId } from "./api/meta_pb";

export function g<T>(v: T | undefined): T {
  return v as T;
}

export function serializeKey(key: Key): string {
  return `ns=${g(key.getPartitionid()).getNamespace()}|${key
    .getPathList()
    .map(pe => {
      if (pe.getIdtypeCase() === PathElement.IdtypeCase.ID) {
        return `id=${pe.getId()}`;
      } else if (pe.getIdtypeCase() === PathElement.IdtypeCase.NAME) {
        return `name=${pe.getName()}`;
      } else {
        return `unset`;
      }
    })
    .join("|")}`;
}

export function deserializeKey(keyString: string): Key {
  const components = keyString.split('|');
  const ns = components[0].substr(3);
  const partitionId = new PartitionId();
  partitionId.setNamespace(ns);
  const key = new Key();
  key.setPartitionid(partitionId);
  for (let i = 1; i < components.length; i++) {
    if (components[i].startsWith("id=")) {
      const pathElement = new PathElement();
      pathElement.setId(parseInt(components[i].substr(3)));
      key.addPath(pathElement);
    } else if (components[i].startsWith("name=")) {
      const pathElement = new PathElement();
      pathElement.setName(components[i].substr(5));
      key.addPath(pathElement);
    } else if (components[i] == "unset") {
      const pathElement = new PathElement();
      key.addPath(pathElement);
    }
  }
  return key;
}
