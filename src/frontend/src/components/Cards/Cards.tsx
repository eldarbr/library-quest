import WordCard from "../WordCard/WordCard.tsx";
import { useEffect, useState } from "react";
import { type QuestWord } from "../../Api";

export default function Cards({
  words,
  onSubmitAnswers,
}: {
  words: QuestWord[];
  onSubmitAnswers: (answers: string[]) => Promise<void>;
}) {
  const [answers, setAnswers] = useState<string[]>(
    Array(words.length).fill("")
  );

  useEffect(() => {
    setAnswers(Array(words.length).fill(""));
  }, [words]);

  return (
    <>
      {words.map((word, idx) => {
        return (
          <WordCard
            key={word.quest_word_idx}
            wordDescription={word.position}
            thisWordAnswer={answers[idx]}
            setThisWordAnswer={(answer: string) => {
              const newAnswers = [...answers];
              newAnswers[idx] = answer;
              setAnswers(newAnswers);
            }}
          />
        );
      })}
      <button onClick={() => onSubmitAnswers(answers)}>Отправить</button>
    </>
  );
}
