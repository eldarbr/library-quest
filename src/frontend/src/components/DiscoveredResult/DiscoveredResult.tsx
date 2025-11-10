import { useState } from "react";

export default function ShowDiscoveredResult({
  discoveredQuestKeyword,
}: {
  discoveredQuestKeyword: string;
}) {
  const [buttonStyle, setButtonStyle] = useState<object | undefined>(undefined);
  const colorSwitchTimeout = 500;

  return (
    <>
      <p>Ответ для формы:</p>
      <button
        className="message success"
        style={buttonStyle}
        onClick={() => {
          setButtonStyle({ backgroundColor: "#fff" });
          navigator.clipboard.writeText(discoveredQuestKeyword);
          setTimeout(() => setButtonStyle(undefined), colorSwitchTimeout);
        }}
      >
        <span>{discoveredQuestKeyword}</span>
        <CopyIcon />
      </button>
    </>
  );
}

const CopyIcon = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="18"
    height="18"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
  </svg>
);
