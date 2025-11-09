import { useState } from "react";

export default function TeamForm({
  onSubmitTeam,
}: {
  onSubmitTeam: (teamId: number) => Promise<void>;
}) {
  const [teamID, setTeamID] = useState("");
  const teamIDMaxLen = 10;

  return (
    <>
      <input
        placeholder="team id"
        onChange={(data) => {
          setTeamID(
            data.currentTarget.value
              .split("")
              .filter((ch) => /^\d$/.test(ch))
              .slice(0, teamIDMaxLen)
              .join("")
          );
        }}
        value={teamID}
      ></input>
      <button
        onClick={() => onSubmitTeam(Number.parseInt(teamID))}
        disabled={teamID === null || teamID.length == 0}
      >
        получить задание
      </button>
    </>
  );
}
