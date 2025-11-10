import WordCard from "../WordCard/WordCard.tsx";
import { useEffect, useMemo, useState } from "react";
import { type QuestWord } from "../../Api";
import "./Cards.css";
import "../Shared.css";

export default function Cards({
  words,
  mistakes,
  onSubmitAnswers,
}: {
  words: QuestWord[];
  mistakes: string[] | null;
  onSubmitAnswers: (answers: string[]) => Promise<void>;
}) {
  const [answers, setAnswers] = useState<string[]>(
    Array(words.length).fill("")
  );

  useEffect(() => {
    setAnswers(Array(words.length).fill(""));
  }, [words]);

  const mistakeSet = useMemo(() => new Set(mistakes), [mistakes]);

  return (
    <div className="cards-container">
      {words.map((word, idx) => {
        return (
          <WordCard
            key={word.quest_word_idx}
            wordPosition={word.position}
            thisWordAnswer={answers[idx]}
            hasMistake={mistakeSet.has((idx + 1).toString())}
            setThisWordAnswer={(answer: string) => {
              const newAnswers = [...answers];
              newAnswers[idx] = answer;
              setAnswers(newAnswers);
            }}
          />
        );
      })}
      <button onClick={() => onSubmitAnswers(answers)}>Отправить</button>
    </div>
  );
}
