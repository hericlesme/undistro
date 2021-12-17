import { useMemo } from 'react'
import { PAGINATION_DOTS } from '@/helpers/constants'
import { range } from '@/helpers/pagination'

export interface usePaginationProps {
  totalCount: number
  pageSize: number
  currentPage: number
  siblingCount?: number
}

const usePagination = ({ totalCount, pageSize, siblingCount = 1, currentPage }: usePaginationProps) => {
  const paginationRange = useMemo(() => {
    const totalPageCount = Math.ceil(totalCount / pageSize)
    const totalPageNumbers = siblingCount + 5

    if (totalPageNumbers > totalPageCount + 2) {
      return range(1, totalPageCount)
    }

    const leftSiblingIndex = Math.max(currentPage - siblingCount, 1)
    const rightSiblingIndex = Math.min(currentPage + siblingCount, totalPageCount)

    const shouldShowLeftDots = leftSiblingIndex > 1
    const shouldShowRightDots = rightSiblingIndex < totalPageCount

    const firstPageIndex = 1
    const lastPageIndex = totalPageCount

    if (!shouldShowLeftDots && shouldShowRightDots) {
      let leftRange = range(Math.max(1, currentPage - 1), currentPage + 1)
      return [...leftRange, PAGINATION_DOTS, totalPageCount]
    }

    if (shouldShowLeftDots && !shouldShowRightDots) {
      let rightRange = range(currentPage - 1, totalPageCount)
      return [firstPageIndex, PAGINATION_DOTS, ...rightRange]
    }

    if (shouldShowLeftDots && shouldShowRightDots) {
      let middleRange = range(leftSiblingIndex, rightSiblingIndex)
      return [firstPageIndex, PAGINATION_DOTS, ...middleRange, PAGINATION_DOTS, lastPageIndex]
    }
    return []
  }, [totalCount, pageSize, siblingCount, currentPage])

  return paginationRange
}

export { usePagination, PAGINATION_DOTS }
