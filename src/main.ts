import fetch from "cross-fetch";

import * as z from "zod";

const resultsSchema = z.array(z.array(z.string()));

async function main() {
  const query = new URLSearchParams({
    keyword: "長野",
  });
  const res = await fetch(
    `http://localhost:5001/postalcode-firebase/us-central1/postalcode?${query}`
  );
  const results: string[][] = resultsSchema.parse(await res.json());
  console.log(results);
}

main().catch((e) => {
  console.error(e);
});
