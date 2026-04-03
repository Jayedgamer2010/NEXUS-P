import { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import StatCard from '../../components/ui/StatCard';
import ConsoleInput from '../../components/console/ConsoleInput';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { useConsole } from '../../hooks/useConsole';
import { adminApi } from '../../api/admin';
import { useServerStats } from '../../hooks/useServerStats';
import { Server, PowerAction } from '../../types';
import './ServerDetail.css';

export default function ServerDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [server, setServer] = useState<Server | null>(null);
  const [loading, setLoading] = useState(true);
  const [sendingAction, setSendingAction] = useState<string | null>(null);
  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermRef = useRef<Terminal | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);

  const { stats, loading: statsLoading } = useServerStats({
    serverUUID: id || '',
  });

  const { connectionStatus, send: wsSend, disconnect } = useConsole({
    serverUUID: id || '',
    onMessage: (message) => {
      if (xtermRef.current) {
        xtermRef.current.writeln(message);
      }
    },
  });

  useEffect(() => {
    loadServer();
  }, [id]);

  useEffect(() => {
    if (!terminalRef.current) return;

    const term = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'JetBrains Mono, Fira Code, monospace',
      theme: {
        background: '#000000',
        foreground: '#10b981',
        cursor: '#10b981',
        selection: '#7c3aed',
      },
      scrollback: 1000,
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    xtermRef.current = term;
    fitAddonRef.current = fitAddon;

    term.writeln('\x1b[36m%s\x1b[0m', `[NEXUS] Console ready`);

    return () => {
      term.dispose();
    };
  }, []);

  useEffect(() => {
    const handleResize = () => {
      if (fitAddonRef.current) {
        fitAddonRef.current.fit();
      }
    };
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  useEffect(() => {
    return () => {
      disconnect();
    };
  }, [disconnect]);

  const loadServer = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const response = await adminApi.getServers(1, 100);
      const found = response.data.find(s => s.uuid === id);
      if (found) {
        setServer(found);
      } else {
        navigate('/admin/servers');
      }
    } catch (error) {
      console.error('Failed to load server:', error);
      navigate('/admin/servers');
    } finally {
      setLoading(false);
    }
  };

  const handlePowerAction = async (action: PowerAction) => {
    if (!server) return;
    setSendingAction(action);
    try {
      await adminApi.powerServer(server.id, action);
      setTimeout(loadServer, 2000);
    } catch (error) {
      console.error('Failed to send power action:', error);
      alert('Failed to send power action');
    } finally {
      setSendingAction(null);
    }
  };

  const handleConsoleSend = (command: string) => {
    wsSend(command);
    if (xtermRef.current) {
      xtermRef.current.writeln(`$ ${command}`);
    }
  };

  const formatBytes = (bytes: number): string => {
    const gb = bytes / (1024 * 1024 * 1024);
    return gb.toFixed(2) + ' GB';
  };

  if (loading) return <div className="loading">Loading server...</div>;
  if (!server) return null;

  const cpuPercent = stats ? (stats.cpu * 100).toFixed(1) : '0.0';
  const ramUsed = stats ? formatBytes(stats.memory) : '0 GB';
  const ramTotal = stats ? formatBytes(stats.memory_max) : formatBytes(server.memory * 1024 * 1024);
  const diskUsed = stats ? formatBytes(stats.disk) : '0 GB';
  const diskTotal = stats ? formatBytes(stats.disk_max) : formatBytes(server.disk * 1024 * 1024);

  const getStatusColor = () => {
    switch (connectionStatus) {
      case 'connected':
        return '#10b981';
      case 'connecting':
        return '#f59e0b';
      default:
        return '#ef4444';
    }
  };

  const getStatusText = () => {
    switch (connectionStatus) {
      case 'connected':
        return 'Connected';
      case 'connecting':
        return 'Connecting...';
      default:
        return 'Disconnected';
    }
  };

  return (
    <div className="server-detail">
      <div className="server-header">
        <div className="server-info">
          <h1>{server.name}</h1>
          <div className="server-meta">
            <span className="uuid">{server.uuid.substring(0, 12)}...</span>
            <span className={`status-badge status-${server.status}`}>{server.status}</span>
          </div>
        </div>
        <div className="server-actions">
          <button
            className="action-btn start"
            onClick={() => handlePowerAction(PowerAction.START)}
            disabled={sendingAction !== null || server.status === 'running'}
          >
            ▶ Start
          </button>
          <button
            className="action-btn restart"
            onClick={() => handlePowerAction(PowerAction.RESTART)}
            disabled={sendingAction !== null || server.status !== 'running'}
          >
            ↻ Restart
          </button>
          <button
            className="action-btn stop"
            onClick={() => handlePowerAction(PowerAction.STOP)}
            disabled={sendingAction !== null || server.status !== 'running'}
          >
            ◼ Stop
          </button>
          <button
            className="action-btn kill"
            onClick={() => handlePowerAction(PowerAction.KILL)}
            disabled={sendingAction !== null}
          >
            ✕ Kill
          </button>
        </div>
      </div>

      <div className="stats-section">
        <h2>Live Statistics</h2>
        <div className="stats-grid">
          <StatCard
            title="CPU Usage"
            value={`${cpuPercent}%`}
            icon="⚡"
            accentColor="#7c3aed"
            subtitle={statsLoading ? 'Fetching...' : 'Live'}
          />
          <StatCard
            title="RAM Usage"
            value={ramUsed}
            subtitle={`of ${ramTotal}`}
            icon="💾"
            accentColor="#3b82f6"
          />
          <StatCard
            title="Disk Usage"
            value={diskUsed}
            subtitle={`of ${diskTotal}`}
            icon="💿"
            accentColor="#10b981"
          />
        </div>
      </div>

      <div className="console-section">
        <h2>Console</h2>
        <div className="console-wrapper">
          <div className="console-container">
            <div className="console-header">
              <div className="console-status">
                <span
                  className="status-indicator"
                  style={{ backgroundColor: getStatusColor() }}
                />
                {getStatusText()}
              </div>
              <div className="console-title">Server Console</div>
            </div>
            <div ref={terminalRef} className="console-terminal" />
          </div>
          <ConsoleInput onSend={handleConsoleSend} />
        </div>
      </div>
    </div>
  );
}
