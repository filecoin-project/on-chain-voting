import React from 'react';
export default function Table ({ title = '', list = [] as any, subTitle = '' }) {
  return (
    <table className='min-w-full bg-[#273141] rounded text-left'>
      <thead>
        <tr>
          <th scope='col' colSpan={2}>
            <h2 className='font-normal text-white px-8 py-7 text-2xl border-b border-[#313D4F]'>
              {title}
              {subTitle && (
                <span className='text-[#8896AA] text-xl px-1'> - {subTitle}</span>
              )} 
            </h2>
          </th>
        </tr>
      </thead>
      <tbody className='divide-y divide-[#313D4F]'>
        {list.map((item: any) => (
          <tr key={item.name} className='divide-x divide-[#313D4F]'>
            <td className='w-[280px] whitespace-nowrap py-9 px-8 text-xl text-[#8896AA]'>
              {item.name}
            </td>
            <td className='py-5 px-4 text-xl text-white'>
              {item.comp}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
