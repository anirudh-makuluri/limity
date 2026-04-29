import crypto from 'crypto'
import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function generateRandomKey(prefix = 'limity_'): string {
  const randomPart = crypto.randomBytes(32).toString('hex')
  return `${prefix}${randomPart}`
}

export function hashKey(key: string): string {
  return crypto.createHash('sha256').update(key).digest('hex')
}
