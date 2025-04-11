import { useState, useRef, useEffect } from "react";
import "./App.css";

// NOTE: USE ENVRIONMENT VARIABLES LATER
const WS_URL = "ws://localhost:4000/ws";

function App() {
  const [isConnected, setIsConnected] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState("Connecting...");
  const [messageToSend, setMessageToSend] = useState("");

  const clientId = useRef(
    `client_${Math.random().toString(36).substring(2, 9)}`
  );
  const [receivedMessages, setReceivedMessages] = useState([]);

  // To hold Websocket instance
  const ws = useRef(null);

  // To manage Websocket connecion lifecycle
  useEffect(() => {
    if (ws.current && ws.current.readyState != WebSocket.CLOSED) {
      console.log("Websocket connection already exists or closing");
      return;
    }

    // Create Websocket connection instance
    console.log(`Attempting to connect to ${WS_URL}...`);
    ws.current = new WebSocket(WS_URL);
    setConnectionStatus("Connecting...");
    setIsConnected(false); //Explicitly set false initially

    // --- Websocket Event Handlers ---
    // Look at MDN Websocket docs if confused

    // Called when the connection is succesfully opened
    ws.current.onopen = () => {
      console.log("Websocket connected");
      setConnectionStatus("Connected");
      setIsConnected(true);
    };

    // Called when a message is received from the server
    ws.current.onmessage = (event) => {
      console.log("Raw websocket Message received:", event.data);
      try {
        const message = JSON.parse(event.data);
        console.log("Parsed WS message: ", message);

        // Ignore message sent by self
        if (message.sender === clientId.current) {
          console.log("Ignoring self-sent message");
          return;
        }
      } catch (error) {
        console.error(
          "Failed to parse WS message or invalid JSON:",
          event.data,
          error
        );
        setReceivedMessages((prev) => [
          ...prev,
          { type: "error", payload: `Parse Error: ${event.data}` },
        ]);
      }
    };

    ws.current.onerror = (error) => {
      console.error("Websocket Error:", error);
      setConnectionStatus("Error");
      setIsConnected(false);
    };

    ws.current.onclose = (event) => {
      console.log("Websocket Disconnected", event);
      setConnectionStatus(
        `Disconnected: ${event.reason || "No reason given"} (Code: ${
          event.code
        })`
      );
      setIsConnected(false);

      // TODO: Implement reconnection logic here if desired
    };

    // Cleanup function
    return () => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        console.log("Closing Websocket connection");
        ws.current.close(1000, "Component unmounting");
      } else {
        console.log("Websocket already closed or closing when unmounting.");
      }
      ws.current = null;
    };
  }, []); // Should only run once on mount/unmount

  // Function to send message
  const sendMessage = (messageObject) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      try {
        // Ensure sender ID is always included
        const messageToSendWithSender = {
          ...messageObject,
          sender: clientId.current, // Add our client ID
        };
        const messageString = JSON.stringify(messageToSendWithSender);
        ws.current.send(messageString);
        console.log("Message Sent:", messageToSendWithSender); // Log the object
      } catch (error) {
        console.error("Failed to stringify or send message:", error);
      }
    } else {
      console.error("WebSocket is not connected. Cannot send message.");
      setConnectionStatus("Disconnected (Cannot Send)");
      setIsConnected(false);
    }
  };
  const handleSendButtonClick = () => {
    if (!messageToSend) {
      console.log("Cannot send empty message.");
      return;
    }
    // Send a structured message object
    sendMessage({
      type: "message", // Example type for simple chat-like messages
      payload: messageToSend,
      // Target and Room are omitted for simple broadcast
    });
    setMessageToSend(""); // Clear input after sending
  };

  return (
    <div className="App">
      <h1> sswan </h1>
      <h2>
        WebSocket Status:{" "}
        <span style={{ color: isConnected ? "green" : "red" }}>
          {connectionStatus}
        </span>
      </h2>

      <div style={{ margin: "20px 0" }}>
        <input
          type="text"
          value={messageToSend}
          onChange={(e) => setMessageToSend(e.target.value)}
          placeholder="Type a message to send..."
          disabled={!isConnected} // Disable input if not connected
          style={{ marginRight: "10px", padding: "5px" }}
        />
        <button
          onClick={handleSendButtonClick}
          disabled={!isConnected} // Disable button if not connected
          style={{ padding: "5px 10px" }}
        >
          Send Message
        </button>
      </div>
      <div
        style={{
          marginTop: "20px",
          border: "1px solid #ccc",
          padding: "10px",
          height: "200px",
          overflowY: "scroll",
        }}
      >
        <h3>Received Messages:</h3>
        {receivedMessages.length === 0 && <p>No messages received yet.</p>}
        <ul>
          {receivedMessages.map((msg, index) => (
            <li key={index}>
              <strong>Type:</strong> {msg.type},<strong>Sender:</strong>{" "}
              {msg.sender ? msg.sender.substring(0, 8) : "N/A"}...,{" "}
              {/* Show partial sender */}
              <strong>Payload:</strong>{" "}
              {typeof msg.payload === "object"
                ? JSON.stringify(msg.payload)
                : msg.payload}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;
