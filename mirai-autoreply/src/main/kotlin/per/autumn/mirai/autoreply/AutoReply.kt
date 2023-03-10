package per.autumn.mirai.autoreply

import net.mamoe.mirai.console.command.Command
import net.mamoe.mirai.console.command.CommandManager.INSTANCE.register
import net.mamoe.mirai.console.command.CommandManager.INSTANCE.unregister
import net.mamoe.mirai.console.plugin.jvm.JvmPluginDescriptionBuilder
import net.mamoe.mirai.console.plugin.jvm.KotlinPlugin
import net.mamoe.mirai.event.GlobalEventChannel
import net.mamoe.mirai.event.events.GroupMessageEvent
import net.mamoe.mirai.event.events.UserMessageEvent
import per.autumn.mirai.autoreply.Config.enabledGroups
import per.autumn.mirai.autoreply.ReplyManager.getResponse
import java.io.File


object AutoReply : KotlinPlugin(
    JvmPluginDescriptionBuilder("per.autumn.mirai.autoreply.plugin", "1.3.2")
        .info("一个简易的自动回复插件")
        .name("auto-reply").build()
) {
    val imgFolder = File(dataFolder, "/img")
    override fun onEnable() {
        Config.reload()
        println("启动auto replay")
        registerCommands()
        println("注册auto replay")
        //添加群组监听
        GlobalEventChannel
            .filterIsInstance<GroupMessageEvent>()
//            .filter { enabledGroups.contains(it.group.id) }
            .subscribeAlways<GroupMessageEvent> {
                println("监听私聊auto replay")
                val content = it.message.contentToString()
                if (content.length >= 5 && content.startsWith("/chat", false)) {
                    it.group.sendMessage(getResponse(it))
                }
            }
        //添加私聊监听
        GlobalEventChannel
            .subscribeAlways<UserMessageEvent> {
                println("监听私聊auto replay")
                val content = it.message.contentToString()
                if (content.length >= 5 && content.startsWith("/chat", false)) {
                    it.subject.sendMessage(getResponse(it))
                }
            }
    }

    override fun onDisable() {
        unregisterCommands()
        Config.save()
    }

    private val commands = getAllCommands()

    private fun registerCommands(commands: List<Command> = AutoReply.commands) {
        commands.forEach {
            it.register()
        }
    }

    private fun unregisterCommands(commands: List<Command> = AutoReply.commands) {
        commands.forEach {
            it.unregister()
        }
    }
}