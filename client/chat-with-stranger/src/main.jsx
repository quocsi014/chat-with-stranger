import ReactDOM from 'react-dom/client'
import {createBrowserRouter, RouterProvider} from "react-router-dom"
import "./index.css"
import HomePage from './pages/HomePage.jsx'
import ChatPage from './pages/ChatPage.jsx'
import { Provider } from 'react-redux'
import { store } from './redux/store.js'

const router = createBrowserRouter([
  {
    path:"/",
    element:<HomePage/>
  },
  {
    path:"/chat",
    element:<ChatPage/>
  }
])

ReactDOM.createRoot(document.getElementById('root')).render(
  <Provider store={store}>
    <RouterProvider router={router}/>
  </Provider>,
)
