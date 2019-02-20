import { Key, PathElement } from "./api/meta_pb";

export function g<T>(v: T | undefined): T {
  return v as T;
}

export function serializeKey(key: Key): string {
  return `ns=${g(key.getPartitionid()).getNamespace()}|${key
    .getPathList()
    .map(pe => {
      if (pe.getId() !== undefined) {
        return `id=${pe.getId()}`;
      } else if (pe.getName() !== undefined) {
        return `name=${pe.getName()}`;
      } else {
        return `unset`;
      }
    })
    .join("|")}`;
}
