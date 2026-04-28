import crypto from 'crypto'

export function generateRandomKey(prefix = 'limity_'): string {
  const randomPart = crypto.randomBytes(32).toString('hex')
  return `${prefix}${randomPart}`
}

export function hashKey(key: string): string {
  return crypto.createHash('sha256').update(key).digest('hex')
}
