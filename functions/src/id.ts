const digits =
  " !#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~";

const base = digits.length;

export function indexToId(index: number): string {
  let id = "";
  while (0 < index) {
    id += digits[index % base];
    index = Math.floor(index / base);
  }
  return id;
}

export function idToIndex(id: string): number {
  let index = 0;
  let t = 1;
  for (let i = 0; i < id.length; i++) {
    const digit = digits.indexOf(id[i]);
    index += digit * t;
    t *= base;
  }
  return index;
}

export function test() {
  for (let i = 0; i < 100; i++) {
    const index = Math.floor(100000 * Math.random());
    const encoded = indexToId(index);
    const decoded = idToIndex(encoded);
    if (index !== decoded) {
      console.log("index", index);
      console.log("encoded", encoded);
      console.log("decoded", decoded);
    }
  }
}

test();
