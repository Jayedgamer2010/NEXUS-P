import { useState } from 'react'

interface ConsoleInputProps {
  onSend: (command: string) => void
}

export default function ConsoleInput({ onSend }: ConsoleInputProps) {
  const [value, setValue] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (value.trim()) {
      onSend(value)
      setValue('')
    }
  }

  return (
    <form className="console-input-bar" onSubmit={handleSubmit}>
      <input
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder="Enter command..."
      />
      <button type="submit">Send</button>
    </form>
  )
}
