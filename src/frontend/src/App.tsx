import { useEffect, useState } from "react";
import {
  Api,
  type QuestResponse,
  type ValidationError,
  type AuthorizationError,
  type InternalServerError,
} from "./Api.ts";
import Cards from "./components/Cards/Cards.tsx";
import "./App.css";
import ShowDiscoveredResult from "./components/DiscoveredResult/DiscoveredResult.tsx";

const api = new Api({
  baseUrl: window.config.API_BASE_URL,
});

function App() {
  const [questData, setQuestData] = useState<QuestResponse | null>(null);
  const [discoveredQuestKeyword, setdiscoveredQuestKeyword] = useState<
    string | null
  >(null);
  const [validationMistakes, setValidationMistakes] = useState<string[] | null>(
    null
  );

  const handleValidate = async (answers: string[]) => {
    if (questData === null) {
      console.log("validate called before get quest request");
      return;
    }
    try {
      const res = await api.api.v1ValidateRandomCreate({
        answer: answers,
        quest_id: questData?.quest_id,
      });
      setValidationMistakes(null);
      setdiscoveredQuestKeyword(res.data.keyword);
    } catch (e) {
      handleValidationError(e, setValidationMistakes);
    }
  };

  // load random quest
  useEffect(() => {
    const handleFetchQuest = async () => {
      try {
        const data = await api.api.v1QuestRandomList();
        setQuestData(data.data);
      } catch (err) {
        setValidationMistakes([String(err)]);
      }
    };
    handleFetchQuest();
  }, []);

  return (
    <div className="App">
      <h1>Library quest 🐣</h1>
      {questData !== null && discoveredQuestKeyword === null && (
        <Cards
          words={questData.words}
          mistakes={validationMistakes}
          onSubmitAnswers={handleValidate}
        />
      )}
      {validationMistakes !== null && (
        <p className="message error">Ошибки: {validationMistakes.join(", ")}</p>
      )}
      {discoveredQuestKeyword !== null && (
        <ShowDiscoveredResult discoveredQuestKeyword={discoveredQuestKeyword} />
      )}
    </div>
  );
}

function handleValidationError(
  err: unknown,
  setResultUI: (res: string[]) => void
) {
  if (
    typeof err === "object" &&
    err !== null &&
    "status" in err &&
    typeof err.status === "number" &&
    "error" in err &&
    typeof err.error === "object"
  ) {
    const { status, error } = err;

    switch (status) {
      case 401: {
        const authErr = error as AuthorizationError;
        setResultUI([authErr.error || "Authorization failed."]);
        return;
      }
      case 400: {
        const valErr = error as ValidationError;
        if (valErr.mistakes === undefined) {
          setResultUI(["mistaks undefined"]);
          return;
        }
        setResultUI(valErr.mistakes?.map((id) => (id + 1).toString()));
        return;
      }
      case 500: {
        const serverErr = error as InternalServerError;
        setResultUI([serverErr.error || "An internal server error occurred."]);
        return;
      }
      default: {
        setResultUI([`Unhandled error with status: ${status}`]);
        return;
      }
    }
  }

  console.log(err);
  setResultUI(["An unknown or network error occurred."]);
}

export default App;
