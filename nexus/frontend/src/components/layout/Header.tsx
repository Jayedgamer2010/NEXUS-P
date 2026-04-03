import './Header.css';

interface HeaderProps {
  title: string;
}

export default function Header({ title }: HeaderProps) {
  return (
    <header className="header">
      <div className="header-left">
        <h1 className="header-title">{title}</h1>
      </div>
      <div className="header-right">
        {/* Future: notifications, search, etc */}
        <div className="header-spacer" />
      </div>
    </header>
  );
}
