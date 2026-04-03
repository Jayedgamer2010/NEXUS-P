import { Link } from 'react-router-dom';
import './Error.css';

export default function Forbidden() {
  return (
    <div className="error-page">
      <div className="error-content">
        <h1 className="error-code">403</h1>
        <h2 className="error-title">Access Forbidden</h2>
        <p className="error-message">
          You don't have permission to access this page.
        </p>
        <Link to="/admin/dashboard" className="back-link">
          ← Back to Dashboard
        </Link>
      </div>
    </div>
  );
}
