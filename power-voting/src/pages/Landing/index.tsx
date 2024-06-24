import { Collapse } from "antd";
import React from "react";
import { useNavigate } from "react-router-dom";
const { Panel } = Collapse;


const Landing = () => {

    const DESC = [
        {
            icon: "/images/landing_ic_1.png",
            title: "Lorem ipsum dolor sit amet1",
            desc: "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut."
        },
        {
            icon: "/images/landing_ic_2.png",

            title: "Lorem ipsum dolor sit amet2",
            desc: "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut."
        },
        {
            icon: "/images/landing_ic_3.png",
            title: "Lorem ipsum dolor sit amet3",
            desc: "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut."
        }

    ]
    const QUESTIONS = [
        {
            title: "Question 1",
            answer: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi utLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut.",
        },
        {
            title: "Question 2",
            answer: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi utLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut.",
        },
        {
            title: "Question 3",
            answer: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi utLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut.",
        },
        {
            title: "Question 4",
            answer: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi utLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut.",
        },
        {
            title: "Question 5",
            answer: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi utLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut.",
        }
    ]
    const navigate = useNavigate();
    const goHome = async () => {
        navigate("/home");
    }

    return <div className="w-full justify-center">
        <div className='flex w-full items-center justify-center'>
            <div className='flex-shrink-0'>
                <img height={"50px"} width={"50px"} src="/images/logo.png" alt="" />
            </div>
            <div className='ml-3 flex items-baseline space-x-20'>
                <span
                    className='text-black text-2xl font-semibold hover:opacity-80'
                >
                    Power Voting
                </span>
            </div>
        </div>
        <div className="mt-5 text-black font-bold text-[54px] text-center">
            Unlock the Power of Decentralized Decision-Making
        </div>
        <div className="mt-10 text-[#445063] text-[24px] px-[60px] text-center">
            Power Voting empowers you to participate in governance and make impactful decisions within the Filecoin ecosystem.
        </div>
        <div className="mt-5 w-full flex items-center justify-center">
            <div className="cursor-pointer flex items-center justify-center text-center rounded w-[128px] h-[31px] border-solid border-[1px] border-[#DFDFDF] bg-white text-[#575757]">
                Learn More
            </div>
            <div onClick={goHome} className="cursor-pointer flex items-center justify-center ml-5 text-center rounded w-[128px] h-[31px] bg-[#0190FF] text-[#ffffff]">
                Get Started
            </div>
        </div>

        <img className="mt-20" width={"100%"} src="/images/landing_1.png" alt="" />


        <div className="mt-40 text-[#005292] text-[20px] px-[60px] text-center mb-[10px]">
            Lorem ipsum dolor sit amet
        </div>
        <div className="mt-5 text-[#000000] text-[40px] px-[60px] text-center mb-[10px]">
            Ut enim ad minim veniam quis nostrud
        </div>
        <div className="mt-10 text-[#445063] text-[24px] px-[60px] text-center mb-[10px]">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut.
        </div>
        <div className="flex mt-[50px]">
            {
                DESC.map((v, i) => {
                    return <div key={i}>
                        <img height={"50px"} width={"50px"} src={v.icon} alt="" />
                        <div className="mt-[10px] text-[#000000] text-[20px]">
                            {v.title}
                        </div>
                        <div className="mt-[10px] text-[#445063] text-[20px]">
                            {v.desc}
                        </div>
                    </div>
                })

            }
        </div>
        <div className="mt-40 text-[#005292] text-[20px] px-[60px] text-center mb-[10px]">
            Lorem ipsum dolor sit amet
        </div>
        <div className="mt-10 text-[#000000] text-[40px] px-[60px] text-center mb-[10px]">
            Frequently Asked Questions
        </div>

        <Collapse
            expandIconPosition={"end"}
            bordered={false}

            style={{
                background: "#F9F9F9",
            }}>
            {
                QUESTIONS.map((v, i) => {
                    return <Panel
                        className="text-[20px]"
                        header={v.title} key={i}>
                        <span className="text-[16px]"> {v.answer}</span>

                    </Panel>
                })
            }

        </Collapse>

        <div className="mt-40 text-[#000000] text-[40px] px-[60px] text-center mb-[10px]">
            Ut enim ad minim veniam quis nostrud
        </div>
        <div className="mt-10 text-[#445063] text-[24px] px-[60px] text-center mb-[10px]">
            Power Voting empowers you to participate in governance and make impactful decisions within the Filecoin ecosystem.
        </div>

        <div className="mt-10 w-full flex items-center justify-center mb-[50px] ">
            <div onClick={goHome} className="cursor-pointer flex items-center justify-center ml-[5px] text-center rounded w-[128px] h-[31px] bg-[#0190FF] text-[#ffffff]">
                Get Started
            </div>
        </div>
    </div>
}

export default Landing