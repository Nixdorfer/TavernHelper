export async function copyToClipboard(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(text)
      return true
    }
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)
    textarea.select()
    const result = document.execCommand('copy')
    document.body.removeChild(textarea)
    return result
  } catch {
    return false
  }
}

export async function readFromClipboard(): Promise<string | null> {
  try {
    if (navigator.clipboard) {
      return await navigator.clipboard.readText()
    }
    return null
  } catch {
    return null
  }
}

export async function copyImageToClipboard(blob: Blob): Promise<boolean> {
  try {
    if (!navigator.clipboard) return false
    await navigator.clipboard.write([
      new ClipboardItem({ [blob.type]: blob })
    ])
    return true
  } catch {
    return false
  }
}
