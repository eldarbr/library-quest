export default function WordCard({
  wordDescription,
  thisWordAnswer,
  setThisWordAnswer,
}: {
  wordDescription: WordDescription;
  thisWordAnswer: string;
  setThisWordAnswer: (newAnswer: string) => void;
}) {
  const maxAnswerLen = 20;

  return (
    <>
      <p>Книга: {wordDescription.book}</p>
      <p>Страница: {wordDescription.page}</p>
      <p>Строка: {wordDescription.line}</p>
      <p>Слово: {wordDescription.word}</p>
      <>
        <input
          value={thisWordAnswer}
          onInput={(e) =>
            setThisWordAnswer(e.currentTarget.value.slice(0, maxAnswerLen))
          }
        />
      </>
    </>
  );
}

interface WordDescription {
  book: string;
  page: string;
  line: string;
  word: string;
}
