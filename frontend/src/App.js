import { useState } from "react";
import ShortenForm from "./components/ShortenForm";
import Result from "./components/Result";
import LinksList from "./components/LinksList";

export default function App() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [shortUrl, setShortUrl] = useState("");

  function handleCreated(newShortUrl) {
    setShortUrl(newShortUrl);
    setRefreshKey(k => k + 1); // trigger reload listy
  }

  return (
    <main className="container">
      <article>
        <h2>Skracacz linków</h2>
        <ShortenForm onCreated={handleCreated} />
        <Result shortUrl={shortUrl} />
      </article>

      <LinksList refreshKey={refreshKey} />
    </main>
  );
}
