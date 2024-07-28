import { useEffect, useRef, useState } from "react";
import { useSelector } from "react-redux";
import SystemMessage from "../components/SystemMessage";
import UserMessage from "../components/UserMessage";

const ChatPage = () => {
  const [msg, setMsg] = useState("");
  const [msgList, setMsgList] = useState([]);
  const socketRef = useRef(null);

  
  const name = useSelector((state) => state.user.name);


  useEffect(() => {
    const socket = new WebSocket(
      `ws://localhost:8080/ws?name=${name}`
    );
    socketRef.current = socket;
    socket.onopen = () => {
      console.log("WebSocket connection opened");
    };

    socket.onmessage = (event) => {
      console.log("Received from server:", JSON.parse(event.data));
      setMsgList((prevMsgList) => [
        ...prevMsgList,
        {...JSON.parse(event.data), is_you: false},
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
  }, [name]);

  const sendMsg = (msg) => {
    socketRef.current.send(msg);
    setMsgList((prevMsgList) => [
      ...prevMsgList,
      { message: msg, user_name: "You", is_system: false, is_you:true },
    ]);
    setMsg("");
  };

  console.log(msgList)
  return (
    <div className="h-dvh bg-gray-100 flex flex-col-reverse items-center">
      <div className="w-full h-full max-w-200 p-2 flex flex-col bg-white">
        <div className="size-full border-b-2 border-gray-400 mb-2 flex flex-col overflow-y-scroll scrollbar-hide">
          {msgList.map((msg) => {
            if(msg.is_system){
              return <SystemMessage message={msg.message}/>
            }else{
              if (msg.you){

                return <UserMessage isYou={msg.is_you} name={msg.user_name} message={msg.message}/>
              }else{
                return <UserMessage isYou={msg.is_you} name={msg.user_name} message={msg.message} />
              }
            }
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
