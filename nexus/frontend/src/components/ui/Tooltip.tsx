import { useState } from 'react';
import './Tooltip.css';

interface TooltipProps {
  content: string;
  children: React.ReactNode;
  position?: 'top' | 'bottom' | 'left' | 'right';
}

export default function Tooltip({
  content,
  children,
  position = 'top',
}: TooltipProps) {
  const [visible, setVisible] = useState(false);

  return (
    <span
      className={`tooltip-wrapper tooltip-wrapper--${position}`}
      onMouseEnter={() => setVisible(true)}
      onMouseLeave={() => setVisible(false)}
    >
      {children}
      {visible && <span className="tooltip-bubble">{content}</span>}
    </span>
  );
}
