
import MarkdownIt from 'markdown-it'
let mdItInstance: MarkdownIt | null = null
export const Mdit = (text: string): string => {
    if (mdItInstance === null) {
        mdItInstance = new MarkdownIt()
    }

    return mdItInstance.render(text)
}
