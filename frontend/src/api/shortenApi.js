export async function shortenUrl(url, customShort) {
  const body = { url };
  if (customShort) body.customShort = customShort;

  const res = await fetch("http://localhost:8080/add", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });

  if (!res.ok) throw new Error("Add failed");
  return res.json();
}


export async function fetchLinks() {
  const res = await fetch("http://localhost:8080/links");

  if (!res.ok) {
    throw new Error("Fetch failed");
  }

  return res.json();
}
