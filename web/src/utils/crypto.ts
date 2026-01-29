/**
 * 密码加密工具
 * 使用RSA公钥加密密码，确保传输安全
 */

// 缓存公钥
let cachedPublicKey: string | null = null

/**
 * 获取RSA公钥
 */
export async function getPublicKey(): Promise<string> {
  if (cachedPublicKey) {
    return cachedPublicKey
  }

  const response = await fetch('/api/v1/auth/public-key')
  const data = await response.json()

  if (data.code === 0 && data.data?.public_key) {
    cachedPublicKey = data.data.public_key
    return cachedPublicKey!
  }

  throw new Error('获取公钥失败')
}

/**
 * 清除缓存的公钥
 */
export function clearPublicKeyCache(): void {
  cachedPublicKey = null
}

/**
 * 将字符串转换为ArrayBuffer
 */
function str2ab(str: string): ArrayBuffer {
  const encoder = new TextEncoder()
  return encoder.encode(str).buffer as ArrayBuffer
}

/**
 * 将ArrayBuffer转换为Base64
 */
function ab2base64(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer)
  let binary = ''
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return btoa(binary)
}

/**
 * 解析PEM格式的公钥
 */
function pemToArrayBuffer(pem: string): ArrayBuffer {
  // 移除PEM头尾和换行符
  const pemContents = pem
    .replace(/-----BEGIN PUBLIC KEY-----/, '')
    .replace(/-----END PUBLIC KEY-----/, '')
    .replace(/\s/g, '')

  // Base64解码
  const binaryString = atob(pemContents)
  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }
  return bytes.buffer
}

/**
 * 导入RSA公钥
 */
async function importPublicKey(pemKey: string): Promise<CryptoKey> {
  const keyData = pemToArrayBuffer(pemKey)

  return await crypto.subtle.importKey(
    'spki',
    keyData,
    {
      name: 'RSA-OAEP',
      hash: 'SHA-256',
    },
    false,
    ['encrypt']
  )
}

/**
 * 使用RSA公钥加密密码
 * @param password 明文密码
 * @returns Base64编码的加密密码
 */
export async function encryptPassword(password: string): Promise<string> {
  try {
    const publicKeyPem = await getPublicKey()
    const publicKey = await importPublicKey(publicKeyPem)

    const passwordBuffer = str2ab(password)
    const encryptedBuffer = await crypto.subtle.encrypt(
      { name: 'RSA-OAEP' },
      publicKey,
      passwordBuffer
    )

    return ab2base64(encryptedBuffer)
  } catch (error) {
    console.error('密码加密失败:', error)
    throw new Error('密码加密失败，请刷新页面重试')
  }
}

/**
 * SHA256哈希（备用方案，当RSA不可用时）
 */
export async function hashPassword(password: string): Promise<string> {
  const encoder = new TextEncoder()
  const data = encoder.encode(password)
  const hashBuffer = await crypto.subtle.digest('SHA-256', data)
  return ab2base64(hashBuffer)
}
