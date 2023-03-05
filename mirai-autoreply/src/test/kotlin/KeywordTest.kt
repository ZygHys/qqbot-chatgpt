package per.autumn.mirai.autoreply

import cn.hutool.json.JSON
import cn.hutool.json.JSONObject
import com.google.gson.Gson
import net.mamoe.mirai.message.data.MessageChainBuilder
import okhttp3.MediaType
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import org.junit.jupiter.api.Test

/**
 * @author SoundOfAutumn
 * @date 2022/5/16 10:21
 */
class KeywordTest {
    @Test
    fun `keyword match with test`() {
        println(Keyword("""123${123}456""").isMatchWith("""12123456"""))
    }

    @Test
    fun `get`(){
        val content = "/chat你好"

        val okHttpClient = OkHttpClient()
        val url = "http://localhost:801/chat?msg="+content.subSequence(5, content.length);
        val respstr = url?.let { Request.Builder().url(it).get() }
            ?.let { it.header("Content-type", "application/json") }
            ?.let { it.build() }
            ?.let { okHttpClient.newCall(it).execute() }.body?.string()

        val jsonObject = JSONObject(respstr)
        println(jsonObject["Msg"] as String?)
    }

}