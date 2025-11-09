import { type WordPosition } from "../../Api";
import "./WordCard.css";
import "../Shared.css";

export default function WordCard({
  wordPosition,
  thisWordAnswer,
  setThisWordAnswer,
}: {
  wordPosition: WordPosition;
  thisWordAnswer: string;
  setThisWordAnswer: (newAnswer: string) => void;
}) {
  const maxAnswerLen = 20;

  return (
    <div className="word-card">
      <div className="word-card-details">
        <p>
          Книга: <span>{wordPosition.book}</span>
        </p>
        <p>
          Страница: <span>{wordPosition.page}</span>
        </p>
        <p>
          Строка: <span>{wordPosition.line}</span>
        </p>
        <p>
          Слово: <span>{wordPosition.word}</span>
        </p>
      </div>
      <>
        <input
          value={thisWordAnswer}
          onInput={(e) =>
            setThisWordAnswer(e.currentTarget.value.slice(0, maxAnswerLen))
          }
          placeholder="Enter word here..."
        />
      </>
    </div>
  );
}
