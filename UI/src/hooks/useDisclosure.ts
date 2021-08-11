import { useState } from "react";

export const useDisclosure = (
  initialIsOpen: boolean = false
): [boolean, () => void, () => void] => {
  const [isOpen, setIsOpen] = useState(initialIsOpen);

  const close = () => {
    setIsOpen(false);
  };

  const open = () => {
    setIsOpen(true);
  };

  return [isOpen, close, open];
};
