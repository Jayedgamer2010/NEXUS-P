import { useState, KeyboardEvent } from 'react';
import './ConsoleInput.css';

interface ConsoleInputProps {
  onSend: (command: string) => void;
  disabled?: boolean;
  placeholder?: string;
}

export default function ConsoleInput({ onSend, disabled = false, placeholder = 'Enter command...' }: ConsoleInputProps) {
  const [input, setInput] = useState('');

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && input.trim() && !disabled) {
      onSend(input.trim());
      setInput('');
    }
  };

  return (
    <div className="console-input">
      <span className="prompt-symbol">❯</span>
      <input
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        disabled={disabled}
        className="console-input-field"
      />
      <button
        className="send-btn"
        onClick={() => {
          if (input.trim() && !disabled) {
            onSend(input.trim());
            setInput('');
          }
        }}
        disabled={disabled || !input.trim()}
      >
        Send
      </button>
    </div>
  );
}
