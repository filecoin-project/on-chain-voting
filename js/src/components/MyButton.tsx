import { Button } from "antd"
import React, { useEffect, useState } from "react"

export default function MyButton(props: any) {
  const { startCounting, handlerNavigate,change } = props
  const [buttonState, setButtonState] = useState(true);
  const [buttonText,setButtonText] = useState('Vote Counting');

  useEffect(()=>{
    if(!buttonState && !change){
      setButtonText('View');
    }
  },[change])

  function handleClick() {
    if (buttonState) {
      setButtonState(false)
    }
    buttonState ? startCounting() : handlerNavigate()
  }

  return (
    <Button
      className="menu_btn"
      type="primary"
      onClick={(pr) => {
        handleClick()
        
      }}
    >
      {buttonText}
    </Button>
  )
}
