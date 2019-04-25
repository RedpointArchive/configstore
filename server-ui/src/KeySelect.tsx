import * as React from "react";
import {
  Key,
  Schema,
  SchemaField,
  MetaEntity,
  ValueType,
  MetaListEntitiesRequest
} from "./api/meta_pb";
import AsyncSelect from "react-select/lib/Async";
import { serializeKey, deserializeKey, prettifyKey, g } from "./core";
import { useAsync } from "react-async";
import { createGrpcPromiseClient } from "./svcHost";
import { ConfigstoreMetaServicePromiseClient } from "./api/meta_grpc_web_pb";

export interface KeySelectProps {
  value: Key | undefined;
  onChange(value: Key | undefined): void;
  schema: Schema;
  field: SchemaField;
  className?: string;
}

async function getUnfilteredOptionsMap(props: {
  schema: Schema;
  field: SchemaField;
}) {
  const client = createGrpcPromiseClient(ConfigstoreMetaServicePromiseClient);
  const kindsMap = props.schema.getKindsMap();
  const results = await Promise.all(
    kindsMap
      .toArray()
      .map(kv => kv[0])
      .filter(kindName => {
        const fieldEditor = props.field.getEditor();
        if (fieldEditor !== undefined) {
          const allowedKinds = fieldEditor.getAllowedkindsList();
          if (allowedKinds.length === 0) {
            return true;
          } else {
            return allowedKinds.includes(kindName);
          }
        }
      })
      .map(kindName =>
        (async () => {
          const kindForLookup = kindsMap.get(kindName);
          if (kindForLookup === undefined) {
            return {
              kindName: kindName,
              selector: () => "error",
              entities: []
            };
          }
          try {
            let selector = (r: MetaEntity) => prettifyKey(g(r.getKey()));
            const kindForLookupEditor = kindForLookup.getEditor();
            if (kindForLookupEditor !== undefined) {
              if (
                kindForLookupEditor.getRendereditordropdownwithfield() !== ""
              ) {
                selector = (r: MetaEntity) => {
                  const values = r.getValuesList();
                  const fields = kindForLookup
                    .getFieldsList()
                    .filter(
                      x =>
                        x.getName() ===
                        kindForLookupEditor.getRendereditordropdownwithfield()
                    );
                  const lookupField =
                    fields.length === 0 ? undefined : fields[0];
                  if (lookupField !== undefined) {
                    const value = values.filter(
                      x => x.getId() === lookupField.getId()
                    )[0];
                    switch (lookupField.getType()) {
                      case ValueType.STRING:
                        return value.getStringvalue();
                      default:
                        return `(unsupported field type ${lookupField.getType()} for 'renderEditorDropdownWithField' setting): ${prettifyKey(
                          g(r.getKey())
                        )}`;
                    }
                  }
                  return prettifyKey(g(r.getKey()));
                };
              }
            }

            const req = new MetaListEntitiesRequest();
            req.setKindname(kindName);
            req.setStart("");
            req.setLimit(0);
            return {
              kindName: kindName,
              selector: selector,
              entities: (await client.svc.metaList(
                req,
                client.meta
              )).getEntitiesList()
            };
          } catch (e) {
            console.error(e);
            return {
              kindName: kindName,
              selector: () => "error",
              entities: []
            };
          }
        })()
      )
  );
  return results.map(result => ({
    label: result.kindName,
    options: result.entities.map(r => ({
      value: serializeKey(g(r.getKey())),
      label: result.selector(r)
    }))
  }));
}

type AsyncResult = {
  label: string;
  options: {
    value: string;
    label: string;
  }[];
}[];

export const KeySelect = (props: KeySelectProps) => {
  const { data, isLoading } = useAsync<AsyncResult>({
    promiseFn: getUnfilteredOptionsMap,
    schema: props.schema,
    field: props.field
  } as any);

  if (isLoading || data === undefined) {
    return (
      <AsyncSelect
        key="loading"
        defaultValue={undefined}
        className={props.className}
        isLoading={true}
        cacheOptions
        defaultOptions
        loadOptions={async () => {
          return [];
        }}
      />
    );
  }

  let value = undefined;
  if (props.value !== undefined) {
    for (const group of data) {
      for (const option of group.options) {
        if (option.value === serializeKey(props.value)) {
          value = option;
        }
      }
    }
  }

  return (
    <AsyncSelect
      isClearable
      key={"loaded"}
      className={props.className}
      value={value}
      onChange={(selectedOption: any) => {
        if (
          selectedOption === undefined ||
          selectedOption === null ||
          selectedOption instanceof Array
        ) {
          props.onChange(undefined);
        } else {
          props.onChange(deserializeKey(selectedOption.value));
        }
      }}
      isLoading={isLoading}
      cacheOptions
      defaultOptions
      loadOptions={async inputValue =>
        data.map(group => ({
          label: group.label,
          options: group.options.filter(vl =>
            vl.label.toLowerCase().includes(inputValue.toLowerCase())
          )
        }))
      }
    />
  );
};
