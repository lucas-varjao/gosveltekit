// frontend/src/lib/utils.ts

import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export type WithElementRef<T, E extends HTMLElement = HTMLElement> = T & {
    ref?: E | null
}

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}
