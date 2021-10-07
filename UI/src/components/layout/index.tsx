import React from 'react'
import MenuTop from '@components/menuTopBar'
import MenuSideBar from '@components/menuSideBar'

type LayoutProps = {
  children: React.ReactNode
}

export const Layout = ({ children }: LayoutProps) => {
  return (
    <>
      <MenuTop />
      <MenuSideBar />
      <div className="layout-content">{children}</div>
    </>
  )
}
