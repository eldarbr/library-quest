import { type WordPosition } from "../../Api";
import "./WordCard.css";
import "../Shared.css";

export default function WordCard({
  wordPosition,
  thisWordAnswer,
  hasMistake,
  onEnterKey,
  setThisWordAnswer,
}: {
  wordPosition: WordPosition;
  thisWordAnswer: string;
  hasMistake: boolean;
  onEnterKey: () => void;
  setThisWordAnswer: (newAnswer: string) => void;
}) {
  const maxAnswerLen = 25;

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
          onKeyDown={(e) => {
            if (e.key !== "Enter") return;
            onEnterKey();
          }}
          placeholder="Enter word here..."
          style={
            hasMistake ? { backgroundColor: "rgba(255, 180, 180, 1)" } : {}
          }
        />
      </>
    </div>
  );
}
