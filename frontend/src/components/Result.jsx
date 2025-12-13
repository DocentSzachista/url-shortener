export default function Result({ shortUrl }) {
  if (!shortUrl) return null;

  return (
    <p>
      👉 <a href={shortUrl} target="_blank" rel="noreferrer">
        {shortUrl}
      </a>
    </p>
  );
}