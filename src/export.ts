import { Index } from "flexsearch";
import fs from "fs";

async function main() {
  const kennallString = fs.readFileSync("./functions/x-ken-all-hiragana.json");
  const kennall: string[][] = JSON.parse(kennallString.toString());

  const index = new Index({
    preset: "memory",
    tokenize: "full",
    optimize: true,
    resolution: 1,
    context: false,
  });
  console.log(`indexing: ${kennall.length} items`);
  const indexStart = performance.now();
  for (let i = 0; i < kennall.length; i++) {
    if (i % 100 === 0) {
      console.log(i);
    }

    const [
      ,
      ,
      ,
      address1Furigana,
      address2Furigana,
      address3Furigana,
      address1,
      address2,
      address3,
    ] = kennall[i];

    index.add(i, address1Furigana);
    index.append(i, address2Furigana);
    index.append(i, address3Furigana);
    index.append(i, address1);
    index.append(i, address2);
    index.append(i, address3);
  }
  const indexEnd = performance.now();
  console.log(`index completed: ${(indexEnd - indexStart).toFixed()}msecs`);

  await index.export((key, data) => {
    // https://github.com/nextapps-de/flexsearch/issues/290
    const k = key.toString().split(".").pop() || "";

    console.log("exporting " + k);
    fs.writeFileSync(`./functions/export/${k}`, data);
  });

  console.log("completed");
}

main().catch((e) => {
  console.error(e);
});
