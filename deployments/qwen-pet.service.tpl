[Unit]
Description=Qwen PET Knowledge Agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart={{BINARY}}
Environment=PET_MODE=standalone
Environment=PET_CONFIG={{CONFIG}}
Environment=PET_DATA_DIR={{DATA_DIR}}
EnvironmentFile={{ENV_FILE}}
Restart=on-failure
RestartSec=5
KillSignal=SIGTERM
TimeoutStopSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=qwen-pet

[Install]
WantedBy=default.target
