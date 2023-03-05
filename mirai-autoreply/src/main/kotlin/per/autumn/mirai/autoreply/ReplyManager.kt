package per.autumn.mirai.autoreply

import cn.hutool.json.JSONObject
import net.mamoe.mirai.event.events.GroupMessageEvent
import net.mamoe.mirai.event.events.MessageEvent
import net.mamoe.mirai.message.data.*
import okhttp3.MediaType
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import per.autumn.mirai.autoreply.Config.replyMap
import java.util.concurrent.TimeUnit


/**
 * @author SoundOfAutumn
 * @date 2022/5/7 22:26
 */
object ReplyManager {
    fun hasKeyword(msg: MessageChain): Boolean {
        return hasKeyword(msg.contentToString())
    }

    private fun hasKeyword(s: String): Boolean {
        if (s.startsWith("/")) return false
        for (keyword in replyMap.keys) {
            if (Keyword(keyword).isMatchWith(s)) {
                return true
            }
        }
        return false
    }

    private fun getKeyword(msg: MessageChain): String? {
        return getKeyword(msg.contentToString())
    }

    private fun getKeyword(s: String): String? {
        for (keyword in replyMap.keys) {
            if (Keyword(keyword).isMatchWith(s)) {
                return keyword
            }
        }
        return null
    }

    fun removeByKeyword(text: String) {
        replyMap.remove(text)
    }

    fun getResponse(event: MessageEvent): Message {
        val msg = event.message
        val content = msg.contentToString()

        val okHttpClient = OkHttpClient.Builder()
            .connectTimeout(60, TimeUnit.SECONDS)
            .writeTimeout(60, TimeUnit.SECONDS)
            .readTimeout(60, TimeUnit.SECONDS)
            .build()
        val url = "http://localhost:801/chat?msg=" + content.subSequence(5, content.length);
        var respstr = url?.let { Request.Builder().url(it).get() }
            ?.let { it.header("Content-type", "text/plain") }
            ?.let { it.build() }
            ?.let { okHttpClient.newCall(it).execute() }.body?.string()

        println(respstr)
        val jsonObject = JSONObject(respstr)
        respstr = jsonObject["Msg"] as String?

        val chainBuilder = MessageChainBuilder()
        val response = respstr?.let { Response(it) }
        response?.let { addExtraMessage(event, it, chainBuilder) }
        respstr?.let { chainBuilder.add(it) }

        return chainBuilder.asMessageChain()
    }

    private fun addExtraMessage(event: MessageEvent, response: Response, cb: MessageChainBuilder) {
        if (event is GroupMessageEvent) {
            addExtraMessage(event, response, cb)
        }
    }

    private fun addExtraMessage(event: GroupMessageEvent, response: Response, cb: MessageChainBuilder) {
        if (response.quote) {
            cb.add(QuoteReply(event.message))
        }
        if (response.atAll) {
            cb.add(AtAll)
        }
        if (response.atSender) {
            cb.add(At(event.sender))
        }
    }
}

