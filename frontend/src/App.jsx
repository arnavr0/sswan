import { useState, useRef, useEffect } from "react";
import "./App.css";

// NOTE: USE ENVRIONMENT VARIABLES LATER
const WS_URL = "ws://localhost:4000/ws";

function App() {
  const [isConnected, setIsConnected] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState("Connecting...");
  const [messageToSend, setMessageToSend] = useState("");

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
      console.log("Websocket Message received:", event.data);
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
  const handleSendMessage = () => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      if (!messageToSend) {
        console.log("Cannot send empty message.");
        return; // Don't send empty messages
      }
      try {
        const basicMessage = {
          type: "echo", // Example type
          payload: messageToSend,
          sender: "react-client-temp-id", // Example sender
        };
        const messageString = JSON.stringify(basicMessage);

        ws.current.send(messageString); // Send the message
        console.log("Message Sent:", messageString);
        // Optional: Clear the input field after sending
        // setMessageToSend('');
      } catch (error) {
        console.error("Failed to send message:", error);
      }
    } else {
      console.error("WebSocket is not connected. Cannot send message.");
      setConnectionStatus("Disconnected (Cannot Send)");
      setIsConnected(false);
    }
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
          onClick={handleSendMessage}
          disabled={!isConnected} // Disable button if not connected
          style={{ padding: "5px 10px" }}
        >
          Send Message
        </button>
      </div>
    </div>
  );
}

export default App;
