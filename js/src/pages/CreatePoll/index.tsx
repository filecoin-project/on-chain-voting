import React, { useState } from "react"
import {
  Form,
  Input,
  Button,
  Radio,
  DatePicker,
  InputNumber,
  Space,
  Modal,
  message,
} from "antd"
import { Buffer } from "buffer"
import { PlusCircleOutlined, PlusOutlined } from "@ant-design/icons"
// @ts-ignore
import nftStorage from "../../utils/storeNFT.js"
import { mainnetClient, timelockEncrypt, roundAt } from 'tlock-js';
import { RangePickerProps } from "antd/es/date-picker/index"
import dayjs from "dayjs"
import { usePowerVotingContract } from "../../hooks/use-power-voting-contract"
import { useNavigate } from "react-router-dom"
import create from "@ant-design/icons/lib/components/IconFont.js"

const { createVotingApi } = usePowerVotingContract()

interface Values {
  title: string
  description: string
  modifier: string
}
interface CollectionCreateFormProps {
  open: boolean
  onCreate: (values: Values) => void
  onCancel: () => void
  loading: boolean
}

const CollectionCreateForm: React.FC<CollectionCreateFormProps> = ({
  open,
  onCreate,
  onCancel,
  loading,
}) => {
  const [form] = Form.useForm()
  return (
    <Modal
      confirmLoading={loading}
      open={open}
      title="Add Option"
      okText="Create"
      cancelText="Cancel"
      onCancel={onCancel}
      onOk={() => {
        form
          .validateFields()
          .then((values) => {
            form.resetFields()
            onCreate(values)
          })
          .catch((info) => {
            console.log("Validate Failed:", info)
          })
      }}
    >
      <Form
        form={form}
        layout="vertical"
        name="form_in_modal"
        initialValues={{ modifier: "public" }}
      >
        <Form.Item
          name="Name"
          label="Option Name"
          rules={[{ required: true, message: "Option Name!" }]}
        >
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  )
}
const { RangePicker } = DatePicker
const { TextArea } = Input

const CreatePoll = () => {
  
  const navigate = useNavigate()
  const onFinish = async (values: any) => {

    console.log("Success:", values)
    if (radio.length <= 0) {
      message.error("Please confirm if you want to add a voting option")
    } else {
      setLoading(true)
      // 调取接口发送数据
      const timestamp = values.Time.valueOf()
      // console.log(values.Time.valueOf());
      // const info = await mainnetClient().chain().info();
      // const roundNumber = roundAt(timestamp, info)

      const _values = { ...values, Time: timestamp, option: radio }

      const cid = await nftStorage(_values)

      if (cid) {
        setLoading(false)
        message.success("Waiting for the transaction to be chained!")
        // console.log(cid)
        if (createVotingApi) {
          const res = await createVotingApi(cid)
          // console.log(res, "res")
          if (res) {
            message.success("Preparing to wind the chain!")
            navigate("/", { state: true })
          }
        }
      }
    }
  }

  const onFinishFailed = (errorInfo: any) => {
    console.log("Failed:", errorInfo)
  }

  // tan

  const [open, setOpen] = useState(false)
  const [radio, setRadio] = useState([] as any[string])
  const [loading, setLoading] = useState<boolean>(false)
  const onCreate = (values: any) => {
    const radioAdd = radio
    radioAdd.push(values.Name)
    setRadio(radioAdd)
    setOpen(false)
  }
  const range = (start: number, end: number) => {
    const result = []
    for (let i = start; i < end; i++) {
      result.push(i)
    }
    return result
  }

  const disabledDate: RangePickerProps["disabledDate"] = (current) => {
    // Can not select days before today and today
    return current && current < dayjs()
  }

  const disabledDateTime = () => ({
    disabledHours: () => range(0, 24).splice(4, 20),
    disabledMinutes: () => range(30, 60),
    disabledSeconds: () => [55, 56],
  })

  return (
    <div>
      <Form
        layout="vertical"
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        labelCol={{ span: 12 }}
        wrapperCol={{ span: 24 }}
      >
        <Form.Item
          name="Name"
          label="Name"
          rules={[{ required: true, message: "Please enter the name!" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="Description"
          label="Description"
          rules={[{ required: true, message: "Please enter the Description!" }]}
        >
          <TextArea rows={4} />
        </Form.Item>

        <div
          style={{ cursor: "pointer" }}
          onClick={() => {
            setOpen(true)
          }}
        >
          Add Option <PlusCircleOutlined />
        </div>
        <Form.Item label="options:">
          <Radio.Group>
            <Space direction="vertical">
              {radio.map((item: any, index: any) => {
                return <div key={index}>{" " + item}</div>
              })}
            </Space>
          </Radio.Group>
        </Form.Item>
        <Form.Item
          name={"Time"}
          label="Closing Time"
          rules={[
            {
              required: true,
              message: "Please enter your Number of Closing Time!",
            },
          ]}
        >
          <DatePicker
            showTime={{ format: 'HH:mm' }}
            format="YYYY-MM-DD HH:mm"
          // disabledDate={disabledDate}

          />
        </Form.Item>

        <Form.Item>
          <Button
            loading={loading}
            htmlType="submit"
            style={{ background: "#e99d42", width: "100%" }}
          >
            Submit
          </Button>
        </Form.Item>
      </Form>

      <>
        <CollectionCreateForm
          open={open}
          onCreate={onCreate}
          onCancel={() => {
            setOpen(false)
          }}
          loading={loading}
        />
      </>
    </div>
  )
}

export default CreatePoll
