import { Collapse } from "antd";
import React from "react";
import { useNavigate } from "react-router-dom";
const { Panel } = Collapse;
import { useTranslation } from 'react-i18next';

const Landing = () => {

    const { t } = useTranslation();

    const DESC = [
        {
            icon: "/images/landing_ic_1.png",
            title: t('content.section1Title'),
            desc: t('content.section1Content'),
        },
        {
            icon: "/images/landing_ic_2.png",
            title: t('content.section2Title'),
            desc: t('content.section2Content'),
        },
        {
            icon: "/images/landing_ic_3.png",
            title: t('content.section3Title'),
            desc: t('content.section3Content'),
        }

    ]
    const QUESTIONS = [
        {
            title: t('content.question1'),
            answer: t('content.answer1'),
        },
        {
            title: t('content.question2'),
            answer: t('content.answer2'),
        },
        {
            title: t('content.question3'),
            answer: t('content.answer3'),
        },
        {
            title: t('content.question4'),
            answer: t('content.answer4'),
        },
        {
            title: t('content.question5'),
            answer: t('content.answer5'),
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
            {
                t('content.headTitle')
            }
        </div>
        <div className="mt-10 text-[#445063] text-[24px] px-[60px] text-center">
            {
                t('content.headContent')
            }
        </div>
        <div className="mt-5 w-full flex items-center justify-center">
            <div className="cursor-pointer flex items-center justify-center text-center rounded w-[128px] h-[31px] border-solid border-[1px] border-[#DFDFDF] bg-white text-[#575757]">
                {
                    t('content.headButtonLeft')
                }
            </div>
            <div onClick={goHome} className="cursor-pointer flex items-center justify-center ml-5 text-center rounded w-[128px] h-[31px] bg-[#0190FF] text-[#ffffff]">
                {
                    t('content.headButtonRight')
                }
            </div>
        </div>

        <img className="mt-20" width={"100%"} src="/images/landing_1.png" alt="" />


        <div className="mt-40 text-[#005292] text-[20px] px-[60px] text-center mb-[10px]">
            {
                t('content.topTitle')
            }
        </div>
        <div className="mt-5 text-[#000000] text-[40px] px-[60px] text-center mb-[10px]">
            {
                t('content.topHead')
            }
        </div>
        <div className="mt-10 text-[#445063] text-[24px] px-[60px] text-center mb-[10px]">
            {
                t('content.topContent')
            }
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
            {
                t('content.questionTitle')
            }
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
            {
                t('content.bottomTitle')
            }
        </div>
        <div className="mt-10 text-[#445063] text-[24px] px-[60px] text-center mb-[10px]">
            {
                t('content.bottomHead')
            }
        </div>

        <div className="mt-10 w-full flex items-center justify-center mb-[50px] ">
            <div onClick={goHome} className="cursor-pointer flex items-center justify-center ml-[5px] text-center rounded w-[128px] h-[31px] bg-[#0190FF] text-[#ffffff]">
                {
                    t('content.bottomButton')
                }
            </div>
        </div>
    </div>
}

export default Landing