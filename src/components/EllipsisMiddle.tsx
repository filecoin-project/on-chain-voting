import React, { CSSProperties } from "react";
import { Typography } from "antd";
import "../common/styles/common.less";

const { Text } = Typography;

/**
 * @description
 * If you want to omit the middle part of the text and display it with an ellipsis instead, you can try using the subcomponent
 * Both ends of the text are displayed in a symmetrical manner, and the end of the text is the same length as the beginning of the displayed characters
 * */
const EllipsisMiddle: React.FC<{ className?: string | undefined; style?: CSSProperties | undefined; suffixCount: number; children: string }> = (
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
  // The data length is greater than 10 before processing, otherwise it is returned directly
  if (children.length < 10) {
    return <span>{children}</span>;
  }
  const start = children.slice(0, suffixCount).trim();
  const suffix = children.slice(-suffixCount).trim();

  return (
    <Text className={`${className || ""} ellipsis-sp`} style={{...style }}>
      {start}
      <span className={"ellipsis"}>{children.slice(suffixCount, -suffixCount).trim()}</span>
      {suffix}
    </Text>
  );
};

export default EllipsisMiddle;
