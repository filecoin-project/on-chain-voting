import React from "react"
import { ConnectButton } from "@rainbow-me/rainbowkit"
import { Button, Col, Row } from "antd"

export const ConnectWeb3Button = (props: any) => {
  return (
    <div
      style={{
        // width: isConnected ? "31rem" : "20rem",
        color: "#fff",
        padding: "8px 8px",
        textAlign: "center",
        borderRadius: "8px",
      }}
      className="connect flex items-center sm:mx-4 "
    >
      <ConnectButton.Custom>
        {({
          account,
          chain,
          openAccountModal,
          openChainModal,
          openConnectModal,
          authenticationStatus,
          mounted,
        }) => {
          // Note: If your app doesn't use authentication, you
          // can remove all 'authenticationStatus' checks
          const ready = mounted && authenticationStatus !== "loading"
          const connected =
            ready &&
            account &&
            chain &&
            (!authenticationStatus || authenticationStatus === "authenticated")
          // handler(connected);
          localStorage.setItem('isConnect',JSON.stringify(connected));
          return (
            <div
              {...(!ready && {
                "aria-hidden": true,
                style: {
                  opacity: 0,
                  pointerEvents: "none",
                  userSelect: "none",
                },
              })}
            >
              {(() => {
                if (!connected) {
                  return (
                    <div>
                      <Button
                        className="menu_btn"
                        onClick={openConnectModal}
                        type="primary"
                      >
                        Connect Wallet
                      </Button>
                    </div>
                  )
                }
                if (chain.unsupported) {
                  return (
                    <div>
                      <button
                        className="connect-pc connect connect-warn"
                        onClick={openChainModal}
                        type="button"
                      >
                        Wrong network
                      </button>
                      <img
                        className="connect-mob"
                        onClick={openChainModal}
                        src="/assets/wallet.svg"
                        width={20}
                      />
                    </div>
                  )
                }
                return (
                  <a>
                    <div
                      onClick={openAccountModal}
                      style={{
                        fontSize: "16px",
                        display: "flex",
                        gap: 12,
                        color: "orange",
                        fontWeight: "bold",
                      }}
                    >
                      <Row gutter={[8, 8]}>
                        <Col>
                          <i className="iconfont icon-qianbao"></i>
                        </Col>
                        <Col>
                          {" "}
                          <span>
                            {account.displayName}
                            {account.displayBalance
                              ? ` (${account.displayBalance})`
                              : ""}
                          </span>
                        </Col>
                      </Row>
                    </div>
                  </a>
                )
              })()}
            </div>
          )
        }}
      </ConnectButton.Custom>
    </div>
  )
}
