import { useState } from "react";
import TeamForm from "./components/TeamForm/TeamForm";
import {
  Api,
  type QuestResponse,
  type ValidationError,
  type AuthorizationError,
  type InternalServerError,
} from "./Api.ts";
import Cards from "./components/Cards/Cards.tsx";
import "./App.css";

function App() {
  const [questData, setQuestData] = useState<QuestResponse | null>(null);
  const [validationResult, setValidationResult] = useState<string | null>(null);
  const [validationMistakes, setValidationMistakes] = useState<string | null>(
    null
  );
  const [questTeamID, setQuestTeamID] = useState<number | null>(null);

  const api = new Api({
    baseUrl: window.config.API_BASE_URL,
  });

  const handleFetchQuest = async (teamId: number) => {
    try {
      const data = await api.api.v1QuestList({ team_id: teamId });
      setQuestData(data.data);
      setQuestTeamID(teamId);
    } catch (e) {
      console.log(e);
    }
  };

  const handleValidate = async (answers: string[]) => {
    if (questTeamID === null || questData === null) {
      console.log("validate called before get quest request");
      return;
    }
    try {
      await api.api.v1ValidateCreate({
        team_id: questTeamID,
        answer: answers,
        quest_id: questData?.quest_id,
      });
      setValidationMistakes(null);
      setValidationResult("Success");
    } catch (e) {
      handleValidationError(e, setValidationMistakes);
    }
  };

  return (
    <div className="App">
      <h2>Library quest</h2>
      {questData === null && <TeamForm onSubmitTeam={handleFetchQuest} />}
      {questData !== null && validationResult === null && (
        <Cards words={questData.words} onSubmitAnswers={handleValidate} />
      )}
      {validationMistakes !== null && (
        <p className="message error">Ошибки в словах: {validationMistakes}</p>
      )}
      {validationResult !== null && (
        <p className="message success">{validationResult}</p>
      )}
    </div>
  );
}

function handleValidationError(
  err: unknown,
  setResultUI: (res: string) => void
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
        setResultUI(authErr.error || "Authorization failed.");
        return;
      }
      case 400: {
        const valErr = error as ValidationError;
        if (valErr.mistakes === undefined) {
          setResultUI("mistaks undefined");
          return;
        }
        setResultUI(valErr.mistakes?.map((id) => id + 1).join(", "));
        return;
      }
      case 500: {
        const serverErr = error as InternalServerError;
        setResultUI(serverErr.error || "An internal server error occurred.");
        return;
      }
      default: {
        setResultUI(`Unhandled error with status: ${status}`);
        return;
      }
    }
  }

  console.log(err);
  setResultUI("An unknown or network error occurred.");
}

export default App;
