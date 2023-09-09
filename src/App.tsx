import React from "react";
import {useRoutes, Link} from "react-router-dom";
import routes from "./router";
import ConnectWeb3Button from "./components/ConnectWeb3Button";
import Footer from './components/Footer';

import "./common/styles/reset.less";
import "tailwindcss/tailwind.css";
import "./common/styles/app.less";

const App: React.FC = () => {
  const element = useRoutes(routes);

  return (
    <div className="layout font-body">
      <header className='h-[96px]  bg-[#273141]'>
        <div className='w-[1000px] h-[96px] mx-auto flex items-center justify-between'>
          <div className='flex items-center'>
            <div className='flex-shrink-0'>
              <Link to='/'>
                <img className="logo" src="/images/logo1.png" alt=""/>
              </Link>
            </div>
            <div className='ml-6 flex items-baseline space-x-20'>
              <Link
                to='/'
                className='text-white text-2xl font-semibold hover:opacity-80'
              >
                Power Voting
              </Link>
            </div>
          </div>
          <div className='flex items-center'>
            <ConnectWeb3Button/>
          </div>
        </div>
      </header>
      <div className='content w-[1000px] mx-auto pt-10 pb-10'>
        {element}
      </div>
      <Footer/>
    </div>
  )
}

export default App
