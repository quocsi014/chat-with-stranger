const SystemMessage = (props)=>{
  const {message} = props
  return(
    <div className="w-full text-center text-md text-gray-300">
      {message}
    </div>
  )
}

export default SystemMessage