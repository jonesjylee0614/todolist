import { useBreakpoints } from '@vueuse/core';

export const BREAKPOINTS = {
  desktop: 1280,
  tablet: 768
};

export const breakpoints = useBreakpoints({
  desktop: `(min-width: ${BREAKPOINTS.desktop}px)`,
  tablet: `(min-width: ${BREAKPOINTS.tablet}px)`
});

export function useDevice() {
  const isDesktop = breakpoints.greater('desktop');
  const isTablet = breakpoints.between('tablet', 'desktop');
  const isMobile = breakpoints.smaller('tablet');
  return {
    isDesktop,
    isTablet,
    isMobile
  };
}

