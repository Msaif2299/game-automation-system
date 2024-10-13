import { useRef } from "react"

export function capitalize(value: string): string {
    return value.slice(0,1).toUpperCase() + value.slice(1, value.length)
}

export const useDebounce = (callback: (...args: any[]) => void, delay: number) => {
    const timeoutRef = useRef<number|null>(null)
    const debounce = (...args: any[]) => {
        if (timeoutRef.current) {
            clearTimeout(timeoutRef.current)
        }
        timeoutRef.current = window.setTimeout(() => {
            callback(...args)
        }, delay)
    }
    return debounce
}