import { useState } from "react";
import { useDispatch } from "react-redux";
import { updateName } from "../redux/userSlice";
import { useNavigate } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";

const HomePage = () => {
  const [name, setName] = useState("");
  const [code, setCode] = useState("");
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const handleNameChange = (e) => {
    setName(e.target.value);
    dispatch(updateName({ name: e.target.value }));
  };

  return (
    <div className="h-dvh bg-slate-300 flex items-center flex-col pt-28">
      <h2 className="mb-10 text-3xl font-extrabold">Chat with stranger</h2>
      <div className="w-96">
        <div className="flex h-8 mb-4">
          <label htmlFor="user_name" className="text-xl font-bold min-w-16">
            Name:{" "}
          </label>
          <input
            type="text"
            className="border-2 border-gray-300 rounded-lg h-8 w-full outline-none pl-1"
            value={name}
            onChange={(e) => handleNameChange(e)}
          />
        </div>
        <div className="flex h-8 mb-4">
          <label htmlFor="user_name" className="text-xl font-bold min-w-16">
            Code:{" "}
          </label>
          <input
            type="text"
            className="border-2 border-gray-300 rounded-lg h-8 w-full outline-none pl-1"
            value={code}
            onChange={(e) => {
              setCode(e.target.value);
            }}
          />
          {code.length == 0 ? (
            <button
              className="px-2 ml-2 bg-blue-300 text-white font-bold rounded-lg"
              disabled
            >
              join
            </button>
          ) : (
            <button
              className="px-2 ml-2 bg-blue-500 hover:bg-blue-600 text-white font-bold rounded-lg"
              onClick={() => {
                navigate(`/chat/${code}`);
              }}
            >
              join
            </button>
          )}
        </div>
        <button
          className="w-full py-2 rounded-lg bg-cyan-400 hover:bg-cyan-500 text-gray-700 font-bold"
          onClick={() => {
            let id = uuidv4();
            navigate(`/chat/${id}`);
          }}
        >
          New Room
        </button>
      </div>
    </div>
  );
};

export default HomePage;
