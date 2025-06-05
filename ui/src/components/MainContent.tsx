"use client"

import { useNavigation } from "./contexts/NavigationContext";
import { Habits } from "./habits/Habits";
import { Statistics } from "./statistics/Statistics";

export function MainContent() {
  const { activeItem } = useNavigation();
  return (
    <>
      {activeItem === 'habits' && <Habits />}
      {activeItem === 'statistics' && <Statistics />}
    </>
  );
}