export default function LinkItem({ originalUrl, shortUrl, visits: visitsNumber }) {
  return (
    <tr>
      <td>
        <a href={shortUrl} target="_blank" rel="noreferrer">
          {shortUrl}
        </a>
      </td>
      <td>{originalUrl}</td>
      <td>{visitsNumber}</td>
      {/* <td>{new Date(createdAt).toLocaleString()}</td> */}
    </tr>
  );
}