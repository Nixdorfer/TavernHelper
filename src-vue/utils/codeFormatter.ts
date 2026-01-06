export function formatHtml(html: string): string {
  if (!html || !html.trim()) return ''
  let formatted = ''
  let indent = 0
  const tab = '  '
  const selfClosingTags = ['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'param', 'source', 'track', 'wbr']
  const inlineTags = ['a', 'span', 'strong', 'em', 'b', 'i', 'u', 'small', 'sub', 'sup', 'code', 'label']
  html = html.replace(/\r\n/g, '\n').replace(/\r/g, '\n')
  html = html.replace(/>\s+</g, '><').trim()
  const tokens: { type: string; content: string; tagName?: string }[] = []
  let i = 0
  while (i < html.length) {
    if (html[i] === '<') {
      let j = i + 1
      while (j < html.length && html[j] !== '>') j++
      const tag = html.slice(i, j + 1)
      const isClosing = tag.startsWith('</')
      const isSelfClosing = tag.endsWith('/>') || selfClosingTags.some(t => tag.match(new RegExp(`^<${t}[\\s/>]`, 'i')))
      const tagMatch = tag.match(/<\/?(\w+)/)
      const tagName = tagMatch ? tagMatch[1].toLowerCase() : ''
      tokens.push({ type: isClosing ? 'close' : (isSelfClosing ? 'self' : 'open'), content: tag, tagName })
      i = j + 1
    } else {
      let j = i
      while (j < html.length && html[j] !== '<') j++
      const text = html.slice(i, j).trim()
      if (text) {
        tokens.push({ type: 'text', content: text })
      }
      i = j
    }
  }
  for (let idx = 0; idx < tokens.length; idx++) {
    const token = tokens[idx]
    if (token.type === 'close') {
      indent = Math.max(0, indent - 1)
      formatted += tab.repeat(indent) + token.content + '\n'
    } else if (token.type === 'open') {
      formatted += tab.repeat(indent) + token.content + '\n'
      indent++
    } else if (token.type === 'self') {
      formatted += tab.repeat(indent) + token.content + '\n'
    } else {
      formatted += tab.repeat(indent) + token.content + '\n'
    }
  }
  return formatted.trim()
}

export function formatCss(css: string): string {
  if (!css || !css.trim()) return ''
  const tab = '  '
  css = css.replace(/\r\n/g, '\n').replace(/\r/g, '\n')
  const comments: string[] = []
  css = css.replace(/\/\*[\s\S]*?\*\//g, match => {
    comments.push(match)
    return `__COMMENT_${comments.length - 1}__`
  })
  css = css.replace(/\s+/g, ' ').trim()
  css = css.replace(/\s*{\s*/g, ' {\n')
  css = css.replace(/\s*}\s*/g, '\n}\n\n')
  css = css.replace(/\s*;\s*/g, ';\n')
  css = css.replace(/\s*:\s*/g, ': ')
  css = css.replace(/__COMMENT_(\d+)__/g, (_, idx) => comments[parseInt(idx)] || '')
  const lines = css.split('\n')
  let formatted = ''
  let indent = 0
  for (let line of lines) {
    line = line.trim()
    if (!line) {
      if (formatted && !formatted.endsWith('\n\n')) {
        formatted += '\n'
      }
      continue
    }
    if (line === '}') {
      indent = Math.max(0, indent - 1)
      formatted += tab.repeat(indent) + line + '\n'
    } else if (line.endsWith('{')) {
      formatted += tab.repeat(indent) + line + '\n'
      indent++
    } else {
      formatted += tab.repeat(indent) + line + '\n'
    }
  }
  formatted = formatted.replace(/\n{3,}/g, '\n\n')
  return formatted.trim()
}
