// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { Collapse } from "antd";
import { useTranslation } from 'react-i18next';
import { useNavigate } from "react-router-dom";
import "../../common/styles/homepage.less";
const { Panel } = Collapse;

const Landing = () => {

    const { t } = useTranslation();

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
            answer: <div>
                <p>{t('content.answer3_1')}</p>
                <p>&nbsp;&nbsp;<span className="font-bold">{t('content.answer4_2')}: </span>{t('content.answer3_2')}</p>
                <p>&nbsp;&nbsp;<span className="font-bold">{t('content.answer4_3')}: </span>{t('content.answer3_3')}</p>
                <p>&nbsp;&nbsp;<span className="font-bold">{t('content.answer4_4')}: </span>{t('content.answer3_4')}</p>
                <p>&nbsp;&nbsp;<span className="font-bold">{t('content.answer4_5')}: </span>{t('content.answer3_5')}</p>
                <p>{t('content.answer3_6')}</p>
            </div>,
        },
        {
            title: t('content.question4'),
            answer: <div>
                <p>{t('content.answer4_1')}</p>
                <p>&nbsp;&nbsp;{t('content.answer4_2')}.</p>
                <p>&nbsp;&nbsp;{t('content.answer4_3')}.</p>
                <p>&nbsp;&nbsp;{t('content.answer4_4')}.</p>
                <p>&nbsp;&nbsp;{t('content.answer4_5')}.</p>
                <p>{t('content.answer4_6')}</p>
            </div>,
        },
        {
            title: t('content.question5'),
            answer: t('content.answer5'),
        },
        {
            title: t('content.question6'),
            answer: t('content.answer6'),
        },
        {
            title: t('content.question7'),
            answer: t('content.answer7'),
        },
        {
            title: t('content.question8'),
            answer: t('content.answer8'),
        },
        {
            title: t('content.question9'),
            answer: t('content.answer9'),
        },
        {
            title: t('content.question10'),
            answer: <div>
                <p>{t('content.answer10_1')}</p>
                <p>&nbsp;&nbsp;● <span>{t('content.answer10_2')}</span></p>
                <p>&nbsp;&nbsp;● <span>{t('content.answer10_3')}</span></p>
                <p>&nbsp;&nbsp;● <span>{t('content.answer10_4')}</span></p>
                <p>&nbsp;&nbsp;● <span>{t('content.answer10_5')}</span></p>
            </div>,
        },
        {
            title: t('content.question11'),
            answer: <span>{t('content.answer11')}&nbsp;<a className="hover:underline" href="https://github.com/filecoin-project/on-chain-voting" style={{ color: 'blue' }}>GitHub repository</a>.</span>,
        },
        {
            title: t('content.question12'),
            answer: t('content.answer12'),
        }
    ]
    const navigate = useNavigate();
    const goHome = async () => {
        navigate("/home");
    }

    return <div className="w-full justify-center space-y-16">
        <div className="space-y-10 mx-auto max-w-2xl">
            <div className='flex w-full items-center justify-center'>
                <div className='flex-shrink-0'>
                    <img height={"50px"} width={"50px"} src="/images/logo.png" alt="" />
                </div>
                <div className='ml-3 flex items-baseline space-x-20'>
                    <span className='text-black text-2xl font-semibold hover:opacity-80'>
                        Power Voting
                    </span>
                </div>
            </div>
            <h1 className="text-black font-bold text-3xl sm:text-4xl md:text-5xl lg:text-[54px] text-center leading-tight">
                {
                    t('content.headTitle')
                }
            </h1>
            <p className="text-[#445063] text-base sm:text-lg md:text-xl lg:text-[24px] text-center leading-relaxed">
                {
                    t('content.headContent')
                }
            </p>
            <div className="w-full flex items-center justify-center">
                <div className="px-6 py-3 cursor-pointer flex items-center justify-center text-center rounded border-solid border-[1px] border-[#DFDFDF] bg-white text-[#575757]">
                    {
                        t('content.headButtonLeft')
                    }
                </div>
                <div onClick={goHome} className="px-6 py-3 cursor-pointer flex items-center justify-center ml-5 text-center rounded bg-primary text-black">
                    {
                        t('content.headButtonRight')
                    }
                </div>
            </div>
        </div>

        <div className="rounded-3xl primary-shadow">
            <img
                className="rounded-3xl"
                src="/images/landing_1.png"
                alt=""
            />
        </div>




        <div className="text-[#000000] text-[40px] text-center leading-tight">
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
                        <span className="text-[16px] leading-relaxed"> {v.answer}</span>

                    </Panel>
                })
            }

        </Collapse>

    </div>
}

export default Landing
