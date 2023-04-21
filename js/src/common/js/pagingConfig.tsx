import { TablePaginationConfig } from "antd/es/table";
import React from 'react'

interface Props{
  count:number
  page:number
  pageSize:number
}
export default function pagingConfig(props:Props) {
  const {count,page,pageSize} = props;
  const pagingConfig: TablePaginationConfig = {
    defaultPageSize: 50, // 如果当pageSize为没有定义时，初始的表格单页展示数据大小为50
    total: count, // 总数据量
    defaultCurrent: 1,
    hideOnSinglePage: true, // 只有一页数据时是否隐藏分页器
    responsive: true, //
    pageSize: pageSize,
    current: page, // 当前的页码
    position: ["bottomCenter"], // 分页器显示位置
    showLessItems: true,
    pageSizeOptions: [20, 50, 100], // 用户可设置的单页显示条数
    showQuickJumper: true,
    showSizeChanger: true,
  }
  return pagingConfig;
}
