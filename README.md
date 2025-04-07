sswan/
├── backend/                  # Go Backend Code
│   ├── cmd/
│   │   └── server/           # Main application entry point directory
│   │       └── main.go       # Main function, HTTP server setup
│   ├── internal/             # Private application code (optional but good practice)
│   │   └── websocket/        # WebSocket handling logic (or just 'ws')
│   │       ├── client.go     # Represents a connected client
│   │       ├── hub.go        # Manages clients, rooms, broadcasting
│   │       ├── handler.go    # The HTTP handler that upgrades to WebSocket
│   │       └── message.go    # Defines signaling message structures
│   ├── pkg/                  # Shared libraries (if any, less common for internal)
│   ├── go.mod                # Go module definition
│   ├── go.sum                # Go module checksums
│   └── Makefile              # (Optional) Build/run commands
│
├── frontend/                 # React Frontend Code (e.g., created with Create React App or Vite)
│   ├── public/               # Static assets (index.html, favicon, etc.)
│   │   └── index.html
│   ├── src/                  # React source code
│   │   ├── components/       # Reusable UI components (e.g., VideoPlayer, Controls)
│   │   │   └── VideoPlayer.js
│   │   │   └── Controls.js
│   │   ├── hooks/            # Custom React hooks (e.g., for WebSocket, WebRTC logic)
│   │   │   └── useWebSocket.js
│   │   │   └── useWebRTC.js  # Encapsulates WebRTC setup/signaling
│   │   ├── services/         # Interacting with backend/APIs (WebSocket connection setup)
│   │   │   └── websocketService.js
│   │   ├── contexts/         # (Optional) React Context for global state if needed
│   │   │   └── AppContext.js
│   │   ├── pages/            # (Optional) Top-level page components (if using routing)
│   │   │   └── ShareScreenPage.js
│   │   │   └── WatchScreenPage.js
│   │   ├── App.js            # Main application component
│   │   ├── index.js          # React entry point
│   │   └── App.css           # Main styles (or other styling approach)
│   ├── package.json          # Frontend dependencies and scripts
│   ├── package-lock.json     # (or yarn.lock) Lockfile for dependencies
│   └── .env.development      # (Optional) Environment variables for development (e.g., WS URL)
│   └── .env.production       # (Optional) Environment variables for production
│
├── .gitignore                # Git ignore rules for both backend and frontend
└── README.md                 # Project overview, setup instructions, etc.

