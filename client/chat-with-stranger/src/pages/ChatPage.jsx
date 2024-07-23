import { useEffect, useRef, useState } from "react";
import { useSelector } from "react-redux";
import { useParams } from "react-router-dom";

const ChatPage = () => {
  const [msg, setMsg] = useState("");
  const [msgList, setMsgList] = useState([]);
  const socketRef = useRef(null);

  let key = 1;
  const name = useSelector((state) => state.user.name);

  const { code } = useParams();

  useEffect(() => {
    const socket = new WebSocket(
      `ws://localhost:8080/ws?name=${name}&room_key=${code}`
    );
    socketRef.current = socket;
    socket.onopen = () => {
      console.log("WebSocket connection opened");
    };

    socket.onmessage = (event) => {
      console.log("Received from server:", event.data);
      setMsgList((prevMsgList) => [
        ...prevMsgList,
        { msg: event.data, isYou: false },
      ]);
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };

    socket.onerror = (error) => {
      console.log("WebSocket error:", error);
    };

    return () => {
      socket.close();
    };
  }, [code, name]);

  const sendMsg = (msg) => {
    socketRef.current.send(msg);
    setMsgList((prevMsgList) => [
      ...prevMsgList,
      { msg: "You: " + msg, isYou: true },
    ]);
    setMsg("");
  };

  return (
    <div className="h-dvh bg-slate-300 flex flex-col-reverse items-center">
      <div className="w-full h-full max-w-200 p-2 flex flex-col">
        <div className="text-3xl border-b-4 border-white mb-4">
          Code: {code}
        </div>
        <div className="size-full border-b-2 border-gray-400 mb-2 flex flex-col overflow-y-scroll scrollbar-hide">
          {msgList.map((msg) => {
            return (
              <div
                key={++key}
                className={`${
                  msg.isYou ? "self-end pl-4" : "pr-4"
                } border-b-2 w-fit mb-2`}
              >
                {msg.msg}
              </div>
            );
          })}
        </div>
        <div className="flex w-full ">
          <input
            type="text"
            className="w-full rounded-lg mr-1 px-2 outline-none text-lg"
            value={msg}
            onChange={(e) => {
              setMsg(e.target.value);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                sendMsg(msg)
              }
            }}
          />
          {msg.length == 0 ? (
            <button
              className="p-2 bg-blue-300 text-white font-bold rounded-lg"
              disabled
            >
              Send
            </button>
          ) : (
            <button
              className="p-2 bg-blue-500 text-white font-bold rounded-lg"
              onClick={() => {
                sendMsg(msg);
              }}
            >
              Send
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default ChatPage;
