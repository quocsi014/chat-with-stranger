const UserMessage = (props)=>{
  const {name, message, isYou} = props
  return(
    <div className={`bg-blue-500 rounded-xl px-4 py-2 text-lg text-white mb-2 max-w-96 w-fit ${isYou? "self-end":"self-start"}`}>
      <span className="text-2xl font-bold">{name}</span>
      <div className="">
        {message}
      </div>
    </div>
  )

}

export default UserMessage