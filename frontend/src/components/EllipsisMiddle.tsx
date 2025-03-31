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

import type { CSSProperties } from "react";
import React from "react";
import { Typography } from "antd";
import "../common/styles/common.less";

const { Text } = Typography;

interface Props {
  className?: string;
  style?: CSSProperties;
  suffixCount: number;
  children: string;
}

/**
 * @description
 * If you want to omit the middle part of the text and display it with an ellipsis instead, you can try using the subcomponent
 * Both ends of the text are displayed in a symmetrical manner, and the end of the text is the same length as the beginning of the displayed characters
 * */
const EllipsisMiddle: React.FC<Props> = (
  /**
   * @param
   * className
   * suffixCount
   * children
   * */
  { className, style, suffixCount, children }
) => {
  // Renders the component only if it has data, otherwise it doesn't
  if (!children) {
    return null;
  }

  const start = children.slice(0, suffixCount).trim();
  const suffix = children.slice(-suffixCount).trim();
  const middle = children.length < 10 ? children : `${children.slice(suffixCount, -suffixCount).trim()}`;

  return (
    <Text className={`${className} ellipsis-sp`} style={{...style }}>
      {start}
      <span className="ellipsis">{middle}</span>
      {suffix}
    </Text>
  );
};

export default EllipsisMiddle;
