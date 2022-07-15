import * as functions from "firebase-functions";
import * as fs from "fs";
import * as z from "zod";
import { performance } from "perf_hooks";
import { Index } from "flexsearch";

const paramSchema = z.object({
  q: z.string().max(100).default(""),
  limit: z
    .string()
    .regex(/^[0-9]{1,3}$/)
    .default("10"),
});

const kennallString = fs.readFileSync("./x-ken-all-hiragana.json");
const kennall: string[][] = JSON.parse(kennallString.toString());
let index: undefined | "Loading" | Index;

async function loadIndex(): Promise<void> {
  const indexStart = performance.now();
  const loadingIndex = new Index();
  for (const key of ["reg", "cfg", "map", "ctx"]) {
    await loadingIndex.import(
      key,
      fs.readFileSync("./export/" + key).toString()
    );
  }
  index = loadingIndex;
  const indexEnd = performance.now();
  functions.logger.log(
    `index imported in ${(indexEnd - indexStart).toFixed()}msecs`
  );
}

function linearSearch(q: string, limit: number): string[][] {
  const results: string[][] = [];

  for (const row of kennall) {
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
    ] = row;

    if (
      address1Furigana.includes(q) ||
      address2Furigana.includes(q) ||
      address3Furigana.includes(q) ||
      address1.includes(q) ||
      address2.includes(q) ||
      address3.includes(q)
    ) {
      results.push(row);
      if (limit <= results.length) {
        return results;
      }
    }
  }

  return results;
}

function search(q: string, limit: number): string[][] {
  if (index === undefined) {
    index = "Loading";
    loadIndex();
    return linearSearch(q, limit);
  } else if (index == "Loading") {
    return linearSearch(q, limit);
  } else {
    const flex = index.search(q, limit);
    const hits = flex.map((i) => kennall[Number(i)]);
    return hits;
  }
}

export const postalcode = functions
  .runWith({
    timeoutSeconds: 300,
    memory: "2GB", // 2GBが最小、4GB以上にしてもインポート速度に大差ない
    // minInstances: 1,
    maxInstances: 4,
  })
  .region("asia-northeast1")
  .https.onRequest(async (request, response) => {
    const params = paramSchema.parse(request.query);
    const limit = Math.max(1, Math.min(100, Number(params.limit)));
    const start = performance.now();
    const hits = search(params.q, limit);
    const end = performance.now();
    response.status(200);
    response.end(JSON.stringify({ hits, time: end - start }));
  });
