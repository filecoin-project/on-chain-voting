import { Table } from "antd"
import React, { useEffect, useState } from "react"
import { ColumnsType, TablePaginationConfig } from "antd/es/table"
import { TableProps as RcTableProps } from "rc-table/lib/Table"
import { GetRowKey } from "antd/lib/table/interface"

/**
 * @description
 * header: 表格标题
 * rowKey: 表格主键，如果需要增加过滤功能，或者排序等，会根据传入的rowKey进行规则判断；rowKey通常是ColumnsType表格配置的每列的key属性
 * tableDataTypeConfig：
 * @see https://ant.design/components/table-cn/#Column
 * dataSource: 需要渲染的数据，渲染数据需与ColumnsType配置中定于的列属性规则相对应
 *
 * */

interface BlockListProps {
    className?: string
    rowKey: string | GetRowKey<any>
    dataConf: { page: number; pageSize: number; count: number; data: any }
    tableDataTypeConfig: () => any
    loading?: boolean
    onChange: (page?: number, pageSize?: number) => void
    scroll?: RcTableProps["scroll"] & {
        scrollToFirstRowOnChange?: boolean
    }
}

// @ts-ignore
const Tabulation = (props: BlockListProps) => {
    const {
        className,
        rowKey,
        tableDataTypeConfig,
        onChange,
        dataConf,
        loading,
        scroll,
    } = props
    // console.log(dataConf.data,dataConf.count);

    /**
     * @constructor
     * 设置表格初始渲染的属性值
     * */

    /**
     * 表格分页器属性配置
     * */

    const pagingConfig: TablePaginationConfig = {
        defaultPageSize: 50, // 如果当pageSize为没有定义时，初始的表格单页展示数据大小为50
        total: dataConf.count, // 总数据量
        defaultCurrent: 1,
        hideOnSinglePage: true, // 只有一页数据时是否隐藏分页器
        responsive: true, //
        pageSize: dataConf.pageSize,
        current: dataConf.page, // 当前的页码
        position: ["bottomCenter"], // 分页器显示位置
        showLessItems: true,
        pageSizeOptions: [20, 50, 100], // 用户可设置的单页显示条数
        showQuickJumper: true,
        showSizeChanger: true,
    }
    return (
        <Table
            rowClassName={className}
            className={className}
            loading={loading}
            style={{ borderRadius: "8px",textAlign:"center",minHeight:"100vh" }}
            scroll={scroll}
            rowKey={rowKey}
            size={"small"}
            columns={tableDataTypeConfig()}
            dataSource={dataConf.data}
            pagination={pagingConfig}
            onChange={(pagination) => {
                const page = pagination.current || 1
                onChange(pagination.pageSize, page)
            }}
        />
    )
}

export default Tabulation
