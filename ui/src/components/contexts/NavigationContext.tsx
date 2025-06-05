"use client"

import React, { createContext, useContext, useState, ReactNode } from 'react'

export type NavigationItem = 'habits' | 'statistics'

interface NavigationContextType {
  activeItem: NavigationItem
  setActiveItem: (item: NavigationItem) => void
}

const NavigationContext = createContext<NavigationContextType | undefined>(undefined)

interface NavigationProviderProps {
  children: ReactNode
}

export function NavigationProvider({ children }: NavigationProviderProps) {
  const [activeItem, setActiveItem] = useState<NavigationItem>('habits')

  return (
    <NavigationContext.Provider value={{ activeItem, setActiveItem }}>
      {children}
    </NavigationContext.Provider>
  )
}

export function useNavigation() {
  const context = useContext(NavigationContext)
  if (context === undefined) {
    throw new Error('useNavigation must be used within a NavigationProvider')
  }
  return context
} 