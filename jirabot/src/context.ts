import Bot from 'keybase-bot'
import {Issue} from './jira'
import {CommentMessage} from './message'
import util from 'util'
import * as BotConfig from './bot-config'
import * as Jira from './jira'
import Aliases from './aliases'
import Configs from './configs'
import logger from './logger'

const setTimeoutPromise = util.promisify(setTimeout)

type CommentContextItem = {
  message: CommentMessage
  issues: Array<Issue>
}

class CommentContext {
  _respMsgIDToCommentMessage = new Map()

  add = (responseID: number, message: CommentMessage, issues: Array<Issue>) => {
    this._respMsgIDToCommentMessage.set(responseID, {message, issues})
    setTimeoutPromise(1000 * 120 /* 2min */).then(() =>
      this._respMsgIDToCommentMessage.delete(responseID)
    )
  }

  get = (responseID: number): null | CommentContextItem =>
    this._respMsgIDToCommentMessage.get(responseID)
}

export type Context = {
  aliases: Aliases
  bot: Bot
  botConfig: BotConfig.BotConfig
  comment: CommentContext
  configs: Configs
  getJiraFromTeamnameAndUsername: typeof Jira.getJiraFromTeamnameAndUsername
}

export const init = (botConfig: BotConfig.BotConfig): Promise<Context> => {
  var bot = new Bot()
  const context = {
    aliases: new Aliases({}),
    bot,
    botConfig,
    comment: new CommentContext(),
    configs: new Configs(bot, botConfig),
    getJiraFromTeamnameAndUsername: Jira.getJiraFromTeamnameAndUsername,
  }
  return context.bot
    .init(
      context.botConfig.keybase.username,
      context.botConfig.keybase.paperkey,
      {
        autoLogSendOnExit: true,
        verbose: true,
      }
    )
    .then(() => {
      logger.info('init done')
      return context
    })
}
