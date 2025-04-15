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

import { LoadingOutlined } from '@ant-design/icons';
import { Spin, Space } from 'antd';

const LoadingButton  = (props: any) => {
  const { loading, handleClick, text, isFull, className } = props;
  return (
    <button onClick={handleClick} className={`h-[40px] bg-sky-500 hover:bg-sky-700 text-white py-2 px-6 rounded-xl ${className} ${loading ? 'cursor-not-allowed' : ''} ${isFull ? 'w-full' : ''}`} type="submit" disabled={loading}>
      <Space size='middle'>
        {
          text
        }
        {
          loading && <Spin indicator={<LoadingOutlined style={{ fontSize: 16, color: '#ffffff' }} spin />} />
        }
      </Space>
    </button>
  )
}

export default LoadingButton;