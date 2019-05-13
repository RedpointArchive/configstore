import * as React from "react";
import BigInt from "big-integer";

interface Props {
  value: string;
  onChange: (value: string, isValid: boolean) => void;
  placeholder?: string;
  readOnly?: boolean;
  className?: string;
}

export const nibblinsToDollarString = (val: BigInt.BigInteger): string => {
  const r = val.divmod(BigInt("1e11"));
  const dollars = r.quotient;
  const nibblins = r.remainder;

  function pad(num: string, size: number) {
    let s = num;
    while (s.length < size) s = "0" + s;
    return s;
  }

  function trimRight(num: string) {
    while (num.length > 2 && num.substr(num.length - 1) === "0") {
      num = num.substr(0, num.length - 1);
    }
    return num;
  }

  return dollars.toString() + "." + trimRight(pad(nibblins.toString(), 11));
};

const dollarStringToNibblins = (str: string): BigInt.BigInteger => {
  if (str.length == 0) {
    return BigInt();
  }

  let decimal = str.indexOf(".");
  if (decimal == -1) {
    decimal = str.length;
  }

  let start = 0;
  if (str[0] === "-") {
    start = 1;
  }

  let dollars = BigInt();
  if (decimal > start) {
    dollars = BigInt(str.substr(start, decimal - start), 10);
  }

  let nibblins = BigInt();
  if (decimal + 1 < str.length) {
    length = str.length - (decimal + 1);
    if (length < 11) {
      length = 11;
    }
    if (length > 11) {
      throw new Error("too many decimal places (limit to 11)");
    }
    for (let i = 0; i < length; i += 1) {
      if (i < 11) {
        nibblins = nibblins.multiply(BigInt("10"));
      }
      const index = decimal + 1 + i;
      if (index < str.length) {
        const char = str.charCodeAt(index);
        if (char < "0".charCodeAt(0) || char > "9".charCodeAt(0)) {
          throw new Error("invalid dollar string: " + str);
        }
        if (i < 11) {
          nibblins = nibblins.add(BigInt(String.fromCharCode(char)));
        }
      }
    }
  }

  if (str[0] === "-") {
    dollars = dollars.multiply(BigInt("-1"));
    nibblins = nibblins.multiply(BigInt("-1"));
  }

  return dollars.multiply(BigInt("1e11")).add(nibblins);
};

export const FinancialInput = (props: Props) => {
  const defaultInput = React.useMemo(() => {
    return props.value.startsWith("ERROR:")
      ? props.value.substr("ERROR:".length)
      : nibblinsToDollarString(BigInt(props.value));
  }, [
    props.value.startsWith("ERROR:")
      ? props.value.substr("ERROR:".length)
      : nibblinsToDollarString(BigInt(props.value))
  ]);

  const [input, setInput] = React.useState(defaultInput);

  return (
    <input
      value={input}
      onChange={e => {
        setInput(e.target.value);
        try {
          const bir = dollarStringToNibblins(e.target.value);
          props.onChange(bir.toString(), true);
        } catch (err) {
          props.onChange("ERROR:" + input, false);
        }
      }}
      type="text"
      placeholder={props.placeholder}
      readOnly={props.readOnly}
      className={props.className}
    />
  );
};
