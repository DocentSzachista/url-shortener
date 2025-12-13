import { shortenUrl } from "../api/shortenApi";
import { useState } from "react";

export default function ShortenForm({ onCreated }) {
  const [custom, setCustom] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    const url = e.target.url.value;

    const data = await shortenUrl(url, custom || undefined);

    const shortUrl = `http://localhost:8080/${data.shortedID}`;
    onCreated(shortUrl);

    e.target.reset();
    setCustom("");
  }

  return (
    <form onSubmit={handleSubmit}>
      <input name="url" type="url" placeholder="https://example.com" required />
      <input
        name="customShort"
        type="text"
        placeholder="Opcjonalny skrót"
        value={custom}
        onChange={e => setCustom(e.target.value)}
      />
      <button>Skróć</button>
    </form>
  );
}