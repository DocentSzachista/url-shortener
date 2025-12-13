import { useEffect, useState } from "react";
import { fetchLinks } from "../api/shortenApi";
import LinkItem from "./LinkItemComp";
import { HOSTNAME } from "../config";

export default function LinksList() {
  const [links, setLinks] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchLinks()
      .then(setLinks)
      .catch(() => alert("Nie udało się pobrać linków"))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p>Ładowanie...</p>;
  if (!links.length) return <p>Brak linków</p>;
    console.log(links);
  return (
    <article>
      <h3>Utworzone linki</h3>
      <table>
        <thead>
          <tr>
            <th>Skrócony</th>
            <th>Oryginalny</th>
            <th>Liczba wejść</th>
          </tr>
        </thead>
        <tbody>
          {links.map(link => (
            <LinkItem key={link.shortedId} originalUrl={link.Url} shortUrl={`${HOSTNAME}/${link.ShortedId}`} visitsNumber={link.Visits} />
          ))}
        </tbody>
      </table>
    </article>
  );
}